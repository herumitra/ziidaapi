package middleware

import (
	fiber "github.com/gofiber/fiber/v2"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
)

func JWTMiddleware(c *fiber.Ctx) error {
	helpers.TokenValidation(c, "sub")

	// Lanjutkan ke handler berikutnya
	return c.Next()
}

func RoleMiddleware(allowedRoles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil user_role dari token
		userRole, _ := helpers.GetClaimsToken(c, "user_role")

		// Periksa apakah user_role termasuk dalam allowedRoles
		for _, role := range allowedRoles {
			if string(role) == userRole {
				return c.Next() // Akses diizinkan
			}
		}

		// Jika role tidak sesuai, tolak akses
		return helpers.JSONResponse(c, fiber.StatusForbidden, "Forbidden", "You don't have permission to access this resource!")
	}
}
