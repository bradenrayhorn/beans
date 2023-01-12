package contract_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/contract"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	cleanup := func() {
		_, err := pool.Exec(context.Background(), "truncate table users, budgets cascade;")
		require.Nil(t, err)
	}

	accountRepository := postgres.NewAccountRepository(pool)
	c := contract.NewAccountContract(accountRepository)

	t.Run("create", func(t *testing.T) {
		t.Run("handles validation error", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)

			_, err := c.Create(context.Background(), testutils.BudgetAuthContext(t, userID, budget), beans.Name(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create account", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)

			account, err := c.Create(context.Background(), testutils.BudgetAuthContext(t, userID, budget), beans.Name("Account"))
			require.Nil(t, err)

			// account was returned
			assert.Equal(t, "Account", string(account.Name))
			assert.Equal(t, budget.ID, account.BudgetID)
			assert.False(t, account.ID.Empty())

			// account was saved
			dbAccount, err := accountRepository.Get(context.Background(), account.ID)
			require.Nil(t, err)
			assert.True(t, reflect.DeepEqual(account, dbAccount))
		})
	})

	t.Run("get all", func(t *testing.T) {
		t.Run("can get all accounts", func(t *testing.T) {
			defer cleanup()
			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			account := testutils.MakeAccount(t, pool, "Account", budget.ID)

			budget2 := testutils.MakeBudget(t, pool, "Budget", userID)
			_ = testutils.MakeAccount(t, pool, "Account", budget2.ID)

			accounts, err := c.GetAll(context.Background(), testutils.BudgetAuthContext(t, userID, budget))
			require.Nil(t, err)
			require.Len(t, accounts, 1)

			assert.True(t, reflect.DeepEqual(account, accounts[0]))
		})
	})
}
