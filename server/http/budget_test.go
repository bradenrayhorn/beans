package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/mocks"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestBudget(t *testing.T) {
	contract := mocks.NewMockBudgetContract()
	sv := Server{budgetContract: contract}

	user := &beans.User{ID: beans.NewBeansID()}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}

	t.Run("create", func(t *testing.T) {
		contract.CreateFunc.PushReturn(budget, nil)

		req := `{"name":"Budget1"}`
		res := testutils.HTTP(t, sv.handleBudgetCreate(), user, nil, req, http.StatusOK)

		expected := fmt.Sprintf(`{"data":{"name":"Budget1","id":"%s"}}`, budget.ID)
		assert.JSONEq(t, expected, res)

		params := contract.CreateFunc.History()[0]
		assert.Equal(t, user.ID, params.Arg1.UserID())
		assert.Equal(t, budget.Name, params.Arg2)
	})

	t.Run("get", func(t *testing.T) {
		contract.GetFunc.PushReturn(budget, nil)

		options := &testutils.HTTPOptions{URLParams: map[string]string{"budgetID": budget.ID.String()}}
		res := testutils.HTTPWithOptions(t, sv.handleBudgetGet(), options, user, nil, nil, http.StatusOK)
		expected := fmt.Sprintf(`{"data":{"name":"Budget1","id":"%s"}}`, budget.ID)

		assert.JSONEq(t, expected, res)
	})

	t.Run("get all", func(t *testing.T) {
		contract.GetAllFunc.PushReturn([]*beans.Budget{budget}, nil)

		options := &testutils.HTTPOptions{URLParams: map[string]string{"budgetID": budget.ID.String()}}
		res := testutils.HTTPWithOptions(t, sv.handleBudgetGetAll(), options, user, nil, nil, http.StatusOK)
		expected := fmt.Sprintf(`{"data":[{"name":"Budget1","id":"%s"}]}`, budget.ID)

		assert.JSONEq(t, expected, res)
	})
}
