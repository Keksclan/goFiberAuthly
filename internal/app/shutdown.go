package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// GracefulShutdown blocks until SIGINT or SIGTERM is received,
// then shuts down the Fiber app gracefully.
func (a *Application) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	slog.Info("shutting down", "signal", sig.String())

	a.Ready = false

	if err := a.Fiber.Shutdown(); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	slog.Info("server stopped")
}
