package postgres

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
	"github.com/bradenrayhorn/beans/server/postgres/mapper"
)

type monthCategoryRepository struct{ repository }

var _ beans.MonthCategoryRepository = (*monthCategoryRepository)(nil)

func (r *monthCategoryRepository) Create(ctx context.Context, tx beans.Tx, monthCategory beans.MonthCategory) error {
	return r.DB(tx).CreateMonthCategory(ctx, db.CreateMonthCategoryParams{
		ID:         monthCategory.ID.String(),
		MonthID:    monthCategory.MonthID.String(),
		CategoryID: monthCategory.CategoryID.String(),
		Amount:     mapper.AmountToNumeric(monthCategory.Amount),
	})
}

func (r *monthCategoryRepository) UpdateAmount(ctx context.Context, monthCategory beans.MonthCategory) error {
	return r.DB(nil).UpdateMonthCategoryAmount(ctx, db.UpdateMonthCategoryAmountParams{
		ID:     monthCategory.ID.String(),
		Amount: mapper.AmountToNumeric(monthCategory.Amount),
	})
}

func (r *monthCategoryRepository) GetForMonth(ctx context.Context, month beans.Month) ([]beans.MonthCategoryWithDetails, error) {
	monthCategories := []beans.MonthCategoryWithDetails{}

	res, err := r.DB(nil).GetMonthCategoriesForMonth(ctx, db.GetMonthCategoriesForMonthParams{
		FromDate: mapper.MonthDateToPg(month.Date),
		ToDate:   mapper.DateToPg(month.Date.LastDay()),
		MonthID:  month.ID.String(),
	})
	if err != nil {
		return monthCategories, mapPostgresError(err)
	}

	previousAssigned, err := r.DB(nil).GetPastMonthCategoriesAvailable(ctx, db.GetPastMonthCategoriesAvailableParams{
		BudgetID:   month.BudgetID.String(),
		BeforeDate: mapper.MonthDateToPg(month.Date),
	})
	if err != nil {
		return monthCategories, mapPostgresError(err)
	}

	previousAssignedByCategory := make(map[string]beans.Amount)
	for _, v := range previousAssigned {
		amount, err := mapper.NumericToAmount(v.Assigned)
		if err != nil {
			return monthCategories, err
		}
		previousAssignedByCategory[v.ID] = amount
	}

	previousActivity, err := r.DB(nil).GetActivityBeforeDateByCategory(ctx, db.GetActivityBeforeDateByCategoryParams{
		BudgetID: month.BudgetID.String(),
		Date:     mapper.MonthDateToPg(month.Date),
	})
	if err != nil {
		return monthCategories, mapPostgresError(err)
	}

	previousActivityByCategory := make(map[string]beans.Amount)
	for _, v := range previousActivity {
		amount, err := mapper.NumericToAmount(v.Activity)
		if err != nil {
			return monthCategories, err
		}
		previousActivityByCategory[v.ID] = amount
	}

	for _, v := range res {
		monthCategory, err := mapper.GetMonthCategoriesForMonthRow(v)
		if err != nil {
			return monthCategories, err
		}

		pastAssigned := previousAssignedByCategory[v.CategoryID]
		pastActivity := previousActivityByCategory[v.CategoryID]

		available, err := beans.Arithmetic.Add(
			pastAssigned.OrZero(),
			pastActivity.OrZero(),
			monthCategory.Amount,
			monthCategory.Activity,
		)
		if err != nil {
			return monthCategories, err
		}

		monthCategory.Available = available

		monthCategories = append(monthCategories, monthCategory)
	}

	return monthCategories, nil
}

func (r *monthCategoryRepository) GetOrCreate(ctx context.Context, tx beans.Tx, month beans.Month, categoryID beans.ID) (beans.MonthCategory, error) {
	res, err := r.DB(tx).GetMonthCategoryByMonthAndCategory(ctx, db.GetMonthCategoryByMonthAndCategoryParams{
		MonthID:    month.ID.String(),
		CategoryID: categoryID.String(),
	})

	if err != nil {
		err = mapPostgresError(err)

		if errors.Is(err, beans.ErrorNotFound) {
			monthCategory := beans.MonthCategory{
				ID:         beans.NewBeansID(),
				MonthID:    month.ID,
				CategoryID: categoryID,
				Amount:     beans.NewAmount(0, 0),
			}

			return monthCategory, r.Create(ctx, tx, monthCategory)
		}
	}

	id, err := beans.BeansIDFromString(res.ID)
	if err != nil {
		return beans.MonthCategory{}, err
	}
	amount, err := mapper.NumericToAmount(res.Amount)
	if err != nil {
		return beans.MonthCategory{}, err
	}

	return beans.MonthCategory{
		ID:         id,
		MonthID:    month.ID,
		CategoryID: categoryID,
		Amount:     amount,
	}, nil
}

func (r *monthCategoryRepository) GetAssignedInMonth(ctx context.Context, month beans.Month) (beans.Amount, error) {
	res, err := r.DB(nil).GetAssignedInMonth(ctx, month.ID.String())
	if err != nil {
		return beans.NewEmptyAmount(), err
	}

	amount, err := mapper.NumericToAmount(res)
	if err != nil {
		return beans.NewEmptyAmount(), err
	}

	return amount.OrZero(), nil
}
