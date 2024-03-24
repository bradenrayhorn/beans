package contract_test

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/service"
	"github.com/bradenrayhorn/beans/server/specification"
	"github.com/bradenrayhorn/beans/server/specification/contractadapter"
)

func TestContracts(t *testing.T) {
	ds, done := testutils.TmpDatasource(t)
	t.Cleanup(done)

	sessionRepository := inmem.NewSessionRepository()
	contracts := contract.NewContracts(ds, sessionRepository)
	services := service.NewServices(ds, sessionRepository)
	adapter := contractadapter.New(contracts, services)

	specification.DoTests(t, adapter)
}
