package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// watchSignals starts a goroutine that cancels ctx on SIGINT or SIGTERM.
// This enables graceful shutdown when the user presses Ctrl+C in watch mode.
func watchSignals(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		cancel()
	}()
}
