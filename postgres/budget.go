package postgres

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BudgetRepository struct {
	db   *db.Queries
	pool *pgxpool.Pool
}

func NewBudgetRepository(pool *pgxpool.Pool) *BudgetRepository {
	return &BudgetRepository{db: db.New(pool), pool: pool}
}

func (r *BudgetRepository) Create(ctx context.Context, id beans.ID, name beans.BudgetName, userID beans.UserID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	q := db.New(tx)

	if err := q.CreateBudget(ctx, db.CreateBudgetParams{ID: id.String(), Name: string(name)}); err != nil {
		return err
	}
	if err := q.CreateBudgetUser(ctx, db.CreateBudgetUserParams{BudgetID: id.String(), UserID: userID.String()}); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *BudgetRepository) Get(ctx context.Context, id beans.ID) (*beans.Budget, error) {
	budget, err := r.db.GetBudget(ctx, id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, beans.WrapError(err, beans.ErrorNotFound)
		}
		return nil, err
	}

	// load users
	userIDStrings, err := r.db.GetBudgetUserIDs(ctx, id.String())
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
		Name:    beans.BudgetName(budget.Name),
		UserIDs: userIDs,
	}, nil
}

func (r *BudgetRepository) GetBudgetsForUser(ctx context.Context, userID beans.UserID) ([]*beans.Budget, error) {
	budgets := []*beans.Budget{}
	dbBudgets, err := r.db.GetBudgetsForUser(ctx, userID.String())
	if err != nil {
		return budgets, err
	}

	for _, b := range dbBudgets {
		id, err := beans.BeansIDFromString(b.ID)
		if err != nil {
			return budgets, err
		}

		budgets = append(budgets, &beans.Budget{ID: id, Name: beans.BudgetName(b.Name)})
	}

	return budgets, nil
}
