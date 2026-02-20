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

	// Initialize structured JSON logger with configured log level.
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Server.SlogLevel(),
	})
	slog.SetDefault(slog.New(logHandler))

	slog.Info("logger initialized", "level", cfg.Server.LogLevel)
	slog.Info("config loaded",
		"port", cfg.Server.Port,
		"log_level", cfg.Server.LogLevel,
		"auth_issuer", cfg.Auth.Issuer,
		"auth_audience", cfg.Auth.Audience,
		"auth_jwks_url", cfg.Auth.JWKSURL,
		"auth_introspection_url", cfg.Auth.IntrospectionURL,
		"auth_client_id", cfg.Auth.ClientID,
		"auth_required_headers", cfg.Auth.RequiredHeaders,
	)

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
