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

	// Edits a transaction.
	Update(ctx context.Context, auth *BudgetAuthContext, params TransactionUpdateParams) error
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error

	Update(ctx context.Context, transaction *Transaction) error

	// Attaches Account, CategoryName fields to Transactions.
	GetForBudget(ctx context.Context, budgetID ID) ([]*Transaction, error)

	// Get transaction. Attaches Account field to Transaction.
	Get(ctx context.Context, id ID) (*Transaction, error)

	// Gets sum of all income transactions between the dates.
	GetIncomeBetween(ctx context.Context, budgetID ID, begin Date, end Date) (Amount, error)
}

type TransactionParams struct {
	AccountID  ID
	CategoryID ID
	Amount     Amount
	Date       Date
	Notes      TransactionNotes
}

type TransactionCreateParams struct {
	TransactionParams
}

type TransactionUpdateParams struct {
	ID ID
	TransactionParams
}

func (t TransactionUpdateParams) ValidateAll() error {
	if err := t.TransactionParams.ValidateAll(); err != nil {
		return err
	}

	return ValidateFields(Field("Transaction ID", Required(t.ID)))
}

func (t TransactionParams) ValidateAll() error {
	return ValidateFields(
		Field("Account ID", Required(t.AccountID)),
		Field("Amount", Required(&t.Amount), MaxPrecision(t.Amount)),
		Field("Date", Required(t.Date)),
		Field("Notes", Max(t.Notes, 255, "characters")),
	)
}
