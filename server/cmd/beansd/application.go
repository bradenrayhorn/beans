package main

import (
	"fmt"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/http"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/postgres"
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
	pool, err := postgres.CreatePool(
		fmt.Sprintf("postgres://%s:%s@%s/%s",
			a.config.Postgres.Username,
			a.config.Postgres.Password,
			a.config.Postgres.Addr,
			a.config.Postgres.Database,
		))

	if err != nil {
		panic(err)
	}
	a.pool = pool

	a.datasource = postgres.NewDataSource(pool)
	a.sessionRepository = inmem.NewSessionRepository()

	a.httpServer = http.NewServer(contract.NewContracts(
		a.datasource,
		a.sessionRepository,
	))
	if err := a.httpServer.Open(":" + a.config.Port); err != nil {
		panic(err)
	}

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
