package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/assert"
)

func TestBudget(t *testing.T) {
	user := &beans.User{ID: beans.NewBeansID()}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}

	t.Run("create", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.budgetContract.CreateFunc.PushReturn(budget, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "POST",
			path:   "/api/v1/budgets",
			body:   `{"name":"Budget1"}`,
			user:   user,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"data":{"id":"%s"}}`, budget.ID,
		), res.body)

		params := test.budgetContract.CreateFunc.History()[0]
		assert.Equal(t, user.ID, params.Arg1.UserID())
		assert.Equal(t, budget.Name, params.Arg2)
	})

	t.Run("get", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.budgetContract.GetFunc.PushReturn(budget, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "GET",
			path:   fmt.Sprintf("/api/v1/budgets/%s", budget.ID),
			user:   user,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"data":{"name":"Budget1","id":"%s"}}`, budget.ID,
		), res.body)
	})

	t.Run("get all", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.budgetContract.GetAllFunc.PushReturn([]*beans.Budget{budget}, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "GET",
			path:   "/api/v1/budgets",
			user:   user,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"data":[{"name":"Budget1","id":"%s"}]}`, budget.ID,
		), res.body)
	})

}
