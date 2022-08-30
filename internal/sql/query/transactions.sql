-- name: CreateTransaction :exec
INSERT INTO transactions (
  id, account_id, date, amount, notes
) VALUES ($1, $2, $3, $4, $5);

-- name: GetTransactionsForAccount :many
SELECT * from transactions WHERE account_id = $1;

