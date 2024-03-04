package main

import (
	"context"
	"service-profile/profile"
	pb "service-profile/proto"
)

func (s *server) GetProfile(ctx context.Context, in *pb.Empty) (*pb.Profile, error) {
	return profile.GetProfile(ctx, s.storage)
}
func (s *server) UpdateProfile(ctx context.Context, in *pb.Profile) (*pb.Profile, error) {
	return profile.UpdateProfile(ctx, s.storage, in)
}
