package beans

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransactionParamsValidation(t *testing.T) {
	params := TransactionParams{
		AccountID:  NewID(),
		CategoryID: NewID(),
		PayeeID:    NewID(),
		Amount:     NewAmount(7, 0),
		Date:       NewDate(time.Now()),
		Notes:      NewTransactionNotes("blah"),
	}

	t.Run("account is required", func(t *testing.T) {
		params := params
		params.AccountID = EmptyID()

		_, msg := params.ValidateAll().(Error).BeansError()
		assert.Equal(t, "Account ID is required.", msg)
	})

	t.Run("amount is required", func(t *testing.T) {
		params := params
		params.Amount = NewEmptyAmount()

		_, msg := params.ValidateAll().(Error).BeansError()
		assert.Equal(t, "Amount is required.", msg)
	})

	t.Run("amount max precision", func(t *testing.T) {
		params := params
		params.Amount = NewAmount(5, -3)

		_, msg := params.ValidateAll().(Error).BeansError()
		assert.Equal(t, "Amount must have at most 2 decimal points.", msg)
	})

	t.Run("date is required", func(t *testing.T) {
		params := params
		params.Date = Date{}

		_, msg := params.ValidateAll().(Error).BeansError()
		assert.Equal(t, "Date is required.", msg)
	})

	t.Run("notes has max character count", func(t *testing.T) {
		params := params
		params.Notes = NewTransactionNotes(strings.Repeat("a", 256))

		_, msg := params.ValidateAll().(Error).BeansError()
		assert.Equal(t, "Notes must be at most 255 characters.", msg)
	})
}

func TestTransactionUpdateParamsValidation(t *testing.T) {
	params := TransactionUpdateParams{
		TransactionParams: TransactionParams{
			AccountID:  NewID(),
			CategoryID: NewID(),
			PayeeID:    NewID(),
			Amount:     NewAmount(7, 0),
			Date:       NewDate(time.Now()),
			Notes:      NewTransactionNotes("blah"),
		},
	}

	t.Run("id is required", func(t *testing.T) {
		params := params
		params.ID = EmptyID()

		_, msg := params.ValidateAll().(Error).BeansError()
		assert.Equal(t, "Transaction ID is required.", msg)
	})

	t.Run("validates other params", func(t *testing.T) {
		params := params
		params.ID = NewID()
		params.AccountID = EmptyID()

		_, msg := params.ValidateAll().(Error).BeansError()
		assert.Equal(t, "Account ID is required.", msg)
	})
}
