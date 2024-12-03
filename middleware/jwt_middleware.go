package middleware

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/helpers"
	"github.com/herumitra/ziidaapi/models"
	"github.com/redis/go-redis/v9"
)

// JWTMiddleware is a middleware to check if the token is valid and the user status is active
func JWTMiddleware(c *fiber.Ctx) error {
	// Get token from header
	token := c.Get("Authorization")

	// Remove "Bearer " prefix
	if strings.HasPrefix(token, "Bearer ") {
		token = token[len("Bearer "):]
	}
	// Check if token is empty
	if token == "" {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Missing token", nil)
	}

	// Verify token JWT using secret key
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET"))
		return secretKey, nil
	})
	// Check if token is valid
	if err != nil || !parsedToken.Valid {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	// Get claim from token (example: user ID)
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Invalid token claims", nil)
	}
	// Get user ID from claims
	userID := claims["sub"].(float64)

	// Verify token JWT using Redis to check if the token is still valid
	ctx := context.Background()
	redisKey := fmt.Sprintf("auth:%s", token) // Gunakan prefix "auth:" untuk key Redis
	rdb := config.RDB                         // Ambil Redis client dari context
	redisValue, err := rdb.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Token not found in Redis", nil)
	}
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Redis error", nil)
	}

	// Cek status user dengan query ke database
	if redisValue == fmt.Sprintf("%v", userID) {
		// Ambil instance DB dari model
		var users []models.User

		// Cek apakah user status aktif
		if config.DB.Where("id = ? AND user_status ='active'", userID).First(&users).Error == nil {
			// User aktif, lanjutkan ke handler berikutnya
			return c.Next()
		}

		// // Jika user status tidak aktif, kembalikan error
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "User not active, please call administrator", err)
		// fmt.Printf("User ID: %v, User Status: %v\n", userID, users)
	}

	// Jika token tidak valid, kembalikan error
	return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Unauthorized access", nil)
}
