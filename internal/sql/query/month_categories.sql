-- name: CreateMonthCategory :exec
INSERT INTO month_categories (
  id, month_id, category_id, amount
) VALUES ($1, $2, $3, $4);

-- name: GetMonthCategoriesForMonth :many
SELECT * FROM month_categories WHERE month_id = $1;
