//go:build slow

package postgres_test

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/internal/tests/datasource"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
)

func TestPostgresDatasource(t *testing.T) {
	_, ds, _, stop := testutils.StartPoolWithDataSource(t)
	defer stop()

	datasource.DoTestDatasource(t, ds)
}
