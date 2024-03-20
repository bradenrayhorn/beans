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
	split := SplitParams{CategoryID: NewID(), Amount: NewAmount(7, 0), Notes: NewTransactionNotes("blah")}

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

	t.Run("splits", func(t *testing.T) {

		t.Run("category is required", func(t *testing.T) {
			params := params
			params.Splits = []SplitParams{split}
			params.Splits[0].CategoryID = EmptyID()

			_, msg := params.ValidateAll().(Error).BeansError()
			assert.Equal(t, "Category ID is required.", msg)
		})

		t.Run("amount is required", func(t *testing.T) {
			params := params
			params.Splits = []SplitParams{split}
			params.Splits[0].Amount = NewEmptyAmount()

			_, msg := params.ValidateAll().(Error).BeansError()
			assert.Equal(t, "Amount is required.", msg)
		})

		t.Run("amount max precision", func(t *testing.T) {
			params := params
			params.Splits = []SplitParams{split}
			params.Splits[0].Amount = NewAmount(5, -3)

			_, msg := params.ValidateAll().(Error).BeansError()
			assert.Equal(t, "Amount must have at most 2 decimal points.", msg)
		})

		t.Run("notes has max character count", func(t *testing.T) {
			params := params
			params.Splits = []SplitParams{split}
			params.Splits[0].Notes = NewTransactionNotes(strings.Repeat("a", 256))

			_, msg := params.ValidateAll().(Error).BeansError()
			assert.Equal(t, "Notes must be at most 255 characters.", msg)
		})

		t.Run("sum to transaction", func(t *testing.T) {
			params := params
			params.Splits = []SplitParams{split, split}

			t.Run("positive", func(t *testing.T) {
				params.Amount = NewAmount(5, 0)

				t.Run("works", func(t *testing.T) {
					params.Splits[0].Amount = NewAmount(3, 0)
					params.Splits[1].Amount = NewAmount(2, 0)

					err := params.ValidateAll()
					assert.NoError(t, err)
				})

				t.Run("too big", func(t *testing.T) {
					params.Splits[0].Amount = NewAmount(3, 0)
					params.Splits[1].Amount = NewAmount(3, 0)

					_, msg := params.ValidateAll().(Error).BeansError()
					assert.Equal(t, "Splits must sum to transaction.", msg)
				})

				t.Run("too small", func(t *testing.T) {
					params.Splits[0].Amount = NewAmount(1, 0)
					params.Splits[1].Amount = NewAmount(3, 0)

					_, msg := params.ValidateAll().(Error).BeansError()
					assert.Equal(t, "Splits must sum to transaction.", msg)
				})

				t.Run("too negative", func(t *testing.T) {
					params.Splits[0].Amount = NewAmount(-1, 0)
					params.Splits[1].Amount = NewAmount(-3, 0)

					_, msg := params.ValidateAll().(Error).BeansError()
					assert.Equal(t, "Splits must sum to transaction.", msg)
				})
			})

			t.Run("negative", func(t *testing.T) {
				params.Amount = NewAmount(-5, 0)

				t.Run("works", func(t *testing.T) {
					params.Splits[0].Amount = NewAmount(-3, 0)
					params.Splits[1].Amount = NewAmount(-2, 0)

					err := params.ValidateAll()
					assert.NoError(t, err)
				})

				t.Run("too big", func(t *testing.T) {
					params.Splits[0].Amount = NewAmount(-3, 0)
					params.Splits[1].Amount = NewAmount(-1, 0)

					_, msg := params.ValidateAll().(Error).BeansError()
					assert.Equal(t, "Splits must sum to transaction.", msg)
				})

				t.Run("too small", func(t *testing.T) {
					params.Splits[0].Amount = NewAmount(-3, 0)
					params.Splits[1].Amount = NewAmount(-3, 0)

					_, msg := params.ValidateAll().(Error).BeansError()
					assert.Equal(t, "Splits must sum to transaction.", msg)
				})

				t.Run("too positive", func(t *testing.T) {
					params.Splits[0].Amount = NewAmount(2, 0)
					params.Splits[1].Amount = NewAmount(3, 0)

					_, msg := params.ValidateAll().(Error).BeansError()
					assert.Equal(t, "Splits must sum to transaction.", msg)
				})
			})

		})
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
		isSplit         bool
		expected        TransactionVariant
	}{
		{"on-budget, no transfer", accountOnBudgetA, Optional[RelatedAccount]{}, false, TransactionStandard},
		{"off-budget, no transfer", accountOffBudgetA, Optional[RelatedAccount]{}, false, TransactionOffBudget},
		{"on-budget from on-budget", accountOnBudgetA, OptionalWrap(accountOnBudgetB), false, TransactionTransfer},
		{"on-budget from off-budget", accountOnBudgetA, OptionalWrap(accountOffBudgetA), false, TransactionStandard},
		{"off-budget from on-budget", accountOffBudgetA, OptionalWrap(accountOnBudgetA), false, TransactionOffBudget},
		{"off-budget from off-budget", accountOffBudgetA, OptionalWrap(accountOffBudgetB), false, TransactionTransfer},
		{"split", accountOnBudgetA, Optional[RelatedAccount]{}, true, TransactionSplit},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variant := GetTransactionVariant(test.account, test.transferAccount, test.isSplit)
			assert.Equal(t, test.expected, variant)
		})
	}
}
