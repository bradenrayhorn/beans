package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
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
		test := newHttpTest(t)
		defer test.Stop(t)

		test.transactionContract.CreateFunc.PushReturn(transaction, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "POST",
			path:   "/api/v1/transactions",
			body: fmt.Sprintf(`{
				"account_id":"%s",
				"category_id":"%s",
				"payee_id":"%s",
				"amount":5,
				"date":"2022-01-09",
				"notes":"hi there"
			}`, account.ID, category.ID, payee.ID),
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"data":{"transaction_id":"%s"}}`, transaction.ID,
		), res.body)

		params := test.transactionContract.CreateFunc.History()[0]
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
		test := newHttpTest(t)
		defer test.Stop(t)

		test.transactionContract.UpdateFunc.PushReturn(nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "PUT",
			path:   fmt.Sprintf("/api/v1/transactions/%s", transaction.ID),
			body: fmt.Sprintf(`{
				"account_id":"%s",
				"category_id":"%s",
				"payee_id":"%s",
				"amount":5,
				"date":"2022-01-09",
				"notes":"hi there"
			}`, account.ID, category.ID, payee.ID),
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Empty(t, res.body)

		params := test.transactionContract.UpdateFunc.History()[0]
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

	t.Run("delete", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		id1 := beans.NewBeansID()
		id2 := beans.NewBeansID()

		test.transactionContract.DeleteFunc.PushReturn(nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "POST",
			path:   "/api/v1/transactions/delete",
			body: fmt.Sprintf(`{
				"ids":["%s","%s"]
			}`, id1, id2),
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Empty(t, res.body)

		params := test.transactionContract.DeleteFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, []beans.ID{
			id1, id2,
		}, params.Arg2)
	})

	t.Run("get", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.transactionContract.GetAllFunc.PushReturn([]*beans.Transaction{transaction}, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "GET",
			path:   "/api/v1/transactions",
			user:   user,
			budget: budget,
		})

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

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, expected, res.body)
	})
}
