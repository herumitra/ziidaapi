package helpers

import "github.com/gofiber/fiber/v2"

// Response represents the standard JSON response structure.
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents a structured way to express API errors.
type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

// JSONResponse sends a formatted JSON response with a standard structure.
func JSONResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	resp := Response{
		Status:  getStatusText(status),
		Message: message,
		Data:    data,
	}
	return c.Status(status).JSON(resp)
}

// getStatusText returns "success" or "error" based on the HTTP status code.
func getStatusText(status int) string {
	if status >= 200 && status < 300 {
		return "success"
	}
	return "error"
}
