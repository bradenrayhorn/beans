package logic

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

var errorInvalidCredentials = beans.NewError(beans.EUNAUTHORIZED, "Invalid username or password")

type UserService struct {
	UserRepository beans.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, username beans.Username, password beans.Password) (*beans.User, error) {
	if err := beans.Validate(username, password); err != nil {
		return nil, err
	}

	id := beans.UserID(ksuid.New())

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
	if err := beans.Validate(username, password); err != nil {
		return nil, err
	}

	user, err := s.UserRepository.Get(ctx, username)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return nil, beans.WrapError(err, errorInvalidCredentials)
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, beans.WrapError(err, errorInvalidCredentials)
	}

	return user, nil
}
