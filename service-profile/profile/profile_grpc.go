package profile

import (
	"context"
	"log/slog"
	pb "service-profile/proto"
	"service-profile/system"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetProfile(ctx context.Context, storage system.Storage) (*pb.Profile, error) {
	defer system.Perf("get_profile", time.Now())
	claims, err := system.ExtractToken(ctx)
	if err != nil {
		slog.Error("Error extracting token", "system.ExtractToken", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	var db = NewProfileDB(storage)
	profile, exists, err := db.selectProfileByUserId(claims.Id)
	if !exists {
		profile, err = db.insertProfile(&pb.Profile{UserId: claims.Id, Active: false})
	}
	if err != nil {
		slog.Error("Error selecting profile by user id", "db.selectProfileByUserId", err)
		return nil, status.Error(codes.NotFound, "Profile not found")
	}
	return profile, nil
}

func UpdateProfile(ctx context.Context, storage system.Storage, profile *pb.Profile) (*pb.Profile, error) {
	defer system.Perf("update_profile", time.Now())
	claims, err := system.ExtractToken(ctx)
	if err != nil {
		slog.Error("Error extracting token", "system.ExtractToken", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	var db = NewProfileDB(storage)
	profile.UserId = claims.Id
	profile, err = db.updateProfile(profile)
	if err != nil {
		slog.Error("Error updating profile", "db.updateProfile", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}
	return profile, nil
}
