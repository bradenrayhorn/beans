package logic

import (
	"context"
	"errors"
	"time"

	"github.com/bradenrayhorn/beans/beans"
)

type monthService struct {
	monthRepository beans.MonthRepository
}

func NewMonthService(monthRepository beans.MonthRepository) *monthService {
	return &monthService{monthRepository}
}
func (s *monthService) GetOrCreate(ctx context.Context, budgetID beans.ID, date time.Time) (*beans.Month, error) {
	normalizedMonth := beans.NormalizeMonth(date)

	res, err := s.monthRepository.GetByDate(ctx, budgetID, normalizedMonth)
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			month := &beans.Month{ID: beans.NewBeansID(), BudgetID: budgetID, Date: beans.NewDate(normalizedMonth)}
			return month, s.monthRepository.Create(ctx, month)
		}
		return nil, err
	}

	return res, nil
}
