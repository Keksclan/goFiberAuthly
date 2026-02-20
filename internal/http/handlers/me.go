package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/keksclan/goAuthly/authly"
)

// Me returns the authenticated user's identity from the goAuthly Result.
// Protected by Auth middleware – authly.Result is expected in c.Locals("authly").
func Me(c fiber.Ctx) error {
	rid, _ := c.Locals("requestid").(string)

	result, ok := c.Locals("authly").(*authly.Result)
	if !ok || result == nil {
		slog.Error("me handler: authly result missing from context",
			"request_id", rid,
		)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Extract aud and iss from claims if available.
	aud, _ := result.Claims["aud"]
	iss, _ := result.Claims["iss"]

	slog.Info("me handler: returning user identity",
		"sub", result.Subject,
		"type", string(result.Type),
		"source", result.Source,
		"request_id", rid,
	)

	return c.JSON(fiber.Map{
		"sub":    result.Subject,
		"scopes": result.Scopes,
		"aud":    aud,
		"iss":    iss,
		"type":   string(result.Type),
		"source": result.Source,
	})
}
