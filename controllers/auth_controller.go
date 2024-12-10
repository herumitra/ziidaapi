package controllers

import (
	"context"
	fmt "fmt"
	"log"
	"os"
	"strings"
	time "time"

	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	config "github.com/herumitra/ziidaapi/config"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
	bcrypt "golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// generateJWT menghasilkan token JWT untuk pengguna
func generateJWT(user models.User) (string, error) {
	// Define JWT claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(8 * time.Hour).Unix(),
		// "exp": time.Now().Add(5 * time.Minute).Unix(),
	}
	// Generate the token using the claims and a signing key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Replace with your actual signing key (e.g., an environment variable)
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(secretKey)
}

// Fungsi untuk menambahkan token ke blacklist di Redis dengan TTL 8 jam
func blacklistToken(token string) error {
	// Verify token JWT using Redis to check if the token is still valid
	ctx := context.Background()
	redisKey := fmt.Sprintf("auth:%s", token) // Gunakan prefix "auth:" untuk key Redis
	rdb := config.RDB                         // Ambil Redis client dari context

	// Set token ke Redis dengan TTL 8 jam
	err := rdb.Set(ctx, redisKey, "blacklisted", time.Hour*8).Err()
	if err != nil {
		log.Printf("Failed to blacklist token: %v", err)
		return err
	}

	return nil
}

// Function Login menangani login pengguna
func Login(c *fiber.Ctx) error {
	// Definisikan variabel loginRequest dan user
	var loginRequest LoginRequest
	var user models.User

	// Parse input JSON menjadi struct LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Cari user berdasarkan username
	if err := config.DB.Where("username = ? AND user_status = 'active'", loginRequest.Username).First(&user).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Login failed", "User is not active, call admin to activated your account !")
	}

	// Bandingkan password input dengan password yang sudah di-hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Login failed", "Invalid username or password")
	}

	// Buat token JWT
	token, err := generateJWT(user)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Login failed", "Failed to generate token")
	}

	// Jika username dan password cocok, lanjutkan proses (misalnya buat token JWT)
	return helpers.JSONResponse(c, fiber.StatusOK, "Login successful", "token:"+token)
}

// Function Logout menangani logout pengguna
func Logout(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	token := c.Get("Authorization")

	// Remove prefix "Bearer " jika ada
	if strings.HasPrefix(token, "Bearer ") {
		token = token[len("Bearer "):]
	}

	if token == "" {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Missing token", "Insert valid token to access this endpoint !")
	}

	// Blacklist token JWT
	if err := blacklistToken(token); err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Logout failed", "Failed to blacklist token")
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "Logout successful", "Logout successful")
}
