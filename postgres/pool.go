package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreatePool() (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig("postgres://postgres:password@127.0.0.1:5432/beans")
	if err != nil {
		return nil, err
	}

	return pgxpool.ConnectConfig(context.Background(), pgxConfig)
}
