package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryCreate(t *testing.T) {
	categoryService := new(mocks.CategoryService)
	sv := &Server{categoryService: categoryService}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}
	group := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "Group1", BudgetID: budget.ID}
	category := &beans.Category{ID: beans.NewBeansID(), BudgetID: budget.ID, GroupID: group.ID, Name: "Category"}

	t.Run("create calls service and returns response", func(t *testing.T) {
		call := categoryService.On("CreateCategory", mock.Anything, category.BudgetID, category.GroupID, category.Name).Return(category, nil)
		defer call.Unset()

		req := fmt.Sprintf(`{
      "group_id": "%s",
      "name": "%s"
    }`, category.GroupID, category.Name)
		resp := testutils.HTTP(t, sv.handleCategoryCreate(), budget, req, http.StatusOK)
		assert.JSONEq(t, resp, fmt.Sprintf(`{"data":{
      "id": "%s",
      "group_id": "%s",
      "name": "Category"
    }}`, category.ID, category.GroupID))
	})
}

func TestCategoryGroupCreate(t *testing.T) {
	categoryService := new(mocks.CategoryService)
	sv := &Server{categoryService: categoryService}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}
	group := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "Group1", BudgetID: budget.ID}

	t.Run("create calls service and returns response", func(t *testing.T) {
		call := categoryService.On("CreateGroup", mock.Anything, group.BudgetID, group.Name).Return(group, nil)
		defer call.Unset()

		req := fmt.Sprintf(`{
      "name": "%s"
    }`, group.Name)
		resp := testutils.HTTP(t, sv.handleCategoryGroupCreate(), budget, req, http.StatusOK)
		assert.JSONEq(t, resp, fmt.Sprintf(`{"data":{
      "id": "%s",
      "name": "%s",
      "categories": []
    }}`, group.ID, group.Name))
	})
}

func TestGetCategories(t *testing.T) {
	categoryRepository := new(mocks.CategoryRepository)
	sv := &Server{categoryRepository: categoryRepository}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}
	group1 := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "Group1", BudgetID: budget.ID}
	group2 := &beans.CategoryGroup{ID: beans.NewBeansID(), Name: "Group2", BudgetID: budget.ID}
	category1 := &beans.Category{ID: beans.NewBeansID(), Name: "Cat1", BudgetID: budget.ID, GroupID: group1.ID}
	category2 := &beans.Category{ID: beans.NewBeansID(), Name: "Cat2", BudgetID: budget.ID, GroupID: group1.ID}

	categoryRepository.On("GetGroupsForBudget", mock.Anything, budget.ID).Return([]*beans.CategoryGroup{group1, group2}, nil)
	categoryRepository.On("GetForBudget", mock.Anything, budget.ID).Return([]*beans.Category{category1, category2}, nil)

	resp := testutils.HTTP(t, sv.handleCategoryGetAll(), budget, nil, http.StatusOK)
	assert.JSONEq(t, resp, fmt.Sprintf(`{"data":[
    {
      "id": "%s",
      "name": "%s",
      "categories": [
        {"id": "%s", "name": "%s"},
        {"id": "%s", "name": "%s"}
      ]
    },
    {
      "id": "%s",
      "name": "%s",
      "categories": []
    }
  ]}`, group1.ID, group1.Name, category1.ID, category1.Name, category2.ID, category2.Name, group2.ID, group2.Name))
}
