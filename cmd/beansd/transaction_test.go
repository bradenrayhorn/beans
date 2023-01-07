package main_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactions(t *testing.T) {
	t.Parallel()
	ta := StartApplication(t)
	defer ta.Stop(t)

	t.Run("can create transaction", func(t *testing.T) {
		user, session := ta.CreateUserAndSession(t)
		budget := ta.CreateBudget(t, "my budget", user)
		_ = ta.CreateMonth(t, budget, testutils.NewDate(t, "2022-05-01"))
		account := ta.CreateAccount(t, "my account", budget)
		categoryGroup := ta.CreateCategoryGroup(t, budget, "group")
		category := ta.CreateCategory(t, budget, categoryGroup, "category")

		r := ta.PostRequest(t, "api/v1/transactions", newOptionsWithBody(session, budget, fmt.Sprintf(`{
      "account_id":"%s",
      "category_id":"%s",
      "amount":55,
      "date":"2022-05-12",
      "notes":"Some notes"
    }`, account.ID, category.ID)))
		assert.Equal(t, http.StatusOK, r.StatusCode)

		transactions, err := ta.application.TransactionRepository().GetForBudget(context.Background(), budget.ID)
		require.Nil(t, err)
		assert.Len(t, transactions, 1)
		assert.Equal(t, transactions[0].AccountID, account.ID)
		assert.Equal(t, transactions[0].CategoryID, category.ID)
		assert.Equal(t, transactions[0].Amount, beans.NewAmount(55, 0))
		assert.Equal(t, transactions[0].Date, testutils.NewDate(t, "2022-05-12"))
		assert.Equal(t, transactions[0].Notes, beans.NewTransactionNotes("Some notes"))

		assert.JSONEq(t, fmt.Sprintf(`{"data":{"transaction_id":"%s"}}`, transactions[0].ID), r.Body)
	})

	t.Run("can get transactions", func(t *testing.T) {
		user, session := ta.CreateUserAndSession(t)
		budget := ta.CreateBudget(t, "my budget", user)
		account := ta.CreateAccount(t, "my account", budget)
		categoryGroup := ta.CreateCategoryGroup(t, budget, "group")
		category := ta.CreateCategory(t, budget, categoryGroup, "category")

		r := ta.GetRequest(t, "api/v1/transactions", newOptions(session, budget))
		assert.Equal(t, http.StatusOK, r.StatusCode)
		assert.JSONEq(t, `{"data":[]}`, r.Body)

		transaction := &beans.Transaction{
			ID:         beans.NewBeansID(),
			AccountID:  account.ID,
			CategoryID: category.ID,
			Amount:     beans.NewAmount(55, 0),
			Date:       testutils.NewDate(t, "2022-05-12"),
			Notes:      beans.NewTransactionNotes("Some notes"),
		}
		require.Nil(t, ta.application.TransactionRepository().Create(context.Background(), transaction))

		r = ta.GetRequest(t, "api/v1/transactions", newOptions(session, budget))
		assert.Equal(t, http.StatusOK, r.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(`{"data":[{
      "id": "%s",
      "account": {
        "id": "%s",
        "name": "my account"
      },
      "category": {
        "id": "%s",
        "name": "category"
      },
      "amount": {
        "coefficient": 55,
        "exponent": 0
      },
      "date": "2022-05-12",
      "notes": "Some notes"
    }]}`, transaction.ID, account.ID, category.ID), r.Body)
	})
}
