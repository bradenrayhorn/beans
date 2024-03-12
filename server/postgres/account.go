package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
	"github.com/bradenrayhorn/beans/server/postgres/mapper"
)

type AccountRepository struct{ repository }

var _ beans.AccountRepository = (*AccountRepository)(nil)

func (r *AccountRepository) Create(ctx context.Context, account beans.Account) error {
	return r.DB(nil).CreateAccount(ctx, db.CreateAccountParams{
		ID:        account.ID.String(),
		Name:      string(account.Name),
		BudgetID:  account.BudgetID.String(),
		OffBudget: account.OffBudget,
	})
}

func (r *AccountRepository) Get(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.Account, error) {
	account, err := r.DB(nil).GetAccount(ctx, db.GetAccountParams{
		ID:       id.String(),
		BudgetID: budgetID.String(),
	})
	if err != nil {
		return beans.Account{}, mapPostgresError(err)
	}

	return mapper.Account(account)
}

func (r *AccountRepository) GetWithBalance(ctx context.Context, budgetID beans.ID) ([]beans.AccountWithBalance, error) {
	dbAccounts, err := r.DB(nil).GetAccountsWithBalance(ctx, budgetID.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.MapSlice(dbAccounts, mapper.AccountWithBalance)
}

func (r *AccountRepository) GetTransactable(ctx context.Context, budgetID beans.ID) ([]beans.Account, error) {
	dbAccounts, err := r.DB(nil).GetTransactableAccounts(ctx, budgetID.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.MapSlice(dbAccounts, mapper.Account)
}
