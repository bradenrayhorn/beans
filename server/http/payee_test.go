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

func TestPayee(t *testing.T) {
	contract := mocks.NewMockPayeeContract()
	sv := Server{payeeContract: contract}

	user := &beans.User{ID: beans.NewBeansID()}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1", UserIDs: []beans.ID{user.ID}}
	payee := &beans.Payee{ID: beans.NewBeansID(), BudgetID: budget.ID, Name: "Payee"}

	t.Run("create", func(t *testing.T) {
		contract.CreatePayeeFunc.PushReturn(payee, nil)

		req := `{"name":"Payee"}`
		res := testutils.HTTP(t, sv.handlePayeeCreate(), user, budget, req, http.StatusOK)

		expected := fmt.Sprintf(`{"data":{"id":"%s"}}`, payee.ID)
		assert.JSONEq(t, expected, res)

		params := contract.CreatePayeeFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, payee.Name, params.Arg2)
	})

	t.Run("get all", func(t *testing.T) {
		contract.GetAllFunc.PushReturn([]*beans.Payee{payee}, nil)

		res := testutils.HTTP(t, sv.handlePayeeGetAll(), user, budget, nil, http.StatusOK)
		expected := fmt.Sprintf(`{"data":[{"name":"Payee","id":"%s"}]}`, payee.ID)

		assert.JSONEq(t, expected, res)
	})
}
