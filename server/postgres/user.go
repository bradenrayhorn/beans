package postgres

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	repository
}

func NewUserRepository(pool *DbPool) *UserRepository {
	return &UserRepository{repository{pool}}
}

func (r *UserRepository) Create(ctx context.Context, id beans.ID, username beans.Username, passwordHash beans.PasswordHash) error {
	return r.DB(nil).CreateUser(ctx, db.CreateUserParams{
		ID:       id.String(),
		Username: string(username),
		Password: string(passwordHash),
	})
}

func (r *UserRepository) Exists(ctx context.Context, username beans.Username) (bool, error) {
	return r.DB(nil).UserExists(ctx, string(username))
}

func (r *UserRepository) Get(ctx context.Context, id beans.ID) (*beans.User, error) {
	res, err := r.DB(nil).GetUserByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, beans.WrapError(err, beans.ErrorNotFound)
		}
		return nil, err
	}

	return &beans.User{
		ID:           id,
		Username:     beans.Username(res.Username),
		PasswordHash: beans.PasswordHash(res.Password),
	}, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username beans.Username) (*beans.User, error) {
	res, err := r.DB(nil).GetUserByUsername(ctx, string(username))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, beans.WrapError(err, beans.ErrorNotFound)
		}

		return nil, err
	}
	id, err := beans.BeansIDFromString(res.ID)
	if err != nil {
		return nil, err
	}

	return &beans.User{
		ID:           id,
		Username:     beans.Username(res.Username),
		PasswordHash: beans.PasswordHash(res.Password),
	}, nil
}
