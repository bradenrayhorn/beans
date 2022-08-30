package beans

import (
	"context"
)

type Transaction struct {
	ID ID

	AccountID ID

	Amount Amount
	Date   Date
	Notes  TransactionNotes
}

type TransactionNotes ValidatableString

type TransactionRepository interface {
	Create(ctx context.Context, id ID, accountID ID, amount Amount, date Date, notes TransactionNotes) error
}

type TransactionCreate struct {
	AccountID ID
	Amount    Amount
	Date      Date
	Notes     TransactionNotes
}

func (t TransactionCreate) ValidateAll() error {
	return ValidateFields(
		Field("Account ID", Required(t.AccountID)),
		Field("Amount", Required(&t.Amount)),
		Field("Date", Required(t.Date)),
		Field("Notes", Max(ValidatableString(t.Notes), 255, "characters")),
	)
}

type TransactionService interface {
	Create(ctx context.Context, t TransactionCreate) (*Transaction, error)
}
