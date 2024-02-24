package contract

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/argon2"
	"github.com/bradenrayhorn/beans/server/beans"
)

var errorInvalidCredentials = beans.NewError(beans.EUNAUTHORIZED, "Invalid username or password")

var _ beans.UserContract = (*userContract)(nil)

type userContract struct {
	contract
}

func (c *userContract) Register(ctx context.Context, username beans.Username, password beans.Password) error {
	if err := beans.ValidateFields(username.ValidatableField(), password.ValidatableField()); err != nil {
		return err
	}

	id := beans.NewID()

	hashedPassword, err := argon2.GenerateHash(string(password))
	if err != nil {
		return err
	}

	usernameTaken, err := c.ds().UserRepository().Exists(ctx, username)
	if err != nil {
		return err
	}
	if usernameTaken {
		return beans.WrapError(errors.New("invalid username"), beans.ErrorInvalid)
	}

	err = c.ds().UserRepository().Create(ctx, id, username, beans.PasswordHash(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (c *userContract) Login(ctx context.Context, username beans.Username, password beans.Password) (beans.Session, error) {
	if err := beans.ValidateFields(username.ValidatableField(), password.ValidatableField()); err != nil {
		return beans.Session{}, err
	}

	user, err := c.ds().UserRepository().GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return beans.Session{}, beans.WrapError(err, errorInvalidCredentials)
		}
		return beans.Session{}, err
	}

	equal, err := argon2.CompareHashAndPassword(string(user.PasswordHash), string(password))
	if err != nil || !equal {
		return beans.Session{}, beans.WrapError(err, errorInvalidCredentials)
	}

	session, err := c.sessionRepository.Create(user.ID)
	if err != nil {
		return beans.Session{}, err
	}

	return session, nil
}

func (c *userContract) Logout(ctx context.Context, auth *beans.AuthContext) error {
	return c.sessionRepository.Delete(auth.SessionID())
}

func (c *userContract) GetMe(ctx context.Context, auth *beans.AuthContext) (beans.UserPublic, error) {
	user, err := c.ds().UserRepository().Get(ctx, auth.UserID())
	if err != nil {
		return beans.UserPublic{}, err
	}

	return beans.UserPublic{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
