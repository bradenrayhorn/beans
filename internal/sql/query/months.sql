-- name: CreateMonth :exec
INSERT INTO months (
  id, budget_id, date
) VALUES ($1, $2, $3);

-- name: GetMonthByID :one
SELECT * FROM months WHERE id = $1;

-- name: GetMonthByDate :one
SELECT * FROM months WHERE budget_id = $1 AND date = $2;

-- name: GetNewestMonth :one
SELECT * FROM months WHERE budget_id = $1 ORDER BY date desc LIMIT 1;

