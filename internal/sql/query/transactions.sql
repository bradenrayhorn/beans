-- name: CreateTransaction :exec
INSERT INTO transactions (
  id, account_id, date, amount, notes
) VALUES ($1, $2, $3, $4, $5);

-- name: GetTransactionsForBudget :many
SELECT transactions.* from transactions
JOIN accounts
  ON accounts.id = transactions.account_id
  AND accounts.budget_id = $1;

