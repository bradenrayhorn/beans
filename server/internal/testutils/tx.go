package testutils

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
)

func MustRollback(t testing.TB, tx beans.Tx) {
	err := tx.Rollback(context.Background())
	// Ignore ErrTxClosed as transaction may have already been committed.
	if err != nil {
		t.Error("Failed to rollback tx", err)
	}
}
