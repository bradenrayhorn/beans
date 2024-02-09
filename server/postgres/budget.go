package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
)

type BudgetRepository struct{ repository }

var _ beans.BudgetRepository = (*BudgetRepository)(nil)

func (r *BudgetRepository) Create(ctx context.Context, tx beans.Tx, id beans.ID, name beans.Name, userID beans.ID) error {
	q := r.DB(tx)

	if err := q.CreateBudget(ctx, db.CreateBudgetParams{ID: id.String(), Name: string(name)}); err != nil {
		return err
	}
	if err := q.CreateBudgetUser(ctx, db.CreateBudgetUserParams{BudgetID: id.String(), UserID: userID.String()}); err != nil {
		return err
	}

	return nil
}

func (r *BudgetRepository) Get(ctx context.Context, id beans.ID) (beans.Budget, error) {
	budget, err := r.DB(nil).GetBudget(ctx, id.String())
	if err != nil {
		return beans.Budget{}, mapPostgresError(err)
	}

	return beans.Budget{
		ID:   id,
		Name: beans.Name(budget.Name),
	}, nil
}

func (r *BudgetRepository) GetBudgetsForUser(ctx context.Context, userID beans.ID) ([]beans.Budget, error) {
	budgets := []beans.Budget{}
	dbBudgets, err := r.DB(nil).GetBudgetsForUser(ctx, userID.String())
	if err != nil {
		return budgets, err
	}

	for _, b := range dbBudgets {
		id, err := beans.BeansIDFromString(b.ID)
		if err != nil {
			return budgets, err
		}

		budgets = append(budgets, beans.Budget{ID: id, Name: beans.Name(b.Name)})
	}

	return budgets, nil
}

func (r *BudgetRepository) GetBudgetUserIDs(ctx context.Context, id beans.ID) ([]beans.ID, error) {
	userIDStrings, err := r.DB(nil).GetBudgetUserIDs(ctx, id.String())
	if err != nil {
		return nil, err
	}
	userIDs := make([]beans.ID, 0, len(userIDStrings))
	for _, v := range userIDStrings {
		userID, err := beans.BeansIDFromString(v)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}
