// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: categories.sql

package db

import (
	"context"
)

const categoryGroupExists = `-- name: CategoryGroupExists :one
SELECT EXISTS (SELECT id FROM category_groups WHERE id = $1 AND budget_id = $2)
`

type CategoryGroupExistsParams struct {
	ID       string
	BudgetID string
}

func (q *Queries) CategoryGroupExists(ctx context.Context, arg CategoryGroupExistsParams) (bool, error) {
	row := q.db.QueryRow(ctx, categoryGroupExists, arg.ID, arg.BudgetID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createCategory = `-- name: CreateCategory :exec
INSERT INTO categories (
  id, budget_id, group_id, name 
) VALUES ($1, $2, $3, $4)
`

type CreateCategoryParams struct {
	ID       string
	BudgetID string
	GroupID  string
	Name     string
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) error {
	_, err := q.db.Exec(ctx, createCategory,
		arg.ID,
		arg.BudgetID,
		arg.GroupID,
		arg.Name,
	)
	return err
}

const createCategoryGroup = `-- name: CreateCategoryGroup :exec
INSERT INTO category_groups (
  id, budget_id, name, is_income
) VALUES ($1, $2, $3, $4)
`

type CreateCategoryGroupParams struct {
	ID       string
	BudgetID string
	Name     string
	IsIncome bool
}

func (q *Queries) CreateCategoryGroup(ctx context.Context, arg CreateCategoryGroupParams) error {
	_, err := q.db.Exec(ctx, createCategoryGroup,
		arg.ID,
		arg.BudgetID,
		arg.Name,
		arg.IsIncome,
	)
	return err
}

const getCategoriesForBudget = `-- name: GetCategoriesForBudget :many
SELECT id, name, budget_id, group_id, created_at FROM categories WHERE budget_id = $1
`

func (q *Queries) GetCategoriesForBudget(ctx context.Context, budgetID string) ([]Category, error) {
	rows, err := q.db.Query(ctx, getCategoriesForBudget, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.BudgetID,
			&i.GroupID,
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

const getCategoryForBudget = `-- name: GetCategoryForBudget :one
SELECT id, name, budget_id, group_id, created_at FROM categories WHERE id = $1 AND budget_id = $2
`

type GetCategoryForBudgetParams struct {
	ID       string
	BudgetID string
}

func (q *Queries) GetCategoryForBudget(ctx context.Context, arg GetCategoryForBudgetParams) (Category, error) {
	row := q.db.QueryRow(ctx, getCategoryForBudget, arg.ID, arg.BudgetID)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.BudgetID,
		&i.GroupID,
		&i.CreatedAt,
	)
	return i, err
}

const getCategoryGroupsForBudget = `-- name: GetCategoryGroupsForBudget :many
SELECT id, name, is_income, budget_id, created_at FROM category_groups WHERE budget_id = $1
`

func (q *Queries) GetCategoryGroupsForBudget(ctx context.Context, budgetID string) ([]CategoryGroup, error) {
	rows, err := q.db.Query(ctx, getCategoryGroupsForBudget, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CategoryGroup
	for rows.Next() {
		var i CategoryGroup
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IsIncome,
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
