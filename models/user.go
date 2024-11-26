package models

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
}

// HashPassword melakukan hashing password pengguna
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword membandingkan password yang diberikan dengan password yang ter-hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// AuthenticateUser checks user credentials and generates a JWT token
func AuthenticateUser(db *gorm.DB, rdb *redis.Client, username, password string) (string, error) {
	// Find the user in the database
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", fmt.Errorf("user not found: %v", err)
	}

	// Check if the provided password matches the stored hashed password
	if !user.CheckPassword(password) {
		return "", fmt.Errorf("incorrect password")
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return "", fmt.Errorf("error generating JWT: %v", err)
	}

	// Store the token in Redis with a 8-hour expiration
	ctx := context.Background()
	err = rdb.Set(ctx, token, user.ID, 8*time.Hour).Err() // Gunakan time.Hour di sini
	if err != nil {
		return "", fmt.Errorf("error saving token to Redis: %v", err)
	}

	return token, nil
}

// generateJWT generates a JWT token for the user
func generateJWT(user User) (string, error) {
	// Define JWT claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(8 * time.Hour).Unix(),
	}
	// Generate the token using the claims and a signing key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Replace with your actual signing key (e.g., an environment variable)
	secretKey := []byte("your-secret-key")
	return token.SignedString(secretKey)
}
