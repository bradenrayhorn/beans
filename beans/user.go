package beans

import (
	"context"
	"errors"
	"strings"
)

type UserID ID

func (u UserID) String() string {
	return ID(u).String()
}

func UserIDFromString(id string) (UserID, error) {
	beansID, err := BeansIDFromString(id)
	return UserID(beansID), err
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
	Get(ctx context.Context, id UserID) (*User, error)
	GetByUsername(ctx context.Context, username Username) (*User, error)
}

type UserService interface {
	CreateUser(ctx context.Context, username Username, password Password) (*User, error)
	Login(ctx context.Context, username Username, password Password) (*User, error)
}
