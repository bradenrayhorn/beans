package postgres

import (
	"context"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool}
}

func (m *TxManager) Create(ctx context.Context) (beans.Tx, error) {
	ptx, err := m.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &Tx{tx: ptx}, nil
}
