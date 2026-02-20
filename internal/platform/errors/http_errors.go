package errors

import "github.com/gofiber/fiber/v3"

// Standard error codes used in JSON error responses.
const (
	CodeUnauthorized          = "unauthorized"
	CodeInvalidToken          = "invalid_token"
	CodeTokenExpired          = "token_expired"
	CodeForbidden             = "forbidden"
	CodeMissingRequiredHeader = "missing_required_header"
	CodeMissingAuthHeader     = "missing_authorization_header"
)

// ErrorResponse is the standard JSON error format.
type ErrorResponse struct {
	Error     string `json:"error"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// NewUnauthorized returns a 401 JSON error response.
func NewUnauthorized(c fiber.Ctx, code, message string) error {
	reqID, _ := c.Locals("requestid").(string)
	return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
		Error:     CodeUnauthorized,
		Code:      code,
		Message:   message,
		RequestID: reqID,
	})
}

// NewForbidden returns a 403 JSON error response.
func NewForbidden(c fiber.Ctx, code, message string) error {
	reqID, _ := c.Locals("requestid").(string)
	return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
		Error:     CodeForbidden,
		Code:      code,
		Message:   message,
		RequestID: reqID,
	})
}

// NewBadRequest returns a 400 JSON error response.
func NewBadRequest(c fiber.Ctx, code, message string) error {
	reqID, _ := c.Locals("requestid").(string)
	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		Error:     "bad_request",
		Code:      code,
		Message:   message,
		RequestID: reqID,
	})
}
