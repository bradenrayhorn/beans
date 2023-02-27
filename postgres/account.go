package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AccountRepository struct {
	db *db.Queries
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db.New(pool)}
}

func (r *AccountRepository) Create(ctx context.Context, id beans.ID, name beans.Name, budgetID beans.ID) error {
	return r.db.CreateAccount(ctx, db.CreateAccountParams{ID: id.String(), Name: string(name), BudgetID: budgetID.String()})
}

func (r *AccountRepository) Get(ctx context.Context, id beans.ID) (*beans.Account, error) {
	account, err := r.db.GetAccount(ctx, id.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}
	budgetID, err := beans.BeansIDFromString(account.BudgetID)
	if err != nil {
		return nil, err
	}
	return &beans.Account{
		ID:       id,
		Name:     beans.Name(account.Name),
		BudgetID: budgetID,
	}, nil
}

func (r *AccountRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]*beans.Account, error) {
	accounts := []*beans.Account{}
	dbAccounts, err := r.db.GetAccountsForBudget(ctx, budgetID.String())
	if err != nil {
		return accounts, err
	}

	for _, a := range dbAccounts {
		id, err := beans.BeansIDFromString(a.ID)
		if err != nil {
			return accounts, err
		}

		budgetID, err := beans.BeansIDFromString(a.BudgetID)
		if err != nil {
			return accounts, err
		}

		balance, err := numericToAmount(a.Balance)
		if err != nil {
			return accounts, err
		}
		if balance.Empty() {
			balance = beans.NewAmount(0, 0)
		}

		accounts = append(accounts, &beans.Account{ID: id, Name: beans.Name(a.Name), BudgetID: budgetID, Balance: balance})
	}

	return accounts, nil
}
