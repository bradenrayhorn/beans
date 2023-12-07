package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
)

type payeeContract struct {
	payeeRepository beans.PayeeRepository
}

func NewPayeeContract(
	payeeRepository beans.PayeeRepository,
) *payeeContract {
	return &payeeContract{
		payeeRepository,
	}
}

func (c *payeeContract) CreatePayee(ctx context.Context, auth *beans.BudgetAuthContext, name beans.Name) (*beans.Payee, error) {
	if err := beans.ValidateFields(
		beans.Field("Name", name),
	); err != nil {
		return nil, err
	}

	payee := &beans.Payee{
		ID:       beans.NewBeansID(),
		BudgetID: auth.BudgetID(),
		Name:     name,
	}

	err := c.payeeRepository.Create(ctx, payee)
	if err != nil {
		return nil, err
	}

	return payee, nil
}

func (c *payeeContract) GetAll(ctx context.Context, auth *beans.BudgetAuthContext) ([]*beans.Payee, error) {
	return c.payeeRepository.GetForBudget(ctx, auth.BudgetID())
}
