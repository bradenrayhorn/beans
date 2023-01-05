package postgres

import (
	"context"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type monthRepository struct {
	repository
}

func NewMonthRepository(pool *pgxpool.Pool) *monthRepository {
	return &monthRepository{repository{pool}}
}

func (r *monthRepository) Create(ctx context.Context, tx beans.Tx, month *beans.Month) error {
	return r.DB(tx).CreateMonth(ctx, db.CreateMonthParams{
		ID:       month.ID.String(),
		BudgetID: month.BudgetID.String(),
		Date:     month.Date.Time,
	})
}

func (r *monthRepository) Get(ctx context.Context, id beans.ID) (*beans.Month, error) {
	res, err := r.DB(nil).GetMonthByID(ctx, id.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	budgetID, err := beans.BeansIDFromString(res.BudgetID)
	if err != nil {
		return nil, err
	}

	return &beans.Month{
		ID:       id,
		BudgetID: budgetID,
		Date:     beans.NewDate(res.Date),
	}, nil
}

func (r *monthRepository) GetByDate(ctx context.Context, budgetID beans.ID, date time.Time) (*beans.Month, error) {
	res, err := r.DB(nil).GetMonthByDate(ctx, db.GetMonthByDateParams{BudgetID: budgetID.String(), Date: date})
	if err != nil {
		return nil, mapPostgresError(err)
	}

	id, err := beans.BeansIDFromString(res.ID)
	if err != nil {
		return nil, err
	}

	return &beans.Month{
		ID:       id,
		BudgetID: budgetID,
		Date:     beans.NewDate(res.Date),
	}, nil
}

func (r *monthRepository) GetLatest(ctx context.Context, budgetID beans.ID) (*beans.Month, error) {
	res, err := r.DB(nil).GetNewestMonth(ctx, budgetID.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	id, err := beans.BeansIDFromString(res.ID)
	if err != nil {
		return nil, err
	}

	return &beans.Month{
		ID:       id,
		BudgetID: budgetID,
		Date:     beans.NewDate(res.Date),
	}, nil
}
