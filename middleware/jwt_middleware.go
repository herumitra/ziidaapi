package middleware

import (
	context "context"
	fmt "fmt"
	log "log"
	os "os"
	strings "strings"

	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	config "github.com/herumitra/ziidaapi/config"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
)

func JWTMiddleware(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	token := c.Get("Authorization")

	// Remove prefix "Bearer " jika ada
	if strings.HasPrefix(token, "Bearer ") {
		token = token[len("Bearer "):]
	}

	// Periksa jika token kosong
	if token == "" {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Missing token", "Insert valid token to access this endpoint!")
	}

	// Cek apakah token masuk dalam blacklist Redis
	ctx := context.Background()
	redisKey := fmt.Sprintf("blacklist:%s", token)
	rdb := config.RDB
	isBlacklisted, err := rdb.Exists(ctx, redisKey).Result()

	if err != nil {
		log.Printf("Error checking token in Redis: %v", err)
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Token verification failed", "Server error!")
	}

	if isBlacklisted > 0 {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Using token failed", "Token was revoked, please login again!")
	}

	// Verifikasi token menggunakan secret key
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET"))
		return secretKey, nil
	})

	if err != nil || !parsedToken.Valid {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Invalid token", "Try to login again!")
	}

	// Periksa klaim token (opsional, misalnya validasi user_id, role, dll.)
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Invalid token claims", "Try to login again!")
	}

	// Lanjutkan ke handler berikutnya
	return c.Next()
}

func RoleMiddleware(allowedRoles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		token := c.Get("Authorization")

		// Remove prefix "Bearer " jika ada
		if strings.HasPrefix(token, "Bearer ") {
			token = token[len("Bearer "):]
		}

		// Periksa jika token kosong
		if token == "" {
			return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Missing token", "Insert valid token to access this endpoint!")
		}

		// Parse token dan ambil klaim
		parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			secretKey := []byte(os.Getenv("JWT_SECRET"))
			return secretKey, nil
		})

		if err != nil || !parsedToken.Valid {
			return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Invalid token", "Try to login again!")
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || claims["user_role"] == nil {
			return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Invalid token claims", "Try to login again!")
		}

		// Ambil user_role dari token
		userRole := claims["user_role"].(string)

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
