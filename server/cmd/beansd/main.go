package main

import (
	"os"
	"os/signal"

	"log/slog"
)

func main() {
	slog.Info("starting beansd...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	application := NewApplication(config)

	if err := application.Start(); err != nil {
		panic(err)
	}

	<-c
	slog.Info("shutting down server")

	if err := application.Stop(); err != nil {
		panic(err)
	}
}
