package sqlite_test

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/internal/tests/datasource"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
)

func TestSQLiteDatasource(t *testing.T) {
	ds, done := testutils.TmpDatasource(t)
	defer done()

	datasource.DoTestDatasource(t, ds)
}
