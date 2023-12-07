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

func TestTransaction(t *testing.T) {
	contract := mocks.NewMockTransactionContract()
	sv := Server{transactionContract: contract}

	user := &beans.User{ID: beans.NewBeansID()}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1", UserIDs: []beans.ID{user.ID}}
	account := &beans.Account{ID: beans.NewBeansID(), Name: "Accounty", BudgetID: budget.ID}
	category := &beans.Category{ID: beans.NewBeansID(), Name: "Cool Category", BudgetID: budget.ID, GroupID: beans.NewBeansID()}
	payee := &beans.Payee{ID: beans.NewBeansID(), Name: "Good Payee", BudgetID: budget.ID}

	transaction := &beans.Transaction{
		ID:         beans.NewBeansID(),
		AccountID:  account.ID,
		CategoryID: category.ID,
		PayeeID:    payee.ID,
		Amount:     beans.NewAmount(5, 0),
		Date:       testutils.NewDate(t, "2022-01-09"),
		Notes:      beans.NewTransactionNotes("hi there"),

		Account:      account,
		CategoryName: beans.NewNullString("Cool Category"),
		PayeeName:    beans.NewNullString(string(payee.Name)),
	}

	t.Run("create", func(t *testing.T) {
		contract.CreateFunc.PushReturn(transaction, nil)

		req := fmt.Sprintf(`{"account_id":"%s","category_id":"%s","payee_id":"%s","amount":5,"date":"2022-01-09","notes":"hi there"}`, account.ID, category.ID, payee.ID)
		res := testutils.HTTP(t, sv.handleTransactionCreate(), user, budget, req, http.StatusOK)

		expected := fmt.Sprintf(`{"data":{"transaction_id":"%s"}}`, transaction.ID)
		assert.JSONEq(t, expected, res)

		params := contract.CreateFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, beans.TransactionCreateParams{
			TransactionParams: beans.TransactionParams{
				AccountID:  transaction.AccountID,
				CategoryID: transaction.CategoryID,
				PayeeID:    transaction.PayeeID,
				Amount:     transaction.Amount,
				Date:       transaction.Date,
				Notes:      transaction.Notes,
			},
		}, params.Arg2)
	})

	t.Run("update", func(t *testing.T) {
		contract.UpdateFunc.PushReturn(nil)

		req := fmt.Sprintf(`{"account_id":"%s","category_id":"%s","payee_id":"%s","amount":5,"date":"2022-01-09","notes":"hi there"}`, account.ID, category.ID, payee.ID)
		options := &testutils.HTTPOptions{URLParams: map[string]string{"transactionID": transaction.ID.String()}}
		_ = testutils.HTTPWithOptions(t, sv.handleTransactionUpdate(), options, user, budget, req, http.StatusOK)

		params := contract.UpdateFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, beans.TransactionUpdateParams{
			ID: transaction.ID,
			TransactionParams: beans.TransactionParams{
				AccountID:  transaction.AccountID,
				CategoryID: transaction.CategoryID,
				PayeeID:    transaction.PayeeID,
				Amount:     transaction.Amount,
				Date:       transaction.Date,
				Notes:      transaction.Notes,
			},
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
			"payee": {
				"id": "%s",
				"name": "Good Payee"
			},
			"amount": {
				"coefficient": 5,
				"exponent": 0
			},
			"date": "2022-01-09",
			"notes": "hi there"
		}]}`, transaction.ID, account.ID, category.ID, payee.ID)

		assert.JSONEq(t, expected, res)
	})
}
