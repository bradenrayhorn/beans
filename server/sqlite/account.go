package sqlite

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"zombiezen.com/go/sqlite"
)

type accountRepository struct{ repository }

var _ beans.AccountRepository = (*accountRepository)(nil)

const accountCreateSQL = `
INSERT INTO accounts
	(id, budget_id, name, off_budget)
	VALUES (:id, :budgetID, :name, :offBudget)
`

func (r *accountRepository) Create(ctx context.Context, account beans.Account) error {
	return db[any](r.pool).execute(ctx, accountCreateSQL, map[string]any{
		":id":        account.ID.String(),
		":budgetID":  account.BudgetID.String(),
		":name":      string(account.Name),
		":offBudget": account.OffBudget,
	})
}

const accountGetOneSQL = `
SELECT * FROM accounts
	WHERE budget_id = :budgetID AND id = :id
`

func (r *accountRepository) Get(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.Account, error) {
	return db[beans.Account](r.pool).
		mapWith(mapAccount).
		one(ctx, accountGetOneSQL, map[string]any{
			":budgetID": budgetID.String(),
			":id":       id.String(),
		})
}

const accountGetWithBalance = `
SELECT accounts.*, sum(transactions.amount) as balance
	FROM accounts
	LEFT JOIN transactions ON
		accounts.id = transactions.account_id
		AND transactions.is_split = false
	WHERE budget_id = :budgetID
	GROUP BY (accounts.id)
`

func (r *accountRepository) GetWithBalance(ctx context.Context, budgetID beans.ID) ([]beans.AccountWithBalance, error) {
	return db[beans.AccountWithBalance](r.pool).
		mapWith(mapAccountWithBalance).
		many(ctx, accountGetWithBalance, map[string]any{
			":budgetID": budgetID.String(),
		})
}

const accountGetTransactableSQL = `
SELECT * FROM accounts
	WHERE budget_id = :budgetID
`

func (r *accountRepository) GetTransactable(ctx context.Context, budgetID beans.ID) ([]beans.Account, error) {
	return db[beans.Account](r.pool).
		mapWith(mapAccount).
		many(ctx, accountGetTransactableSQL, map[string]any{
			":budgetID": budgetID.String(),
		})
}

// mappers

func mapAccount(stmt *sqlite.Stmt) (beans.Account, error) {
	id, err := mapID(stmt, "id")
	if err != nil {
		return beans.Account{}, err
	}
	budgetID, err := mapID(stmt, "budget_id")
	if err != nil {
		return beans.Account{}, err
	}

	return beans.Account{
		ID:        id,
		Name:      beans.Name(stmt.GetText("name")),
		OffBudget: stmt.GetBool("off_budget"),
		BudgetID:  budgetID,
	}, nil
}

func mapAccountWithBalance(stmt *sqlite.Stmt) (beans.AccountWithBalance, error) {
	account, err := mapAccount(stmt)
	if err != nil {
		return beans.AccountWithBalance{}, err
	}

	return beans.AccountWithBalance{
		Account: account,
		Balance: mapAmount(stmt, "balance"),
	}, nil
}
