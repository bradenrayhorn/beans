package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
	"github.com/bradenrayhorn/beans/server/postgres/mapper"
	"github.com/jackc/pgx/v4/pgxpool"
)

type payeeRepository struct {
	repository
}

func NewPayeeRepository(pool *pgxpool.Pool) *payeeRepository {
	return &payeeRepository{repository{pool}}
}

func (r *payeeRepository) Create(ctx context.Context, payee *beans.Payee) error {
	return r.DB(nil).CreatePayee(ctx, db.CreatePayeeParams{
		ID:       payee.ID.String(),
		BudgetID: payee.BudgetID.String(),
		Name:     string(payee.Name),
	})
}

func (r *payeeRepository) Get(ctx context.Context, id beans.ID) (*beans.Payee, error) {
	res, err := r.DB(nil).GetPayee(ctx, id.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.Payee(res)
}

func (r *payeeRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]*beans.Payee, error) {
	res, err := r.DB(nil).GetPayeesForBudget(ctx, budgetID.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.MapSlice(res, mapper.Payee)
}
