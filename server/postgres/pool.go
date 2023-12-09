package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbPool pgxpool.Pool

func CreatePool(url string) (*DbPool, error) {
	pgxConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	pgpool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, err
	}

	return (*DbPool)(pgpool), nil
}

func (p *DbPool) Close() {
	(*pgxpool.Pool)(p).Close()
}
