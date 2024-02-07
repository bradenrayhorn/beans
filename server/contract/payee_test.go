package contract_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayee(t *testing.T) {
	t.Parallel()
	pool, ds, factory, stop := testutils.StartPoolWithDataSource(t)
	defer stop()

	cleanup := func() {
		testutils.MustExec(t, pool, "truncate table users, budgets cascade;")
	}

	payeeRepository := postgres.NewPayeeRepository(pool)
	c := contract.NewContracts(ds, inmem.NewSessionRepository()).Payee

	t.Run("create", func(t *testing.T) {
		t.Run("handles validation error", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			_, err := c.CreatePayee(context.Background(), auth, beans.Name(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
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

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			payee := factory.MakePayee("Payee1", budget.ID)

			budget2 := factory.MakeBudget("Budget", userID)
			_ = factory.MakePayee("Payee2", budget2.ID)

			payees, err := c.GetAll(context.Background(), auth)
			require.Nil(t, err)
			require.Len(t, payees, 1)
			assert.True(t, reflect.DeepEqual(payee, payees[0]))
		})
	})
}
