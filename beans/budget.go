package beans

import (
	"context"
)

type Budget struct {
	ID   ID
	Name Name

	// Must be explicitly loaded.
	UserIDs []ID
}

func (b *Budget) UserHasAccess(userID ID) bool {
	for _, id := range b.UserIDs {
		if id == userID {
			return true
		}
	}

	return false
}

type BudgetContract interface {
	// Creates a budget.
	Create(ctx context.Context, auth *AuthContext, name Name) (*Budget, error)

	// Gets a budget by the budget ID.
	// Ensures the user has access to the budget.
	Get(ctx context.Context, auth *AuthContext, id ID) (*Budget, error)

	// Gets all budgets accessible to the user.
	GetAll(ctx context.Context, auth *AuthContext) ([]*Budget, error)
}

type BudgetRepository interface {
	// Creates a budget and assigns user to the budget.
	Create(ctx context.Context, tx Tx, id ID, name Name, userID ID) error
	// Gets budget by ID. Attaches UserIDs field.
	Get(ctx context.Context, id ID) (*Budget, error)
	GetBudgetsForUser(ctx context.Context, userID ID) ([]*Budget, error)
}
