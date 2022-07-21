package beans

import (
	"context"
	"errors"
	"strings"

	"github.com/segmentio/ksuid"
)

type UserID ksuid.KSUID

func (u UserID) String() string {
	return ksuid.KSUID(u).String()
}

func UserIDFromString(id string) (UserID, error) {
	userID, err := ksuid.Parse(id)
	if err != nil {
		return UserID(ksuid.Nil), err
	}
	return UserID(userID), nil
}

type Username string

func (u Username) Validate() error {
	if strings.TrimSpace(string(u)) == "" {
		return errors.New("username is required")
	}
	return nil
}

type Password string

func (p Password) Validate() error {
	if string(p) == "" {
		return errors.New("password is required")
	}
	return nil
}

type PasswordHash string

type User struct {
	ID           UserID
	Username     Username
	PasswordHash PasswordHash
}

type UserRepository interface {
	Create(ctx context.Context, id UserID, username Username, passwordHash PasswordHash) error
	Exists(ctx context.Context, username Username) (bool, error)
	Get(ctx context.Context, username Username) (*User, error)
}

type UserService interface {
	CreateUser(ctx context.Context, username Username, password Password) (*User, error)
	Login(ctx context.Context, username Username, password Password) (*User, error)
}
