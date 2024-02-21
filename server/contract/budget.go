package contract

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
)

type budgetContract struct{ contract }

func (c *budgetContract) Create(ctx context.Context, auth *beans.AuthContext, name beans.Name) (beans.Budget, error) {
	if err := beans.ValidateFields(beans.Field("Budget name", name)); err != nil {
		return beans.Budget{}, err
	}

	budgetID := beans.NewID()

	return beans.ExecTx(ctx, c.ds().TxManager(), func(tx beans.Tx) (beans.Budget, error) {
		// create budget
		if err := c.ds().BudgetRepository().Create(ctx, tx, budgetID, name, auth.UserID()); err != nil {
			return beans.Budget{}, err
		}

		// create income group and category
		categoryGroup := beans.CategoryGroup{
			ID:       beans.NewID(),
			BudgetID: budgetID,
			Name:     "Income",
			IsIncome: true,
		}
		if err := c.ds().CategoryRepository().CreateGroup(ctx, tx, categoryGroup); err != nil {
			return beans.Budget{}, err
		}

		category := beans.Category{
			ID:       beans.NewID(),
			GroupID:  categoryGroup.ID,
			BudgetID: budgetID,
			Name:     "Income",
		}
		if err := c.ds().CategoryRepository().Create(ctx, tx, category); err != nil {
			return beans.Budget{}, err
		}

		return beans.Budget{
			ID:   budgetID,
			Name: name,
		}, nil
	})
}

func (c *budgetContract) Get(ctx context.Context, auth *beans.AuthContext, id beans.ID) (beans.Budget, error) {
	budget, err := c.ds().BudgetRepository().Get(ctx, id)
	if err != nil {
		return beans.Budget{}, err
	}

	userIDs, err := c.ds().BudgetRepository().GetBudgetUserIDs(ctx, id)
	if err != nil {
		return beans.Budget{}, err
	}

	// only return budget if the authed user is in the list
	for _, userID := range userIDs {
		if userID == auth.UserID() {
			return budget, nil
		}
	}

	return beans.Budget{}, beans.ErrorNotFound
}

func (c *budgetContract) GetAll(ctx context.Context, auth *beans.AuthContext) ([]beans.Budget, error) {
	return c.ds().BudgetRepository().GetBudgetsForUser(ctx, auth.UserID())
}
