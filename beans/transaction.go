package beans

import (
	"context"
)

type Transaction struct {
	ID ID

	AccountID  ID
	CategoryID ID

	Amount Amount
	Date   Date
	Notes  TransactionNotes

	// Must be explicitly loaded.
	Account *Account
	// Must be explicitly loaded.
	CategoryName NullString
}

type TransactionNotes struct{ NullString }

func NewTransactionNotes(string string) TransactionNotes {
	return TransactionNotes{NullString: NewNullString(string)}
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	// Attaches Account, CategoryName fields to Transactions.
	GetForBudget(ctx context.Context, budgetID ID) ([]*Transaction, error)
}

type TransactionCreate struct {
	AccountID  ID
	CategoryID ID
	Amount     Amount
	Date       Date
	Notes      TransactionNotes
}

func (t TransactionCreate) ValidateAll() error {
	return ValidateFields(
		Field("Account ID", Required(t.AccountID)),
		Field("Amount", Required(&t.Amount), MaxPrecision(t.Amount)),
		Field("Date", Required(t.Date)),
		Field("Notes", Max(t.Notes, 255, "characters")),
	)
}

type TransactionService interface {
	// Attaches Account field to Transaction.
	Create(ctx context.Context, activeBudget *Budget, t TransactionCreate) (*Transaction, error)
}
