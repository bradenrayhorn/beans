package sqlite

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"zombiezen.com/go/sqlite"
)

type payeeRepository struct{ repository }

var _ beans.PayeeRepository = (*payeeRepository)(nil)

const payeeCreateSQL = `
INSERT INTO payees (id, budget_id, name) VALUES (:id, :budgetID, :name)
`

func (r *payeeRepository) Create(ctx context.Context, payee beans.Payee) error {
	return db[any](r.pool).
		execute(ctx, payeeCreateSQL, map[string]any{
			":id":       payee.ID.String(),
			":budgetID": payee.BudgetID.String(),
			":name":     string(payee.Name),
		})
}

const payeeGetSQL = `
SELECT * FROM payees WHERE budget_id = :budgetID and id = :id
`

func (r *payeeRepository) Get(ctx context.Context, budgetID beans.ID, id beans.ID) (beans.Payee, error) {
	return db[beans.Payee](r.pool).
		mapWith(mapPayee).
		one(ctx, payeeGetSQL, map[string]any{
			":id":       id.String(),
			":budgetID": budgetID.String(),
		})
}

const payeeGetForBudgetSQL = `
SELECT * FROM payees WHERE budget_id = :budgetID
`

func (r *payeeRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]beans.Payee, error) {
	return db[beans.Payee](r.pool).
		mapWith(mapPayee).
		many(ctx, payeeGetForBudgetSQL, map[string]any{
			":budgetID": budgetID.String(),
		})
}

func mapPayee(stmt *sqlite.Stmt) (beans.Payee, error) {
	id, err := mapID(stmt, "id")
	if err != nil {
		return beans.Payee{}, err
	}
	budgetID, err := mapID(stmt, "budget_id")
	if err != nil {
		return beans.Payee{}, err
	}

	return beans.Payee{
		ID:       id,
		BudgetID: budgetID,
		Name:     beans.Name(stmt.GetText("name")),
	}, nil
}
