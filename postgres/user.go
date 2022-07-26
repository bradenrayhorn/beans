package postgres

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	db *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db.New(pool)}
}

func (r *UserRepository) Create(ctx context.Context, id beans.UserID, username beans.Username, passwordHash beans.PasswordHash) error {
	return r.db.CreateUser(ctx, db.CreateUserParams{
		ID:       id.String(),
		Username: string(username),
		Password: string(passwordHash),
	})
}

func (r *UserRepository) Exists(ctx context.Context, username beans.Username) (bool, error) {
	return r.db.UserExists(ctx, string(username))
}

func (r *UserRepository) Get(ctx context.Context, id beans.UserID) (*beans.User, error) {
	res, err := r.db.GetUserByID(ctx, id.String())
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
	res, err := r.db.GetUserByUsername(ctx, string(username))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, beans.WrapError(err, beans.ErrorNotFound)
		}

		return nil, err
	}
	id, err := beans.UserIDFromString(res.ID)
	if err != nil {
		return nil, err
	}

	return &beans.User{
		ID:           id,
		Username:     beans.Username(res.Username),
		PasswordHash: beans.PasswordHash(res.Password),
	}, nil
}
