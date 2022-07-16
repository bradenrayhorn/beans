package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreatePool(url string) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	return pgxpool.ConnectConfig(context.Background(), pgxConfig)
}
