package main

import (
	"os"
	"os/signal"

	"log/slog"
)

func main() {
	slog.Info("starting beans server")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	application := NewApplication(Config{
		Postgres: PostgresConfig{
			Addr:     "127.0.0.1:5432",
			Username: "postgres",
			Password: "password",
			Database: "beans",
		},
		Port: "8000",
	})

	if err := application.Start(); err != nil {
		panic(err)
	}

	<-c
	slog.Info("shutting down server")

	if err := application.Stop(); err != nil {
		panic(err)
	}
}
