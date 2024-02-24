-- name: CreateMonthCategory :exec
INSERT INTO month_categories (
  id, month_id, category_id, amount
) VALUES ($1, $2, $3, $4);

-- name: GetMonthCategoriesForMonth :many
SELECT month_categories.*
  FROM month_categories
  WHERE month_id = @month_id;

-- name: GetPastMonthCategoriesAssigned :many
SELECT
    mc.category_id,
    sum(mc.amount)::numeric as assigned
  FROM month_categories mc
  JOIN months m on m.id = mc.month_id
    AND m.budget_id = @budget_id
    AND m.date < @before_date
  GROUP BY (
    mc.category_id
  );

-- name: GetMonthCategoryByMonthAndCategory :one
SELECT * FROM month_categories WHERE month_id = $1 and category_id = $2;

-- name: UpdateMonthCategoryAmount :exec
UPDATE month_categories SET amount = $1 WHERE id = $2;

-- name: GetAssignedInMonth :one
SELECT sum(month_categories.amount)::numeric as amount
  FROM month_categories
  JOIN months m on m.id = month_categories.month_id
    AND m.id = $1;

