package logic

import (
	"context"
	"time"

	"github.com/bradenrayhorn/beans/beans"
)

type budgetService struct {
	budgetRepository beans.BudgetRepository
	monthRepository  beans.MonthRepository
	txManager        beans.TxManager
}

func NewBudgetService(txManager beans.TxManager, budgetRepository beans.BudgetRepository, monthRepository beans.MonthRepository) *budgetService {
	return &budgetService{budgetRepository, monthRepository, txManager}
}

func (s *budgetService) CreateBudget(ctx context.Context, name beans.Name, userID beans.UserID) (*beans.Budget, error) {
	if err := beans.ValidateFields(beans.Field("Budget name", name)); err != nil {
		return nil, err
	}

	budgetID := beans.NewBeansID()

	tx, err := s.txManager.Create(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.budgetRepository.Create(ctx, tx, budgetID, name, userID); err != nil {
		return nil, err
	}

	if err := s.monthRepository.Create(ctx, &beans.Month{
		ID:       beans.NewBeansID(),
		BudgetID: budgetID,
		Date:     beans.NewDate(beans.NormalizeMonth(time.Now())),
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &beans.Budget{
		ID:   budgetID,
		Name: name,
	}, nil
}
