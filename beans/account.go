package beans

import "context"

type Account struct {
	ID   ID
	Name Name

	BudgetID ID
}

type AccountRepository interface {
	Create(ctx context.Context, id ID, name Name, budgetID ID) error
	Get(ctx context.Context, id ID) (*Account, error)
	GetForBudget(ctx context.Context, budgetID ID) ([]*Account, error)
}

type AccountService interface {
	Create(ctx context.Context, name Name, budgetID ID) (*Account, error)
}
