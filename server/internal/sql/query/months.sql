-- name: CreateMonth :exec
INSERT INTO months (
  id, budget_id, date, carryover
) VALUES ($1, $2, $3, $4);

-- name: GetMonthByID :one
SELECT * FROM months WHERE id = $1;

-- name: GetMonthByDate :one
SELECT * FROM months WHERE budget_id = $1 AND date = $2;

-- name: GetMonthsByBudget :many
SELECT * FROM months WHERE budget_id = $1;

-- name: UpdateMonth :exec
UPDATE months
  SET carryover = $1
  WHERE id = $2;
