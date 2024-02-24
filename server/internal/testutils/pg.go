package testutils

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/sql/migrations"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orlangure/gnomock"
	pg "github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/require"
)

func StartPool(tb testing.TB) (*postgres.DbPool, func()) {
	p := pg.Preset(
		pg.WithVersion("16.0"),
		pg.WithDatabase("beans"),
		pg.WithQueries(getMigrationQueries(tb)),
	)

	container, err := gnomock.Start(p)
	if err != nil {
		tb.Fatal(err)
	}

	pool, err := postgres.CreatePool(
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			"postgres",
			"password",
			"127.0.0.1",
			container.DefaultPort(),
			"beans",
		))

	if err != nil {
		tb.Fatal(err)
	}

	return pool, func() {
		err := gnomock.Stop(container)
		if err != nil {
			tb.Error("failed to stop container", err)
		}
	}
}

func StartPoolWithDataSource(tb testing.TB) (*postgres.DbPool, beans.DataSource, *Factory, func()) {
	pool, stop := StartPool(tb)
	ds := postgres.NewDataSource(pool)
	factory := NewFactory(tb, ds)
	return pool, ds, factory, stop
}

func getMigrationQueries(tb testing.TB) string {
	queries := ""

	files, err := migrations.MigrationsFS.ReadDir(".")
	if err != nil {
		tb.Fatal(err)
	}

	for _, file := range files {
		content, err := migrations.MigrationsFS.ReadFile(file.Name())
		if err != nil {
			tb.Fatal(err)
		}

		queries += string(content)
	}

	return queries
}

func MustExec(t testing.TB, pool *postgres.DbPool, sql string) {
	_, err := (*pgxpool.Pool)(pool).Exec(context.Background(), sql)
	require.Nil(t, err)
}

func MustRollback(t testing.TB, tx beans.Tx) {
	err := tx.Rollback(context.Background())
	// Ignore ErrTxClosed as transaction may have already been committed.
	if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		t.Error("Failed to rollback tx", err)
	}
}
