package contract_test

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/specification"
)

func TestAccount(t *testing.T) {
	t.Parallel()
	_, ds, _, stop := testutils.StartPoolWithDataSource(t)
	defer stop()

	contracts := contract.NewContracts(ds, inmem.NewSessionRepository())
	adapter := ContractsAdapter{contracts}

	specification.TestAccounts(t, &adapter)
}
