package beans

import (
	"context"
	"strings"
)

type Username string

func (u Username) ValidatableField() validatableField {
	return Field("Username", Required(u), Max(u, 32, "characters"))
}

func (u Username) Empty() bool {
	return strings.TrimSpace(string(u)) == ""
}

func (u Username) Length() int {
	return len(u)
}

type Password string

func (p Password) ValidatableField() validatableField {
	return Field("Password", Required(p), Max(p, 255, "characters"))
}

func (p Password) Empty() bool {
	return string(p) == ""
}

func (p Password) Length() int {
	return len(p)
}

type PasswordHash string

type User struct {
	ID           ID
	Username     Username
	PasswordHash PasswordHash
}

type UserContract interface {
	CreateUser(ctx context.Context, username Username, password Password) (*User, error)
	Login(ctx context.Context, username Username, password Password) (*User, error)
}

type UserRepository interface {
	Create(ctx context.Context, id ID, username Username, passwordHash PasswordHash) error
	Exists(ctx context.Context, username Username) (bool, error)
	Get(ctx context.Context, id ID) (*User, error)
	GetByUsername(ctx context.Context, username Username) (*User, error)
}
