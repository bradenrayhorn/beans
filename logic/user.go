package logic

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/argon2"
	"github.com/bradenrayhorn/beans/beans"
	"github.com/segmentio/ksuid"
)

var errorInvalidCredentials = beans.NewError(beans.EUNAUTHORIZED, "Invalid username or password")

type UserService struct {
	UserRepository beans.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, username beans.Username, password beans.Password) (*beans.User, error) {
	if err := beans.ValidateFields(username.ValidatableField(), password.ValidatableField()); err != nil {
		return nil, err
	}

	id := beans.UserID(ksuid.New())

	hashedPassword, err := argon2.GenerateHash(string(password))
	if err != nil {
		return nil, err
	}

	usernameTaken, err := s.UserRepository.Exists(ctx, username)
	if err != nil {
		return nil, err
	}
	if usernameTaken {
		return nil, beans.WrapError(errors.New("invalid username"), beans.ErrorInvalid)
	}

	err = s.UserRepository.Create(ctx, id, username, beans.PasswordHash(hashedPassword))
	if err != nil {
		return nil, err
	}

	return &beans.User{
		ID:           id,
		Username:     username,
		PasswordHash: beans.PasswordHash(hashedPassword),
	}, nil
}

func (s *UserService) Login(ctx context.Context, username beans.Username, password beans.Password) (*beans.User, error) {
	if err := beans.ValidateFields(username.ValidatableField(), password.ValidatableField()); err != nil {
		return nil, err
	}

	user, err := s.UserRepository.GetByUsername(ctx, username)
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
