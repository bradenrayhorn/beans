package sqlite

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type repository struct {
	pool *Pool
}

type executor[T any] struct {
	pool *Pool

	conn   *sqlite.Conn // not nil if in a transaction
	mapper func(stmt *sqlite.Stmt) (T, error)
}

func db[T any](pool *Pool) *executor[T] {
	return &executor[T]{pool: pool}
}

func (d *executor[T]) inTx(tx beans.Tx) *executor[T] {
	if tx != nil {
		d.conn = tx.(*Tx).conn
	}
	return d
}

func (d *executor[T]) mapWith(mapper func(stmt *sqlite.Stmt) (T, error)) *executor[T] {
	d.mapper = mapper
	return d
}

func (d *executor[T]) execute(ctx context.Context, query string, args map[string]any) error {
	conn, done, err := d.connOrNew(ctx)
	if err != nil {
		return err
	}
	defer done()

	return sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Named: args,
	})
}

func (d *executor[T]) executeWithArgs(ctx context.Context, query string, args []any) error {
	conn, done, err := d.connOrNew(ctx)
	if err != nil {
		return err
	}
	defer done()

	return sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: args,
	})
}

func (d *executor[T]) one(ctx context.Context, query string, args map[string]any) (T, error) {
	var result T
	conn, done, err := d.connOrNew(ctx)
	if err != nil {
		return result, err
	}
	defer done()

	found := false
	err = sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Named: args,
		ResultFunc: func(stmt *sqlite.Stmt) error {
			found = true
			mapped, err := d.mapper(stmt)
			if err != nil {
				return err
			}

			result = mapped
			return nil
		},
	})
	if !found {
		return result, beans.NewError(beans.ENOTFOUND, "not found")
	}

	return result, err
}

func (d *executor[T]) oneWithArgs(ctx context.Context, query string, args []any) (T, error) {
	var result T
	conn, done, err := d.connOrNew(ctx)
	if err != nil {
		return result, err
	}
	defer done()

	found := false
	err = sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: args,
		ResultFunc: func(stmt *sqlite.Stmt) error {
			found = true
			mapped, err := d.mapper(stmt)
			if err != nil {
				return err
			}

			result = mapped
			return nil
		},
	})
	if !found {
		return result, beans.NewError(beans.ENOTFOUND, "not found")
	}

	return result, err
}

func (d *executor[T]) many(ctx context.Context, query string, args map[string]any) ([]T, error) {
	result := []T{}
	conn, done, err := d.connOrNew(ctx)
	if err != nil {
		return result, err
	}
	defer done()

	err = sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Named: args,
		ResultFunc: func(stmt *sqlite.Stmt) error {
			mapped, err := d.mapper(stmt)
			if err != nil {
				return err
			}

			result = append(result, mapped)
			return nil
		},
	})

	return result, err
}

func (d *executor[T]) manyWithArgs(ctx context.Context, query string, args []any) ([]T, error) {
	result := []T{}
	conn, done, err := d.connOrNew(ctx)
	if err != nil {
		return result, err
	}
	defer done()

	err = sqlitex.Execute(conn, query, &sqlitex.ExecOptions{
		Args: args,
		ResultFunc: func(stmt *sqlite.Stmt) error {
			mapped, err := d.mapper(stmt)
			if err != nil {
				return err
			}

			result = append(result, mapped)
			return nil
		},
	})

	return result, err
}

func (d *executor[T]) connOrNew(ctx context.Context) (*sqlite.Conn, func(), error) {
	if d.conn == nil {
		conn, done, err := d.pool.Conn(ctx)
		if err != nil {
			return nil, func() {}, err
		}
		return conn, done, nil
	} else {
		return d.conn, func() {}, nil
	}
}
