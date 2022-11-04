package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	err := run(ctx, os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	var command string
	if len(args) > 0 {
		command = args[0]
	}

	switch command {
	case "migrate":
		return runMigrations(ctx)
	default:
		return errors.New("unknown command")
	}
}
