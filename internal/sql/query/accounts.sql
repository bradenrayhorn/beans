-- name: CreateAccount :exec
INSERT INTO accounts (
  id, name, budget_id
) VALUES ($1, $2, $3);

-- name: GetAccount :one
SELECT * from accounts WHERE id = $1;

-- name: GetAccountsForBudget :many
SELECT * from accounts WHERE budget_id = $1;

