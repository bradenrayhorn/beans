package logic_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/stretchr/testify/assert"
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

	categoryRepository := mocks.NewMockCategoryRepository()
	svc := logic.NewCategoryService(categoryRepository)

	t.Run("fields are required", func(t *testing.T) {
		nilID, _ := beans.BeansIDFromString("")
		_, err := svc.CreateCategory(context.Background(), nilID, nilID, "")
		testutils.AssertError(t, err, "Budget ID is required. Group ID is required. Name is required.")
	})

	t.Run("can create", func(t *testing.T) {
		categoryRepository.GroupExistsFunc.SetDefaultReturn(true, nil)

		category, err := svc.CreateCategory(context.Background(), budget.ID, group.ID, "My Cat")
		require.Nil(t, err)
		assert.Equal(t, budget.ID, category.BudgetID)
		assert.Equal(t, group.ID, category.GroupID)
		assert.Equal(t, beans.Name("My Cat"), category.Name)
	})

	t.Run("cannot create with invalid group", func(t *testing.T) {
		categoryRepository.GroupExistsFunc.SetDefaultReturn(false, nil)

		_, err := svc.CreateCategory(context.Background(), budget.ID, group.ID, "My Cat")
		testutils.AssertError(t, err, "Invalid Group ID.")
	})

	t.Run("cannot create with group check error", func(t *testing.T) {
		categoryRepository.GroupExistsFunc.SetDefaultReturn(false, errors.New("no"))

		_, err := svc.CreateCategory(context.Background(), budget.ID, group.ID, "My Cat")
		assert.Error(t, err, "no")
	})
}

func TestCreateCategoryGroup(t *testing.T) {
	budget := &beans.Budget{
		ID:   beans.NewBeansID(),
		Name: "Budget1",
	}

	categoryRepository := mocks.NewMockCategoryRepository()
	svc := logic.NewCategoryService(categoryRepository)

	t.Run("fields are required", func(t *testing.T) {
		nilID, _ := beans.BeansIDFromString("")
		_, err := svc.CreateGroup(context.Background(), nilID, "")
		testutils.AssertError(t, err, "Budget ID is required. Name is required.")
	})

	t.Run("can create", func(t *testing.T) {
		group, err := svc.CreateGroup(context.Background(), budget.ID, "My Group")
		require.Nil(t, err)
		assert.Equal(t, budget.ID, group.BudgetID)
		assert.Equal(t, beans.Name("My Group"), group.Name)
	})
}
