package postgres_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccounts(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	accountRepository := postgres.NewAccountRepository(pool)

	userID := testutils.MakeUser(t, pool, "user")
	budgetID := testutils.MakeBudget(t, pool, "budget", userID).ID

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
		assertPgError(t, pgerrcode.UniqueViolation, err)
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
