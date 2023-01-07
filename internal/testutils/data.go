package testutils

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func MakeUser(tb testing.TB, pool *pgxpool.Pool, username string) beans.UserID {
	userID := beans.UserID(beans.NewBeansID())
	err := postgres.NewUserRepository(pool).Create(context.Background(), userID, beans.Username(username), beans.PasswordHash("x"))
	require.Nil(tb, err)
	return userID
}

func MakeBudget(tb testing.TB, pool *pgxpool.Pool, name string, userID beans.UserID) *beans.Budget {
	id := beans.NewBeansID()
	err := postgres.NewBudgetRepository(pool).Create(context.Background(), nil, id, beans.Name(name), userID)
	require.Nil(tb, err)
	return &beans.Budget{
		ID:      id,
		Name:    beans.Name(name),
		UserIDs: []beans.UserID{userID},
	}
}

func MakeMonth(tb testing.TB, pool *pgxpool.Pool, budgetID beans.ID, date beans.Date) *beans.Month {
	month := &beans.Month{ID: beans.NewBeansID(), BudgetID: budgetID, Date: date}
	err := postgres.NewMonthRepository(pool).Create(context.Background(), nil, month)
	require.Nil(tb, err)
	return month
}
