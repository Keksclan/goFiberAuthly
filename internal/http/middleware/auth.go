package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/keksclan/goAuthly/authly"

	httperrors "goFiberAuthly/internal/platform/errors"
)

// Auth returns a Fiber v3 middleware that validates Bearer tokens using goAuthly.
//
// Flow:
//  1. Check required headers (if configured)
//  2. Extract Authorization header (Bearer ...)
//  3. Verify token via goAuthly Engine (JWT validation, optional introspection fallback)
//  4. Store Result in c.Locals("authly") with sub, scopes, claims
//  5. On failure: return structured JSON error (401/403)
func Auth(engine *authly.Engine, requiredHeaders []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Check required headers.
		for _, h := range requiredHeaders {
			if c.Get(h) == "" {
				return httperrors.NewBadRequest(c,
					httperrors.CodeMissingRequiredHeader,
					"missing required header: "+h,
				)
			}
		}

		// Extract Authorization header.
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return httperrors.NewUnauthorized(c,
				httperrors.CodeMissingAuthHeader,
				"missing authorization header",
			)
		}

		// Only Bearer scheme supported.
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return httperrors.NewUnauthorized(c,
				httperrors.CodeInvalidToken,
				"unsupported authorization scheme, expected Bearer",
			)
		}

		token := authHeader[7:]
		if token == "" {
			return httperrors.NewUnauthorized(c,
				httperrors.CodeInvalidToken,
				"empty bearer token",
			)
		}

		// Verify token via goAuthly engine.
		result, err := engine.Verify(context.Background(), token)
		if err != nil {
			return httperrors.NewUnauthorized(c,
				httperrors.CodeInvalidToken,
				"token invalid or expired",
			)
		}

		// Store result in locals for downstream handlers.
		c.Locals("authly", result)
		c.Locals("sub", result.Subject)
		c.Locals("scopes", result.Scopes)
		c.Locals("claims", result.Claims)

		return c.Next()
	}
}
