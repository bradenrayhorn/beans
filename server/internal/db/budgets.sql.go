// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: budgets.sql

package db

import (
	"context"
)

const createBudget = `-- name: CreateBudget :exec
INSERT INTO budgets (
  id, name
) VALUES ($1, $2)
`

type CreateBudgetParams struct {
	ID   string
	Name string
}

func (q *Queries) CreateBudget(ctx context.Context, arg CreateBudgetParams) error {
	_, err := q.db.Exec(ctx, createBudget, arg.ID, arg.Name)
	return err
}

const createBudgetUser = `-- name: CreateBudgetUser :exec
INSERT INTO budgets_users (budget_id, user_id) VALUES ($1, $2)
`

type CreateBudgetUserParams struct {
	BudgetID string
	UserID   string
}

func (q *Queries) CreateBudgetUser(ctx context.Context, arg CreateBudgetUserParams) error {
	_, err := q.db.Exec(ctx, createBudgetUser, arg.BudgetID, arg.UserID)
	return err
}

const getBudget = `-- name: GetBudget :one
SELECT id, name, created_at FROM budgets WHERE id = $1
`

func (q *Queries) GetBudget(ctx context.Context, id string) (Budget, error) {
	row := q.db.QueryRow(ctx, getBudget, id)
	var i Budget
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const getBudgetUserIDs = `-- name: GetBudgetUserIDs :many
SELECT user_id from budgets_users WHERE budget_id = $1
`

func (q *Queries) GetBudgetUserIDs(ctx context.Context, budgetID string) ([]string, error) {
	rows, err := q.db.Query(ctx, getBudgetUserIDs, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var user_id string
		if err := rows.Scan(&user_id); err != nil {
			return nil, err
		}
		items = append(items, user_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBudgetsForUser = `-- name: GetBudgetsForUser :many
SELECT budgets.id, budgets.name, budgets.created_at FROM budgets
JOIN budgets_users ON budgets_users.budget_id = budgets.id
                      AND budgets_users.user_id = $1
`

func (q *Queries) GetBudgetsForUser(ctx context.Context, userID string) ([]Budget, error) {
	rows, err := q.db.Query(ctx, getBudgetsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Budget
	for rows.Next() {
		var i Budget
		if err := rows.Scan(&i.ID, &i.Name, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}