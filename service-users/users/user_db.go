package users

import (
	"fmt"
	"powerit/system"

	pb "powerit/proto"

	"github.com/google/uuid"
)

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

func selectUserById(id string) (*pb.User, error) {
	row := system.Db.QueryRow("update users set updated = current_timestamp where id = $1 returning *", id)
    var user pb.User
    err := row.Scan(dest(&user)...)
	if err != nil {
		return nil, fmt.Errorf("scanUser: %w", err)
	}

	return &user, nil
}

func selectUserByEmailAndSub(email string, sub string) (*pb.User, error) {
	row := system.Db.QueryRow("select * from users where email = $1 and sub = $2", email, sub)
    var user pb.User
    err := row.Scan(dest(&user)...)
	if err != nil {
		return nil, fmt.Errorf("scanUser: %w", err)
	}

	return &user, nil
}

func insertUser(email string, sub string, avatar string) (*pb.User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("uuid.NewRandom: %w", err)
	}
	row := system.Db.QueryRow("insert into users (id, email, sub, role, avatar) values ($1, $2, $3, $4, $5) returning *",
		id, email, sub, pb.UserRole_ROLE_USER, avatar)
    var user pb.User
    err = row.Scan(dest(&user)...)
	if err != nil {
		return nil, fmt.Errorf("scanUser: %w", err)
	}

	return &user, nil
}

func updateSubscriptionId(userId string, subscriptionId string) error {
	_, err := system.Db.Exec("update users set subscription_id = $1 where id = $2", subscriptionId, userId)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func updateSubscriptionCheck(userId string, subscriptionCheck string) error {
	_, err := system.Db.Exec("update users set subscription_check = $1 where id = $2", subscriptionCheck, userId)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func updateSubscriptionEnd(userId string, subscriptionEnd string) error {
	_, err := system.Db.Exec("update users set subscription_end = $1 where id = $2", subscriptionEnd, userId)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
