package postgres_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccounts(t *testing.T) {
	pool, container := StartPool(t)
	defer StopPool(t, container)

	accountRepository := postgres.NewAccountRepository(pool)

	userID := makeUser(t, pool, "user")
	budgetID := makeBudget(t, pool, "budget", userID)

	t.Run("can create and get account", func(t *testing.T) {
		defer pool.Exec(context.Background(), "truncate accounts;")
		accountID := beans.NewBeansID()
		err := accountRepository.Create(context.Background(), accountID, "Account1", budgetID)
		require.Nil(t, err)

		account, err := accountRepository.Get(context.Background(), accountID)
		require.Nil(t, err)
		assert.Equal(t, accountID, account.ID)
		assert.Equal(t, "Account1", string(account.Name))
		assert.Equal(t, budgetID, account.BudgetID)
	})

	t.Run("cannot create duplicate account", func(t *testing.T) {
		defer pool.Exec(context.Background(), "truncate accounts;")
		accountID := beans.NewBeansID()
		err := accountRepository.Create(context.Background(), accountID, "Account1", budgetID)
		require.Nil(t, err)

		err = accountRepository.Create(context.Background(), accountID, "Account1", budgetID)
		require.NotNil(t, err)
		var pgErr *pgconn.PgError
		require.ErrorAs(t, err, &pgErr)
		assert.Equal(t, pgerrcode.UniqueViolation, pgErr.Code)
	})

	t.Run("cannot get fictitious account", func(t *testing.T) {
		defer pool.Exec(context.Background(), "truncate accounts;")
		accountID := beans.NewBeansID()
		_, err := accountRepository.Get(context.Background(), accountID)
		require.NotNil(t, err)
		var beansError beans.Error
		require.ErrorAs(t, err, &beansError)
		code, _ := beansError.BeansError()
		assert.Equal(t, beans.ENOTFOUND, code)
	})
}
