package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/helpers"
	"github.com/herumitra/ziidaapi/services"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login mengautentikasi pengguna dan menghasilkan JWT
func Login(c *fiber.Ctx) error {
	var req LoginRequest

	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	// Verifikasi kredensial dan hasilkan token JWT
	_, resp, err := services.AuthenticateUser(config.DB, config.RDB, req.Username, req.Password)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Incorrect username or password", nil)
	}

	// Kembalikan token sebagai response
	return helpers.JSONResponse(c, fiber.StatusOK, "Login successful", resp.Data)
}

// Logout menangani proses logout pengguna dengan menghapus token di Redis
func Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if err := config.RDB.Del(c.Context(), token).Err(); err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to logout", nil)
	}
	return helpers.JSONResponse(c, fiber.StatusOK, "Logout successful", nil)
}
