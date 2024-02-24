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

func (r *monthCategoryRepository) GetForMonth(ctx context.Context, month beans.Month) ([]beans.MonthCategory, error) {
	monthCategories := []beans.MonthCategory{}

	res, err := r.DB(nil).GetMonthCategoriesForMonth(ctx, month.ID.String())
	if err != nil {
		return monthCategories, mapPostgresError(err)
	}

	return mapper.MapSlice(res, mapper.MonthCategory)
}

func (r *monthCategoryRepository) GetAssignedByCategory(ctx context.Context, budgetID beans.ID, before beans.Date) (map[beans.ID]beans.Amount, error) {
	res, err := r.DB(nil).GetPastMonthCategoriesAssigned(ctx, db.GetPastMonthCategoriesAssignedParams{
		BudgetID:   budgetID.String(),
		BeforeDate: mapper.DateToPg(before),
	})
	if err != nil {
		return nil, mapPostgresError(err)
	}

	previousAssignedByCategory := make(map[beans.ID]beans.Amount)
	for _, v := range res {
		id, amount, err := mapper.GetPastMonthCategoriesAssigned(v)
		if err != nil {
			return nil, err
		}
		previousAssignedByCategory[id] = amount
	}

	return previousAssignedByCategory, nil
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
				ID:         beans.NewID(),
				MonthID:    month.ID,
				CategoryID: categoryID,
				Amount:     beans.NewAmount(0, 0),
			}

			return monthCategory, r.Create(ctx, tx, monthCategory)
		}
	}

	id, err := beans.IDFromString(res.ID)
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
