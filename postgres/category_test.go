package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategories(t *testing.T) {
	pool, stop := testutils.StartPool(t)
	defer stop()

	categoryRepository := postgres.NewCategoryRepository(pool)

	userID := makeUser(t, pool, "user")
	budgetID := makeBudget(t, pool, "budget", userID)

	cleanup := func() {
		pool.Exec(context.Background(), "truncate categories; truncate category_groups;")
	}

	t.Run("can create", func(t *testing.T) {
		defer cleanup()
		group1 := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budgetID}
		group2 := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group2", BudgetID: budgetID}
		require.Nil(t, categoryRepository.CreateGroup(context.Background(), group1))
		require.Nil(t, categoryRepository.CreateGroup(context.Background(), group2))

		groups, err := categoryRepository.GetGroupsForBudget(context.Background(), budgetID)
		require.Nil(t, err)
		require.Len(t, groups, 2)
		assert.True(t, reflect.DeepEqual(group1, groups[0]))
		assert.True(t, reflect.DeepEqual(group2, groups[1]))

		category1 := &beans.Category{ID: beans.NewBeansID(), GroupID: group1.ID, Name: "cat 1", BudgetID: budgetID}
		category2 := &beans.Category{ID: beans.NewBeansID(), GroupID: group2.ID, Name: "cat 2", BudgetID: budgetID, IsIncome: true}
		require.Nil(t, categoryRepository.Create(context.Background(), category1))
		require.Nil(t, categoryRepository.Create(context.Background(), category2))

		categories, err := categoryRepository.GetForBudget(context.Background(), budgetID)
		require.Nil(t, err)
		require.Len(t, categories, 2)
		assert.True(t, reflect.DeepEqual(category1, categories[0]))
		assert.True(t, reflect.DeepEqual(category2, categories[1]))
	})

	t.Run("can get single category", func(t *testing.T) {
		defer cleanup()
		group := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budgetID}
		require.Nil(t, categoryRepository.CreateGroup(context.Background(), group))

		_, err := categoryRepository.GetSingleForBudget(context.Background(), beans.NewBeansID(), budgetID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		category := &beans.Category{ID: beans.NewBeansID(), GroupID: group.ID, Name: "cat 1", BudgetID: budgetID, IsIncome: true}
		require.Nil(t, categoryRepository.Create(context.Background(), category))

		res, err := categoryRepository.GetSingleForBudget(context.Background(), category.ID, budgetID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(category, res))
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		defer cleanup()
		group := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budgetID}
		require.Nil(t, categoryRepository.CreateGroup(context.Background(), group))
		assertPgError(t, pgerrcode.UniqueViolation, categoryRepository.CreateGroup(context.Background(), group))

		category := &beans.Category{ID: beans.NewBeansID(), GroupID: group.ID, Name: "cat 1", BudgetID: budgetID}
		require.Nil(t, categoryRepository.Create(context.Background(), category))
		assertPgError(t, pgerrcode.UniqueViolation, categoryRepository.Create(context.Background(), category))
	})

	t.Run("group exists", func(t *testing.T) {
		defer cleanup()
		group := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budgetID}
		res, err := categoryRepository.GroupExists(context.Background(), budgetID, group.ID)
		require.Nil(t, err)
		require.False(t, res)

		require.Nil(t, categoryRepository.CreateGroup(context.Background(), group))

		res, err = categoryRepository.GroupExists(context.Background(), budgetID, group.ID)
		require.Nil(t, err)
		require.True(t, res)
	})

	t.Run("group exists checks budget id", func(t *testing.T) {
		defer cleanup()
		group := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budgetID}
		res, err := categoryRepository.GroupExists(context.Background(), beans.NewBeansID(), group.ID)
		require.Nil(t, err)
		require.False(t, res)
	})
}
