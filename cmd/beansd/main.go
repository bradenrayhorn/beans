package main

import (
	"fmt"

	"github.com/bradenrayhorn/beans/http"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/bradenrayhorn/beans/postgres"
)

func main() {
	fmt.Println("Starting beans server")

	pool, err := postgres.CreatePool()
	if err != nil {
		panic(err)
	}
	userRepository := postgres.NewUserRepository(pool)
	userService := &logic.UserService{UserRepository: userRepository}

	httpServer := http.NewServer(userService)
	httpServer.Start()
}
