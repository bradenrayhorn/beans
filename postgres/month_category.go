package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type monthCategoryRepository struct {
	db *db.Queries
}

func NewMonthCategoryRepository(pool *pgxpool.Pool) *monthCategoryRepository {
	return &monthCategoryRepository{db: db.New(pool)}
}

func (r *monthCategoryRepository) Create(ctx context.Context, monthCategory *beans.MonthCategory) error {
	return r.db.CreateMonthCategory(ctx, db.CreateMonthCategoryParams{
		ID:         monthCategory.ID.String(),
		MonthID:    monthCategory.MonthID.String(),
		CategoryID: monthCategory.CategoryID.String(),
		Amount:     amountToNumeric(monthCategory.Amount),
	})
}

func (r *monthCategoryRepository) UpdateAmount(ctx context.Context, monthCategoryID beans.ID, amount beans.Amount) error {
	return r.db.UpdateMonthCategoryAmount(ctx, db.UpdateMonthCategoryAmountParams{
		ID:     monthCategoryID.String(),
		Amount: amountToNumeric(amount),
	})
}

func (r *monthCategoryRepository) GetForMonth(ctx context.Context, monthID beans.ID) ([]*beans.MonthCategory, error) {
	monthCategories := []*beans.MonthCategory{}
	res, err := r.db.GetMonthCategoriesForMonth(ctx, monthID.String())
	if err != nil {
		return monthCategories, err
	}

	for _, v := range res {
		id, err := beans.BeansIDFromString(v.ID)
		if err != nil {
			return monthCategories, err
		}
		categoryID, err := beans.BeansIDFromString(v.CategoryID)
		if err != nil {
			return monthCategories, err
		}
		amount, err := numericToAmount(v.Amount)
		if err != nil {
			return monthCategories, err
		}
		spent, err := numericToAmount(v.Spent)
		if err != nil {
			return monthCategories, err
		}
		if spent.Empty() {
			spent = beans.NewAmount(0, 0)
		}

		monthCategories = append(monthCategories, &beans.MonthCategory{
			ID:         id,
			MonthID:    monthID,
			CategoryID: categoryID,
			Amount:     amount,
			Spent:      spent,
		})
	}

	return monthCategories, nil
}

func (r *monthCategoryRepository) GetByMonthAndCategory(ctx context.Context, monthID beans.ID, categoryID beans.ID) (*beans.MonthCategory, error) {
	res, err := r.db.GetMonthCategoryByMonthAndCategory(ctx, db.GetMonthCategoryByMonthAndCategoryParams{
		MonthID:    monthID.String(),
		CategoryID: categoryID.String(),
	})

	if err != nil {
		return nil, mapPostgresError(err)
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
