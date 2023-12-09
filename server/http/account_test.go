package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	user := &beans.User{ID: beans.NewBeansID()}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1", UserIDs: []beans.ID{user.ID}}
	account := &beans.Account{ID: beans.NewBeansID(), BudgetID: budget.ID, Name: "Account1", Balance: beans.NewAmount(4, 0)}

	t.Run("create", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.accountContract.CreateFunc.PushReturn(account, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "POST",
			path:   "/api/v1/accounts",
			body:   `{"name":"Account1"}`,
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"data":{"id":"%s"}}`,
			account.ID,
		), res.body)

		params := test.accountContract.CreateFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, account.Name, params.Arg2)
	})

	t.Run("get all", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.accountContract.GetAllFunc.PushReturn([]*beans.Account{account}, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "GET",
			path:   "/api/v1/accounts",
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"data":[{"name":"Account1","id":"%s","balance":{"coefficient":4,"exponent":0}}]}`,
			account.ID,
		), res.body)
	})
}
