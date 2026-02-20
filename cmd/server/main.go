package main

import (
	"log/slog"
	"os"

	"goFiberAuthly/internal/app"
	"goFiberAuthly/internal/config"
)

func main() {
	// Load configuration via goConfy (YAML + ENV macro expansion).
	cfg, err := config.Load("")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize application (Fiber + goAuthly engine).
	application, err := app.New(cfg)
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Start server in a goroutine.
	go func() {
		addr := ":" + cfg.Server.Port
		slog.Info("starting server", "addr", addr)
		if err := application.Fiber.Listen(addr); err != nil {
			slog.Error("server error", "error", err)
		}
	}()

	// Block until shutdown signal.
	application.GracefulShutdown()
}
