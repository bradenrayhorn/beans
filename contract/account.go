package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
)

type accountContract struct {
	accountRepository beans.AccountRepository
}

func NewAccountContract(
	accountRepository beans.AccountRepository,
) *accountContract {
	return &accountContract{accountRepository}
}

func (c *accountContract) Create(ctx context.Context, budgetID beans.ID, name beans.Name) (*beans.Account, error) {
	if err := beans.ValidateFields(beans.Field("Account name", name)); err != nil {
		return nil, err
	}

	accountID := beans.NewBeansID()
	if err := c.accountRepository.Create(ctx, accountID, name, budgetID); err != nil {
		return nil, err
	}

	return &beans.Account{
		ID:       accountID,
		Name:     name,
		BudgetID: budgetID,
	}, nil
}

func (c *accountContract) GetAll(ctx context.Context, budgetID beans.ID) ([]*beans.Account, error) {
	return c.accountRepository.GetForBudget(ctx, budgetID)
}
