package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"service-auth/system"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v76"
	portal_session "github.com/stripe/stripe-go/v76/billingportal/session"
	checkout_session "github.com/stripe/stripe-go/v76/checkout/session"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "service-auth/proto"
)

type AuthService interface {
	// gRPC
	Auth(ctx context.Context, storage system.Storage) (*pb.AuthResponse, error)
	CreateStripeCheckout(ctx context.Context, storage system.Storage) (*pb.StripeUrlResponse, error)
	CreateStripePortal(ctx context.Context, storage system.Storage) (*pb.StripeUrlResponse, error)
	// HTTP
	OauthLogin(c echo.Context, storage system.Storage) error
	OauthCallback(c echo.Context, storage system.Storage) error
    // Task
    CleanTokens(ctx context.Context, storage system.Storage) error
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

/**
 * 1. Extract phantom token from context
 * 2. Using it's id, get oauth token from database
 * 3. Check if oauth token is valid
 * 4. Refresh oauth token if it's expired
 * 5. Create new phantom token
 * 6. Get user from database
 * 7. Return user and new phantom token
 */
func (a *authService) Auth(ctx context.Context, storage system.Storage) (*pb.AuthResponse, error) {
	defer system.Perf("auth", time.Now())
	var authDB = newAuthDB(&storage)
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

func (a *authService) CreateStripeCheckout(ctx context.Context, storage system.Storage) (*pb.StripeUrlResponse, error) {
	defer system.Perf("CreateStripeCheckout", time.Now())
	user, err := getUser(ctx, storage)
	if err != nil {
		slog.Error("Error authorizing user", "auth.GetUser", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	var authDB = newAuthDB(&storage)
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

func (a *authService) CreateStripePortal(ctx context.Context, storage system.Storage) (*pb.StripeUrlResponse, error) {
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

func (a *authService) OauthLogin(c echo.Context, storage system.Storage) error {
	defer system.Perf("oauth_login", time.Now())
	var authDB = newAuthDB(&storage)
    var OAuthConfig = newOAuthConfig()
	config, err := OAuthConfig.getOAuthConfig(c.Param("provider"))
	if err != nil {
		slog.Error("Error getting provider", "getOAuthConfig", err)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}

	// generate random state and verifier
	state := system.GenerateRandomState(32)
	verifier := oauth2.GenerateVerifier()
	// store state and verifier
	_, err = authDB.insertToken(time.Now().Add(10*time.Second).Format(time.RFC3339), "", state, verifier)
	if err != nil {
		slog.Error("Error inserting token", "insertToken", err)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}
	// Redirect user to consent page to ask for permission
	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *authService) OauthCallback(c echo.Context, storage system.Storage) error {
	defer system.Perf("oauth_callback", time.Now())
	var authDB = newAuthDB(&storage)
    var OAuthConfig = newOAuthConfig()
	provider := c.Param("provider")
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	// get verifier from state
	token, err := authDB.seleteTokenByState(state)
	if err != nil {
		slog.Error("Error getting token by state", "seleteTokenByState", err)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}
	expires, err := time.Parse(time.RFC3339, token.Expires)
	if err != nil || time.Now().After(expires) {
		slog.Error("Token expired", "token.Expires", token.Expires)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}

	// get oauth config
	config, err := OAuthConfig.getOAuthConfig(provider)
	if err != nil {
		slog.Error("Error getting provider", "getOAuthConfig", err)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}

	// get oauth token
	oauthToken, err := config.Exchange(context.Background(), code, oauth2.VerifierOption(token.Verifier))
	if err != nil {
		slog.Error("Error exchanging code for token", "config.Exchange", err)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}

	// fetch user info from github
	userInfo, err := OAuthConfig.getUserInfo(provider, oauthToken.AccessToken)
	if err != nil {
		slog.Error("Error fetching user info", "configProvider.getUserInfo", err)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}

	// get user, create if not exists
	user, err := authDB.selectUserByEmailAndSub(userInfo.email, userInfo.sub)
	if err != nil {
		user, err = authDB.insertUser(userInfo.email, userInfo.sub, userInfo.avatar)
		if err != nil {
			slog.Error("Error inserting user", "insertUser", err)
			return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=invalid_user")
		}
	}

	// create oauth token with a 10 seconds expiration
	token, err = authDB.insertToken(time.Now().Add(10*time.Second).Format(time.RFC3339), user.Id, "", "")
	if err != nil {
		slog.Error("Error inserting token", "insertToken", err)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}

	// redirect to home page
	return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/token/"+token.Id)
}

func (a *authService) CleanTokens(ctx context.Context, storage system.Storage) error {
	authDb := newAuthDB(&storage)
	err := authDb.CleanTokens()
	if err != nil {
		return fmt.Errorf("cleanTokens: %w", err)
	}
	return nil
}
