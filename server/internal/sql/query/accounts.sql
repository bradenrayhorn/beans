-- name: CreateAccount :exec
INSERT INTO accounts (
  id, name, budget_id
) VALUES ($1, $2, $3);

-- name: GetAccount :one
SELECT * from accounts WHERE id = $1 AND budget_id = $2;

-- name: GetAccountsForBudget :many
SELECT accounts.*, sum(transactions.amount)::numeric as balance
  FROM accounts
  LEFT JOIN transactions ON
    accounts.id = transactions.account_id
  WHERE budget_id = $1
  GROUP BY (
    accounts.id,
    accounts.name,
    accounts.budget_id
  );

