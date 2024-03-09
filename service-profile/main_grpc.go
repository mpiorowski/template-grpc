package main

import (
	"context"
	"service-profile/profile"
	pb "service-profile/proto"
)

func (s *server) GetProfile(ctx context.Context, in *pb.Empty) (*pb.Profile, error) {
	profileDB := profile.NewProfileDB(&s.storage)
	profileService := profile.NewProfileService(profileDB)
	return profileService.GetProfile(ctx)
}
func (s *server) UpdateProfile(ctx context.Context, in *pb.Profile) (*pb.Profile, error) {
	profileDB := profile.NewProfileDB(&s.storage)
	profileService := profile.NewProfileService(profileDB)
	return profileService.UpdateProfile(ctx, in)
}
