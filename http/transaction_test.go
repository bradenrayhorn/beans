package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransaction(t *testing.T) {
	transactionService := new(mocks.TransactionService)
	sv := &Server{transactionService: transactionService}
	user := &beans.User{ID: beans.UserID(beans.NewBeansID())}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}
	account := &beans.Account{ID: beans.NewBeansID(), Name: "Account1"}

	transaction := beans.Transaction{
		ID:        beans.NewBeansID(),
		AccountID: account.ID,
		Amount:    beans.NewAmount(1456, -2),
		Date:      testutils.NewDate(t, "2022-08-29"),
		Notes:     beans.TransactionNotes{NullString: beans.NewNullString("My Notes")},
		Account:   account,
	}

	t.Run("create returns response", func(t *testing.T) {
		call := transactionService.On("Create", mock.Anything, budget, mock.Anything).Return(&transaction, nil)
		defer call.Unset()

		resp := testutils.HTTP(t, sv.handleTransactionCreate(), user, budget, nil, http.StatusOK)
		assert.JSONEq(t, resp, fmt.Sprintf(`{"data":{
    "id": "%s",
    "account": {
      "id": "%s",
      "name": "Account1"
    },
    "amount": {
      "coefficient": 1456,
      "exponent": -2
    },
    "date": "2022-08-29",
    "notes": "My Notes"
    }}`, transaction.ID, transaction.AccountID))
	})

	t.Run("create sends data to service", func(t *testing.T) {
		call := transactionService.On("Create", mock.Anything, budget, beans.TransactionCreate{
			AccountID: transaction.AccountID,
			Amount:    transaction.Amount,
			Date:      transaction.Date,
			Notes:     transaction.Notes,
		}).Return(&transaction, nil)
		defer call.Unset()

		testutils.HTTP(t, sv.handleTransactionCreate(), user, budget, fmt.Sprintf(`{
      "account_id": "%s",
      "amount": 14.56,
      "date": "2022-08-29",
      "notes": "My Notes"
      }`, transaction.AccountID), http.StatusOK)
	})
}

func TestGetTransactions(t *testing.T) {
	transactionRepository := new(mocks.TransactionRepository)
	accountRepository := new(mocks.AccountRepository)
	sv := &Server{transactionRepository: transactionRepository, accountRepository: accountRepository}
	user := &beans.User{ID: beans.UserID(beans.NewBeansID())}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}
	account := &beans.Account{ID: beans.NewBeansID(), Name: "Account1"}

	transaction1 := &beans.Transaction{
		ID:        beans.NewBeansID(),
		AccountID: account.ID,
		Amount:    beans.NewAmount(1456, -2),
		Date:      testutils.NewDate(t, "2022-08-29"),
		Notes:     beans.NewTransactionNotes("My notes"),
		Account:   account,
	}
	transaction2 := &beans.Transaction{
		ID:        beans.NewBeansID(),
		AccountID: account.ID,
		Amount:    beans.NewAmount(1494191, 0),
		Date:      testutils.NewDate(t, "2022-08-29"),
		Account:   account,
	}
	call := transactionRepository.On("GetForBudget", mock.Anything, budget.ID).Return([]*beans.Transaction{transaction1, transaction2}, nil)
	defer call.Unset()

	resp := testutils.HTTP(t, sv.handleTransactionGetAll(), user, budget, nil, http.StatusOK)
	assert.JSONEq(t, resp, fmt.Sprintf(`{"data":[
    {
      "id": "%s",
      "account": {
        "id": "%s",
        "name": "Account1"
      },
      "amount": {
        "coefficient": 1456,
        "exponent": -2
      },
      "date": "2022-08-29",
      "notes": "My notes"
    },
    {
      "id": "%s",
      "account": {
        "id": "%s",
        "name": "Account1"
      },
      "amount": {
        "coefficient": 1494191,
        "exponent": 0 
      },
      "date": "2022-08-29",
      "notes": null
    }
    ]}`, transaction1.ID, transaction1.AccountID, transaction2.ID, transaction2.AccountID))
}
