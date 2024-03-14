CREATE TABLE transactions (
    id CHAR(27) PRIMARY KEY,
    account_id CHAR(27) NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    payee_id CHAR(27) REFERENCES payees(id) ON DELETE CASCADE,
    category_id CHAR(27) REFERENCES categories(id) ON DELETE CASCADE,
    transfer_id CHAR(27) REFERENCES transactions(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    amount NUMERIC NOT NULL,
    notes VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

