package postgres

import (
	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func (r *repository) DB(tx beans.Tx) *db.Queries {
	if tx == nil {
		return db.New(r.pool)
	}

	ptx := tx.(*Tx)
	return db.New(ptx.tx)
}
