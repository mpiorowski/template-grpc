package users

import (
	"context"
	"fmt"
	pb "powerit/proto"
	"powerit/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/metadata"
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
func Auth(ctx context.Context) (*pb.User, string, error) {
	claims, err := extractToken(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("extractToken: %w", err)
	}
	// get oauth token from redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.REDIS_URL,
		Password: utils.REDIS_PASSWORD,
	})
	userId, err := rdb.Get(context.Background(), claims.Id).Result()
	if err != nil {
		return nil, "", fmt.Errorf("rdb.Get: %w", err)
	}
	// get user from database
	user, err := selectUserById(userId)
	if err != nil {
		return nil, "", fmt.Errorf("selectUserById: %w", err)
	}
	// create new phantom token with a 7 day expiration
	tokenId, err := uuid.NewV7()
	if err != nil {
		return nil, "", fmt.Errorf("uuid.NewV7: %w", err)
	}
	err = rdb.Set(context.Background(), tokenId.String(), userId, 7*24*time.Hour).Err()
	if err != nil {
		return nil, "", fmt.Errorf("rdb.Set: %w", err)
	}
	subscribed := checkIfSubscribed(user)
	user.SubscriptionActive = subscribed
	return user, tokenId.String(), nil
}

func GetUser(ctx context.Context) (*pb.User, error) {
	claims, err := extractToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("extractToken: %w", err)
	}
	user, err := selectUserById(claims.Id)
	if err != nil {
		return nil, fmt.Errorf("selectUserById: %w", err)
	}
	subscribed := checkIfSubscribed(user)
	user.SubscriptionActive = subscribed
	return user, nil
}

type Claims struct {
	Id string
}

func extractToken(ctx context.Context) (Claims, error) {
	claims := Claims{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return claims, fmt.Errorf("Missing context metadata")
	}

	token := md.Get("x-authorization")
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
