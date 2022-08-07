package beans

import (
	"context"
	"errors"
	"strings"
)

type Budget struct {
	ID   ID
	Name BudgetName

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

type BudgetName string

func (n BudgetName) Validate() error {
	trimmedName := strings.TrimSpace(string(n))
	if trimmedName == "" {
		return errors.New("Budget name is required")
	}
	if len(trimmedName) > 255 {
		return errors.New("Budget name must be at most 255 characters")
	}
	return nil
}

type BudgetRepository interface {
	Create(ctx context.Context, id ID, name BudgetName, userID UserID) error
	// Gets budget by ID. Attaches UserIDs field.
	Get(ctx context.Context, id ID) (*Budget, error)
	GetBudgetsForUser(ctx context.Context, userID UserID) ([]*Budget, error)
}

type BudgetService interface {
	CreateBudget(ctx context.Context, name BudgetName, userID UserID) (*Budget, error)
}
