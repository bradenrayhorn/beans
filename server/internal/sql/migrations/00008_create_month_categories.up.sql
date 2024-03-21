-- +goose Up
CREATE TABLE month_categories (
    id CHAR(27) PRIMARY KEY,
    month_id CHAR(27) NOT NULL REFERENCES months(id) ON DELETE CASCADE,
    category_id CHAR(27) NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    amount NUMERIC,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX month_category_idx on month_categories (month_id, category_id);

