package logic

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

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
		return nil, errors.New("invalid username")
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
