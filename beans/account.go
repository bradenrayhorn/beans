package beans

import "context"

type Account struct {
	ID   ID
	Name Name

	BudgetID ID
}

type AccountContract interface {
	// Creates an account.
	Create(ctx context.Context, budgetID ID, name Name) (*Account, error)

	// Gets all accounts associated with the budget.
	GetAll(ctx context.Context, budgetID ID) ([]*Account, error)
}

type AccountRepository interface {
	Create(ctx context.Context, id ID, name Name, budgetID ID) error
	Get(ctx context.Context, id ID) (*Account, error)
	GetForBudget(ctx context.Context, budgetID ID) ([]*Account, error)
}
