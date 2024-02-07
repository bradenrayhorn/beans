package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
)

type accountContract struct {
	contract
}

func (c *accountContract) Create(ctx context.Context, auth *beans.BudgetAuthContext, name beans.Name) (beans.ID, error) {
	if err := beans.ValidateFields(beans.Field("Account name", name)); err != nil {
		return beans.ID{}, err
	}

	accountID := beans.NewBeansID()
	if err := c.ds().AccountRepository().Create(ctx, accountID, name, auth.BudgetID()); err != nil {
		return beans.ID{}, err
	}

	return accountID, nil
}

func (c *accountContract) GetAll(ctx context.Context, auth *beans.BudgetAuthContext) ([]beans.Account, error) {
	return c.ds().AccountRepository().GetForBudget(ctx, auth.BudgetID())
}

func (c *accountContract) Get(ctx context.Context, auth *beans.BudgetAuthContext, id beans.ID) (beans.Account, error) {
	account, err := c.ds().AccountRepository().Get(ctx, id)
	if err != nil {
		return beans.Account{}, err
	}

	if account.BudgetID != auth.BudgetID() {
		return beans.Account{}, beans.NewError(beans.ENOTFOUND, "account not found")
	}

	return account, nil
}
