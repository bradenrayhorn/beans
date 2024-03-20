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

	TransferID ID
	SplitID    ID
	IsSplit    bool
}

type Split struct {
	ID       ID
	Amount   Amount
	Notes    TransactionNotes
	Category RelatedCategory
}

type TransactionAsSplit struct {
	Transaction
	Split
}

type TransactionWithRelations struct {
	ID     ID
	Amount Amount
	Date   Date
	Notes  TransactionNotes

	Variant TransactionVariant

	Account         RelatedAccount
	Category        Optional[RelatedCategory]
	Payee           Optional[RelatedPayee]
	TransferAccount Optional[RelatedAccount]
}

type TransactionNotes struct{ NullString }

func NewTransactionNotes(string string) TransactionNotes {
	return TransactionNotes{NullString: NewNullString(string)}
}

type TransactionVariant string

const (
	TransactionStandard  TransactionVariant = "standard"
	TransactionOffBudget TransactionVariant = "off_budget"
	TransactionTransfer  TransactionVariant = "transfer"
	TransactionSplit     TransactionVariant = "split"
)

type TransactionContract interface {
	// Creates a transaction.
	Create(ctx context.Context, auth *BudgetAuthContext, params TransactionCreateParams) (ID, error)

	// Gets all transactions for budget. Excludes splits.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]TransactionWithRelations, error)

	// Get splits for a transaction.
	GetSplits(ctx context.Context, auth *BudgetAuthContext, id ID) ([]Split, error)

	// Edits a transaction.
	Update(ctx context.Context, auth *BudgetAuthContext, params TransactionUpdateParams) error

	// Deletes transactions.
	Delete(ctx context.Context, auth *BudgetAuthContext, transactionIDs []ID) error

	// Gets a transaction details.
	Get(ctx context.Context, auth *BudgetAuthContext, id ID) (TransactionWithRelations, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, transactions []Transaction) error

	Update(ctx context.Context, transactions []Transaction) error

	Delete(ctx context.Context, budgetID ID, transactionIDs []ID) error

	// Gets all transactions for budget. Excludes splits.
	GetForBudget(ctx context.Context, budgetID ID) ([]TransactionWithRelations, error)

	// Gets a single transaction for budget.
	GetWithRelations(ctx context.Context, budgetID ID, id ID) (TransactionWithRelations, error)

	// Gets all splits for a transaction.
	GetSplits(ctx context.Context, budgetID ID, transactionID ID) ([]TransactionAsSplit, error)

	// Get transaction.
	Get(ctx context.Context, budgetID ID, id ID) (Transaction, error)

	// Gets sum of all income transactions between the dates.
	GetIncomeBetween(ctx context.Context, budgetID ID, begin Date, end Date) (Amount, error)

	// Gets sum of transactions grouped by category between the dates.
	GetActivityByCategory(ctx context.Context, budgetID ID, from Date, to Date) (map[ID]Amount, error)
}

type TransactionParams struct {
	AccountID  ID
	CategoryID ID
	PayeeID    ID
	Amount     Amount
	Date       Date
	Notes      TransactionNotes
	Splits     []SplitParams
}

type SplitParams struct {
	Amount     Amount
	CategoryID ID
	Notes      TransactionNotes
}

type TransactionCreateParams struct {
	TransferAccountID ID
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
	err := ValidateFields(
		Field("Account ID", Required(t.AccountID)),
		Field("Amount", Required(&t.Amount), MaxPrecision(t.Amount)),
		Field("Date", Required(t.Date)),
		Field("Notes", Max(t.Notes, 255, "characters")),
	)
	if err != nil {
		return err
	}

	sum := NewAmount(0, 0)
	for _, s := range t.Splits {
		err := ValidateFields(
			Field("Category ID", Required(s.CategoryID)),
			Field("Amount", Required(&s.Amount), MaxPrecision(s.Amount)),
			Field("Notes", Max(s.Notes, 255, "characters")),
		)
		if err != nil {
			return err
		}
		sum, err = Arithmetic.Add(sum, s.Amount)
		if err != nil {
			return err
		}
	}

	if len(t.Splits) > 0 && sum.Compare(t.Amount) != 0 {
		return NewError(EINVALID, "Splits must sum to transaction.")
	}

	return nil
}

// helpers

func GetTransactionVariant(
	account RelatedAccount,
	transferAccount Optional[RelatedAccount],
	isSplit bool,
) TransactionVariant {
	if isSplit {
		return TransactionSplit
	}

	if transferAccount, ok := transferAccount.Value(); ok {

		// only a transfer variant if both accounts have same on/off budget
		if transferAccount.OffBudget == account.OffBudget {
			return TransactionTransfer
		}
	}

	if account.OffBudget {
		return TransactionOffBudget
	} else {
		return TransactionStandard
	}
}
