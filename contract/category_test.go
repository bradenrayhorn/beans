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

func TestCategory(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	cleanup := func() {
		_, err := pool.Exec(context.Background(), "truncate table users, budgets cascade;")
		require.Nil(t, err)
	}

	categoryRepository := postgres.NewCategoryRepository(pool)
	c := contract.NewCategoryContract(categoryRepository)

	t.Run("create category", func(t *testing.T) {
		t.Run("handles validation error", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			_, err := c.CreateCategory(context.Background(), auth, beans.NewBeansID(), beans.Name(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot create if group does not exist", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			_, err := c.CreateCategory(context.Background(), auth, beans.NewBeansID(), beans.Name("Category"))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create category", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			group := testutils.MakeCategoryGroup(t, pool, "Group", budget.ID)

			category, err := c.CreateCategory(context.Background(), auth, group.ID, beans.Name("Category"))
			require.Nil(t, err)

			// category was returned
			assert.Equal(t, "Category", string(category.Name))
			assert.Equal(t, budget.ID, category.BudgetID)
			assert.Equal(t, group.ID, category.GroupID)
			assert.Equal(t, false, category.IsIncome)

			// category was saved
			dbCategory, err := categoryRepository.GetSingleForBudget(context.Background(), category.ID, budget.ID)
			require.Nil(t, err)
			assert.True(t, reflect.DeepEqual(category, dbCategory))
		})
	})

	t.Run("get all", func(t *testing.T) {
		t.Run("can get all categories and groups", func(t *testing.T) {
			defer cleanup()

			userID := testutils.MakeUser(t, pool, "user")
			budget := testutils.MakeBudget(t, pool, "Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			group := testutils.MakeCategoryGroup(t, pool, "Group", budget.ID)
			category := testutils.MakeCategory(t, pool, "Category", group.ID, budget.ID)

			budget2 := testutils.MakeBudget(t, pool, "Budget", userID)
			_ = testutils.MakeCategoryGroup(t, pool, "Group", budget2.ID)
			_ = testutils.MakeCategory(t, pool, "Category", group.ID, budget2.ID)

			groups, categories, err := c.GetAll(context.Background(), auth)
			require.Nil(t, err)
			require.Len(t, groups, 1)
			require.Len(t, categories, 1)

			assert.True(t, reflect.DeepEqual(group, groups[0]))
			assert.True(t, reflect.DeepEqual(category, categories[0]))
		})
	})
}
