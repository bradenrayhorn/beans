package beans

import (
	"context"

	"github.com/segmentio/ksuid"
)

type UserID ksuid.KSUID

func (u UserID) String() string {
	return ksuid.KSUID(u).String()
}

type Username string

type Password string

type PasswordHash string

type User struct {
	ID           UserID
	Username     Username
	PasswordHash PasswordHash
}

type UserRepository interface {
	Create(ctx context.Context, id UserID, username Username, passwordHash PasswordHash) error
	Exists(ctx context.Context, username Username) (bool, error)
}

type UserService interface {
	CreateUser(ctx context.Context, username Username, password Password) (*User, error)
}
