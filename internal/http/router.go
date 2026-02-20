package http

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/keksclan/goAuthly/authly"

	"goFiberAuthly/internal/config"
	"goFiberAuthly/internal/http/handlers"
	"goFiberAuthly/internal/http/middleware"
)

// SetupRoutes registers all routes on the Fiber app.
func SetupRoutes(app *fiber.App, engine *authly.Engine, cfg *config.Config, ready *bool) {
	// Global middleware.
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	slog.Info("global middleware registered", "middleware", []string{"RequestID", "Logger"})

	// Public endpoints.
	app.Get("/healthz", handlers.Healthz)
	app.Get("/readyz", handlers.Readyz(ready))
	slog.Info("public routes registered", "routes", []string{"/healthz", "/readyz"})

	// Protected endpoints (goAuthly auth middleware).
	protected := app.Group("", middleware.Auth(engine, cfg.Auth.RequiredHeaders))
	protected.Get("/me", handlers.Me)
	slog.Info("protected routes registered",
		"routes", []string{"/me"},
		"required_headers", cfg.Auth.RequiredHeaders,
	)
}
