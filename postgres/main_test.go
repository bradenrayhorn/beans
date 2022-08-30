package postgres_test

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orlangure/gnomock"
	pg "github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/require"
)

func StartPool(tb testing.TB) (*pgxpool.Pool, *gnomock.Container) {
	p := pg.Preset(
		pg.WithVersion("13.4"),
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

	return pool, container
}

func StopPool(tb testing.TB, container *gnomock.Container) {
	gnomock.Stop(container)
}

func getMigrationQueries(tb testing.TB) string {
	queries := ""
	err := filepath.WalkDir("../internal/sql/migrations/", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		queries += string(content)

		return nil
	})
	if err != nil {
		tb.Fatal(err)
	}

	return queries
}

func makeUser(tb testing.TB, pool *pgxpool.Pool, username string) beans.UserID {
	userID := beans.UserID(beans.NewBeansID())
	err := postgres.NewUserRepository(pool).Create(context.Background(), userID, beans.Username(username), beans.PasswordHash("x"))
	require.Nil(tb, err)
	return userID
}

func makeBudget(tb testing.TB, pool *pgxpool.Pool, name string, userID beans.UserID) beans.ID {
	id := beans.NewBeansID()
	err := postgres.NewBudgetRepository(pool).Create(context.Background(), id, beans.Name(name), userID)
	require.Nil(tb, err)
	return id
}
