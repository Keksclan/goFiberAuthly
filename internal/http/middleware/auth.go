package middleware

import (
	"context"
	"log/slog"
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
		rid, _ := c.Locals("requestid").(string)

		// Check required headers.
		for _, h := range requiredHeaders {
			if c.Get(h) == "" {
				slog.Warn("missing required header",
					"header", h,
					"method", c.Method(),
					"path", c.Path(),
					"request_id", rid,
				)
				return httperrors.NewBadRequest(c,
					httperrors.CodeMissingRequiredHeader,
					"missing required header: "+h,
				)
			}
		}

		// Extract Authorization header.
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			slog.Warn("missing authorization header",
				"method", c.Method(),
				"path", c.Path(),
				"request_id", rid,
			)
			return httperrors.NewUnauthorized(c,
				httperrors.CodeMissingAuthHeader,
				"missing authorization header",
			)
		}

		// Only Bearer scheme supported.
		if !strings.HasPrefix(authHeader, "Bearer ") {
			slog.Warn("unsupported authorization scheme",
				"method", c.Method(),
				"path", c.Path(),
				"request_id", rid,
			)
			return httperrors.NewUnauthorized(c,
				httperrors.CodeInvalidToken,
				"unsupported authorization scheme, expected Bearer",
			)
		}

		token := authHeader[7:]
		if token == "" {
			slog.Warn("empty bearer token",
				"method", c.Method(),
				"path", c.Path(),
				"request_id", rid,
			)
			return httperrors.NewUnauthorized(c,
				httperrors.CodeInvalidToken,
				"empty bearer token",
			)
		}

		slog.Debug("verifying token",
			"method", c.Method(),
			"path", c.Path(),
			"request_id", rid,
		)

		// Verify token via goAuthly engine.
		result, err := engine.Verify(context.Background(), token)
		if err != nil {
			slog.Warn("token verification failed",
				"error", err.Error(),
				"method", c.Method(),
				"path", c.Path(),
				"request_id", rid,
			)
			return httperrors.NewUnauthorized(c,
				httperrors.CodeInvalidToken,
				"token invalid or expired",
			)
		}

		slog.Info("token verified",
			"sub", result.Subject,
			"type", string(result.Type),
			"source", result.Source,
			"scopes", result.Scopes,
			"method", c.Method(),
			"path", c.Path(),
			"request_id", rid,
		)

		// Store result in locals for downstream handlers.
		c.Locals("authly", result)
		c.Locals("sub", result.Subject)
		c.Locals("scopes", result.Scopes)
		c.Locals("claims", result.Claims)

		return c.Next()
	}
}
