package contract_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayee(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	cleanup := func() {
		_, err := pool.Exec(context.Background(), "truncate table users, budgets cascade;")
		require.Nil(t, err)
	}

	payeeRepository := postgres.NewPayeeRepository(pool)
	c := contract.NewPayeeContract(
		payeeRepository,
	)

	t.Run("create", func(t *testing.T) {
		t.Run("handles validation error", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			_, err := c.CreatePayee(context.Background(), auth, beans.Name(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			payee, err := c.CreatePayee(context.Background(), auth, beans.Name("Payee"))
			require.Nil(t, err)

			// payee was returned
			assert.Equal(t, "Payee", string(payee.Name))
			assert.Equal(t, budget.ID, payee.BudgetID)

			// payee was saved
			res, err := payeeRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, res, 1)
			assert.True(t, reflect.DeepEqual(payee, res[0]))
		})
	})

	t.Run("get all", func(t *testing.T) {
		t.Run("can get all payees", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			payee := testutils.MakePayee(t, pool, "Payee1", budget.ID)

			budget2 := testutils.MakeBudget(t, pool, "Budget", userID)
			_ = testutils.MakePayee(t, pool, "Payee2", budget2.ID)

			payees, err := c.GetAll(context.Background(), auth)
			require.Nil(t, err)
			require.Len(t, payees, 1)
			assert.True(t, reflect.DeepEqual(payee, payees[0]))
		})
	})
}
