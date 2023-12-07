package contract

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/argon2"
	"github.com/bradenrayhorn/beans/beans"
)

var errorInvalidCredentials = beans.NewError(beans.EUNAUTHORIZED, "Invalid username or password")

var _ beans.UserContract = (*userContract)(nil)

type userContract struct {
	sessionRepository beans.SessionRepository
	userRepository    beans.UserRepository
}

func NewUserContract(
	sessionRepository beans.SessionRepository,
	userRepository beans.UserRepository,
) *userContract {
	return &userContract{sessionRepository, userRepository}
}

func (c *userContract) Register(ctx context.Context, username beans.Username, password beans.Password) error {
	if err := beans.ValidateFields(username.ValidatableField(), password.ValidatableField()); err != nil {
		return err
	}

	id := beans.NewBeansID()

	hashedPassword, err := argon2.GenerateHash(string(password))
	if err != nil {
		return err
	}

	usernameTaken, err := c.userRepository.Exists(ctx, username)
	if err != nil {
		return err
	}
	if usernameTaken {
		return beans.WrapError(errors.New("invalid username"), beans.ErrorInvalid)
	}

	err = c.userRepository.Create(ctx, id, username, beans.PasswordHash(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (c *userContract) Login(ctx context.Context, username beans.Username, password beans.Password) (*beans.Session, error) {
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

	session, err := c.sessionRepository.Create(user.ID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (c *userContract) Logout(ctx context.Context, auth *beans.AuthContext) error {
	return c.sessionRepository.Delete(auth.SessionID())
}

func (c *userContract) GetMe(ctx context.Context, auth *beans.AuthContext) (*beans.User, error) {
	user, err := c.userRepository.Get(ctx, auth.UserID())
	if err != nil {
		return nil, err
	}

	return user, nil
}
