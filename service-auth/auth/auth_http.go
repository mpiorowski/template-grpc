package auth

import (
	"context"
	"log/slog"
	"net/http"
	"service-auth/system"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

/**
 *  Oauth login
 *  @api {get} /oauth-login/:provider Oauth login
 */
func OauthLogin(c echo.Context, storage system.Storage) error {
	defer system.Perf("oauth_login", time.Now())
	var authDB = NewAuthDB(&storage)
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

/**
 *  Oauth callback
 *  @api {get} /oauth-callback/:provider Oauth callback
 */
func OauthCallback(c echo.Context, storage system.Storage) error {
	defer system.Perf("oauth_callback", time.Now())
	var authDB = NewAuthDB(&storage)
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
	return c.Redirect(http.StatusTemporaryRedirect, system.CLIENT_URL+"/token/"+token.Id)
}
