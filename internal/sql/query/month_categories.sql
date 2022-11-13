-- name: CreateMonthCategory :exec
INSERT INTO month_categories (
  id, month_id, category_id, amount
) VALUES ($1, $2, $3, $4);

-- name: GetMonthCategoriesForMonth :many
SELECT * FROM month_categories WHERE month_id = $1;

-- name: GetMonthCategoryByMonthAndCategory :one
SELECT * FROM month_categories WHERE month_id = $1 and category_id = $2;

-- name: UpdateMonthCategoryAmount :exec
UPDATE month_categories SET amount = $1 WHERE id = $2;
