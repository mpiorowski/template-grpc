package auth

import (
	"context"
	"fmt"
	"log/slog"
	pb "service-auth/proto"
	"service-auth/system"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v76"
	portal_session "github.com/stripe/stripe-go/v76/billingportal/session"
	checkout_session "github.com/stripe/stripe-go/v76/checkout/session"
)

/**
 * 1. Extract phantom token from context
 * 2. Using it's id, get oauth token from database
 * 3. Check if oauth token is valid
 * 4. Refresh oauth token if it's expired
 * 5. Create new phantom token
 * 6. Get user from database
 * 7. Return user and new phantom token
 */
func Auth(ctx context.Context, storage system.Storage) (*pb.User, string, error) {
	var authDB = NewAuthDB(&storage)
	claims, err := system.ExtractToken(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("extractToken: %w", err)
	}
	// get oauth token from redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     system.REDIS_URL,
		Password: system.REDIS_PASSWORD,
	})
	userId, err := rdb.Get(context.Background(), claims.Id).Result()
	if err != nil {
		return nil, "", fmt.Errorf("rdb.Get: %w", err)
	}
	// get user from database
	user, err := authDB.selectUserById(userId)
	if err != nil {
		return nil, "", fmt.Errorf("selectUserById: %w", err)
	}
	// create new phantom token with a 7 day expiration
	tokenId, err := uuid.NewV7()
	if err != nil {
		return nil, "", fmt.Errorf("uuid.NewV7: %w", err)
	}
	go func() {
		err = rdb.Set(context.Background(), tokenId.String(), userId, 7*24*time.Hour).Err()
		if err != nil {
			slog.Error("Error setting token in redis", "rdb.Set", err)
		}
	}()
	subscribed := checkIfSubscribed(user, authDB)
	user.SubscriptionActive = subscribed
	return user, tokenId.String(), nil
}

func GetUser(ctx context.Context, storage system.Storage) (*pb.User, error) {
	var authDB = NewAuthDB(&storage)
	claims, err := system.ExtractToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("extractToken: %w", err)
	}
	user, err := authDB.selectUserById(claims.Id)
	if err != nil {
		return nil, fmt.Errorf("selectUserById: %w", err)
	}
	subscribed := checkIfSubscribed(user, authDB)
	user.SubscriptionActive = subscribed
	return user, nil
}

func CreateStripeCheckout(user *pb.User, storage system.Storage) (string, error) {
	var authDB = NewAuthDB(&storage)
	customerId := user.SubscriptionId
	if customerId == "" {
		var err error
		customerId, err = createStripeUser(user.Id, user.Email, authDB)
		if err != nil {
			return "", fmt.Errorf("createStripeUser: %w", err)
		}
	}

	stripe.Key = system.STRIPE_API_KEY

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(system.CLIENT_URL + "/billing?success"),
		CancelURL:  stripe.String(system.CLIENT_URL + "/billing?cancel"),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String("subscription"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(system.STRIPE_PRICE_ID),
				Quantity: stripe.Int64(1),
			},
		},
		Customer: stripe.String(customerId),
	}

	session, err := checkout_session.New(params)
	if err != nil {
		return "", err
	}

	err = authDB.updateSubscriptionCheck(user.Id, "1970-01-01T00:00:00Z")
	if err != nil {
		slog.Error("Error updating subscription check date", "updateSubscriptionCheck", err)
	}
	return session.URL, nil
}

func CreateStripePortal(user *pb.User) (string, error) {
	stripe.Key = system.STRIPE_API_KEY

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(user.SubscriptionId),
		ReturnURL: stripe.String(system.CLIENT_URL + "/billing"),
	}
	session, err := portal_session.New(params)
	if err != nil {
		return "", err
	}

	return session.URL, nil
}
