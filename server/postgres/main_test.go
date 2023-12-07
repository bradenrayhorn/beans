package postgres_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/postgres"
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

func makeTransaction(tb testing.TB, pool *pgxpool.Pool, transaction *beans.Transaction) *beans.Transaction {
	err := postgres.NewTransactionRepository(pool).Create(context.Background(), transaction)
	require.Nil(tb, err)
	return transaction
}

func findResult[K comparable](s []*K, compare func(*K) bool) *K {
	for _, v := range s {
		if compare(v) {
			return v
		}
	}
	return nil
}
