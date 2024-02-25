package service_test

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/internal/fake"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/service"
)

func makeServices(t *testing.T) (*service.All, *testutils.Factory, beans.DataSource, beans.SessionRepository) {
	ds := fake.NewDataSource()
	factory := testutils.NewFactory(t, ds)
	sessionRepository := inmem.NewSessionRepository()

	return service.NewServices(ds, sessionRepository), factory, ds, sessionRepository
}
