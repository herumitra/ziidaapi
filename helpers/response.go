package helpers

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents the standard JSON response format / structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents the standard JSON error response format / structure
type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

// getStatusText returns "success" or "error" based on the HTTP status code
func getStatusText(status int) string {
	if status >= 200 && status < 300 {
		return "success"
	}
	return "error"
}

// JSONResponse sends a standard JSON response format / structure
func JSONResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	resp := Response{
		Status:  getStatusText(status),
		Message: message,
		Data:    data,
	}
	return c.Status(status).JSON(resp)
}
