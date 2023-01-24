-- name: CreateTransaction :exec
INSERT INTO transactions (
  id, account_id, category_id, date, amount, notes
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateTransaction :exec
UPDATE transactions
  SET account_id=$1, category_id=$2, date=$3, amount=$4, notes=$5
  WHERE id=$6;

-- name: GetTransaction :one
SELECT transactions.*, accounts.name as account_name, accounts.budget_id as budget_id
  FROM transactions
  JOIN accounts
    ON accounts.id = transactions.account_id
  WHERE transactions.id = $1;

-- name: GetTransactionsForBudget :many
SELECT transactions.*, accounts.name as account_name, categories.name as category_name from transactions
JOIN accounts
  ON accounts.id = transactions.account_id
  AND accounts.budget_id = $1
LEFT JOIN categories
  ON categories.id = transactions.category_id
ORDER BY date desc;

-- name: GetIncomeBeforeOrOnDate :one
SELECT sum(transactions.amount)::numeric
FROM transactions
JOIN categories
  ON categories.id = transactions.category_id
  AND categories.is_income = true
WHERE transactions.date <= $1;

