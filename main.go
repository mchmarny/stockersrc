package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger = log.New(os.Stdout, "[SOURCE] ", 0)
)

func main() {

	logger.Println("Initializing...")

	// context to cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		logger.Println(<-ch)
		cancel()
		os.Exit(0)
	}()

	// TODO: Handle graceful shutdown
	go getContent(ctx)

	// wait
	for {
		select {
		case <-ctx.Done():
			break
		default:
			//do nothing here
		}
	}

}
