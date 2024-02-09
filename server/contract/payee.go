package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
)

type payeeContract struct{ contract }

var _ beans.PayeeContract = (*payeeContract)(nil)

func (c *payeeContract) CreatePayee(ctx context.Context, auth *beans.BudgetAuthContext, name beans.Name) (beans.Payee, error) {
	if err := beans.ValidateFields(
		beans.Field("Name", name),
	); err != nil {
		return beans.Payee{}, err
	}

	payee := beans.Payee{
		ID:       beans.NewBeansID(),
		BudgetID: auth.BudgetID(),
		Name:     name,
	}

	err := c.ds().PayeeRepository().Create(ctx, payee)
	if err != nil {
		return beans.Payee{}, err
	}

	return payee, nil
}

func (c *payeeContract) GetAll(ctx context.Context, auth *beans.BudgetAuthContext) ([]beans.Payee, error) {
	return c.ds().PayeeRepository().GetForBudget(ctx, auth.BudgetID())
}

func (c *payeeContract) Get(ctx context.Context, auth *beans.BudgetAuthContext, id beans.ID) (beans.Payee, error) {
	payee, err := c.ds().PayeeRepository().Get(ctx, id)
	if err != nil {
		return beans.Payee{}, err
	}

	if payee.BudgetID != auth.BudgetID() {
		return beans.Payee{}, beans.NewError(beans.ENOTFOUND, "payee not found")
	}

	return payee, nil
}
