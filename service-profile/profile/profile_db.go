package profile

import (
	"database/sql"
	"fmt"
	pb "service-profile/proto"
	"service-profile/system"

	"github.com/google/uuid"
)

type ProfileDBProvider interface {
	selectProfileByUserId(userId string) (*pb.Profile, bool, error)
	insertProfile(profile *pb.Profile) (*pb.Profile, error)
	updateProfile(profile *pb.Profile) (*pb.Profile, error)
}

type ProfileDBImpl struct {
	system.Storage
}

var _ ProfileDBProvider = ProfileDBImpl{}

func NewProfileDB(s system.Storage) ProfileDBProvider {
	return ProfileDBImpl{s}
}

func dest(profile *pb.Profile) []interface{} {
	return []interface{}{
		&profile.Id,
		&profile.Created,
		&profile.Updated,
		&profile.UserId,
		&profile.Active,
		&profile.Username,
		&profile.About,
		&profile.FirstName,
		&profile.LastName,
		&profile.Email,
		&profile.Country,
		&profile.StreetAddress,
		&profile.City,
		&profile.State,
		&profile.Zip,
		&profile.EmailNotifications,
		&profile.PushNotification,
		&profile.Resume,
		&profile.Cover,
		&profile.Position,
		&profile.Skills,
	}
}

func (db ProfileDBImpl) selectProfileByUserId(userId string) (*pb.Profile, bool, error) {
	row := db.Conn.QueryRow("select * from profiles where user_id = ?", userId)
	var profile pb.Profile
	err := row.Scan(dest(&profile)...)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, true, fmt.Errorf("scanProfile: %w", err)
	}
	return &profile, true, nil
}

func (db ProfileDBImpl) insertProfile(profile *pb.Profile) (*pb.Profile, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("uuid.NewV7: %w", err)
	}
	row := db.Conn.QueryRow(`insert into profiles (
        id,
        user_id,
        active,
        username,
        about,
        first_name,
        last_name,
        email,
        country,
        street_address,
        city,
        state,
        zip,
        email_notifications,
        push_notification,
        resume,
        cover,
        position,
        skills
    ) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) returning *`,
        id.String(),
		profile.UserId,
		profile.Active,
		profile.Username,
		profile.About,
		profile.FirstName,
		profile.LastName,
		profile.Email,
		profile.Country,
		profile.StreetAddress,
		profile.City,
		profile.State,
		profile.Zip,
		profile.EmailNotifications,
		profile.PushNotification,
		profile.Resume,
		profile.Cover,
		profile.Position,
		profile.Skills,
	)
	err = row.Scan(dest(profile)...)
	if err != nil {
		return nil, fmt.Errorf("insertProfile: %w", err)
	}
	return profile, nil
}

func (db ProfileDBImpl) updateProfile(profile *pb.Profile) (*pb.Profile, error) {
	row := db.Conn.QueryRow(`update profiles set
        active = ?,
        username = ?,
        about = ?,
        first_name = ?,
        last_name = ?,
        email = ?,
        country = ?,
        street_address = ?,
        city = ?,
        state = ?,
        zip = ?,
        email_notifications = ?,
        push_notification = ?,
        resume = ?,
        cover = ?,
        position = ?,
        skills = ?
    where id = ? and user_id = ? returning *`,
		profile.Active,
		profile.Username,
		profile.About,
		profile.FirstName,
		profile.LastName,
		profile.Email,
		profile.Country,
		profile.StreetAddress,
		profile.City,
		profile.State,
		profile.Zip,
		profile.EmailNotifications,
		profile.PushNotification,
		profile.Resume,
		profile.Cover,
		profile.Position,
        profile.Skills,
		profile.Id,
		profile.UserId,
	)
	err := row.Scan(dest(profile)...)
	if err != nil {
		return nil, fmt.Errorf("updateProfile: %w", err)
	}
	return profile, nil
}
