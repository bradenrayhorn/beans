package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tx struct {
	tx pgx.Tx
}

func (t *Tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *Tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

type TxManager struct {
	pool *DbPool
}

func NewTxManager(pool *DbPool) *TxManager {
	return &TxManager{pool}
}

func (m *TxManager) Create(ctx context.Context) (beans.Tx, error) {
	ptx, err := (*pgxpool.Pool)(m.pool).Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &Tx{tx: ptx}, nil
}
