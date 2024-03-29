package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
)

type payeeContract struct{ contract }

var _ beans.PayeeContract = (*payeeContract)(nil)

func (c *payeeContract) CreatePayee(ctx context.Context, auth *beans.BudgetAuthContext, name beans.Name) (beans.ID, error) {
	if err := beans.ValidateFields(
		beans.Field("Name", name),
	); err != nil {
		return beans.EmptyID(), err
	}

	payee := beans.Payee{
		ID:       beans.NewID(),
		BudgetID: auth.BudgetID(),
		Name:     name,
	}

	err := c.ds().PayeeRepository().Create(ctx, payee)
	if err != nil {
		return beans.EmptyID(), err
	}

	return payee.ID, nil
}

func (c *payeeContract) GetAll(ctx context.Context, auth *beans.BudgetAuthContext) ([]beans.Payee, error) {
	return c.ds().PayeeRepository().GetForBudget(ctx, auth.BudgetID())
}

func (c *payeeContract) Get(ctx context.Context, auth *beans.BudgetAuthContext, id beans.ID) (beans.Payee, error) {
	return c.ds().PayeeRepository().Get(ctx, auth.BudgetID(), id)
}
