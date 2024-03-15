package postgres

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
	pgDB "github.com/bradenrayhorn/beans/server/postgres/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *DbPool
}

func (r *repository) DB(tx beans.Tx) *db.Queries {
	if tx == nil {
		return db.New((*pgxpool.Pool)(r.pool))
	}

	ptx := tx.(*Tx)
	return db.New(ptx.tx)
}

func (r *repository) db(tx beans.Tx) *pgDB.Executor {
	if tx == nil {
		return pgDB.NewExecutor((*pgxpool.Pool)(r.pool))
	}

	ptx := tx.(*Tx)
	return pgDB.NewExecutor(ptx.tx)
}
