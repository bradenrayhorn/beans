// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: accounts.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAccount = `-- name: CreateAccount :exec
INSERT INTO accounts (
  id, name, budget_id, off_budget
) VALUES ($1, $2, $3, $4)
`

type CreateAccountParams struct {
	ID        string
	Name      string
	BudgetID  string
	OffBudget bool
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) error {
	_, err := q.db.Exec(ctx, createAccount,
		arg.ID,
		arg.Name,
		arg.BudgetID,
		arg.OffBudget,
	)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, name, budget_id, off_budget, created_at from accounts WHERE id = $1 AND budget_id = $2
`

type GetAccountParams struct {
	ID       string
	BudgetID string
}

func (q *Queries) GetAccount(ctx context.Context, arg GetAccountParams) (Account, error) {
	row := q.db.QueryRow(ctx, getAccount, arg.ID, arg.BudgetID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.BudgetID,
		&i.OffBudget,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountsWithBalance = `-- name: GetAccountsWithBalance :many
SELECT accounts.id, accounts.name, accounts.budget_id, accounts.off_budget, accounts.created_at, sum(transactions.amount)::numeric as balance
  FROM accounts
  LEFT JOIN transactions ON
    accounts.id = transactions.account_id
  WHERE budget_id = $1
  GROUP BY (
    accounts.id,
    accounts.name,
    accounts.budget_id
  )
`

type GetAccountsWithBalanceRow struct {
	Account Account
	Balance pgtype.Numeric
}

func (q *Queries) GetAccountsWithBalance(ctx context.Context, budgetID string) ([]GetAccountsWithBalanceRow, error) {
	rows, err := q.db.Query(ctx, getAccountsWithBalance, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAccountsWithBalanceRow
	for rows.Next() {
		var i GetAccountsWithBalanceRow
		if err := rows.Scan(
			&i.Account.ID,
			&i.Account.Name,
			&i.Account.BudgetID,
			&i.Account.OffBudget,
			&i.Account.CreatedAt,
			&i.Balance,
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

const getTransactableAccounts = `-- name: GetTransactableAccounts :many
SELECT accounts.id, accounts.name, accounts.budget_id, accounts.off_budget, accounts.created_at
  FROM accounts
  WHERE
    budget_id = $1
`

func (q *Queries) GetTransactableAccounts(ctx context.Context, budgetID string) ([]Account, error) {
	rows, err := q.db.Query(ctx, getTransactableAccounts, budgetID)
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
			&i.OffBudget,
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
