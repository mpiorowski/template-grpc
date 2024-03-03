package main

import (
	"context"
	"log/slog"
	"time"

	"service-auth/auth"
	pb "service-auth/proto"
	"service-auth/system"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Auth(ctx context.Context, in *pb.Empty) (*pb.AuthResponse, error) {
	defer system.Perf("auth", time.Now())
	user, token, err := auth.Auth(ctx, s.storage)
	if err != nil {
		slog.Error("Error authorizing user", "auth.Auth", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	return &pb.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *server) CreateStripeCheckout(ctx context.Context, in *pb.Empty) (*pb.StripeUrlResponse, error) {
	defer system.Perf("create_stripe_checkout", time.Now())
	user, err := auth.GetUser(ctx, s.storage)
	if err != nil {
		slog.Error("Error authorizing user", "auth.GetUser", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	// Create a new checkout session
	url, err := auth.CreateStripeCheckout(user, s.storage)
	if err != nil {
		slog.Error("Error creating stripe checkout", "auth.CreateStripeCheckout", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &pb.StripeUrlResponse{
		Url: url,
	}, nil
}

func (s *server) CreateStripePortal(ctx context.Context, in *pb.Empty) (*pb.StripeUrlResponse, error) {
	defer system.Perf("create_stripe_portal", time.Now())
	user, err := auth.GetUser(ctx, s.storage)
	if err != nil {
		slog.Error("Error authorizing user", "auth.GetUser", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	// Create a new portal session
	url, err := auth.CreateStripePortal(user)
	if err != nil {
		slog.Error("Error creating stripe portal", "auth.CreateStripePortal", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &pb.StripeUrlResponse{
		Url: url,
	}, nil
}
