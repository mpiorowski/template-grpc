package main

import (
	"context"
	"service-profile/profile"
	pb "service-profile/proto"
)

func (s *server) GetProfile(ctx context.Context, in *pb.Empty) (*pb.Profile, error) {
	return profile.GetProfile(ctx, s.storage)
}
func (s *server) InsertProfile(ctx context.Context, in *pb.Profile) (*pb.Profile, error) {
	return profile.InsertProfile(ctx, s.storage, in)
}
func (s *server) UpdateProfile(ctx context.Context, in *pb.Profile) (*pb.Profile, error) {
	return profile.UpdateProfile(ctx, s.storage, in)
}
