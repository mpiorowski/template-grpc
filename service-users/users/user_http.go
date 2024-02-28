package users

import (
	"context"
	"log/slog"
	"net/http"
	"powerit/system"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
)

/**
 *  Oauth login
 *  @api {get} /oauth-login/:provider Oauth login
 */
func OauthLogin(c echo.Context) error {
	config, err := OAuthConfig.getOAuthConfig(c.Param("provider"))
	if err != nil {
		slog.Error("Error getting provider", "configProvider.getOAuthConfig", err)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}

	// generate random state and verifier
	state := system.GenerateRandomState(32)
	verifier := oauth2.GenerateVerifier()
	// store state and verifier
	rdb := redis.NewClient(&redis.Options{
		Addr:     system.REDIS_URL,
        Password: system.REDIS_PASSWORD,
	})
	err = rdb.Set(context.Background(), state, verifier, 5*time.Minute).Err()
    if err != nil {
        slog.Error("Error setting state and verifier", "rdb.Set", err)
        return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
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
		Addr:     system.REDIS_URL,
        Password: system.REDIS_PASSWORD,
	})
	verifier, err := rdb.Get(context.Background(), state).Result()
	if err != nil {
		slog.Error("Error getting verifier", "rdb.Get", err)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}

	// get oauth config
	config, err := OAuthConfig.getOAuthConfig(provider)
	if err != nil {
		slog.Error("Error getting provider", "configProvider.getOAuthConfig", err)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}

	// get oauth token
	oauthToken, err := config.Exchange(context.Background(), code, oauth2.VerifierOption(verifier))
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
	user, err := selectUserByEmailAndSub(userInfo.email, userInfo.sub)
	if err != nil {
		user, err = insertUser(userInfo.email, userInfo.sub, userInfo.avatar)
		if err != nil {
			slog.Error("Error inserting user", "insertUser", err)
			return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=invalid_user")
		}
	}

	// create oauth token with a 10 seconds expiration
	tokenId, err := uuid.NewV7()
	if err != nil {
		slog.Error("Error creating token id", "uuid.NewV7", err)
		return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
	}
	err = rdb.Set(context.Background(), tokenId.String(), user.Id, 10*time.Second).Err()
    if err != nil {
        slog.Error("Error setting token id", "rdb.Set", err)
        return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/auth?error=unauthorized")
    }

	// set cookie
	// cookie := &http.Cookie{}
	// cookie.Domain = system.COOKIE_DOMAIN
	// cookie.Name = "token"
	// cookie.Value = tokenId.String()
	// cookie.Path = "/"
	// cookie.Secure = true
	// cookie.SameSite = http.SameSiteLaxMode
	// cookie.HttpOnly = true
	// // 7 days
	// cookie.MaxAge = 7 * 24 * 60 * 60
	// c.SetCookie(cookie)

	// redirect to home page
	return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/token/"+tokenId.String())
}
