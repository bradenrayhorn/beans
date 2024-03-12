package beans

import "context"

// models

type Account struct {
	ID        ID
	Name      Name
	OffBudget bool

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

// repository

type AccountRepository interface {
	Create(ctx context.Context, account Account) error
	Get(ctx context.Context, budgetID ID, id ID) (Account, error)
	GetWithBalance(ctx context.Context, budgetID ID) ([]AccountWithBalance, error)
	GetTransactable(ctx context.Context, budgetID ID) ([]Account, error)
}

// contract

type AccountContract interface {
	// Creates an account.
	Create(ctx context.Context, auth *BudgetAuthContext, params AccountCreate) (ID, error)

	// Gets all accounts associated with the budget.
	GetAll(ctx context.Context, auth *BudgetAuthContext) ([]AccountWithBalance, error)

	// Gets all accounts that can be used in a transaction.
	GetTransactable(ctx context.Context, auth *BudgetAuthContext) ([]Account, error)

	// Gets an account's details.
	Get(ctx context.Context, auth *BudgetAuthContext, id ID) (Account, error)
}

type AccountCreate struct {
	Name      Name
	OffBudget bool
}
