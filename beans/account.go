package beans

import "context"

type Account struct {
	ID   ID
	Name Name

	BudgetID ID

	// Must be loaded explicitly.
	Balance Amount
}

type AccountContract interface {
	// Creates an account.
	Create(ctx context.Context, auth *BudgetAuthContext, name Name) (*Account, error)

	// Gets all accounts associated with the budget. Loads Balance.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]*Account, error)
}

type AccountRepository interface {
	Create(ctx context.Context, id ID, name Name, budgetID ID) error
	Get(ctx context.Context, id ID) (*Account, error)
	// Loads Balance.
	GetForBudget(ctx context.Context, budgetID ID) ([]*Account, error)
}
