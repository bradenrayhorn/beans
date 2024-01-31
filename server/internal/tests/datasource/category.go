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
		group1 := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budget.ID}
		group2 := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group2", BudgetID: budget.ID, IsIncome: true}
		require.Nil(t, categoryRepository.CreateGroup(ctx, nil, group1))
		require.Nil(t, categoryRepository.CreateGroup(ctx, nil, group2))

		groups, err := categoryRepository.GetGroupsForBudget(ctx, budget.ID)
		require.Nil(t, err)
		testutils.IsEqualInAnyOrder(t, []*beans.CategoryGroup{group1, group2}, groups, testutils.CmpCategoryGroup)

		category1 := &beans.Category{ID: beans.NewBeansID(), GroupID: group1.ID, Name: "cat 1", BudgetID: budget.ID}
		category2 := &beans.Category{ID: beans.NewBeansID(), GroupID: group2.ID, Name: "cat 2", BudgetID: budget.ID}
		require.Nil(t, categoryRepository.Create(ctx, nil, category1))
		require.Nil(t, categoryRepository.Create(ctx, nil, category2))

		categories, err := categoryRepository.GetForBudget(ctx, budget.ID)
		require.Nil(t, err)
		testutils.IsEqualInAnyOrder(t, []*beans.Category{category1, category2}, categories, testutils.CmpCategory)
	})

	t.Run("create respects tx", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		txManager := ds.TxManager()

		group := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budget.ID}
		category := &beans.Category{ID: beans.NewBeansID(), GroupID: group.ID, Name: "cat 1", BudgetID: budget.ID}

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

		group := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budget.ID}
		require.Nil(t, categoryRepository.CreateGroup(ctx, nil, group))

		_, err := categoryRepository.GetSingleForBudget(ctx, beans.NewBeansID(), budget.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

		category := &beans.Category{ID: beans.NewBeansID(), GroupID: group.ID, Name: "cat 1", BudgetID: budget.ID}
		require.Nil(t, categoryRepository.Create(ctx, nil, category))

		res, err := categoryRepository.GetSingleForBudget(ctx, category.ID, budget.ID)
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(category, res))
	})

	t.Run("cannot create duplicate IDs", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		group := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budget.ID}
		require.Nil(t, categoryRepository.CreateGroup(ctx, nil, group))
		assert.NotNil(t, categoryRepository.CreateGroup(ctx, nil, group))

		category := &beans.Category{ID: beans.NewBeansID(), GroupID: group.ID, Name: "cat 1", BudgetID: budget.ID}
		require.Nil(t, categoryRepository.Create(ctx, nil, category))
		assert.NotNil(t, categoryRepository.Create(ctx, nil, category))
	})

	t.Run("group exists", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		group := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budget.ID}
		res, err := categoryRepository.GroupExists(ctx, budget.ID, group.ID)
		require.Nil(t, err)
		require.False(t, res)

		require.Nil(t, categoryRepository.CreateGroup(ctx, nil, group))

		res, err = categoryRepository.GroupExists(ctx, budget.ID, group.ID)
		require.Nil(t, err)
		require.True(t, res)
	})

	t.Run("group exists checks budget id", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()

		group := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "group1", BudgetID: budget.ID}
		res, err := categoryRepository.GroupExists(ctx, beans.NewBeansID(), group.ID)
		require.Nil(t, err)
		require.False(t, res)
	})
}
