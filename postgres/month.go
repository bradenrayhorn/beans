package postgres

import (
	"context"

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
