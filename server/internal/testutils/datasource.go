package testutils

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TmpDatasource(tb testing.TB) (beans.DataSource, func()) {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	require.NoError(tb, err)
	suffix := fmt.Sprintf("%x", b)[2:12]
	path := "/tmp/testdb-" + suffix

	pool, err := sqlite.CreatePool(context.Background(), "file:"+path)
	if err != nil {
		tb.Fatal(err)
	}

	return sqlite.NewDataSource(pool),
		func() {
			assert.NoError(tb, pool.Close(context.Background()))
			assert.NoError(tb, os.Remove(path))
		}
}
