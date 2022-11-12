package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

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
	err := postgres.NewBudgetRepository(pool).Create(context.Background(), id, beans.Name(name), userID, time.Now())
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
