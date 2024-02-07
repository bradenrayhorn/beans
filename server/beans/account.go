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
	Create(ctx context.Context, auth *BudgetAuthContext, name Name) (ID, error)

	// Gets all accounts associated with the budget. Loads Balance.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]Account, error)

	// Gets an account's details. Loads Balance.
	Get(ctx context.Context, auth *BudgetAuthContext, id ID) (Account, error)
}

type AccountRepository interface {
	Create(ctx context.Context, id ID, name Name, budgetID ID) error
	Get(ctx context.Context, id ID) (Account, error)
	// Loads Balance.
	GetForBudget(ctx context.Context, budgetID ID) ([]Account, error)
}
