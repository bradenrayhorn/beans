package beans

import (
	"context"
)

type Budget struct {
	ID   ID
	Name Name

	// Must be explicitly loaded.
	UserIDs []UserID
}

func (b *Budget) UserHasAccess(userID UserID) bool {
	for _, id := range b.UserIDs {
		if id == userID {
			return true
		}
	}

	return false
}

type BudgetRepository interface {
	// Creates a budget and assigns user to the budget.
	Create(ctx context.Context, tx Tx, id ID, name Name, userID UserID) error
	// Gets budget by ID. Attaches UserIDs field.
	Get(ctx context.Context, id ID) (*Budget, error)
	GetBudgetsForUser(ctx context.Context, userID UserID) ([]*Budget, error)
}

type BudgetService interface {
	CreateBudget(ctx context.Context, name Name, userID UserID) (*Budget, error)
}
