-- name: CreateMonth :exec
INSERT INTO months (
  id, budget_id, date
) VALUES ($1, $2, $3);
