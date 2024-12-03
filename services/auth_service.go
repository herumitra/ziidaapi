package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/herumitra/ziidaapi/helpers"
	"github.com/herumitra/ziidaapi/models"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthenticateUser mengautentikasi pengguna dan menghasilkan JWT token
func AuthenticateUser(db *gorm.DB, rdb *redis.Client, username, password string) (string, *helpers.Response, error) {
	// Temukan pengguna di database
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", nil, fmt.Errorf("user not found: %v", err)
	}

	// Periksa apakah password cocok
	if !user.CheckPassword(password) {
		return "", nil, fmt.Errorf("incorrect password")
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("error generating JWT: %v", err)
	}

	// Simpan token di Redis dengan prefix "auth:"
	redisKey := "auth:" + token // Menambahkan prefix "auth:" pada key Redis
	ctx := context.Background()
	err = rdb.Set(ctx, redisKey, user.ID, 8*time.Hour).Err() // Set token dengan durasi kadaluarsa 8 jam
	if err != nil {
		return "", nil, fmt.Errorf("error saving token to Redis: %v", err)
	}

	// Return token, success response, and no error
	return token, &helpers.Response{
		Status:  "success",
		Message: "Login successful",
		Data:    map[string]string{"token": token},
	}, nil
}

// generateJWT menghasilkan token JWT untuk pengguna
func generateJWT(user models.User) (string, error) {
	// Define JWT claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(8 * time.Hour).Unix(),
	}
	// Generate the token using the claims and a signing key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Replace with your actual signing key (e.g., an environment variable)
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(secretKey)
}

func HashingPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err.Error()
	}
	password = string(hashedPassword)
	return password
}

func SetBranchIdToRedis(rdb *redis.Client, token string, branchID string) (string, *helpers.Response, error) {
	redisKey := "auth:" + token
	ctx := context.Background()
	err := rdb.HSet(ctx, redisKey, "branchID", branchID).Err()
	if err != nil {
		return "", nil, fmt.Errorf("error saving branchID to Redis: %v", err)
	}

	return branchID, &helpers.Response{
		Status:  "success",
		Message: "Branch telah dipilih",
		Data:    "branch_id: " + branchID,
	}, nil
}

func JWTDecodeID(token string) string {
	// Remove "Bearer " prefix
	if strings.HasPrefix(token, "Bearer ") {
		token = token[len("Bearer "):]
	}

	// Verify token JWT using secret key
	parsedToken, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET"))
		return secretKey, nil
	})

	// Get claim from token (example: user ID)
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return ""
	}

	// Get user ID from claims
	userID := claims["sub"].(float64)

	return fmt.Sprintf("%v", userID)

}
