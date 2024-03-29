package datasource

import (
	"context"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testPayee(t *testing.T, ds beans.DataSource) {
	factory := testutils.NewFactory(t, ds)

	payeeRepository := ds.PayeeRepository()
	ctx := context.Background()

	t.Run("can create and get", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		payee := beans.Payee{ID: beans.NewID(), Name: "payee1", BudgetID: budget.ID}
		require.Nil(t, payeeRepository.Create(ctx, payee))

		res, err := payeeRepository.Get(ctx, budget.ID, payee.ID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(res, payee))
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		payee := factory.Payee(beans.Payee{BudgetID: budget.ID})

		assert.NotNil(t, payeeRepository.Create(ctx, payee))
	})

	t.Run("cannot get non existant payee", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		_, err := payeeRepository.Get(ctx, budget.ID, beans.NewID())
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("cannot get payee for other budget", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()
		payee := factory.Payee(beans.Payee{BudgetID: budget.ID})

		_, err := payeeRepository.Get(ctx, budget2.ID, payee.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("can get payees for budget", func(t *testing.T) {
		budget1, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()

		payee1 := factory.Payee(beans.Payee{BudgetID: budget1.ID})
		payee2 := factory.Payee(beans.Payee{BudgetID: budget2.ID})

		// budget 1 contains a payee
		res, err := payeeRepository.GetForBudget(ctx, budget1.ID)
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.True(t, reflect.DeepEqual(res[0], payee1))

		// budget 2 contains a payee
		res, err = payeeRepository.GetForBudget(ctx, budget2.ID)
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.True(t, reflect.DeepEqual(res[0], payee2))

		// random budget contains no payee
		res, err = payeeRepository.GetForBudget(ctx, beans.NewID())
		require.Nil(t, err)
		require.Len(t, res, 0)
	})

	t.Run("sorts by name", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		payeeB := factory.Payee(beans.Payee{BudgetID: budget.ID, Name: "Bert"})
		payeeA := factory.Payee(beans.Payee{BudgetID: budget.ID, Name: "Ally"})
		payeeC := factory.Payee(beans.Payee{BudgetID: budget.ID, Name: "Charlie"})

		// check payees are returned in alphabetical order
		res, err := payeeRepository.GetForBudget(ctx, budget.ID)
		require.Nil(t, err)
		require.Len(t, res, 3)

		assert.Equal(t, payeeA.ID, res[0].ID)
		assert.Equal(t, payeeB.ID, res[1].ID)
		assert.Equal(t, payeeC.ID, res[2].ID)
	})
}
