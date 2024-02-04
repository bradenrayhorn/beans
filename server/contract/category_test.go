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

func TestCategory(t *testing.T) {
	t.Parallel()
	pool, _, factory, stop := testutils.StartPoolWithDataSource(t)
	defer stop()

	cleanup := func() {
		testutils.MustExec(t, pool, "truncate table users, budgets cascade;")
	}

	categoryRepository := postgres.NewCategoryRepository(pool)
	monthCategoryRepository := postgres.NewMonthCategoryRepository(pool)
	monthRepository := postgres.NewMonthRepository(pool)
	txManager := postgres.NewTxManager(pool)
	c := contract.NewCategoryContract(
		categoryRepository,
		monthCategoryRepository,
		monthRepository,
		txManager,
	)

	t.Run("create group", func(t *testing.T) {
		t.Run("handles validation error", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			_, err := c.CreateGroup(context.Background(), auth, beans.Name(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			group, err := c.CreateGroup(context.Background(), auth, beans.Name("Group"))
			require.Nil(t, err)

			// group was returned
			assert.Equal(t, "Group", string(group.Name))
			assert.Equal(t, budget.ID, group.BudgetID)
			assert.Equal(t, false, group.IsIncome)

			// group was saved
			res, err := categoryRepository.GetGroupsForBudget(context.Background(), budget.ID)
			require.Nil(t, err)
			require.Len(t, res, 1)
			assert.True(t, reflect.DeepEqual(group, res[0]))
		})
	})

	t.Run("create category", func(t *testing.T) {
		t.Run("handles validation error", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			_, err := c.CreateCategory(context.Background(), auth, beans.NewBeansID(), beans.Name(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot create if group does not exist", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)

			_, err := c.CreateCategory(context.Background(), auth, beans.NewBeansID(), beans.Name("Category"))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create category", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			group := factory.MakeCategoryGroup("Group", budget.ID)

			category, err := c.CreateCategory(context.Background(), auth, group.ID, beans.Name("Category"))
			require.Nil(t, err)

			// category was returned
			assert.Equal(t, "Category", string(category.Name))
			assert.Equal(t, budget.ID, category.BudgetID)
			assert.Equal(t, group.ID, category.GroupID)

			// category was saved
			dbCategory, err := categoryRepository.GetSingleForBudget(context.Background(), category.ID, budget.ID)
			require.Nil(t, err)
			assert.True(t, reflect.DeepEqual(category, dbCategory))
		})

		t.Run("creates category for existing months", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			group := factory.MakeCategoryGroup("Group", budget.ID)
			month := factory.MakeMonth(budget.ID, testutils.NewDate(t, "2022-05-01"))

			category, err := c.CreateCategory(context.Background(), auth, group.ID, beans.Name("Category"))
			require.Nil(t, err)

			res, err := monthCategoryRepository.GetForMonth(context.Background(), month)
			require.Nil(t, err)
			require.Len(t, res, 1)
			assert.Equal(t, category.ID, res[0].CategoryID)
		})
	})

	t.Run("get all", func(t *testing.T) {
		t.Run("can get all categories and groups", func(t *testing.T) {
			defer cleanup()

			userID := factory.MakeUser("user")
			budget := factory.MakeBudget("Budget", userID)
			auth := testutils.BudgetAuthContext(t, userID, budget)
			group := factory.MakeCategoryGroup("Group", budget.ID)
			category := factory.MakeCategory("Category", group.ID, budget.ID)

			budget2 := factory.MakeBudget("Budget", userID)
			_ = factory.MakeCategoryGroup("Group", budget2.ID)
			_ = factory.MakeCategory("Category", group.ID, budget2.ID)

			result, err := c.GetAll(context.Background(), auth)
			require.Nil(t, err)

			assert.True(t, reflect.DeepEqual(result, []beans.CategoryGroupWithCategories{
				{
					CategoryGroup: *group,
					Categories:    []beans.Category{*category},
				},
			}))
		})
	})
}
