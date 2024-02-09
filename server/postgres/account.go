package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
	"github.com/bradenrayhorn/beans/server/postgres/mapper"
)

type AccountRepository struct{ repository }

var _ beans.AccountRepository = (*AccountRepository)(nil)

func (r *AccountRepository) Create(ctx context.Context, id beans.ID, name beans.Name, budgetID beans.ID) error {
	return r.DB(nil).CreateAccount(ctx, db.CreateAccountParams{ID: id.String(), Name: string(name), BudgetID: budgetID.String()})
}

func (r *AccountRepository) Get(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.Account, error) {
	account, err := r.DB(nil).GetAccount(ctx, db.GetAccountParams{
		ID:       id.String(),
		BudgetID: budgetID.String(),
	})
	if err != nil {
		return beans.Account{}, mapPostgresError(err)
	}

	return beans.Account{
		ID:       id,
		Name:     beans.Name(account.Name),
		BudgetID: budgetID,
	}, nil
}

func (r *AccountRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]beans.AccountWithBalance, error) {
	accounts := []beans.AccountWithBalance{}
	dbAccounts, err := r.DB(nil).GetAccountsForBudget(ctx, budgetID.String())
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

		balance, err := mapper.NumericToAmount(a.Balance)
		if err != nil {
			return accounts, err
		}
		if balance.Empty() {
			balance = beans.NewAmount(0, 0)
		}

		accounts = append(accounts, beans.AccountWithBalance{
			Account: beans.Account{ID: id, Name: beans.Name(a.Name), BudgetID: budgetID},
			Balance: balance,
		})
	}

	return accounts, nil
}
