package main

import (
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// TODO - read this from config
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("starting beans server")

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
	log.Info().Msg("shutting down server")

	if err := application.Stop(); err != nil {
		panic(err)
	}
}
