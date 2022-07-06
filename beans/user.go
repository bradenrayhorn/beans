package beans

import (
	"context"

	"github.com/segmentio/ksuid"
)

type UserID ksuid.KSUID

type User struct {
	ID       UserID
	Username string
	Password string
}

type UserRepository interface {
	CreateUser(ctx context.Context, id UserID, username string, passwordHash string)
}
