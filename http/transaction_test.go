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

func TestTransaction(t *testing.T) {
	contract := mocks.NewMockTransactionContract()
	sv := Server{transactionContract: contract}

	user := &beans.User{ID: beans.UserID(beans.NewBeansID())}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}
	account := &beans.Account{ID: beans.NewBeansID(), Name: "Accounty", BudgetID: budget.ID}
	category := &beans.Category{ID: beans.NewBeansID(), Name: "Cool Category", BudgetID: budget.ID, GroupID: beans.NewBeansID()}

	transaction := &beans.Transaction{
		ID:         beans.NewBeansID(),
		AccountID:  account.ID,
		CategoryID: category.ID,
		Amount:     beans.NewAmount(5, 0),
		Date:       testutils.NewDate(t, "2022-01-09"),
		Notes:      beans.NewTransactionNotes("hi there"),

		Account:      account,
		CategoryName: beans.NewNullString("Cool Category"),
	}

	t.Run("create", func(t *testing.T) {
		contract.CreateFunc.PushReturn(transaction, nil)

		req := fmt.Sprintf(`{"account_id":"%s","category_id":"%s","amount":5,"date":"2022-01-09","notes":"hi there"}`, account.ID, category.ID)
		res := testutils.HTTP(t, sv.handleTransactionCreate(), user, budget, req, http.StatusOK)

		expected := fmt.Sprintf(`{"data":{"transaction_id":"%s"}}`, transaction.ID)
		assert.JSONEq(t, expected, res)

		params := contract.CreateFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1)
		assert.Equal(t, beans.TransactionCreateParams{
			AccountID:  transaction.AccountID,
			CategoryID: transaction.CategoryID,
			Amount:     transaction.Amount,
			Date:       transaction.Date,
			Notes:      transaction.Notes,
		}, params.Arg2)
	})

	t.Run("get", func(t *testing.T) {
		contract.GetAllFunc.PushReturn([]*beans.Transaction{transaction}, nil)

		res := testutils.HTTP(t, sv.handleTransactionGetAll(), user, budget, nil, http.StatusOK)

		expected := fmt.Sprintf(`{"data": [{
			"id": "%s",
			"account": {
				"id": "%s",
				"name": "Accounty"
			},
			"category": {
				"id": "%s",
				"name": "Cool Category"
			},
			"amount": {
				"coefficient": 5,
				"exponent": 0
			},
			"date": "2022-01-09",
			"notes": "hi there"
		}]}`, transaction.ID, account.ID, category.ID)

		assert.JSONEq(t, expected, res)
	})
}
