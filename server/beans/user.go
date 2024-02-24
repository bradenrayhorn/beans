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

type UserPublic struct {
	ID       ID
	Username Username
}

type UserContract interface {
	// Creates a new account
	Register(ctx context.Context, username Username, password Password) error

	// Logs in and returns a session
	Login(ctx context.Context, username Username, password Password) (Session, error)

	// Logs out and deletes the active session
	Logout(ctx context.Context, auth *AuthContext) error

	// Gets the currently authenticated user
	GetMe(ctx context.Context, auth *AuthContext) (UserPublic, error)
}

type UserService interface {
	// Builds auth context
	GetAuth(ctx context.Context, sessionID SessionID) (*AuthContext, error)
}

type UserRepository interface {
	Create(ctx context.Context, id ID, username Username, passwordHash PasswordHash) error
	Exists(ctx context.Context, username Username) (bool, error)
	Get(ctx context.Context, id ID) (User, error)
	GetByUsername(ctx context.Context, username Username) (User, error)
}
