package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Moon1706/ginkgo2allure/cmd"
)

func main() {
	ctx := setupSignalHandler()
	cmd.ExecuteContext(ctx)
}

func setupSignalHandler() context.Context {
	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM}
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, len(shutdownSignals))
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1)
	}()
	return ctx
}
