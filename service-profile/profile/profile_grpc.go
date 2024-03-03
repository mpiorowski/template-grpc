package profile

import (
	"context"
	"log/slog"
	pb "service-profile/proto"
	"service-profile/system"
	"time"

	"github.com/google/uuid"
)

func GetProfile(ctx context.Context, storage system.Storage) (*pb.Profile, error) {
	defer system.Perf("GetProfile", time.Now())
	claims, err := system.ExtractToken(ctx)
	if err != nil {
		slog.Error("Error extracting token", "system.ExtractToken", err)
		return nil, err
	}

	var db = NewProfileDB(storage)
	profile, err := db.selectProfileByUserId(claims.Id)
	if err != nil {
		slog.Error("Error selecting profile by user id", "db.selectProfileByUserId", err)
		return nil, err
	}
	return profile, nil
}

func InsertProfile(ctx context.Context, storage system.Storage, profile *pb.Profile) (*pb.Profile, error) {
	defer system.Perf("InsertProfile", time.Now())
	claims, err := system.ExtractToken(ctx)
	if err != nil {
		slog.Error("Error extracting token", "system.ExtractToken", err)
		return nil, err
	}

	var db = NewProfileDB(storage)
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("Error generating uuid", "uuid.NewV7", err)
		return nil, err
	}
	profile.Id = id.String()
	profile.UserId = claims.Id
	profile, err = db.insertProfile(profile)
	if err != nil {
		slog.Error("Error inserting profile", "db.insertProfile", err)
		return nil, err
	}
	return profile, nil
}

func UpdateProfile(ctx context.Context, storage system.Storage, profile *pb.Profile) (*pb.Profile, error) {
	defer system.Perf("UpdateProfile", time.Now())
	claims, err := system.ExtractToken(ctx)
	if err != nil {
		slog.Error("Error extracting token", "system.ExtractToken", err)
		return nil, err
	}

	var db = NewProfileDB(storage)
	profile.UserId = claims.Id
	profile, err = db.updateProfile(profile)
	if err != nil {
		slog.Error("Error updating profile", "db.updateProfile", err)
		return nil, err
	}
	return profile, nil
}
