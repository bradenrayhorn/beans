package main

import (
	"fmt"

	"github.com/bradenrayhorn/beans/http"
)

func main() {
	fmt.Println("Starting beans server")

	httpServer := http.NewServer()
	httpServer.Start()
}
