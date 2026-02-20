package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

// Healthz always returns 200 OK.
func Healthz(c fiber.Ctx) error {
	rid, _ := c.Locals("requestid").(string)
	slog.Debug("healthz check", "request_id", rid)
	return c.JSON(fiber.Map{"status": "ok"})
}

// Readyz returns 200 if config is loaded and goAuthly engine is initialized.
func Readyz(ready *bool) fiber.Handler {
	return func(c fiber.Ctx) error {
		rid, _ := c.Locals("requestid").(string)
		if ready == nil || !*ready {
			slog.Warn("readyz check failed: not ready", "request_id", rid)
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not ready",
			})
		}
		slog.Debug("readyz check: ready", "request_id", rid)
		return c.JSON(fiber.Map{"status": "ready"})
	}
}
