package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
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
