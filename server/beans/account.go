package beans

import "context"

type Account struct {
	ID   ID
	Name Name

	BudgetID ID
}

type AccountWithBalance struct {
	Account
	Balance Amount
}

type RelatedAccount struct {
	ID   ID
	Name Name
}

type AccountContract interface {
	// Creates an account.
	Create(ctx context.Context, auth *BudgetAuthContext, name Name) (ID, error)

	// Gets all accounts associated with the budget.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]AccountWithBalance, error)

	// Gets an account's details.
	Get(ctx context.Context, auth *BudgetAuthContext, id ID) (Account, error)
}

type AccountRepository interface {
	Create(ctx context.Context, id ID, name Name, budgetID ID) error
	Get(ctx context.Context, id ID) (Account, error)
	GetForBudget(ctx context.Context, budgetID ID) ([]AccountWithBalance, error)
}
