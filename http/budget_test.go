package http

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestBudget(t *testing.T) {
	contract := mocks.NewMockBudgetContract()
	sv := Server{budgetContract: contract}

	user := &beans.User{ID: beans.UserID(beans.NewBeansID())}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}

	t.Run("create", func(t *testing.T) {
		contract.CreateFunc.PushReturn(budget, nil)

		req := `{"name":"Budget1"}`
		res := testutils.HTTP(t, sv.handleBudgetCreate(), user, nil, req, http.StatusOK)

		expected := fmt.Sprintf(`{"data":{"name":"Budget1","id":"%s"}}`, budget.ID)
		assert.JSONEq(t, expected, res)

		params := contract.CreateFunc.History()[0]
		assert.Equal(t, user.ID, params.Arg2)
		assert.Equal(t, budget.Name, params.Arg1)
	})

	t.Run("get", func(t *testing.T) {
		month := &beans.Month{ID: beans.NewBeansID(), BudgetID: budget.ID, Date: beans.NewDate(time.Now())}
		contract.GetFunc.PushReturn(budget, month, nil)

		options := &testutils.HTTPOptions{URLParams: map[string]string{"budgetID": budget.ID.String()}}
		res := testutils.HTTPWithOptions(t, sv.handleBudgetGet(), options, user, nil, nil, http.StatusOK)
		expected := fmt.Sprintf(`{"data":{"name":"Budget1","id":"%s","latest_month_id":"%s"}}`, budget.ID, month.ID)

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
