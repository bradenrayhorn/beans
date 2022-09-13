// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: months.sql

package db

import (
	"context"
	"time"
)

const createMonth = `-- name: CreateMonth :exec
INSERT INTO months (
  id, budget_id, date
) VALUES ($1, $2, $3)
`

type CreateMonthParams struct {
	ID       string
	BudgetID string
	Date     time.Time
}

func (q *Queries) CreateMonth(ctx context.Context, arg CreateMonthParams) error {
	_, err := q.db.Exec(ctx, createMonth, arg.ID, arg.BudgetID, arg.Date)
	return err
}

const getMonthByDate = `-- name: GetMonthByDate :one
SELECT id, budget_id, date, created_at FROM months WHERE budget_id = $1 AND date = $2
`

type GetMonthByDateParams struct {
	BudgetID string
	Date     time.Time
}

func (q *Queries) GetMonthByDate(ctx context.Context, arg GetMonthByDateParams) (Month, error) {
	row := q.db.QueryRow(ctx, getMonthByDate, arg.BudgetID, arg.Date)
	var i Month
	err := row.Scan(
		&i.ID,
		&i.BudgetID,
		&i.Date,
		&i.CreatedAt,
	)
	return i, err
}
