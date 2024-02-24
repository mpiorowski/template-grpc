package main

import (
	"context"
	"log/slog"
	"time"

	pb "powerit/proto"
	"powerit/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Auth(ctx context.Context, in *pb.Empty) (*pb.AuthResponse, error) {
	start := time.Now()
	user, token, err := users.Auth(ctx)
	if err != nil {
		slog.Error("Error authorizing user", "users.UserAuth", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	slog.Info("auth", "time", time.Since(start))
	return &pb.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *server) CreateStripeCheckout(ctx context.Context, in *pb.Empty) (*pb.StripeUrlResponse, error) {
	start := time.Now()
	user, err := users.GetUser(ctx)
	if err != nil {
		slog.Error("Error authorizing user", "users.UserAuth", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	// Create a new checkout session
	url, err := users.CreateStripeCheckout(user)
	if err != nil {
		slog.Error("Error creating stripe checkout", "users.CreateStripeCheckout", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	slog.Info("create_stripe_checkout", "time", time.Since(start))
	return &pb.StripeUrlResponse{
		Url: url,
	}, nil
}

func (s *server) CreateStripePortal(ctx context.Context, in *pb.Empty) (*pb.StripeUrlResponse, error) {
	start := time.Now()
	user, err := users.GetUser(ctx)
	if err != nil {
		slog.Error("Error authorizing user", "users.UserAuth", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	// Create a new portal session
	url, err := users.CreateStripePortal(user)
	if err != nil {
		slog.Error("Error creating stripe portal", "users.CreateStripePortal", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	slog.Info("create_stripe_portal", "time", time.Since(start))
	return &pb.StripeUrlResponse{
		Url: url,
	}, nil
}
