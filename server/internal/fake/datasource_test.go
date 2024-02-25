package fake_test

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/internal/fake"
	"github.com/bradenrayhorn/beans/server/internal/tests/datasource"
)

func TestFakeDataSource(t *testing.T) {
	ds := fake.NewDataSource()

	datasource.DoTestDatasource(t, ds)
}
