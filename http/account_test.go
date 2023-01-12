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

func TestAccount(t *testing.T) {
	contract := mocks.NewMockAccountContract()
	sv := Server{accountContract: contract}

	user := &beans.User{ID: beans.NewBeansID()}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1", UserIDs: []beans.ID{user.ID}}
	account := &beans.Account{ID: beans.NewBeansID(), BudgetID: budget.ID, Name: "Account1"}

	t.Run("create", func(t *testing.T) {
		contract.CreateFunc.PushReturn(account, nil)

		req := `{"name":"Account1"}`
		res := testutils.HTTP(t, sv.handleAccountCreate(), user, budget, req, http.StatusOK)

		expected := fmt.Sprintf(`{"data":{"id":"%s"}}`, account.ID)
		assert.JSONEq(t, expected, res)

		params := contract.CreateFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, account.Name, params.Arg2)
	})

	t.Run("get all", func(t *testing.T) {
		contract.GetAllFunc.PushReturn([]*beans.Account{account}, nil)

		res := testutils.HTTP(t, sv.handleAccountsGet(), user, budget, nil, http.StatusOK)
		expected := fmt.Sprintf(`{"data":[{"name":"Account1","id":"%s"}]}`, account.ID)

		assert.JSONEq(t, expected, res)
	})
}
