package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/assert"
)

func TestCategory(t *testing.T) {
	user := &beans.User{ID: beans.NewBeansID()}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1", UserIDs: []beans.ID{user.ID}}
	group := &beans.CategoryGroup{ID: beans.NewBeansID(), BudgetID: budget.ID, Name: "Group1"}
	category := &beans.Category{ID: beans.NewBeansID(), BudgetID: budget.ID, Name: "Category1", GroupID: group.ID}

	t.Run("create category", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.categoryContract.CreateCategoryFunc.PushReturn(category, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "POST",
			path:   "/api/v1/categories",
			body:   fmt.Sprintf(`{"name":"Category1","group_id":"%s"}`, category.GroupID),
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"data":{"id":"%s"}}`, category.ID,
		), res.body)

		params := test.categoryContract.CreateCategoryFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, category.GroupID, params.Arg2)
		assert.Equal(t, "Category1", string(params.Arg3))
	})

	t.Run("create group", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.categoryContract.CreateGroupFunc.PushReturn(group, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "POST",
			path:   "/api/v1/categories/groups",
			body:   `{"name":"Group1"}`,
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"data":{"id":"%s"}}`, group.ID,
		), res.body)

		params := test.categoryContract.CreateGroupFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, "Group1", string(params.Arg2))
	})

	t.Run("get all", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.categoryContract.GetAllFunc.PushReturn([]*beans.CategoryGroup{group}, []*beans.Category{category}, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "GET",
			path:   "/api/v1/categories",
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"data":[{"name":"Group1","id":"%s","is_income":false,"categories":[{"id":"%s","name":"Category1"}]}]}`, group.ID, category.ID,
		), res.body)
	})
}
