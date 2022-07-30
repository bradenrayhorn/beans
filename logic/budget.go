package logic

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
)

type budgetService struct {
	budgetRepository beans.BudgetRepository
}

func NewBudgetService(br beans.BudgetRepository) *budgetService {
	return &budgetService{budgetRepository: br}
}

func (s *budgetService) CreateBudget(ctx context.Context, name beans.BudgetName, userID beans.UserID) (*beans.Budget, error) {
	if err := beans.Validate(name); err != nil {
		return nil, err
	}

	budgetID := beans.NewBeansID()

	if err := s.budgetRepository.Create(ctx, budgetID, name, userID); err != nil {
		return nil, err
	}

	return &beans.Budget{
		ID:   budgetID,
		Name: name,
	}, nil
}
