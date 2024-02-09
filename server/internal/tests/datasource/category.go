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

func TestCategoryRepository(t *testing.T, ds beans.DataSource) {
	factory := testutils.Factory(t, ds)
	categoryRepository := ds.CategoryRepository()
	ctx := context.Background()

	t.Run("can create", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		group1 := beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budget.ID}
		group2 := beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group2", BudgetID: budget.ID, IsIncome: true}
		require.Nil(t, categoryRepository.CreateGroup(ctx, nil, group1))
		require.Nil(t, categoryRepository.CreateGroup(ctx, nil, group2))

		groups, err := categoryRepository.GetGroupsForBudget(ctx, budget.ID)
		require.Nil(t, err)
		testutils.IsEqualInAnyOrder(t, []beans.CategoryGroup{group1, group2}, groups, testutils.CmpCategoryGroup)

		category1 := beans.Category{ID: beans.NewBeansID(), GroupID: group1.ID, Name: "cat 1", BudgetID: budget.ID}
		category2 := beans.Category{ID: beans.NewBeansID(), GroupID: group2.ID, Name: "cat 2", BudgetID: budget.ID}
		require.Nil(t, categoryRepository.Create(ctx, nil, category1))
		require.Nil(t, categoryRepository.Create(ctx, nil, category2))

		categories, err := categoryRepository.GetForBudget(ctx, budget.ID)
		require.Nil(t, err)
		testutils.IsEqualInAnyOrder(t, []beans.Category{category1, category2}, categories, testutils.CmpCategory)
	})

	t.Run("create respects tx", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		txManager := ds.TxManager()

		group := beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budget.ID}
		category := beans.Category{ID: beans.NewBeansID(), GroupID: group.ID, Name: "cat 1", BudgetID: budget.ID}

		tx, err := txManager.Create(ctx)
		require.Nil(t, err)
		defer testutils.MustRollback(t, tx)

		require.Nil(t, categoryRepository.CreateGroup(ctx, tx, group))
		require.Nil(t, categoryRepository.Create(ctx, tx, category))

		groups, err := categoryRepository.GetGroupsForBudget(ctx, budget.ID)
		require.Nil(t, err)
		require.Len(t, groups, 0)

		categories, err := categoryRepository.GetForBudget(ctx, budget.ID)
		require.Nil(t, err)
		require.Len(t, categories, 0)

		require.Nil(t, tx.Commit(ctx))

		groups, err = categoryRepository.GetGroupsForBudget(ctx, budget.ID)
		require.Nil(t, err)
		require.Len(t, groups, 1)

		categories, err = categoryRepository.GetForBudget(ctx, budget.ID)
		require.Nil(t, err)
		require.Len(t, categories, 1)
	})

	t.Run("can get single category", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		group := beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budget.ID}
		require.Nil(t, categoryRepository.CreateGroup(ctx, nil, group))

		_, err := categoryRepository.GetSingleForBudget(ctx, beans.NewBeansID(), budget.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		category := beans.Category{ID: beans.NewBeansID(), GroupID: group.ID, Name: "cat 1", BudgetID: budget.ID}
		require.Nil(t, categoryRepository.Create(ctx, nil, category))

		res, err := categoryRepository.GetSingleForBudget(ctx, category.ID, budget.ID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(category, res))
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		group := beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budget.ID}
		require.Nil(t, categoryRepository.CreateGroup(ctx, nil, group))
		assert.NotNil(t, categoryRepository.CreateGroup(ctx, nil, group))

		category := beans.Category{ID: beans.NewBeansID(), GroupID: group.ID, Name: "cat 1", BudgetID: budget.ID}
		require.Nil(t, categoryRepository.Create(ctx, nil, category))
		assert.NotNil(t, categoryRepository.Create(ctx, nil, category))
	})

	t.Run("get group", func(t *testing.T) {
		t.Run("cannot get non-existent group", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			_, err := categoryRepository.GetCategoryGroup(ctx, beans.NewBeansID(), budget.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot get group for another budget", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()
			group := factory.CategoryGroup(beans.CategoryGroup{BudgetID: budget2.ID})

			_, err := categoryRepository.GetCategoryGroup(ctx, group.ID, budget.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("can get", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			group := factory.CategoryGroup(beans.CategoryGroup{BudgetID: budget.ID})

			res, err := categoryRepository.GetCategoryGroup(ctx, group.ID, budget.ID)
			require.Nil(t, err)

			assert.Equal(t, group, res)
		})
	})

	t.Run("get categories for group", func(t *testing.T) {
		t.Run("empty for non-existent group", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			res, err := categoryRepository.GetCategoriesForGroup(ctx, beans.NewBeansID(), budget.ID)
			require.Nil(t, err)
			assert.Empty(t, res)
		})

		t.Run("empty for group of another budget", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()
			category := factory.Category(beans.Category{BudgetID: budget2.ID})

			res, err := categoryRepository.GetCategoriesForGroup(ctx, category.GroupID, budget.ID)
			require.Nil(t, err)
			assert.Empty(t, res)
		})

		t.Run("can get", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			category := factory.Category(beans.Category{BudgetID: budget.ID})

			// this category belongs to a different group and should not be returned
			factory.Category(beans.Category{BudgetID: budget.ID})

			res, err := categoryRepository.GetCategoriesForGroup(ctx, category.GroupID, budget.ID)
			require.Nil(t, err)

			testutils.IsEqualInAnyOrder(t, []beans.Category{category}, res, testutils.CmpCategory)
		})
	})
}
