package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
)

type accountContract struct{ contract }

var _ beans.AccountContract = (*accountContract)(nil)

func (c *accountContract) Create(ctx context.Context, auth *beans.BudgetAuthContext, params beans.AccountCreate) (beans.ID, error) {
	if err := beans.ValidateFields(beans.Field("Account name", params.Name)); err != nil {
		return beans.ID{}, err
	}

	account := beans.Account{
		ID:        beans.NewID(),
		Name:      params.Name,
		BudgetID:  auth.BudgetID(),
		OffBudget: params.OffBudget,
	}

	if err := c.ds().AccountRepository().Create(ctx, account); err != nil {
		return beans.ID{}, err
	}

	return account.ID, nil
}

func (c *accountContract) GetAll(ctx context.Context, auth *beans.BudgetAuthContext) ([]beans.AccountWithBalance, error) {
	return c.ds().AccountRepository().GetWithBalance(ctx, auth.BudgetID())
}

func (c *accountContract) GetTransactable(ctx context.Context, auth *beans.BudgetAuthContext) ([]beans.Account, error) {
	return c.ds().AccountRepository().GetTransactable(ctx, auth.BudgetID())
}

func (c *accountContract) Get(ctx context.Context, auth *beans.BudgetAuthContext, id beans.ID) (beans.Account, error) {
	return c.ds().AccountRepository().Get(ctx, auth.BudgetID(), id)
}
