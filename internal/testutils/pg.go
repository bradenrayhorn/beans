package testutils

import (
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/internal/sql/migrations"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orlangure/gnomock"
	pg "github.com/orlangure/gnomock/preset/postgres"
)

func StartPool(tb testing.TB) (*pgxpool.Pool, func()) {
	p := pg.Preset(
		pg.WithVersion("15.1"),
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
		gnomock.Stop(container)
	}
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
