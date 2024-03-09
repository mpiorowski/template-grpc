package auth

import (
	"context"
	"fmt"
	pb "service-auth/proto"
	"service-auth/system"
)

func getUser(ctx context.Context, storage system.Storage) (*pb.User, error) {
	var authDB = newAuthDB(&storage)
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

