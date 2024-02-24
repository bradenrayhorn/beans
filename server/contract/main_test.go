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
	_, ds, _, stop := testutils.StartPoolWithDataSource(t)
	t.Cleanup(stop)

	sessionRepository := inmem.NewSessionRepository()
	contracts := contract.NewContracts(ds, sessionRepository)
	services := service.NewServices(ds, sessionRepository)
	adapter := contractadapter.New(contracts, services)

	specification.DoTests(t, adapter)
}
