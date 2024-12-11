package services

import (
	fmt "fmt"
	os "os"
	strings "strings"

	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

func GetUserID(c *fiber.Ctx) (string, error) {
	// Ambil token dari header Authorization
	token := c.Get("Authorization")

	// Hapus prefix "Bearer " jika ada
	if strings.HasPrefix(token, "Bearer ") {
		token = token[len("Bearer "):]
	}

	// Periksa jika token kosong
	if token == "" {
		return "", fmt.Errorf("missing token")
	}

	// Verifikasi token JWT
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET"))
		return secretKey, nil
	})

	if err != nil || !parsedToken.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// Ambil claims dari token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		return "", fmt.Errorf("invalid token claims")
	}

	// Ambil user_id dari claims
	userID := string(claims["sub"].(string))

	return userID, nil
}

func GetBranchID(c *fiber.Ctx) (string, error) {
	// Ambil token dari header Authorization
	token := c.Get("Authorization")

	// Hapus prefix "Bearer " jika ada
	if strings.HasPrefix(token, "Bearer ") {
		token = token[len("Bearer "):]
	}

	// Periksa jika token kosong
	if token == "" {
		return "", fmt.Errorf("missing token")
	}

	// Verifikasi token JWT
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET"))
		return secretKey, nil
	})

	if err != nil || !parsedToken.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// Ambil claims dari token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || claims["branch_id"] == nil {
		return "", fmt.Errorf("invalid token claims")
	}

	// Ambil branch_id dari claims
	branchID := string(claims["branch_id"].(string))

	return branchID, nil
}
