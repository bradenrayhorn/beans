package postgres_test

import (
	"context"
	"testing"

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

func makeCategoryGroup(tb testing.TB, pool *pgxpool.Pool, name string, budgetID beans.ID) beans.ID {
	id := beans.NewBeansID()
	err := postgres.NewCategoryRepository(pool).CreateGroup(context.Background(), nil, &beans.CategoryGroup{ID: id, BudgetID: budgetID, Name: beans.Name(name)})
	require.Nil(tb, err)
	return id
}

func makeCategory(tb testing.TB, pool *pgxpool.Pool, name string, groupID beans.ID, budgetID beans.ID) beans.ID {
	id := beans.NewBeansID()
	err := postgres.NewCategoryRepository(pool).Create(context.Background(), nil, &beans.Category{ID: id, BudgetID: budgetID, GroupID: groupID, Name: beans.Name(name)})
	require.Nil(tb, err)
	return id
}

func makeTransaction(tb testing.TB, pool *pgxpool.Pool, transaction *beans.Transaction) *beans.Transaction {
	err := postgres.NewTransactionRepository(pool).Create(context.Background(), transaction)
	require.Nil(tb, err)
	return transaction
}
