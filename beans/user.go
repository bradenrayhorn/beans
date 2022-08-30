package beans

import (
	"context"
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
