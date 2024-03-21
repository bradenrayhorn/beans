package main

import (
	"fmt"
	"log/slog"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/http"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/bradenrayhorn/beans/server/service"
)

type Application struct {
	httpServer *http.Server
	pool       *postgres.DbPool

	config Config

	datasource        beans.DataSource
	sessionRepository beans.SessionRepository
}

func NewApplication(c Config) *Application {
	return &Application{
		config: c,
	}
}

func (a *Application) Start() error {
	pool, err := postgres.CreatePool(a.config.PostgresURL)

	if err != nil {
		panic(err)
	}
	a.pool = pool

	a.datasource = postgres.NewDataSource(pool)
	a.sessionRepository = inmem.NewSessionRepository()

	a.httpServer = http.NewServer(
		contract.NewContracts(a.datasource, a.sessionRepository),
		service.NewServices(a.datasource, a.sessionRepository),
	)
	if err := a.httpServer.Open(":" + a.config.Port); err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("http listening on port :%s", a.config.Port))

	return nil
}

func (a *Application) Stop() error {
	if err := a.httpServer.Close(); err != nil {
		return err
	}

	a.pool.Close()

	return nil
}

func (a *Application) HttpServer() *http.Server {
	return a.httpServer
}
