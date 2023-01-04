package beans

import "context"

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TxManager interface {
	Create(ctx context.Context) (Tx, error)
}
