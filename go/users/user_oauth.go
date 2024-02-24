package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"powerit/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserInfo struct {
	email  string
	sub    string
	avatar string
}

var githubOAuthConfig = oauth2.Config{
	ClientID:     utils.GITHUB_CLIENT_ID,
	ClientSecret: utils.GITHUB_CLIENT_SECRET,
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
	RedirectURL: utils.SERVER_HTTP + "/oauth-callback/github",
	Scopes:      []string{"user:email"},
}

var googleOAuthConfig = oauth2.Config{
	ClientID:     utils.GOOGLE_CLIENT_ID,
	ClientSecret: utils.GOOGLE_CLIENT_SECRET,
	Endpoint:     google.Endpoint,
	RedirectURL:  utils.SERVER_HTTP + "/oauth-callback/google",
	Scopes:       []string{"profile", "email", "openid"},
}

type OAuthConfigProvider interface {
	getOAuthConfig(provider string) (*oauth2.Config, error)
	getUserInfo(provider string, accessToken string) (*UserInfo, error)
}
type OAuthConfigImpl struct {}
var OAuthConfig OAuthConfigProvider = &OAuthConfigImpl{}

func (o OAuthConfigImpl) getOAuthConfig(provider string) (*oauth2.Config, error) {
	if provider == "github" {
		return &githubOAuthConfig, nil
	} else if provider == "google" {
		return &googleOAuthConfig, nil
	}
	return nil, fmt.Errorf("Invalid provider")
}

func (o *OAuthConfigImpl) getUserInfo(provider string, accessToken string) (*UserInfo, error) {
	var url string
	if provider == "github" {
		url = "https://api.github.com/user"
	} else if provider == "google" {
		url = "https://www.googleapis.com/oauth2/v2/userinfo"
	} else {
		return nil, fmt.Errorf("Invalid provider")
	}

	// Create a GET request to fetch user information from GitHub
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	// Set the "Authorization" header with the access token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var userInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, fmt.Errorf("json.NewDecoder: %w", err)
	}

	var sub string
	var email string
	var avatar string
	var ok bool
	if provider == "github" {
		userId, ok := userInfo["id"].(float64)
		if !ok {
			return nil, fmt.Errorf("Invalid user id")
		}
		sub = fmt.Sprintf("%.0f", userId)
		email, ok = userInfo["email"].(string)
		if !ok {
			email = ""
		}
		avatar, ok = userInfo["avatar_url"].(string)
		if !ok {
			avatar = ""
		}

	} else if provider == "google" {
		sub, ok = userInfo["id"].(string)
		if !ok {
			return nil, fmt.Errorf("Invalid user id")
		}
		email, ok = userInfo["email"].(string)
		if !ok {
			email = ""
		}
		avatar, ok = userInfo["picture"].(string)
		if !ok {
			avatar = ""
		}
	} else {
		return nil, fmt.Errorf("Invalid provider")
	}

	return &UserInfo{
		email:  email,
		sub:    sub,
		avatar: avatar,
	}, nil
}
