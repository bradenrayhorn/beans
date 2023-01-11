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

func TestBudget(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	cleanup := func() {
		_, err := pool.Exec(context.Background(), "truncate table users, budgets cascade;")
		require.Nil(t, err)
	}

	budgetRepository := postgres.NewBudgetRepository(pool)
	categoryRepository := postgres.NewCategoryRepository(pool)
	monthRepository := postgres.NewMonthRepository(pool)
	c := contract.NewBudgetContract(
		budgetRepository,
		categoryRepository,
		monthRepository,
		postgres.NewTxManager(pool),
	)

	t.Run("create", func(t *testing.T) {
		t.Run("handles validation error", func(t *testing.T) {
			defer cleanup()

			_, err := c.Create(context.Background(), beans.Name(""), beans.NewBeansID())
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create budget", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")

			budget, err := c.Create(context.Background(), beans.Name("Test"), userID)
			require.Nil(t, err)

			// budget was returned
			assert.Equal(t, "Test", string(budget.Name))
			assert.False(t, budget.ID.Empty())

			// budget was saved
			dbBudget, err := budgetRepository.Get(context.Background(), budget.ID)
			require.Nil(t, err)
			assert.Equal(t, budget.Name, dbBudget.Name)
			assert.Equal(t, budget.ID, dbBudget.ID)

			// month was created
			_, err = monthRepository.GetLatest(context.Background(), budget.ID)
			require.Nil(t, err)

			// income category was created
			groups, err := categoryRepository.GetGroupsForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			assert.Len(t, groups, 1)
			assert.Equal(t, "Income", string(groups[0].Name))

			categories, err := categoryRepository.GetForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			assert.Len(t, categories, 1)
			assert.Equal(t, "Income", string(categories[0].Name))
			assert.Equal(t, true, categories[0].IsIncome)
		})
	})

	t.Run("get", func(t *testing.T) {
		t.Run("cannot get nonexistant budget", func(t *testing.T) {
			defer cleanup()
			userID := testutils.MakeUser(t, pool, "user")

			_, _, err := c.Get(context.Background(), beans.NewBeansID(), userID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot get budget for wrong user", func(t *testing.T) {
			defer cleanup()
			userID1 := testutils.MakeUser(t, pool, "user1")
			userID2 := testutils.MakeUser(t, pool, "user2")

			budget := testutils.MakeBudget(t, pool, "Budget", userID1)

			_, _, err := c.Get(context.Background(), budget.ID, userID2)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("missing month gives internal error", func(t *testing.T) {
			defer cleanup()
			userID := testutils.MakeUser(t, pool, "user")

			budget := testutils.MakeBudget(t, pool, "Budget", userID)

			_, _, err := c.Get(context.Background(), budget.ID, userID)
			testutils.AssertErrorCode(t, err, beans.EINTERNAL)
		})

		t.Run("can get budget", func(t *testing.T) {
			defer cleanup()
			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			month := testutils.MakeMonth(t, pool, budget.ID, testutils.NewDate(t, "2022-05-01"))

			rBudget, rMonth, err := c.Get(context.Background(), budget.ID, userID)
			require.Nil(t, err)
			assert.True(t, reflect.DeepEqual(budget, rBudget))
			assert.True(t, reflect.DeepEqual(month, rMonth))
		})
	})

	t.Run("get all", func(t *testing.T) {
		t.Run("can get all budgets", func(t *testing.T) {
			defer cleanup()
			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)

			userID2 := testutils.MakeUser(t, pool, "user")
			_ = testutils.MakeBudget(t, pool, "Budget", userID2)

			budgets, err := c.GetAll(context.Background(), userID)
			require.Nil(t, err)
			require.Len(t, budgets, 1)

			budget.UserIDs = nil
			assert.True(t, reflect.DeepEqual(budget, budgets[0]))
		})
	})
}
