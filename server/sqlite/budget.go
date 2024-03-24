package sqlite

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"zombiezen.com/go/sqlite"
)

type budgetRepository struct{ repository }

var _ beans.BudgetRepository = (*budgetRepository)(nil)

const budgetCreateSQL = `
INSERT INTO budgets (id, name) VALUES (:id, :name)
`

const budgetUserCreateSQL = `
INSERT INTO budget_users (budget_id, user_id) VALUES (:budgetID, :userID)
`

func (r *budgetRepository) Create(ctx context.Context, tx beans.Tx, id beans.ID, name beans.Name, userID beans.ID) error {
	err := db[any](r.pool).
		inTx(tx).
		execute(ctx, budgetCreateSQL, map[string]any{
			":id":   id.String(),
			":name": string(name),
		})
	if err != nil {
		return err
	}

	return db[any](r.pool).
		inTx(tx).
		execute(ctx, budgetUserCreateSQL, map[string]any{
			":budgetID": id.String(),
			":userID":   userID.String(),
		})
}

const budgetGetOneSQL = `
SELECT * FROM budgets WHERE id = :id
`

func (r *budgetRepository) Get(ctx context.Context, id beans.ID) (beans.Budget, error) {
	return db[beans.Budget](r.pool).
		mapWith(mapBudget).
		one(ctx, budgetGetOneSQL, map[string]any{
			":id": id.String(),
		})
}

const budgetGetForUserSQL = `
SELECT budgets.* FROM budgets
	JOIN budget_users ON budget_users.budget_id = budgets.id
	AND budget_users.user_id = :userID
`

func (r *budgetRepository) GetBudgetsForUser(ctx context.Context, userID beans.ID) ([]beans.Budget, error) {
	return db[beans.Budget](r.pool).
		mapWith(mapBudget).
		many(ctx, budgetGetForUserSQL, map[string]any{
			":userID": userID.String(),
		})
}

const budgetGetUserIDsSQL = `
SELECT user_id from budget_users WHERE budget_id = :budgetID
`

func (r *budgetRepository) GetBudgetUserIDs(ctx context.Context, id beans.ID) ([]beans.ID, error) {
	return db[beans.ID](r.pool).
		mapWith(func(stmt *sqlite.Stmt) (beans.ID, error) { return mapID(stmt, "user_id") }).
		many(ctx, budgetGetUserIDsSQL, map[string]any{
			":budgetID": id.String(),
		})
}

// mappers

func mapBudget(stmt *sqlite.Stmt) (beans.Budget, error) {
	id, err := mapID(stmt, "id")
	if err != nil {
		return beans.Budget{}, err
	}

	return beans.Budget{
		ID:   id,
		Name: beans.Name(stmt.GetText("name")),
	}, nil
}
