-- name: CreateCategory :exec
INSERT INTO categories (
  id, budget_id, group_id, name 
) VALUES ($1, $2, $3, $4);

-- name: GetCategoryForBudget :one
SELECT * FROM categories WHERE id = $1 AND budget_id = $2;

-- name: GetCategoryGroup :one
SELECT * FROM category_groups WHERE id = $1 AND budget_id = $2;

-- name: GetCategoriesForGroup :many
SELECT categories.* FROM categories
JOIN category_groups ON category_groups.id = categories.group_id
  AND category_groups.id = $1
  AND category_groups.budget_id = $2;

-- name: GetCategoriesForBudget :many
SELECT * FROM categories WHERE budget_id = $1;

-- name: CreateCategoryGroup :exec
INSERT INTO category_groups (
  id, budget_id, name, is_income
) VALUES ($1, $2, $3, $4);

-- name: GetCategoryGroupsForBudget :many
SELECT * FROM category_groups WHERE budget_id = $1;

