package service

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
)

type monthCategoryService struct{ service }

var _ beans.MonthCategoryService = (*monthCategoryService)(nil)

func (s *monthCategoryService) GetForMonth(ctx context.Context, month beans.Month) ([]beans.MonthCategoryWithDetails, error) {
	monthCategories, err := s.ds.MonthCategoryRepository().GetForMonth(ctx, month)
	if err != nil {
		return nil, err
	}

	previousAssigned, err := s.ds.MonthCategoryRepository().GetAssignedByCategory(ctx, month.BudgetID, month.Date.FirstDay())
	if err != nil {
		return nil, err
	}

	previousActivity, err := s.ds.TransactionRepository().GetActivityByCategory(ctx, month.BudgetID, beans.Date{}, month.Date.FirstDay().Previous())
	if err != nil {
		return nil, err
	}

	activity, err := s.ds.TransactionRepository().GetActivityByCategory(ctx, month.BudgetID, month.Date.FirstDay(), month.Date.LastDay())
	if err != nil {
		return nil, err
	}

	res := make([]beans.MonthCategoryWithDetails, len(monthCategories))
	for i, v := range monthCategories {
		// find activity
		activity := activity[v.CategoryID].OrZero()

		// calculate available
		available, err := beans.Arithmetic.Add(
			previousAssigned[v.CategoryID].OrZero(),
			previousActivity[v.CategoryID].OrZero(),
			v.Amount,
			activity,
		)
		if err != nil {
			return nil, err
		}

		// build result
		res[i] = beans.MonthCategoryWithDetails{
			ID:         v.ID,
			CategoryID: v.CategoryID,
			Amount:     v.Amount,
			Activity:   activity,
			Available:  available,
		}
	}

	return res, nil
}
