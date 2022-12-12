-- name: CreateTransaction :exec
INSERT INTO transactions (
  id, account_id, category_id, date, amount, notes
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetTransactionsForBudget :many
SELECT transactions.*, accounts.name as account_name from transactions
JOIN accounts
  ON accounts.id = transactions.account_id
  AND accounts.budget_id = $1;

