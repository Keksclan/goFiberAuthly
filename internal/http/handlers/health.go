package handlers

import "github.com/gofiber/fiber/v3"

// Healthz always returns 200 OK.
func Healthz(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

// Readyz returns 200 if config is loaded and goAuthly engine is initialized.
func Readyz(ready *bool) fiber.Handler {
	return func(c fiber.Ctx) error {
		if ready == nil || !*ready {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not ready",
			})
		}
		return c.JSON(fiber.Map{"status": "ready"})
	}
}
