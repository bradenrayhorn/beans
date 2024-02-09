package contract

import (
	"context"
	"errors"
	"fmt"

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

	id := beans.NewBeansID()

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

func (c *userContract) GetMe(ctx context.Context, auth *beans.AuthContext) (beans.User, error) {
	user, err := c.ds().UserRepository().Get(ctx, auth.UserID())
	if err != nil {
		return beans.User{}, err
	}

	return user, nil
}

func (c *userContract) GetAuth(ctx context.Context, sessionID beans.SessionID) (*beans.AuthContext, error) {
	session, err := c.sessionRepository.Get(sessionID)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			return nil, beans.ErrorUnauthorized
		}

		return nil, fmt.Errorf("UserContract.GetAuth get session: %w", err)
	}

	return beans.NewAuthContext(session.UserID, session.ID), nil
}
