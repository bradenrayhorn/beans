package sqlite

var migrations = []string{
	`CREATE TABLE users (
		id CHAR(27) PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,
	`CREATE TABLE budgets (
		id CHAR(27) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,
	`CREATE TABLE budget_users (
		budget_id CHAR(27) NOT NULL,
		user_id CHAR(27) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (budget_id, user_id),
		FOREIGN KEY (budget_id) REFERENCES budgets (id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);`,
	`CREATE TABLE accounts (
		id CHAR(27) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		budget_id CHAR(27) NOT NULL,
		off_budget BOOLEAN NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (budget_id) REFERENCES budgets (id) ON DELETE CASCADE
	);`,
	`CREATE TABLE category_groups (
		id CHAR(27) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		is_income boolean NOT NULL,
		budget_id CHAR(27) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY (budget_id) REFERENCES budgets (id) ON DELETE CASCADE
	);`,
	`CREATE TABLE categories (
		id CHAR(27) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		budget_id CHAR(27) NOT NULL,
		group_id CHAR(27) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY (budget_id) REFERENCES budgets (id) ON DELETE CASCADE,
		FOREIGN KEY (group_id) REFERENCES category_groups (id) ON DELETE CASCADE
	);`,
	`CREATE TABLE payees (
		id CHAR(27) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		budget_id CHAR(27) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY (budget_id) REFERENCES budgets (id) ON DELETE CASCADE
	);`,
	`CREATE TABLE transactions (
		id CHAR(27) PRIMARY KEY,
		account_id CHAR(27) NOT NULL,
		payee_id CHAR(27),
		category_id CHAR(27),
		transfer_id CHAR(27),
		split_id CHAR(27),
		is_split BOOLEAN NOT NULL,
		date DATE NOT NULL,
		amount INTEGER NOT NULL,
		notes VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE,
		FOREIGN KEY (payee_id) REFERENCES payees (id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE,
		FOREIGN KEY (transfer_id) REFERENCES transactions (id) ON DELETE CASCADE,
		FOREIGN KEY (split_id) REFERENCES transactions (id) ON DELETE CASCADE
	);`,
	`CREATE TABLE months (
		id CHAR(27) PRIMARY KEY,
		budget_id CHAR(27) NOT NULL,
		date DATE NOT NULL,
		carryover INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		UNIQUE (budget_id, date),
		FOREIGN KEY (budget_id) REFERENCES budgets (id) ON DELETE CASCADE
	);`,
	`CREATE TABLE month_categories (
		id CHAR(27) PRIMARY KEY,
		month_id CHAR(27) NOT NULL,
		category_id CHAR(27) NOT NULL,
		amount INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY (month_id) REFERENCES months (id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE,
		UNIQUE (month_id, category_id)
	);`,
}
