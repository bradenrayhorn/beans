package main

import (
	"fmt"

	"github.com/bradenrayhorn/beans/http"
	"github.com/bradenrayhorn/beans/inmem"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Application struct {
	httpServer *http.Server
	pool       *pgxpool.Pool

	config Config
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

	userRepository := postgres.NewUserRepository(pool)
	userService := &logic.UserService{UserRepository: userRepository}
	sessionRepository := inmem.NewSessionRepository()

	a.httpServer = http.NewServer(userRepository, userService, sessionRepository)
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
