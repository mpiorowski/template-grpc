package users

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/stripe/stripe-go/v76"
	portal_session "github.com/stripe/stripe-go/v76/billingportal/session"
	checkout_session "github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/subscription"

	pb "powerit/proto"
	"powerit/system"
)

func checkIfSubscribed(user *pb.User) bool {
	if user.SubscriptionId == "" {
		return false
	}
	subEnd, _ := time.Parse(time.RFC3339, user.SubscriptionEnd)
	if subEnd.After(time.Now()) {
		return true
	}
	subCheck, _ := time.Parse(time.RFC3339, user.SubscriptionCheck)
	if subCheck.After(time.Now()) {
		return false
	}

	stripe.Key = system.STRIPE_API_KEY

	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(user.SubscriptionId),
		Status:   stripe.String("active"),
	}

	i := subscription.List(params)
	for i.Next() {
		if i.Subscription().Status == "active" {
			// get the subscription end date
			subEnd := time.Unix(i.Subscription().CurrentPeriodEnd, 0).Format(time.RFC3339)
			// update the user's subscription end date
			err := updateSubscriptionEnd(user.Id, subEnd)
			if err != nil {
				slog.Error("Error updating subscription end date", "updateSubscriptionEnd", err)
				return false
			}
			return true
		}
	}
	err := updateSubscriptionCheck(user.Id, time.Now().Add(time.Hour).Format(time.RFC3339))
	if err != nil {
		slog.Error("Error updating subscription check date", "updateSubscriptionCheck", err)
	}
	return false
}

func createStripeUser(userId string, email string) (string, error) {
	stripe.Key = system.STRIPE_API_KEY

	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	customer, err := customer.New(params)
	if err != nil {
		return "", fmt.Errorf("customer.New: %w", err)
	}
	err = updateSubscriptionId(userId, customer.ID)
	if err != nil {
		return "", fmt.Errorf("updateSubscriptionId: %w", err)
	}
	return customer.ID, nil
}

func CreateStripeCheckout(user *pb.User) (string, error) {
	customerId := user.SubscriptionId
	if customerId == "" {
		var err error
		customerId, err = createStripeUser(user.Id, user.Email)
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

	err = updateSubscriptionCheck(user.Id, "1970-01-01T00:00:00Z")
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
