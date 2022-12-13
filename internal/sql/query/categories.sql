-- name: CreateCategory :exec
INSERT INTO categories (
  id, budget_id, group_id, name
) VALUES ($1, $2, $3, $4);

-- name: GetCategoryForBudget :one
SELECT * FROM categories WHERE id = $1 AND budget_id = $2;

-- name: GetCategoriesForBudget :many
SELECT * FROM categories WHERE budget_id = $1;

-- name: CreateCategoryGroup :exec
INSERT INTO category_groups (
  id, budget_id, name
) VALUES ($1, $2, $3);

-- name: GetCategoryGroupsForBudget :many
SELECT * FROM category_groups WHERE budget_id = $1;

-- name: CategoryGroupExists :one
SELECT EXISTS (SELECT id FROM category_groups WHERE id = $1 AND budget_id = $2);
