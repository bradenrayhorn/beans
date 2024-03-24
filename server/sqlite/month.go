package sqlite

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/bradenrayhorn/beans/server/beans"
	"zombiezen.com/go/sqlite"
)

type monthRepository struct{ repository }

var _ beans.MonthRepository = (*monthRepository)(nil)

func (r *monthRepository) Create(ctx context.Context, tx beans.Tx, month beans.Month) error {
	carryover, err := serializeAmount(month.Carryover.OrZero())
	if err != nil {
		return err
	}

	q := squirrel.Insert("months").
		Columns("id", "budget_id", "date", "carryover").
		Values(
			serializeID(month.ID),
			serializeID(month.BudgetID),
			serializeDate(month.Date.FirstDay()),
			carryover,
		)
	sql, params, err := q.ToSql()
	if err != nil {
		return err
	}

	return db[any](r.pool).
		inTx(tx).
		executeWithArgs(ctx, sql, params)
}

const monthGetSQL = `
SELECT * FROM months WHERE id = :id AND budget_id = :budgetID
`

func (r *monthRepository) Get(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.Month, error) {
	return db[beans.Month](r.pool).
		mapWith(mapMonth).
		one(ctx, monthGetSQL, map[string]any{
			":id":       id.String(),
			":budgetID": budgetID.String(),
		})
}

const monthUpdateSQL = `
UPDATE months
	SET carryover = :carryover
	WHERE id = :id
`

func (r *monthRepository) Update(ctx context.Context, month beans.Month) error {
	carryover, err := serializeAmount(month.Carryover.OrZero())
	if err != nil {
		return err
	}

	return db[any](r.pool).
		execute(ctx, monthUpdateSQL, map[string]any{
			":id":        month.ID.String(),
			":carryover": carryover,
		})
}

const monthGetByDateSQL = `
SELECT * FROM months WHERE date = :date AND budget_id = :budgetID
`

func (r *monthRepository) GetOrCreate(ctx context.Context, tx beans.Tx, budgetID beans.ID, date beans.MonthDate) (beans.Month, error) {
	month, err := db[beans.Month](r.pool).
		mapWith(mapMonth).
		one(ctx, monthGetByDateSQL, map[string]any{
			":date":     serializeDate(date.FirstDay()),
			":budgetID": budgetID.String(),
		})
	if err != nil {
		if errors.Is(err, beans.ErrorNotFound) {
			month := beans.Month{
				ID:        beans.NewID(),
				BudgetID:  budgetID,
				Date:      date,
				Carryover: beans.NewAmount(0, 0),
			}
			return month, r.Create(ctx, tx, month)
		}

		return beans.Month{}, err
	}

	return month, err
}

const monthGetForBudget = `
SELECT * FROM months WHERE budget_id = :budgetID
`

func (r *monthRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]beans.Month, error) {
	return db[beans.Month](r.pool).
		mapWith(mapMonth).
		many(ctx, monthGetForBudget, map[string]any{
			":budgetID": budgetID.String(),
		})
}

// mappers

func mapMonth(stmt *sqlite.Stmt) (beans.Month, error) {
	id, err := mapID(stmt, "id")
	if err != nil {
		return beans.Month{}, err
	}
	budgetID, err := mapID(stmt, "budget_id")
	if err != nil {
		return beans.Month{}, err
	}
	date, err := mapDate(stmt, "date")
	if err != nil {
		return beans.Month{}, err
	}

	return beans.Month{
		ID:        id,
		BudgetID:  budgetID,
		Date:      beans.NewMonthDate(date),
		Carryover: mapAmount(stmt, "carryover"),
	}, nil
}
