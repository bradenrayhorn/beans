package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/sql/migrations"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orlangure/gnomock"
	pg "github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/require"
)

func StartPool(tb testing.TB) (*pgxpool.Pool, *gnomock.Container) {
	p := pg.Preset(
		pg.WithVersion("15.0"),
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

func assertPgError(tb testing.TB, code string, err error) {
	require.NotNil(tb, err)
	var pgErr *pgconn.PgError
	require.ErrorAs(tb, err, &pgErr)
	require.Equal(tb, code, pgErr.Code)
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

func makeAccount(tb testing.TB, pool *pgxpool.Pool, name string, budgetID beans.ID) beans.Account {
	id := beans.NewBeansID()
	err := postgres.NewAccountRepository(pool).Create(context.Background(), id, beans.Name(name), budgetID)
	require.Nil(tb, err)
	return beans.Account{
		ID:       id,
		Name:     beans.Name(name),
		BudgetID: budgetID,
	}
}

func makeMonth(tb testing.TB, pool *pgxpool.Pool, budgetID beans.ID, date beans.Date) beans.ID {
	id := beans.NewBeansID()
	err := postgres.NewMonthRepository(pool).Create(context.Background(), &beans.Month{ID: id, BudgetID: budgetID, Date: date})
	require.Nil(tb, err)
	return id
}

func makeCategoryGroup(tb testing.TB, pool *pgxpool.Pool, name string, budgetID beans.ID) beans.ID {
	id := beans.NewBeansID()
	err := postgres.NewCategoryRepository(pool).CreateGroup(context.Background(), &beans.CategoryGroup{ID: id, BudgetID: budgetID, Name: beans.Name(name)})
	require.Nil(tb, err)
	return id
}

func makeCategory(tb testing.TB, pool *pgxpool.Pool, name string, groupID beans.ID, budgetID beans.ID) beans.ID {
	id := beans.NewBeansID()
	err := postgres.NewCategoryRepository(pool).Create(context.Background(), &beans.Category{ID: id, BudgetID: budgetID, GroupID: groupID, Name: beans.Name(name)})
	require.Nil(tb, err)
	return id
}
