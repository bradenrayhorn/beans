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

func TestTransactionCreateParamsValidation(t *testing.T) {
	params := TransactionCreateParams{
		TransactionParams: TransactionParams{
			AccountID:  NewID(),
			CategoryID: NewID(),
			PayeeID:    NewID(),
			Amount:     NewAmount(7, 0),
			Date:       NewDate(time.Now()),
			Notes:      NewTransactionNotes("blah"),
		},
	}

	t.Run("validates other params", func(t *testing.T) {
		params := params
		params.AccountID = EmptyID()

		_, msg := params.ValidateAll().(Error).BeansError()
		assert.Equal(t, "Account ID is required.", msg)
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

func TestGetTransactionVariant(t *testing.T) {
	accountOnBudgetA := RelatedAccount{ID: NewID(), Name: "onA", OffBudget: false}
	accountOnBudgetB := RelatedAccount{ID: NewID(), Name: "onB", OffBudget: false}
	accountOffBudgetA := RelatedAccount{ID: NewID(), Name: "offA", OffBudget: true}
	accountOffBudgetB := RelatedAccount{ID: NewID(), Name: "offB", OffBudget: true}

	var tests = []struct {
		name            string
		account         RelatedAccount
		transferAccount Optional[RelatedAccount]
		expected        TransactionVariant
	}{
		{"on-budget, no transfer", accountOnBudgetA, Optional[RelatedAccount]{}, TransactionStandard},
		{"off-budget, no transfer", accountOffBudgetA, Optional[RelatedAccount]{}, TransactionOffBudget},
		{"on-budget from on-budget", accountOnBudgetA, OptionalWrap(accountOnBudgetB), TransactionTransfer},
		{"on-budget from off-budget", accountOnBudgetA, OptionalWrap(accountOffBudgetA), TransactionStandard},
		{"off-budget from on-budget", accountOffBudgetA, OptionalWrap(accountOnBudgetA), TransactionOffBudget},
		{"off-budget from off-budget", accountOffBudgetA, OptionalWrap(accountOffBudgetB), TransactionTransfer},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variant := GetTransactionVariant(test.account, test.transferAccount)
			assert.Equal(t, test.expected, variant)
		})
	}
}
