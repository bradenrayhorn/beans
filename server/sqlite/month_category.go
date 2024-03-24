package sqlite

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/bradenrayhorn/beans/server/beans"
	"zombiezen.com/go/sqlite"
)

type monthCategoryRepository struct{ repository }

var _ beans.MonthCategoryRepository = (*monthCategoryRepository)(nil)

func (r *monthCategoryRepository) Create(ctx context.Context, tx beans.Tx, monthCategory beans.MonthCategory) error {
	amount, err := serializeAmount(monthCategory.Amount)
	if err != nil {
		return err
	}

	q := squirrel.Insert("month_categories").
		Columns("id", "month_id", "category_id", "amount").
		Values(
			serializeID(monthCategory.ID),
			serializeID(monthCategory.MonthID),
			serializeID(monthCategory.CategoryID),
			amount,
		)
	sql, params, err := q.ToSql()
	if err != nil {
		return err
	}

	return db[any](r.pool).
		inTx(tx).
		executeWithArgs(ctx, sql, params)
}

const monthCategoryUpdateAmountSQL = `
UPDATE month_categories SET amount = :amount WHERE id = :id
`

func (r *monthCategoryRepository) UpdateAmount(ctx context.Context, monthCategory beans.MonthCategory) error {
	amount, err := serializeAmount(monthCategory.Amount)
	if err != nil {
		return err
	}

	return db[any](r.pool).
		execute(ctx, monthCategoryUpdateAmountSQL, map[string]any{
			":id":     monthCategory.ID.String(),
			":amount": amount,
		})
}

const monthCategoryGetForMonthSQL = `
SELECT * FROM month_categories WHERE month_id = :monthID
`

func (r *monthCategoryRepository) GetForMonth(ctx context.Context, month beans.Month) ([]beans.MonthCategory, error) {
	return db[beans.MonthCategory](r.pool).
		mapWith(mapMonthCategory).
		many(ctx, monthCategoryGetForMonthSQL, map[string]any{
			":monthID": month.ID,
		})
}

const monthCategoryAssignedByCategorySQL = `
SELECT mc.category_id, sum(mc.amount) as assigned
FROM month_categories mc
	JOIN months m on m.id = mc.month_id
	AND m.budget_id = :budgetID
	AND m.date < :beforeDate
GROUP BY (mc.category_id)
`

type monthCategoryAssignedByCategoryRow struct {
	CategoryID beans.ID
	Assigned   beans.Amount
}

func (r *monthCategoryRepository) GetAssignedByCategory(ctx context.Context, budgetID beans.ID, before beans.Date) (map[beans.ID]beans.Amount, error) {
	rows, err := db[monthCategoryAssignedByCategoryRow](r.pool).
		mapWith(mapMonthCategoryAssignedByCategoryRow).
		many(ctx, monthCategoryAssignedByCategorySQL, map[string]any{
			":budgetID":   budgetID.String(),
			":beforeDate": serializeDate(before),
		})
	if err != nil {
		return nil, err
	}

	previousAssignedByCategory := make(map[beans.ID]beans.Amount)
	for _, v := range rows {
		previousAssignedByCategory[v.CategoryID] = v.Assigned
	}
	return previousAssignedByCategory, nil
}

const monthCategoryGetByMonthAndCategorySQL = `
SELECT * FROM month_categories WHERE month_id = :monthID and category_id = :categoryID
`

func (r *monthCategoryRepository) GetOrCreate(ctx context.Context, tx beans.Tx, month beans.Month, categoryID beans.ID) (beans.MonthCategory, error) {
	monthCategory, err := db[beans.MonthCategory](r.pool).
		mapWith(mapMonthCategory).
		one(ctx, monthCategoryGetByMonthAndCategorySQL, map[string]any{
			":monthID":    month.ID.String(),
			":categoryID": categoryID.String(),
		})
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			monthCategory := beans.MonthCategory{
				ID:         beans.NewID(),
				MonthID:    month.ID,
				CategoryID: categoryID,
				Amount:     beans.NewAmount(0, 0),
			}
			return monthCategory, r.Create(ctx, tx, monthCategory)
		}

		return beans.MonthCategory{}, err
	}

	return monthCategory, err
}

const monthCategoryGetAssignedInMonthSQL = `
SELECT sum(month_categories.amount) as amount
FROM month_categories
JOIN months m on m.id = month_categories.month_id
	AND m.id = :monthID
`

func (r *monthCategoryRepository) GetAssignedInMonth(ctx context.Context, month beans.Month) (beans.Amount, error) {
	return db[beans.Amount](r.pool).
		mapWith(func(stmt *sqlite.Stmt) (beans.Amount, error) { return mapAmount(stmt, "amount"), nil }).
		one(ctx, monthCategoryGetAssignedInMonthSQL, map[string]any{
			":monthID": month.ID,
		})
}

// mappers

func mapMonthCategory(stmt *sqlite.Stmt) (beans.MonthCategory, error) {
	id, err := mapID(stmt, "id")
	if err != nil {
		return beans.MonthCategory{}, err
	}
	monthID, err := mapID(stmt, "month_id")
	if err != nil {
		return beans.MonthCategory{}, err
	}
	categoryID, err := mapID(stmt, "category_id")
	if err != nil {
		return beans.MonthCategory{}, err
	}

	return beans.MonthCategory{
		ID:         id,
		MonthID:    monthID,
		CategoryID: categoryID,
		Amount:     mapAmount(stmt, "amount"),
	}, nil
}

func mapMonthCategoryAssignedByCategoryRow(stmt *sqlite.Stmt) (monthCategoryAssignedByCategoryRow, error) {
	categoryID, err := mapID(stmt, "category_id")
	if err != nil {
		return monthCategoryAssignedByCategoryRow{}, err
	}

	return monthCategoryAssignedByCategoryRow{
		CategoryID: categoryID,
		Assigned:   mapAmount(stmt, "assigned"),
	}, nil
}
