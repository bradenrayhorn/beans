package postgres

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
	"github.com/bradenrayhorn/beans/postgres/mapper"
	"github.com/jackc/pgx/v4/pgxpool"
)

type monthCategoryRepository struct {
	repository
}

func NewMonthCategoryRepository(pool *pgxpool.Pool) *monthCategoryRepository {
	return &monthCategoryRepository{repository: repository{pool}}
}

func (r *monthCategoryRepository) Create(ctx context.Context, tx beans.Tx, monthCategory *beans.MonthCategory) error {
	return r.DB(tx).CreateMonthCategory(ctx, db.CreateMonthCategoryParams{
		ID:         monthCategory.ID.String(),
		MonthID:    monthCategory.MonthID.String(),
		CategoryID: monthCategory.CategoryID.String(),
		Amount:     amountToNumeric(monthCategory.Amount),
	})
}

func (r *monthCategoryRepository) UpdateAmount(ctx context.Context, monthCategoryID beans.ID, amount beans.Amount) error {
	return r.DB(nil).UpdateMonthCategoryAmount(ctx, db.UpdateMonthCategoryAmountParams{
		ID:     monthCategoryID.String(),
		Amount: amountToNumeric(amount),
	})
}

func (r *monthCategoryRepository) GetForMonth(ctx context.Context, month *beans.Month) ([]*beans.MonthCategory, error) {
	monthCategories := []*beans.MonthCategory{}

	res, err := r.DB(nil).GetMonthCategoriesForMonth(ctx, db.GetMonthCategoriesForMonthParams{
		FromDate: month.Date.Time(),
		ToDate:   month.Date.LastDay().Time,
		MonthID:  month.ID.String(),
	})
	if err != nil {
		return monthCategories, mapPostgresError(err)
	}

	previousAssigned, err := r.DB(nil).GetPastMonthCategoriesAvailable(ctx, db.GetPastMonthCategoriesAvailableParams{
		BudgetID:   month.BudgetID.String(),
		BeforeDate: month.Date.Time(),
	})
	if err != nil {
		return monthCategories, mapPostgresError(err)
	}

	previousAssignedByCategory := make(map[string]beans.Amount)
	for _, v := range previousAssigned {
		amount, err := numericToAmount(v.Assigned)
		if err != nil {
			return monthCategories, err
		}
		previousAssignedByCategory[v.ID] = amount
	}

	previousActivity, err := r.DB(nil).GetActivityBeforeDateByCategory(ctx, db.GetActivityBeforeDateByCategoryParams{
		BudgetID: month.BudgetID.String(),
		Date:     month.Date.Time(),
	})
	if err != nil {
		return monthCategories, mapPostgresError(err)
	}

	previousActivityByCategory := make(map[string]beans.Amount)
	for _, v := range previousActivity {
		amount, err := numericToAmount(v.Activity)
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
		if pastAssigned.Empty() {
			pastAssigned = beans.NewAmount(0, 0)
		}
		pastActivity := previousActivityByCategory[v.CategoryID]
		if pastActivity.Empty() {
			pastActivity = beans.NewAmount(0, 0)
		}

		available, err := beans.Add(pastAssigned, pastActivity, monthCategory.Amount, monthCategory.Activity)
		if err != nil {
			return monthCategories, err
		}

		monthCategory.Available = available

		monthCategories = append(monthCategories, monthCategory)
	}

	return monthCategories, nil
}

func (r *monthCategoryRepository) GetOrCreate(ctx context.Context, tx beans.Tx, monthID beans.ID, categoryID beans.ID) (*beans.MonthCategory, error) {
	res, err := r.DB(tx).GetMonthCategoryByMonthAndCategory(ctx, db.GetMonthCategoryByMonthAndCategoryParams{
		MonthID:    monthID.String(),
		CategoryID: categoryID.String(),
	})

	if err != nil {
		err = mapPostgresError(err)

		if errors.Is(err, beans.ErrorNotFound) {
			monthCategory := &beans.MonthCategory{
				ID:         beans.NewBeansID(),
				MonthID:    monthID,
				CategoryID: categoryID,
				Amount:     beans.NewAmount(0, 0),
			}

			return monthCategory, r.Create(ctx, tx, monthCategory)
		}
	}

	id, err := beans.BeansIDFromString(res.ID)
	if err != nil {
		return nil, err
	}
	amount, err := numericToAmount(res.Amount)
	if err != nil {
		return nil, err
	}

	return &beans.MonthCategory{
		ID:         id,
		MonthID:    monthID,
		CategoryID: categoryID,
		Amount:     amount,
	}, nil
}

func (r *monthCategoryRepository) GetAmountInBudget(ctx context.Context, budgetID beans.ID) (beans.Amount, error) {
	res, err := r.DB(nil).GetAmountInBudget(ctx, budgetID.String())
	if err != nil {
		return beans.NewEmptyAmount(), err
	}

	amount, err := numericToAmount(res)
	if err != nil {
		return beans.NewEmptyAmount(), err
	}

	return amount, nil
}
