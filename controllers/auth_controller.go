package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/herumitra/ziidaapi/helpers"
	"github.com/herumitra/ziidaapi/konfig"
	"github.com/herumitra/ziidaapi/models"
)

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
		BranchID uint   `json:"branchID"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	var user models.User
	konfig.DB.Where("username = ? AND branch_id = ?", input.Username, input.BranchID).First(&user)

	if user.Username == "" || !user.CheckPassword(input.Password) {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Invalid username, password or branch", nil)
	}

	// Create JWT Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Username
	claims["branch_id"] = user.BranchID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to generate token", nil)
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "Login successful", fiber.Map{
		"token": t,
	})
}
