package middleware

import (
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
		}
		c.Locals("requestid", rid)
		c.Set("X-Request-Id", rid)
		return c.Next()
	}
}
