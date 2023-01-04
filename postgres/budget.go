package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BudgetRepository struct {
	repository
}

func NewBudgetRepository(pool *pgxpool.Pool) *BudgetRepository {
	return &BudgetRepository{repository{pool}}
}

func (r *BudgetRepository) Create(ctx context.Context, tx beans.Tx, id beans.ID, name beans.Name, userID beans.UserID) error {
	q := r.DB(tx)

	if err := q.CreateBudget(ctx, db.CreateBudgetParams{ID: id.String(), Name: string(name)}); err != nil {
		return err
	}
	if err := q.CreateBudgetUser(ctx, db.CreateBudgetUserParams{BudgetID: id.String(), UserID: userID.String()}); err != nil {
		return err
	}

	return nil
}

func (r *BudgetRepository) Get(ctx context.Context, id beans.ID) (*beans.Budget, error) {
	budget, err := r.DB(nil).GetBudget(ctx, id.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	// load users
	userIDStrings, err := r.DB(nil).GetBudgetUserIDs(ctx, id.String())
	if err != nil {
		return nil, err
	}
	userIDs := make([]beans.UserID, 0, len(userIDStrings))
	for _, v := range userIDStrings {
		userID, err := beans.UserIDFromString(v)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return &beans.Budget{
		ID:      id,
		Name:    beans.Name(budget.Name),
		UserIDs: userIDs,
	}, nil
}

func (r *BudgetRepository) GetBudgetsForUser(ctx context.Context, userID beans.UserID) ([]*beans.Budget, error) {
	budgets := []*beans.Budget{}
	dbBudgets, err := r.DB(nil).GetBudgetsForUser(ctx, userID.String())
	if err != nil {
		return budgets, err
	}

	for _, b := range dbBudgets {
		id, err := beans.BeansIDFromString(b.ID)
		if err != nil {
			return budgets, err
		}

		budgets = append(budgets, &beans.Budget{ID: id, Name: beans.Name(b.Name)})
	}

	return budgets, nil
}
