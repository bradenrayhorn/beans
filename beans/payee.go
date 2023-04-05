package beans

import "context"

type Payee struct {
	ID       ID
	BudgetID ID
	Name     Name
}

type PayeeContract interface {
	// Creates a payee.
	CreatePayee(ctx context.Context, auth *BudgetAuthContext, name Name) (*Payee, error)

	// Gets all payees for a budget.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]*Payee, error)
}

type PayeeRepository interface {
	Create(ctx context.Context, payee *Payee) error
	GetForBudget(ctx context.Context, budgetID ID) ([]*Payee, error)
}
