-- name: CreateMonthCategory :exec
INSERT INTO month_categories (
  id, month_id, category_id, amount
) VALUES ($1, $2, $3, $4);

-- name: GetMonthCategoriesForMonth :many
SELECT month_categories.*, sum(t.amount)::numeric as spent
  FROM month_categories
  LEFT JOIN transactions t on t.category_id = month_categories.category_id
    AND t.date >= @from_date AND t.date <= @to_date
  WHERE month_id = @month_id
  GROUP BY (
    month_categories.id,
    month_categories.month_id,
    month_categories.category_id,
    month_categories.amount
  );

-- name: GetMonthCategoryByMonthAndCategory :one
SELECT * FROM month_categories WHERE month_id = $1 and category_id = $2;

-- name: UpdateMonthCategoryAmount :exec
UPDATE month_categories SET amount = $1 WHERE id = $2;

-- name: GetAmountInBudget :one
SELECT sum(month_categories.amount)::numeric as amount
  FROM month_categories
  JOIN months m on m.id = month_categories.month_id
    AND m.budget_id = $1;

