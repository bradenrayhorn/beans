package beans

import (
	"context"
)

type Transaction struct {
	ID ID

	AccountID  ID
	CategoryID ID
	PayeeID    ID

	Amount Amount
	Date   Date
	Notes  TransactionNotes
}

type TransactionWithRelations struct {
	Transaction

	Account  RelatedAccount
	Category Optional[RelatedCategory]
	Payee    Optional[RelatedPayee]
}

type TransactionNotes struct{ NullString }

func NewTransactionNotes(string string) TransactionNotes {
	return TransactionNotes{NullString: NewNullString(string)}
}

type TransactionContract interface {
	// Creates a transaction.
	Create(ctx context.Context, auth *BudgetAuthContext, params TransactionCreateParams) (ID, error)

	// Gets all transactions for budget.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]TransactionWithRelations, error)

	// Edits a transaction.
	Update(ctx context.Context, auth *BudgetAuthContext, params TransactionUpdateParams) error

	// Deletes transactions.
	Delete(ctx context.Context, auth *BudgetAuthContext, transactionIDs []ID) error

	// Gets a transaction details.
	Get(ctx context.Context, auth *BudgetAuthContext, id ID) (TransactionWithRelations, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction Transaction) error

	Update(ctx context.Context, transaction Transaction) error

	Delete(ctx context.Context, budgetID ID, transactionIDs []ID) error

	GetForBudget(ctx context.Context, budgetID ID) ([]TransactionWithRelations, error)

	// Get transaction.
	Get(ctx context.Context, budgetID ID, id ID) (Transaction, error)

	// Gets sum of all income transactions between the dates.
	GetIncomeBetween(ctx context.Context, budgetID ID, begin Date, end Date) (Amount, error)

	GetActivityByCategory(ctx context.Context, budgetID ID, from Date, to Date) (map[ID]Amount, error)
}

type TransactionParams struct {
	AccountID  ID
	CategoryID ID
	PayeeID    ID
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
