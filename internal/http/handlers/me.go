package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/keksclan/goAuthly/authly"
)

// Me returns the authenticated user's identity from the goAuthly Result.
// Protected by Auth middleware – authly.Result is expected in c.Locals("authly").
func Me(c fiber.Ctx) error {
	result, ok := c.Locals("authly").(*authly.Result)
	if !ok || result == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Extract aud and iss from claims if available.
	aud, _ := result.Claims["aud"]
	iss, _ := result.Claims["iss"]

	return c.JSON(fiber.Map{
		"sub":    result.Subject,
		"scopes": result.Scopes,
		"aud":    aud,
		"iss":    iss,
		"type":   string(result.Type),
		"source": result.Source,
	})
}
