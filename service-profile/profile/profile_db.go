package profile

import (
	"fmt"
	pb "service-profile/proto"
	"service-profile/system"
	"strings"
)

type ProfileDBProvider interface {
	selectProfileByUserId(userId string) (*pb.Profile, error)
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
	}
}

func (db ProfileDBImpl) selectProfileByUserId(userId string) (*pb.Profile, error) {
	row := db.Conn.QueryRow("select * from profiles where user_id = ?", userId)
	var profile pb.Profile
	err := row.Scan(dest(&profile)...)
	if err != nil {
		return nil, fmt.Errorf("scanProfile: %w", err)
	}
	return &profile, nil
}

func (db ProfileDBImpl) insertProfile(profile *pb.Profile) (*pb.Profile, error) {
	row := db.Conn.QueryRow(`insert into profiles (
        id,
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
        position
    ) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
        profile.Id,
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
		strings.Join(profile.EmailNotifications, ","),
		profile.PushNotification,
		profile.Resume,
		profile.Cover,
		profile.Position,
	)
    err := row.Scan(dest(profile)...)
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
        position = ?
    where id = ?`,
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
        profile.Id,
    )
    err := row.Scan(dest(profile)...)
    if err != nil {
        return nil, fmt.Errorf("updateProfile: %w", err)
    }
    return profile, nil
}
