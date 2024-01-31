package postgres_test

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/internal/tests/datasource"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
)

func TestPayees(t *testing.T) {
	t.Parallel()
	_, ds, _, stop := testutils.StartPoolWithDataSource(t)
	defer stop()

	datasource.TestPayeeRepository(t, ds)
}
