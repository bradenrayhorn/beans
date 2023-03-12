CREATE TABLE months (
    id CHAR(27) PRIMARY KEY,
    budget_id CHAR(27) NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    carryover NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (budget_id, date)
);

