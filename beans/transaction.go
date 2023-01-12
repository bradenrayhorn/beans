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

type TransactionContract interface {
	// Creates a transaction. Attaches Account field.
	Create(ctx context.Context, auth *BudgetAuthContext, params TransactionCreateParams) (*Transaction, error)

	// Gets all transactions for budget. Attaches Account, CategoryName fields.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]*Transaction, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	// Attaches Account, CategoryName fields to Transactions.
	GetForBudget(ctx context.Context, budgetID ID) ([]*Transaction, error)

	// Gets sum of all income transactions before or on the date.
	GetIncomeBeforeOrOnDate(ctx context.Context, date Date) (Amount, error)
}

type TransactionCreateParams struct {
	AccountID  ID
	CategoryID ID
	Amount     Amount
	Date       Date
	Notes      TransactionNotes
}

func (t TransactionCreateParams) ValidateAll() error {
	return ValidateFields(
		Field("Account ID", Required(t.AccountID)),
		Field("Amount", Required(&t.Amount), MaxPrecision(t.Amount)),
		Field("Date", Required(t.Date)),
		Field("Notes", Max(t.Notes, 255, "characters")),
	)
}
