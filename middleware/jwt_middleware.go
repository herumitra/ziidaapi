package middleware

import (
	context "context"
	fmt "fmt"
	strings "strings"

	os "os"

	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	config "github.com/herumitra/ziidaapi/config"
	helpers "github.com/herumitra/ziidaapi/helpers"
	redis "github.com/redis/go-redis/v9"
)

func JWTMiddleware(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	token := c.Get("Authorization")

	// Remove prefix "Bearer " jika ada
	if strings.HasPrefix(token, "Bearer ") {
		token = token[len("Bearer "):]
	}

	// Check token jika kosong
	if token == "" {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Missing token", nil)
	}

	// Verifikasi token menggunakan secret key
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET"))
		return secretKey, nil
	})

	// Jika token tidak valid
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Invalid token", "Try to login again")
	}

	// Verify token JWT using Redis to check if the token is still valid
	ctx := context.Background()
	redisKey := fmt.Sprintf("auth:%s", token) // Gunakan prefix "auth:" untuk key Redis
	rdb := config.RDB                         // Ambil Redis client dari context
	redisValue, err := rdb.Get(ctx, redisKey).Result()
	if err != redis.Nil {
		// Jika token ditemukan di Redis, maka gagal login
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Using token failed, token was revoked", nil)
	}

	fmt.Println(redisValue)

	// Jika token valid
	if parsedToken.Valid {
		return c.Next()
	}

	// Jika token tidak valid
	return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Unauthorized access", err)
}
