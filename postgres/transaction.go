package postgres

import (
	"context"
	"database/sql"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionRepository struct {
	db *db.Queries
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db.New(pool)}
}

func (r *TransactionRepository) Create(
	ctx context.Context,
	id beans.ID,
	accountID beans.ID,
	amount beans.Amount,
	date beans.Date,
	notes beans.TransactionNotes,
) error {
	return r.db.CreateTransaction(ctx, db.CreateTransactionParams{
		ID:        id.String(),
		AccountID: accountID.String(),
		Date:      date.Time,
		Amount:    amountToNumeric(amount),
		Notes:     sql.NullString{String: string(notes), Valid: true},
	})
}
