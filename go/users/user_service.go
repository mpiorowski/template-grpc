package users

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"powerit/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
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
func Auth(c echo.Context) error {
	claims, err := extractToken(c)
	if err != nil {
		slog.Error("Error extracting token", "extractToken", err)
		return c.JSON(401, "Unauthorized")
	}
	// get oauth token from redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.REDIS_URL,
		Password: utils.REDIS_PASSWORD,
	})
	userId, err := rdb.Get(context.Background(), claims.Id).Result()
	if err != nil {
		slog.Error("Error getting user id from redis", "rdb.Get", err)
		return c.JSON(401, "Unauthorized")
	}
	// get user from database
	user, err := selectUserById(userId)
	if err != nil {
		slog.Error("Error getting user from database", "selectUserById", err)
		return c.JSON(500, "Error getting user from database")
	}
	// create new phantom token with a 7 day expiration
	tokenId, err := uuid.NewV7()
	if err != nil {
		slog.Error("Error creating new token id", "uuid.NewV7", err)
		return c.JSON(500, "Error creating new token id")
	}
	err = rdb.Set(context.Background(), tokenId.String(), userId, 7*24*time.Hour).Err()
	if err != nil {
		slog.Error("Error setting new token id in redis", "rdb.Set", err)
		return c.JSON(500, "Error setting new token id in redis")
	}
	subscribed := checkIfSubscribed(user)
	user.SubscriptionActive = subscribed
	return c.JSON(200, map[string]interface{}{
		"user":  user,
		"token": tokenId.String(),
	})
}

/**
 *  Oauth login
 *  @api {get} /oauth-login/:provider Oauth login
 */
func OauthLogin(c echo.Context) error {
	config, err := OAuthConfig.getOAuthConfig(c.Param("provider"))
	if err != nil {
		slog.Error("Error getting provider", "configProvider.getOAuthConfig", err)
		return c.Redirect(http.StatusTemporaryRedirect, utils.CLIENT_URL+"/auth?error=unauthorized")
	}

	// generate random state and verifier
	state := utils.GenerateRandomState(32)
	verifier := oauth2.GenerateVerifier()
	// store state and verifier
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.REDIS_URL,
		Password: utils.REDIS_PASSWORD,
	})
	err = rdb.Set(context.Background(), state, verifier, 5*time.Minute).Err()
	if err != nil {
		slog.Error("Error setting state and verifier", "rdb.Set", err)
		return c.Redirect(http.StatusTemporaryRedirect, utils.CLIENT_URL+"/auth?error=unauthorized")
	}
	// Redirect user to consent page to ask for permission
	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

/**
 *  Oauth callback
 *  @api {get} /oauth-callback/:provider Oauth callback
 */
func OauthCallback(c echo.Context) error {
	provider := c.Param("provider")
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	// get verifier from state
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.REDIS_URL,
		Password: utils.REDIS_PASSWORD,
	})
	verifier, err := rdb.Get(context.Background(), state).Result()
	if err != nil {
		slog.Error("Error getting verifier", "rdb.Get", err)
		return c.Redirect(http.StatusTemporaryRedirect, utils.CLIENT_URL+"/auth?error=unauthorized")
	}

	// get oauth config
	config, err := OAuthConfig.getOAuthConfig(provider)
	if err != nil {
		slog.Error("Error getting provider", "configProvider.getOAuthConfig", err)
		return c.Redirect(http.StatusTemporaryRedirect, utils.CLIENT_URL+"/auth?error=unauthorized")
	}

	// get oauth token
	oauthToken, err := config.Exchange(context.Background(), code, oauth2.VerifierOption(verifier))
	if err != nil {
		slog.Error("Error exchanging code for token", "config.Exchange", err)
		return c.Redirect(http.StatusTemporaryRedirect, utils.CLIENT_URL+"/auth?error=unauthorized")
	}

	// fetch user info from github
	userInfo, err := OAuthConfig.getUserInfo(provider, oauthToken.AccessToken)
	if err != nil {
		slog.Error("Error fetching user info", "configProvider.getUserInfo", err)
		return c.Redirect(http.StatusTemporaryRedirect, utils.CLIENT_URL+"/auth?error=unauthorized")
	}

	// get user, create if not exists
	user, err := selectUserByEmailAndSub(userInfo.email, userInfo.sub)
	if err != nil {
		user, err = insertUser(userInfo.email, userInfo.sub, userInfo.avatar)
		if err != nil {
			slog.Error("Error inserting user", "insertUser", err)
			return c.Redirect(http.StatusTemporaryRedirect, utils.CLIENT_URL+"/auth?error=invalid_user")
		}
	}

	// create oauth token with a 7 day expiration
	tokenId, err := uuid.NewV7()
	if err != nil {
		slog.Error("Error creating token id", "uuid.NewV7", err)
		return c.Redirect(http.StatusTemporaryRedirect, utils.CLIENT_URL+"/auth?error=unauthorized")
	}
	err = rdb.Set(context.Background(), tokenId.String(), user.Id, 7*24*time.Hour).Err()
	if err != nil {
		slog.Error("Error setting token id", "rdb.Set", err)
		return c.Redirect(http.StatusTemporaryRedirect, utils.CLIENT_URL+"/auth?error=unauthorized")
	}

	// set cookie
	cookie := &http.Cookie{}
	cookie.Domain = utils.COOKIE_DOMAIN
	cookie.Name = "token"
	cookie.Value = tokenId.String()
	cookie.Path = "/"
	cookie.Secure = true
	cookie.SameSite = http.SameSiteLaxMode
	cookie.HttpOnly = true
	// 7 days
	cookie.MaxAge = 7 * 24 * 60 * 60
	c.SetCookie(cookie)

	// redirect to home page
	return c.Redirect(http.StatusTemporaryRedirect, utils.CLIENT_URL)
}

type Claims struct {
	Id string
}

func extractToken(c echo.Context) (Claims, error) {
	claims := Claims{}
	token := c.Request().Header["Authorization"]
	if len(token) == 0 {
		return claims, fmt.Errorf("Missing authorization header")
	}

	// Validate the token
	tokenParts := strings.SplitN(token[0], " ", 2)
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return claims, fmt.Errorf("Invalid authorization header")
	}

	// Decode the token
	tokenString := tokenParts[1]
	jwtClaims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(tokenString, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.JWT_SECRET), nil
	})
	if err != nil {
		return claims, fmt.Errorf("jwt.ParseWithClaims: %w", err)
	}
	if !t.Valid {
		return claims, fmt.Errorf("Invalid token")
	}
	if _, ok := jwtClaims["id"]; !ok {
		return claims, fmt.Errorf("Missing id in token")
	}

	claims.Id = jwtClaims["id"].(string)
	return claims, nil
}

// TODO: not used
// func oauthRefresh(token Token, configProvider OAuthConfigProvider) (*oauth2.Token, error) {
// 	oauthToken := oauth2.Token{
// 		AccessToken:  token.AccessToken,
// 		RefreshToken: token.RefreshToken,
// 		TokenType:    token.TokenType,
// 		Expiry:       token.Expires,
// 	}
//
// 	config, err := configProvider.getOAuthConfig(token.Provider)
// 	if err != nil {
// 		return nil, fmt.Errorf("getOAuthConfig: %w", err)
// 	}
//
// 	newOauthToken, err := config.TokenSource(context.Background(), &oauthToken).Token()
// 	if err != nil {
// 		return nil, fmt.Errorf("config.TokenSource: %w", err)
// 	}
// 	return newOauthToken, nil
// }
