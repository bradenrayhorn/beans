package postgres_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBudgets(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	budgetRepository := postgres.NewBudgetRepository(pool)
	txManager := postgres.NewTxManager(pool)

	userID := testutils.MakeUser(t, pool, "user")

	cleanup := func() {
		testutils.MustExec(t, pool, "truncate budgets cascade")
	}

	t.Run("can create and get budget", func(t *testing.T) {
		defer cleanup()
		budgetID := beans.NewBeansID()
		err := budgetRepository.Create(context.Background(), nil, budgetID, "Budget1", userID)
		require.Nil(t, err)

		budget, err := budgetRepository.Get(context.Background(), budgetID)
		require.Nil(t, err)
		assert.Equal(t, budgetID, budget.ID)
		assert.Equal(t, "Budget1", string(budget.Name))
		assert.Equal(t, []beans.ID{userID}, budget.UserIDs)
	})

	t.Run("create respects transaction", func(t *testing.T) {
		defer cleanup()
		budgetID1 := beans.NewBeansID()

		tx, err := txManager.Create(context.Background())
		require.Nil(t, err)
		defer testutils.MustRollback(t, tx)

		err = budgetRepository.Create(context.Background(), tx, budgetID1, "Budget1", userID)
		require.Nil(t, err)

		_, err = budgetRepository.Get(context.Background(), budgetID1)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		err = tx.Commit(context.Background())
		require.Nil(t, err)

		_, err = budgetRepository.Get(context.Background(), budgetID1)
		require.Nil(t, err)
	})
}
