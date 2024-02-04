package datasource

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBudgetRepository(t *testing.T, ds beans.DataSource) {
	factory := testutils.Factory(t, ds)
	budgetRepository := ds.BudgetRepository()
	ctx := context.Background()

	user1 := factory.User(beans.User{})

	t.Run("can create and get budget", func(t *testing.T) {
		budgetID := beans.NewBeansID()
		err := budgetRepository.Create(ctx, nil, budgetID, "Budget1", user1.ID)
		require.Nil(t, err)

		budget, err := budgetRepository.Get(context.Background(), budgetID)
		require.Nil(t, err)
		assert.Equal(t, budgetID, budget.ID)
		assert.Equal(t, "Budget1", string(budget.Name))
		assert.Equal(t, []beans.ID{user1.ID}, budget.UserIDs)
	})

	t.Run("create respects transaction", func(t *testing.T) {
		budgetID1 := beans.NewBeansID()

		tx, err := ds.TxManager().Create(context.Background())
		require.Nil(t, err)
		defer testutils.MustRollback(t, tx)

		err = budgetRepository.Create(context.Background(), tx, budgetID1, "Budget1", user1.ID)
		require.Nil(t, err)

		_, err = budgetRepository.Get(context.Background(), budgetID1)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		err = tx.Commit(context.Background())
		require.Nil(t, err)

		_, err = budgetRepository.Get(context.Background(), budgetID1)
		require.Nil(t, err)
	})
}