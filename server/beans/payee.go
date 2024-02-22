package beans

import "context"

type Payee struct {
	ID       ID
	BudgetID ID
	Name     Name
}

type RelatedPayee struct {
	ID   ID
	Name Name
}

type PayeeContract interface {
	// Creates a payee.
	CreatePayee(ctx context.Context, auth *BudgetAuthContext, name Name) (ID, error)

	// Gets all payees for a budget.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]Payee, error)

	// Gets a payee details.
	Get(ctx context.Context, auth *BudgetAuthContext, id ID) (Payee, error)
}

type PayeeRepository interface {
	Create(ctx context.Context, payee Payee) error
	Get(ctx context.Context, budgetID ID, id ID) (Payee, error)
	GetForBudget(ctx context.Context, budgetID ID) ([]Payee, error)
}
