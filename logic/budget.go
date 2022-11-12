package logic

import (
	"context"
	"time"

	"github.com/bradenrayhorn/beans/beans"
)

type budgetService struct {
	budgetRepository beans.BudgetRepository
}

func NewBudgetService(br beans.BudgetRepository) *budgetService {
	return &budgetService{budgetRepository: br}
}

func (s *budgetService) CreateBudget(ctx context.Context, name beans.Name, userID beans.UserID) (*beans.Budget, error) {
	if err := beans.ValidateFields(beans.Field("Budget name", name)); err != nil {
		return nil, err
	}

	budgetID := beans.NewBeansID()

	if err := s.budgetRepository.Create(ctx, budgetID, name, userID, time.Now()); err != nil {
		return nil, err
	}

	return &beans.Budget{
		ID:   budgetID,
		Name: name,
	}, nil
}
