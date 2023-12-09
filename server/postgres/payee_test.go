package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayees(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	payeeRepository := postgres.NewPayeeRepository(pool)

	userID := testutils.MakeUser(t, pool, "user")
	budgetID := testutils.MakeBudget(t, pool, "budget", userID).ID
	budgetID2 := testutils.MakeBudget(t, pool, "budget2", userID).ID

	cleanup := func() {
		testutils.MustExec(t, pool, "truncate payees cascade;")
	}

	t.Run("can create and get", func(t *testing.T) {
		defer cleanup()
		payee := &beans.Payee{ID: beans.NewBeansID(), Name: "payee1", BudgetID: budgetID}
		require.Nil(t, payeeRepository.Create(context.Background(), payee))

		res, err := payeeRepository.Get(context.Background(), payee.ID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(res, payee))
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		defer cleanup()
		payee := &beans.Payee{ID: beans.NewBeansID(), Name: "payee1", BudgetID: budgetID}
		require.Nil(t, payeeRepository.Create(context.Background(), payee))

		assertPgError(t, pgerrcode.UniqueViolation, payeeRepository.Create(context.Background(), payee))
	})

	t.Run("cannot get non existant payee", func(t *testing.T) {
		defer cleanup()

		_, err := payeeRepository.Get(context.Background(), beans.NewBeansID())
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("can get payees for budget", func(t *testing.T) {
		defer cleanup()

		payee1 := &beans.Payee{ID: beans.NewBeansID(), Name: "payee1", BudgetID: budgetID}
		require.Nil(t, payeeRepository.Create(context.Background(), payee1))

		payee2 := &beans.Payee{ID: beans.NewBeansID(), Name: "payee2", BudgetID: budgetID2}
		require.Nil(t, payeeRepository.Create(context.Background(), payee2))

		// budget 1 contains a payee
		res, err := payeeRepository.GetForBudget(context.Background(), budgetID)
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.True(t, reflect.DeepEqual(res[0], payee1))

		// budget 2 contains a payee
		res, err = payeeRepository.GetForBudget(context.Background(), budgetID2)
		require.Nil(t, err)
		require.Len(t, res, 1)
		require.True(t, reflect.DeepEqual(res[0], payee2))

		// random budget contains no payee
		res, err = payeeRepository.GetForBudget(context.Background(), beans.NewBeansID())
		require.Nil(t, err)
		require.Len(t, res, 0)
	})
}
