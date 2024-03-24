package sqlite

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type Tx struct {
	conn *sqlite.Conn

	release    func(*error)
	returnConn func()

	released bool
}

func (t *Tx) Commit(ctx context.Context) error {
	if !t.released {
		var err error
		err = nil
		t.release(&err)
		t.returnConn()
		t.released = true
	}
	return nil
}

func (t *Tx) Rollback(ctx context.Context) error {
	if !t.released {
		err := errors.New("rollback tx")
		t.release(&err)
		t.returnConn()
		t.released = true
	}
	return nil
}

type txManager struct {
	pool *Pool
}

func (m *txManager) Create(ctx context.Context) (beans.Tx, error) {
	conn, done, err := m.pool.Conn(ctx)
	if err != nil {
		return nil, err
	}

	release := sqlitex.Transaction(conn)
	return &Tx{
		conn:       conn,
		release:    release,
		returnConn: done,
	}, nil
}
