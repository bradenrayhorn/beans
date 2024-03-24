package sqlite_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/bradenrayhorn/beans/server/internal/tests/datasource"
	"github.com/bradenrayhorn/beans/server/sqlite"
)

func TestSQLiteDatasource(t *testing.T) {
	b := make([]byte, 12)
	rand.Read(b)
	suffix := fmt.Sprintf("%x", b)[2:12]
	path := "/tmp/testdb-" + suffix
	defer func() {
		os.Remove(path)
	}()

	pool, err := sqlite.CreatePool(context.Background(), "file:"+path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		pool.Close(context.Background())
	}()

	ds := sqlite.NewDataSource(pool)

	datasource.DoTestDatasource(t, ds)
}
