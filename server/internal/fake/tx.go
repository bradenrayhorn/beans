package fake

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
)

type fTx struct {
	do []func()
}

func (t *fTx) Commit(ctx context.Context) error {
	if t.do != nil {
		for _, d := range t.do {
			d()
		}
	}
	return nil
}

func (t *fTx) Rollback(ctx context.Context) error {
	return nil
}

type txManager struct {
}

func (m *txManager) Create(ctx context.Context) (beans.Tx, error) {
	return &fTx{do: make([]func(), 0)}, nil
}
