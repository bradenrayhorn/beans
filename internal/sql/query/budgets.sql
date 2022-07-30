-- name: CreateBudget :exec
INSERT INTO budgets (
  id, name
) VALUES ($1, $2);

-- name: GetBudget :one
SELECT * FROM budgets WHERE id = $1;

-- name: CreateBudgetUser :exec
INSERT INTO budgets_users (budget_id, user_id) VALUES ($1, $2);

-- name: GetBudgetsForUser :many
SELECT budgets.* FROM budgets
JOIN budgets_users ON budgets_users.budget_id = budgets.id
                      AND budgets_users.user_id = $1;

