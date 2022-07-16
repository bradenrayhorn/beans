package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bradenrayhorn/beans/http"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/bradenrayhorn/beans/postgres"
)

func main() {
	fmt.Println("Starting beans server")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	pool, err := postgres.CreatePool("postgres://postgres:password@127.0.0.1:5432/beans")
	if err != nil {
		panic(err)
	}
	userRepository := postgres.NewUserRepository(pool)
	userService := &logic.UserService{UserRepository: userRepository}

	httpServer := http.NewServer(userService)
	if err := httpServer.Open(); err != nil {
		panic(err)
	}

	<-c
	fmt.Println("shutting down server")

	if err := httpServer.Close(); err != nil {
		panic(err)
	}
}
