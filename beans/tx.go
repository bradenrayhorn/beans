package beans

import (
	"context"
	"errors"
)

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TxManager interface {
	Create(ctx context.Context) (Tx, error)
}

func ExecTx[T any](ctx context.Context, m TxManager, callback func(tx Tx) (T, error)) (T, error) {
	var empty T

	tx, err := m.Create(ctx)
	if err != nil {
		return empty, err
	}

	// This deferred rollback will only actually rollback the database if a
	// panic has occurred. Else the rollback will be handled by normal logic
	// below. For this reason we ignore the error.
	defer func() { _ = tx.Rollback(ctx) }()

	res, err := callback(tx)
	if err != nil {
		return empty, errors.Join(err, tx.Rollback(ctx))
	}

	if err = tx.Commit(ctx); err != nil {
		return empty, errors.Join(err, tx.Rollback(ctx))
	}

	return res, nil
}

func ExecTxNil(ctx context.Context, m TxManager, callback func(tx Tx) error) error {
	_, err := ExecTx(ctx, m, func(tx Tx) (interface{}, error) {
		return nil, callback(tx)
	})
	return err
}
