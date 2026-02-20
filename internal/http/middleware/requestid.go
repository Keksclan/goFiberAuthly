package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// RequestID ensures every request has an X-Request-Id header.
// If the header is missing, a new UUID is generated.
func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		rid := c.Get("X-Request-Id")
		if rid == "" {
			rid = uuid.New().String()
			slog.Debug("generated request id", "request_id", rid)
		} else {
			slog.Debug("using provided request id", "request_id", rid)
		}
		c.Locals("requestid", rid)
		c.Set("X-Request-Id", rid)
		return c.Next()
	}
}
