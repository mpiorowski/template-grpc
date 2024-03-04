package main

import (
	"context"

	"service-auth/auth"
	pb "service-auth/proto"
)

func (s *server) Auth(ctx context.Context, in *pb.Empty) (*pb.AuthResponse, error) {
	return auth.Auth(ctx, s.storage)
}

func (s *server) CreateStripeCheckout(ctx context.Context, in *pb.Empty) (*pb.StripeUrlResponse, error) {
	return auth.CreateStripeCheckout(ctx, s.storage)
}

func (s *server) CreateStripePortal(ctx context.Context, in *pb.Empty) (*pb.StripeUrlResponse, error) {
	return auth.CreateStripePortal(ctx, s.storage)
}
