package auth

import (
	"context"
	"log/slog"
	pb "service-auth/proto"
	"service-auth/system"
	"time"

	"github.com/stripe/stripe-go/v76"
	portal_session "github.com/stripe/stripe-go/v76/billingportal/session"
	checkout_session "github.com/stripe/stripe-go/v76/checkout/session"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
func Auth(ctx context.Context, storage system.Storage) (*pb.AuthResponse, error) {
	defer system.Perf("auth", time.Now())
	var authDB = NewAuthDB(&storage)
	claims, err := system.ExtractToken(ctx)
	if err != nil {
		slog.Error("Error extracting token", "system.ExtractToken", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	// get token
	token, err := authDB.selectTokenById(claims.Id)
	if err != nil {
		slog.Error("Error selecting token by id", "authDB.selectTokenById", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	expires, err := time.Parse(time.RFC3339, token.Expires)
	if err != nil || time.Now().After(expires) {
		slog.Error("Token expired", "token.Expires", token.Expires)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	// get user from database
	user, err := authDB.selectUserById(token.UserId)
	if err != nil {
		slog.Error("Error selecting user by id", "authDB.selectUserById", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	// create new phantom token with a 7 day expiration
	// tokenId, err := uuid.NewV7()
	// if err != nil {
	// 	slog.Error("Error creating new token", "uuid.NewV7", err)
	// 	return nil, status.Error(codes.Internal, "Internal error")
	// }
	// go func() {
	// 	err = rdb.Set(context.Background(), tokenId.String(), userId, 7*24*time.Hour).Err()
	// 	if err != nil {
	// 		slog.Error("Error setting token in redis", "rdb.Set", err)
	// 	}
	// }()

	// Update token expiration and user
	go func() {
		err = authDB.updateToken(claims.Id, time.Now().Add(7*24*time.Hour).Format(time.RFC3339))
		if err != nil {
			slog.Error("Error updating token", "authDB.updateToken", err)
		}
		err = authDB.updateUser(user.Id)
		if err != nil {
			slog.Error("Error updating user", "authDB.updateUser", err)
		}
	}()

	subscribed := checkIfSubscribed(user, authDB)
	user.SubscriptionActive = subscribed
	return &pb.AuthResponse{
		Token: claims.Id,
		User:  user,
	}, nil
}

func CreateStripeCheckout(ctx context.Context, storage system.Storage) (*pb.StripeUrlResponse, error) {
	defer system.Perf("CreateStripeCheckout", time.Now())
	user, err := getUser(ctx, storage)
	if err != nil {
		slog.Error("Error authorizing user", "auth.GetUser", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	var authDB = NewAuthDB(&storage)
	customerId := user.SubscriptionId
	if customerId == "" {
		var err error
		customerId, err = createStripeUser(user.Id, user.Email, authDB)
		if err != nil {
			slog.Error("Error creating stripe user", "createStripeUser", err)
			return nil, status.Error(codes.Internal, "Internal error")
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
		slog.Error("Error creating stripe checkout", "checkout_session.New", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	err = authDB.updateSubscriptionCheck(user.Id, "1970-01-01T00:00:00Z")
	if err != nil {
		slog.Error("Error updating subscription check date", "updateSubscriptionCheck", err)
	}
	return &pb.StripeUrlResponse{
		Url: session.URL,
	}, nil
}

func CreateStripePortal(ctx context.Context, storage system.Storage) (*pb.StripeUrlResponse, error) {
	defer system.Perf("CreateStripePortal", time.Now())
	user, err := getUser(ctx, storage)
	if err != nil {
		slog.Error("Error authorizing user", "auth.GetUser", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	stripe.Key = system.STRIPE_API_KEY

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(user.SubscriptionId),
		ReturnURL: stripe.String(system.CLIENT_URL + "/billing"),
	}
	session, err := portal_session.New(params)
	if err != nil {
		slog.Error("Error creating stripe portal", "portal_session.New", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &pb.StripeUrlResponse{
		Url: session.URL,
	}, nil
}
