package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
)

// Logger provides structured request logging middleware.
func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		rid, _ := c.Locals("requestid").(string)
		status := c.Response().StatusCode()

		attrs := []any{
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"duration", time.Since(start).String(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
			"request_id", rid,
			"bytes_sent", len(c.Response().Body()),
		}

		// Include sub if available (set by auth middleware).
		if sub, ok := c.Locals("sub").(string); ok && sub != "" {
			attrs = append(attrs, "sub", sub)
		}

		switch {
		case status >= 500:
			slog.Error("request", attrs...)
		case status >= 400:
			slog.Warn("request", attrs...)
		default:
			slog.Info("request", attrs...)
		}

		return err
	}
}
