// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: accounts.sql

package db

import (
	"context"
)

const createAccount = `-- name: CreateAccount :exec
INSERT INTO accounts (
  id, name, budget_id
) VALUES ($1, $2, $3)
`

type CreateAccountParams struct {
	ID       string
	Name     string
	BudgetID string
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) error {
	_, err := q.db.Exec(ctx, createAccount, arg.ID, arg.Name, arg.BudgetID)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, name, budget_id, created_at from accounts WHERE id = $1
`

func (q *Queries) GetAccount(ctx context.Context, id string) (Account, error) {
	row := q.db.QueryRow(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.BudgetID,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountsForBudget = `-- name: GetAccountsForBudget :many
SELECT id, name, budget_id, created_at from accounts WHERE budget_id = $1
`

func (q *Queries) GetAccountsForBudget(ctx context.Context, budgetID string) ([]Account, error) {
	rows, err := q.db.Query(ctx, getAccountsForBudget, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.BudgetID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
