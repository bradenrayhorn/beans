package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bradenrayhorn/beans/server/internal/sql/migrations"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func runMigrations(ctx context.Context) error {
	db, err := sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s/%s",
		"postgres",
		"password",
		"127.0.0.1:5432",
		"beans"))
	if err != nil {
		return err
	}

	fs, err := iofs.New(migrations.MigrationsFS, ".")
	if err != nil {
		return err
	}

	driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iosfs", fs, "postgres", driver)
	if err != nil {
		return err
	}

	return m.Up()
}
