package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	fmt.Println("Starting beans server")

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
	fmt.Println("shutting down server")

	if err := application.Stop(); err != nil {
		panic(err)
	}
}
