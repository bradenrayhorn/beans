package postgres

import (
	"context"
	"database/sql"

	"github.com/bradenrayhorn/beans/server/internal/sql/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type DbPool pgxpool.Pool

func CreatePool(url string) (*DbPool, error) {
	ctx := context.Background()
	pgxConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	pgpool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	if err := migrate(ctx, pgpool); err != nil {
		return nil, err
	}

	return (*DbPool)(pgpool), nil
}

func (p *DbPool) Close() {
	(*pgxpool.Pool)(p).Close()
}

func migrate(ctx context.Context, pool *pgxpool.Pool) error {
	db, err := sql.Open("pgx/v5", pool.Config().ConnString())
	if err != nil {
		return err
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		migrations.MigrationsFS,
	)
	if err != nil {
		return err
	}

	defer func() { _ = provider.Close() }()
	_, err = provider.Up(ctx)
	return err
}
