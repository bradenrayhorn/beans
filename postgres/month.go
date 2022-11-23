package postgres

import (
	"context"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type monthRepository struct {
	db *db.Queries
}

func NewMonthRepository(pool *pgxpool.Pool) *monthRepository {
	return &monthRepository{db: db.New(pool)}
}

func (r *monthRepository) Create(ctx context.Context, month *beans.Month) error {
	return r.db.CreateMonth(ctx, db.CreateMonthParams{
		ID:       month.ID.String(),
		BudgetID: month.BudgetID.String(),
		Date:     month.Date.Time,
	})
}

func (r *monthRepository) Get(ctx context.Context, id beans.ID) (*beans.Month, error) {
	res, err := r.db.GetMonthByID(ctx, id.String())
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
	res, err := r.db.GetMonthByDate(ctx, db.GetMonthByDateParams{BudgetID: budgetID.String(), Date: date})
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
	res, err := r.db.GetNewestMonth(ctx, budgetID.String())
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
