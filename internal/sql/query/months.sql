-- name: CreateMonth :exec
INSERT INTO months (
  id, budget_id, date
) VALUES ($1, $2, $3);

-- name: GetMonthByDate :one
SELECT * FROM months WHERE budget_id = $1 AND date = $2;

