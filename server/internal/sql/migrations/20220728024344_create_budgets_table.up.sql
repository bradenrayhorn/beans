CREATE TABLE budgets (
    id CHAR(27) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE budgets_users (
    budget_id CHAR(27) NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    user_id CHAR(27) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (budget_id, user_id)
);

CREATE INDEX user_id_budget_id_idx ON budgets_users(user_id, budget_id);
