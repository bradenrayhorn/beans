package contract

import (
	"context"
	"time"

	"github.com/bradenrayhorn/beans/server/beans"
)

type budgetContract struct {
	contract
}

func (c *budgetContract) Create(ctx context.Context, auth *beans.AuthContext, name beans.Name) (*beans.Budget, error) {
	if err := beans.ValidateFields(beans.Field("Budget name", name)); err != nil {
		return nil, err
	}

	budgetID := beans.NewBeansID()

	return beans.ExecTx(ctx, c.ds().TxManager(), func(tx beans.Tx) (*beans.Budget, error) {
		// create budget
		if err := c.ds().BudgetRepository().Create(ctx, tx, budgetID, name, auth.UserID()); err != nil {
			return nil, err
		}

		// create month
		if err := c.ds().MonthRepository().Create(ctx, tx, &beans.Month{
			ID:       beans.NewBeansID(),
			BudgetID: budgetID,
			Date:     beans.NewMonthDate(beans.NewDate(time.Now())),
		}); err != nil {
			return nil, err
		}

		// create income group and category
		categoryGroup := &beans.CategoryGroup{
			ID:       beans.NewBeansID(),
			BudgetID: budgetID,
			Name:     "Income",
			IsIncome: true,
		}
		if err := c.ds().CategoryRepository().CreateGroup(ctx, tx, categoryGroup); err != nil {
			return nil, err
		}

		category := &beans.Category{
			ID:       beans.NewBeansID(),
			GroupID:  categoryGroup.ID,
			BudgetID: budgetID,
			Name:     "Income",
		}
		if err := c.ds().CategoryRepository().Create(ctx, tx, category); err != nil {
			return nil, err
		}

		return &beans.Budget{
			ID:   budgetID,
			Name: name,
		}, nil
	})
}

func (c *budgetContract) Get(ctx context.Context, auth *beans.AuthContext, id beans.ID) (*beans.Budget, error) {
	budget, err := c.ds().BudgetRepository().Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !budget.UserHasAccess(auth.UserID()) {
		return nil, beans.ErrorNotFound
	}

	return budget, nil
}

func (c *budgetContract) GetAll(ctx context.Context, auth *beans.AuthContext) ([]*beans.Budget, error) {
	return c.ds().BudgetRepository().GetBudgetsForUser(ctx, auth.UserID())
}
