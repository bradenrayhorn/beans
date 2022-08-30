package logic_test

import (
	"context"
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
	t.Run("fields are required", func(t *testing.T) {
		transactionRepository := new(mocks.TransactionRepository)
		svc := logic.NewTransactionService(transactionRepository)

		_, err := svc.Create(context.Background(), beans.TransactionCreate{})
		testutils.AssertError(t, err, "Account ID is required. Amount is required. Date is required.")
		transactionRepository.AssertNotCalled(t, "Create")
	})

	t.Run("create transaction", func(t *testing.T) {
		transactionRepository := new(mocks.TransactionRepository)
		svc := logic.NewTransactionService(transactionRepository)

		c := beans.TransactionCreate{
			AccountID: beans.NewBeansID(),
			Amount:    beans.NewAmount(10, 1),
			Date:      beans.NewDate(time.Now()),
			Notes:     "My notes",
		}
		transactionRepository.On("Create", mock.Anything, mock.Anything, c.AccountID, c.Amount, c.Date, c.Notes).Return(nil)

		transaction, err := svc.Create(context.Background(), c)
		require.Nil(t, err)
		assert.False(t, transaction.ID.Empty())
		assert.Equal(t, c.AccountID, transaction.AccountID)
		assert.Equal(t, c.Amount, transaction.Amount)
		assert.Equal(t, c.Date, transaction.Date)
		assert.Equal(t, c.Notes, transaction.Notes)
	})
}
