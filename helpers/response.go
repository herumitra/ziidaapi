package helpers

import "github.com/gofiber/fiber/v2"

// Response is a struct for standard response formatting
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse is a struct to format error message
type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

// JSONResponse formats the response for success or error
func JSONResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	resp := Response{
		Status:  getStatusText(status),
		Message: message,
		Data:    data,
	}
	return c.Status(status).JSON(resp)
}

// Helper function to get status text (success/error)
func getStatusText(status int) string {
	if status >= 200 && status < 300 {
		return "success"
	}
	return "error"
}
