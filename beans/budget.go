package beans

import (
	"context"
	"errors"
	"strings"
)

type Budget struct {
	ID   ID
	Name BudgetName
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
	Get(ctx context.Context, id ID) (*Budget, error)
	GetBudgetsForUser(ctx context.Context, userID UserID) ([]*Budget, error)
}

type BudgetService interface {
	CreateBudget(ctx context.Context, name BudgetName, userID UserID) (*Budget, error)
}
