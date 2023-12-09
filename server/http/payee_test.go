package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/assert"
)

func TestPayee(t *testing.T) {
	user := &beans.User{ID: beans.NewBeansID()}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1", UserIDs: []beans.ID{user.ID}}
	payee := &beans.Payee{ID: beans.NewBeansID(), BudgetID: budget.ID, Name: "Payee"}

	t.Run("create", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.payeeContract.CreatePayeeFunc.PushReturn(payee, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "POST",
			path:   "/api/v1/payees",
			body:   `{"name":"Payee"}`,
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"data":{"id":"%s"}}`, payee.ID,
		), res.body)

		params := test.payeeContract.CreatePayeeFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, payee.Name, params.Arg2)
	})

	t.Run("get all", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.payeeContract.GetAllFunc.PushReturn([]*beans.Payee{payee}, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "GET",
			path:   "/api/v1/payees",
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"data":[{"name":"Payee","id":"%s"}]}`, payee.ID,
		), res.body)
	})
}
