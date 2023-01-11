package contract

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/argon2"
	"github.com/bradenrayhorn/beans/beans"
)

var errorInvalidCredentials = beans.NewError(beans.EUNAUTHORIZED, "Invalid username or password")

type userContract struct {
	userRepository beans.UserRepository
}

func NewUserContract(userRepository beans.UserRepository) *userContract {
	return &userContract{userRepository}
}

func (c *userContract) CreateUser(ctx context.Context, username beans.Username, password beans.Password) (*beans.User, error) {
	if err := beans.ValidateFields(username.ValidatableField(), password.ValidatableField()); err != nil {
		return nil, err
	}

	id := beans.NewBeansID()

	hashedPassword, err := argon2.GenerateHash(string(password))
	if err != nil {
		return nil, err
	}

	usernameTaken, err := c.userRepository.Exists(ctx, username)
	if err != nil {
		return nil, err
	}
	if usernameTaken {
		return nil, beans.WrapError(errors.New("invalid username"), beans.ErrorInvalid)
	}

	err = c.userRepository.Create(ctx, id, username, beans.PasswordHash(hashedPassword))
	if err != nil {
		return nil, err
	}

	return &beans.User{
		ID:           id,
		Username:     username,
		PasswordHash: beans.PasswordHash(hashedPassword),
	}, nil
}

func (c *userContract) Login(ctx context.Context, username beans.Username, password beans.Password) (*beans.User, error) {
	if err := beans.ValidateFields(username.ValidatableField(), password.ValidatableField()); err != nil {
		return nil, err
	}

	user, err := c.userRepository.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return nil, beans.WrapError(err, errorInvalidCredentials)
		}
		return nil, err
	}

	equal, err := argon2.CompareHashAndPassword(string(user.PasswordHash), string(password))
	if err != nil || !equal {
		return nil, beans.WrapError(err, errorInvalidCredentials)
	}

	return user, nil
}
