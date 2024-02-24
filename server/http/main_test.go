package http_test

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/http"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/service"
	"github.com/bradenrayhorn/beans/server/specification"
	"github.com/bradenrayhorn/beans/server/specification/httpadapter"
)

func TestHTTP(t *testing.T) {
	_, ds, _, stop := testutils.StartPoolWithDataSource(t)
	defer stop()
	sessionRepository := inmem.NewSessionRepository()
	httpServer := http.NewServer(
		contract.NewContracts(ds, sessionRepository),
		service.NewServices(ds, sessionRepository),
	)

	if err := httpServer.Open(":0"); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := httpServer.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	adapter := httpadapter.New("http://" + httpServer.GetBoundAddr())

	specification.DoTests(t, adapter)
}
