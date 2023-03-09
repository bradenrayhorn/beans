package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCategory(t *testing.T) {
	contract := mocks.NewMockCategoryContract()
	sv := Server{categoryContract: contract}

	user := &beans.User{ID: beans.NewBeansID()}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1", UserIDs: []beans.ID{user.ID}}
	group := &beans.CategoryGroup{ID: beans.NewBeansID(), BudgetID: budget.ID, Name: "Group1"}
	category := &beans.Category{ID: beans.NewBeansID(), BudgetID: budget.ID, Name: "Category1", GroupID: group.ID}

	t.Run("create category", func(t *testing.T) {
		contract.CreateCategoryFunc.PushReturn(category, nil)

		req := fmt.Sprintf(`{"name":"Category1","group_id":"%s"}`, category.GroupID)
		res := testutils.HTTP(t, sv.handleCategoryCreate(), user, budget, req, http.StatusOK)

		expected := fmt.Sprintf(`{"data":{"id":"%s"}}`, category.ID)
		assert.JSONEq(t, expected, res)

		params := contract.CreateCategoryFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, category.GroupID, params.Arg2)
		assert.Equal(t, "Category1", string(params.Arg3))
	})

	t.Run("create group", func(t *testing.T) {
		contract.CreateGroupFunc.PushReturn(group, nil)

		req := `{"name":"Group1"}`
		res := testutils.HTTP(t, sv.handleCategoryGroupCreate(), user, budget, req, http.StatusOK)

		expected := fmt.Sprintf(`{"data":{"id":"%s"}}`, group.ID)
		assert.JSONEq(t, expected, res)

		params := contract.CreateGroupFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, "Group1", string(params.Arg2))
	})

	t.Run("get all", func(t *testing.T) {
		contract.GetAllFunc.PushReturn([]*beans.CategoryGroup{group}, []*beans.Category{category}, nil)

		res := testutils.HTTP(t, sv.handleCategoryGetAll(), user, budget, nil, http.StatusOK)
		expected := fmt.Sprintf(`{"data":[{"name":"Group1","id":"%s","is_income":false,"categories":[{"id":"%s","name":"Category1"}]}]}`, group.ID, category.ID)

		assert.JSONEq(t, expected, res)
	})
}
