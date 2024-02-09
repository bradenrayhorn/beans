-- name: CreatePayee :exec
INSERT INTO payees (
  id, budget_id, name 
) VALUES ($1, $2, $3);

-- name: GetPayee :one
SELECT * FROM payees WHERE id = $1 AND budget_id = $2;

-- name: GetPayeesForBudget :many
SELECT * FROM payees WHERE budget_id = $1;

