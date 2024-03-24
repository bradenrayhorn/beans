package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/http"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/service"
	"github.com/bradenrayhorn/beans/server/sqlite"
)

type Application struct {
	httpServer *http.Server

	pool *sqlite.Pool

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
	pool, err := sqlite.CreatePool(context.Background(), a.config.DbFilePath)

	if err != nil {
		return err
	}

	a.pool = pool

	a.datasource = sqlite.NewDataSource(pool)
	a.sessionRepository = inmem.NewSessionRepository()

	a.httpServer = http.NewServer(
		contract.NewContracts(a.datasource, a.sessionRepository),
		service.NewServices(a.datasource, a.sessionRepository),
	)
	if err := a.httpServer.Open(":" + a.config.Port); err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("http listening on port :%s", a.config.Port))

	return nil
}

func (a *Application) Stop() error {
	if err := a.httpServer.Close(); err != nil {
		return err
	}

	if err := a.pool.Close(context.Background()); err != nil {
		return err
	}

	return nil
}

func (a *Application) HttpServer() *http.Server {
	return a.httpServer
}
