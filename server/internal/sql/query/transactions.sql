-- name: CreateTransaction :copyfrom
INSERT INTO transactions (
  id, account_id, payee_id, category_id, date, amount, notes, transfer_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateTransaction :exec
UPDATE transactions
  SET account_id=$1, category_id=$2, payee_id=$3, date=$4, amount=$5, notes=$6
  WHERE id=$7;

-- name: DeleteTransactions :exec
DELETE FROM transactions
  USING accounts
  WHERE
    accounts.id = transactions.account_id
    AND accounts.budget_id=$1
    AND transactions.id = ANY(sqlc.arg(IDs)::varchar[]);

-- name: GetTransaction :one
SELECT transactions.*
  FROM transactions
  JOIN accounts
    ON accounts.id = transactions.account_id
    AND accounts.budget_id = @budget_id
  WHERE transactions.id = @id;

-- name: GetIncomeBetween :one
SELECT sum(transactions.amount)::numeric
FROM transactions
JOIN categories
  ON categories.id = transactions.category_id
JOIN category_groups
  ON category_groups.id = categories.group_id
  AND category_groups.is_income = true
JOIN accounts
  ON accounts.id = transactions.account_id
  AND accounts.budget_id = @budget_id
WHERE
  transactions.date <= @end_date
  AND transactions.date >= @begin_date;

-- name: GetActivityByCategory :many
SELECT categories.id, sum(transactions.amount)::numeric as activity
  FROM transactions
  JOIN categories
    ON transactions.category_id = categories.id
  JOIN accounts
    ON accounts.id = transactions.account_id
    AND accounts.budget_id = @budget_id
  WHERE
    (transactions.date >= @from_date OR NOT @filter_from_date)
    AND (transactions.date <= @to_date OR NOT @filter_to_date)
  GROUP BY (
    categories.id
  );

