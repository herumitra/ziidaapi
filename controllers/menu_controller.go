package controllers

import (
	os "os"

	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	config "github.com/herumitra/ziidaapi/config"
	models "github.com/herumitra/ziidaapi/models"
)

func GetMenu(c *fiber.Ctx) error {
	// Ambil user_role dari klaim JWT
	token := c.Get("Authorization")
	parsedToken, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET"))
		return secretKey, nil
	})

	claims := parsedToken.Claims.(jwt.MapClaims)
	userRole := claims["user_role"].(string)

	// Query menu dari database berdasarkan role
	var menus []models.Menu
	err := config.DB.Where("? = ANY (allowed_roles)", userRole).Find(&menus).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch menu")
	}

	return c.JSON(menus)
}
