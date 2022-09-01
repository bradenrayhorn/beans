package logic_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory(t *testing.T) {
	budget := &beans.Budget{
		ID:   beans.NewBeansID(),
		Name: "Budget1",
	}
	group := &beans.CategoryGroup{
		ID:       beans.NewBeansID(),
		Name:     "Group1",
		BudgetID: budget.ID,
	}

	t.Run("fields are required", func(t *testing.T) {
		categoryRepository := new(mocks.CategoryRepository)
		svc := logic.NewCategoryService(categoryRepository)

		nilID, _ := beans.BeansIDFromString("")
		_, err := svc.CreateCategory(context.Background(), nilID, nilID, "")
		testutils.AssertError(t, err, "Budget ID is required. Group ID is required. Name is required.")
	})

	t.Run("can create", func(t *testing.T) {
		categoryRepository := new(mocks.CategoryRepository)
		svc := logic.NewCategoryService(categoryRepository)
		categoryRepository.On("GroupExists", mock.Anything, budget.ID, group.ID).Return(true, nil)

		var category *beans.Category
		categoryRepository.On("Create", mock.Anything, mock.MatchedBy(func(c *beans.Category) bool {
			require.Equal(t, c.BudgetID, budget.ID)
			require.Equal(t, c.GroupID, group.ID)
			require.Equal(t, c.Name, beans.Name("My Cat"))
			category = c
			return true
		})).Return(nil)

		res, err := svc.CreateCategory(context.Background(), budget.ID, group.ID, "My Cat")
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(res, category))
	})

	t.Run("cannot create with invalid group", func(t *testing.T) {
		categoryRepository := new(mocks.CategoryRepository)
		svc := logic.NewCategoryService(categoryRepository)
		categoryRepository.On("GroupExists", mock.Anything, budget.ID, group.ID).Return(false, nil)

		_, err := svc.CreateCategory(context.Background(), budget.ID, group.ID, "My Cat")
		testutils.AssertError(t, err, "Invalid Group ID.")
	})

	t.Run("cannot create with group check error", func(t *testing.T) {
		categoryRepository := new(mocks.CategoryRepository)
		svc := logic.NewCategoryService(categoryRepository)
		categoryRepository.On("GroupExists", mock.Anything, budget.ID, group.ID).Return(true, errors.New("no"))

		_, err := svc.CreateCategory(context.Background(), budget.ID, group.ID, "My Cat")
		assert.Error(t, err, "no")
	})
}

func TestCreateCategoryGroup(t *testing.T) {
	budget := &beans.Budget{
		ID:   beans.NewBeansID(),
		Name: "Budget1",
	}

	t.Run("fields are required", func(t *testing.T) {
		categoryRepository := new(mocks.CategoryRepository)
		svc := logic.NewCategoryService(categoryRepository)

		nilID, _ := beans.BeansIDFromString("")
		_, err := svc.CreateGroup(context.Background(), nilID, "")
		testutils.AssertError(t, err, "Budget ID is required. Name is required.")
	})

	t.Run("can create", func(t *testing.T) {
		categoryRepository := new(mocks.CategoryRepository)
		svc := logic.NewCategoryService(categoryRepository)

		var group *beans.CategoryGroup
		categoryRepository.On("CreateGroup", mock.Anything, mock.MatchedBy(func(g *beans.CategoryGroup) bool {
			require.Equal(t, g.BudgetID, budget.ID)
			require.Equal(t, g.Name, beans.Name("My Group"))
			group = g
			return true
		})).Return(nil)

		res, err := svc.CreateGroup(context.Background(), budget.ID, "My Group")
		require.Nil(t, err)
		assert.True(t, reflect.DeepEqual(res, group))
	})
}
