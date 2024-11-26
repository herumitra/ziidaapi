package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwt "github.com/gofiber/jwt/v3"
	"github.com/herumitra/ziidaapi/helpers"
)

// JWTMiddleware memvalidasi token JWT pada setiap request yang membutuhkannya
func JWTMiddleware() fiber.Handler {
	return jwt.New(jwt.Config{
		SigningKey:  []byte("secret"),       // Gunakan key yang aman untuk signing JWT
		TokenLookup: "header:Authorization", // Token akan diambil dari header Authorization
		AuthScheme:  "Bearer",               // Menentukan bahwa token akan menggunakan skema "Bearer"
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Jika token tidak valid atau kadaluarsa, return custom error message
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ErrorResponse{
				FailedField: "Authorization",
				Tag:         "invalid_token",
				Value:       err.Error(), // Menampilkan error detail
			})
		},
	})
}
