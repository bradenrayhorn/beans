package logic_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	budget := &beans.Budget{
		ID:   beans.NewBeansID(),
		Name: "Budget1",
	}
	account := &beans.Account{
		ID:       beans.NewBeansID(),
		Name:     "Account1",
		BudgetID: budget.ID,
	}
	t.Run("fields are required", func(t *testing.T) {
		transactionRepository := new(mocks.TransactionRepository)
		accountRepository := new(mocks.AccountRepository)
		svc := logic.NewTransactionService(transactionRepository, accountRepository)

		_, err := svc.Create(context.Background(), budget, beans.TransactionCreate{})
		testutils.AssertError(t, err, "Account ID is required. Amount is required. Date is required.")
		transactionRepository.AssertNotCalled(t, "Create")
	})

	t.Run("cannot create transaction with amount more than 2 decimals", func(t *testing.T) {
		transactionRepository := new(mocks.TransactionRepository)
		accountRepository := new(mocks.AccountRepository)
		svc := logic.NewTransactionService(transactionRepository, accountRepository)

		c := beans.TransactionCreate{
			AccountID: account.ID,
			Amount:    beans.NewAmount(10, -3),
			Date:      beans.NewDate(time.Now()),
		}
		_, err := svc.Create(context.Background(), budget, c)
		testutils.AssertError(t, err, "Amount must have at most 2 decimal points.")
	})

	t.Run("can create", func(t *testing.T) {
		transactionRepository := new(mocks.TransactionRepository)
		accountRepository := new(mocks.AccountRepository)
		svc := logic.NewTransactionService(transactionRepository, accountRepository)

		c := beans.TransactionCreate{
			AccountID: account.ID,
			Amount:    beans.NewAmount(10, 1),
			Date:      beans.NewDate(time.Now()),
			Notes:     beans.NewTransactionNotes("My Notes"),
		}
		var transaction *beans.Transaction
		transactionRepository.On("Create", mock.Anything, mock.MatchedBy(func(tr *beans.Transaction) bool {
			require.Equal(t, account.ID, tr.AccountID)
			require.Equal(t, c.Amount, tr.Amount)
			require.Equal(t, c.Date, tr.Date)
			require.Equal(t, c.Notes, tr.Notes)
			transaction = tr
			return true
		})).Return(nil)
		accountRepository.On("Get", mock.Anything, account.ID).Return(account, nil)

		createdTransaction, err := svc.Create(context.Background(), budget, c)
		require.Nil(t, err)
		assert.False(t, transaction.ID.Empty())
		assert.True(t, reflect.DeepEqual(transaction, createdTransaction))
	})

	t.Run("cannot create after account error", func(t *testing.T) {
		transactionRepository := new(mocks.TransactionRepository)
		accountRepository := new(mocks.AccountRepository)
		svc := logic.NewTransactionService(transactionRepository, accountRepository)

		c := beans.TransactionCreate{
			AccountID: account.ID,
			Amount:    beans.NewAmount(10, 1),
			Date:      beans.NewDate(time.Now()),
			Notes:     beans.NewTransactionNotes("My Notes"),
		}
		accountRepository.On("Get", mock.Anything, c.AccountID).Return(nil, errors.New("account not found"))

		_, err := svc.Create(context.Background(), budget, c)
		require.NotNil(t, err)
		assert.Errorf(t, err, "account not found")
	})

	t.Run("cannot create with account from other budget", func(t *testing.T) {
		transactionRepository := new(mocks.TransactionRepository)
		accountRepository := new(mocks.AccountRepository)
		svc := logic.NewTransactionService(transactionRepository, accountRepository)

		c := beans.TransactionCreate{
			AccountID: beans.NewBeansID(),
			Amount:    beans.NewAmount(10, 1),
			Date:      beans.NewDate(time.Now()),
			Notes:     beans.NewTransactionNotes("My notes"),
		}
		badAccount := &beans.Account{
			ID:       c.AccountID,
			Name:     "bad account",
			BudgetID: beans.NewBeansID(),
		}
		accountRepository.On("Get", mock.Anything, c.AccountID).Return(badAccount, nil)

		_, err := svc.Create(context.Background(), budget, c)
		require.NotNil(t, err)
		testutils.AssertError(t, err, "Invalid Account ID")
	})
}
