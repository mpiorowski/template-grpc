package auth

import (
	"fmt"
	"service-auth/system"

	pb "service-auth/proto"

	"github.com/google/uuid"
)

type AuthDBProvider interface {
	selectUserById(id string) (*pb.User, error)
	selectUserByEmailAndSub(email string, sub string) (*pb.User, error)
	insertUser(email string, sub string, avatar string) (*pb.User, error)
	updateSubscriptionId(userId string, subscriptionId string) error
	updateSubscriptionCheck(userId string, subscriptionCheck string) error
	updateSubscriptionEnd(userId string, subscriptionEnd string) error
}

type AuthDBImpl struct {
	*system.Storage
}

var _ AuthDBProvider = AuthDBImpl{}

func NewAuthDB(s *system.Storage) AuthDBProvider {
    return AuthDBImpl{s}
}

func dest(user *pb.User) []interface{} {
	return []interface{}{
		&user.Id,
		&user.Created,
		&user.Updated,
		&user.Deleted,
		&user.Email,
		&user.Role,
		&user.Sub,
		&user.Avatar,
		&user.SubscriptionId,
		&user.SubscriptionEnd,
		&user.SubscriptionCheck,
	}
}

func (db AuthDBImpl) selectUserById(id string) (*pb.User, error) {
	row := db.Conn.QueryRow("update users set updated = current_timestamp where id = ? returning *", id)
	var user pb.User
	err := row.Scan(dest(&user)...)
	if err != nil {
		return nil, fmt.Errorf("scanUser: %w", err)
	}

	return &user, nil
}

func (db AuthDBImpl) selectUserByEmailAndSub(email string, sub string) (*pb.User, error) {
	row := db.Conn.QueryRow("select * from users where email = ? and sub = ?", email, sub)
	var user pb.User
	err := row.Scan(dest(&user)...)
	if err != nil {
		return nil, fmt.Errorf("scanUser: %w", err)
	}

	return &user, nil
}

func (db AuthDBImpl) insertUser(email string, sub string, avatar string) (*pb.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("uuid.NewRandom: %w", err)
	}
	row := db.Conn.QueryRow("insert into users (id, email, sub, role, avatar) values (?, ?, ?, ?, ?) returning *",
		id, email, sub, pb.UserRole_ROLE_USER, avatar)
	var user pb.User
	err = row.Scan(dest(&user)...)
	if err != nil {
		return nil, fmt.Errorf("scanUser: %w", err)
	}

	return &user, nil
}

func (db AuthDBImpl) updateSubscriptionId(userId string, subscriptionId string) error {
	_, err := db.Conn.Exec("update users set subscription_id = ? where id = ?", subscriptionId, userId)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (db AuthDBImpl) updateSubscriptionCheck(userId string, subscriptionCheck string) error {
	_, err := db.Conn.Exec("update users set subscription_check = ? where id = ?", subscriptionCheck, userId)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (db AuthDBImpl) updateSubscriptionEnd(userId string, subscriptionEnd string) error {
	_, err := db.Conn.Exec("update users set subscription_end = ? where id = ?", subscriptionEnd, userId)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
