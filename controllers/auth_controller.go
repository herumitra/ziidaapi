package controllers

import (
	context "context"
	fmt "fmt"
	log "log"
	os "os"
	strings "strings"
	time "time"

	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	config "github.com/herumitra/ziidaapi/config"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
	services "github.com/herumitra/ziidaapi/services"
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
	// Parse token untuk mendapatkan waktu kedaluwarsa (exp)
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET"))
		return secretKey, nil
	})

	if err != nil || !parsedToken.Valid {
		log.Printf("Failed to parse token: %v", err)
		return err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || claims["exp"] == nil {
		log.Println("Invalid token claims, no exp found")
		return fmt.Errorf("invalid token claims")
	}

	// Hitung waktu kedaluwarsa token
	expiryUnix := int64(claims["exp"].(float64)) // Klaim `exp` adalah float64
	expiryTime := time.Unix(expiryUnix, 0)
	ttl := time.Until(expiryTime)

	// Pastikan TTL valid
	if ttl <= 0 {
		log.Println("Token already expired")
		return fmt.Errorf("token already expired")
	}

	// Tambahkan token ke Redis dengan TTL
	ctx := context.Background()
	redisKey := fmt.Sprintf("blacklist:%s", token)
	err = config.RDB.Set(ctx, redisKey, "blacklisted", ttl).Err()
	if err != nil {
		log.Printf("Failed to blacklist token: %v", err)
		return err
	}

	log.Printf("Token blacklisted successfully with TTL: %v", ttl)
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

func SetBranch(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	token := c.Get("Authorization")

	// Hapus prefix "Bearer " jika ada
	if strings.HasPrefix(token, "Bearer ") {
		token = token[len("Bearer "):]
	}

	// Periksa jika token kosong
	if token == "" {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Missing token", "Insert valid token to access this endpoint!")
	}

	// Verifikasi token JWT untuk mendapatkan user ID
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET"))
		return secretKey, nil
	})

	if err != nil || !parsedToken.Valid {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Invalid token", "Try to login again!")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Invalid token claims", "Try to login again!")
	}

	// Ambil user ID dari klaim token
	userID := string(claims["sub"].(string))

	// Parse input JSON untuk mendapatkan branch ID
	var request struct {
		BranchID string `json:"branch_id" validate:"required"`
	}
	if err := c.BodyParser(&request); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Periksa apakah branch_id valid untuk user ini
	var userBranch models.UserBranch
	if err := config.DB.Where("user_id = ? AND branch_id = ?", userID, request.BranchID).First(&userBranch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusForbidden, "Invalid branch ID", "Branch not associated with this user!")
	}

	// Ambil user_role dari tabel users berdasarkan user_id
	var user models.User
	if err := config.DB.Select("user_role").Where("id = ?", userID).First(&user).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to set branch", "Unable to retrieve user role")
	}

	// Buat token JWT baru dengan klaim branch_id dan user_role
	newToken, err := generateBranchJWTWithRole(userID, request.BranchID, string(user.UserRole))
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to set branch", "Failed to generate new token")
	}

	// Tambahkan token lama ke Redis blacklist
	if err := blacklistToken(token); err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to set branch", "Failed to blacklist old token")
	}

	// Berikan token baru ke pengguna
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch set successfully", fiber.Map{
		"new_token": newToken,
	})
}

func generateBranchJWTWithRole(userID string, branchID string, userRole string) (string, error) {
	// Definisikan klaim untuk token baru
	claims := jwt.MapClaims{
		"sub":       userID,                               // User ID
		"branch_id": branchID,                             // Branch ID
		"user_role": userRole,                             // User Role
		"exp":       time.Now().Add(8 * time.Hour).Unix(), // Expired dalam 8 jam
	}

	// Buat token baru dengan klaim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Gunakan secret key untuk menandatangani token
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(secretKey)
}

// GetProfile menangani penampilan branch sesuai branch_id dari tokenJWT
func GetProfile(c *fiber.Ctx) error {
	branchID, _ := services.GetBranchID(c)
	userID, _ := services.GetUserID(c)
	userRole, _ := services.GetUserRole(c)
	var profilStruct models.ProfileStruct

	// Melakukan LEFT OUTER JOIN menggunakan GORM
	if err := config.DB.
		Table("user_branches").
		Select("user_branches.user_id AS user_id, users.name AS profile_name, user_branches.branch_id AS branch_id, branches.branch_name AS branch_name, branches.address, branches.phone, branches.email, branches.sia_id, branches.sia_name, branches.psa_id, branches.psa_name, branches.sipa, branches.sipa_name, branches.aping_id, branches.aping_name, branches.bank_name, branches.account_name, branches.account_number, branches.tax_percentage, branches.journal_method").
		Joins("LEFT JOIN users ON users.id = user_branches.user_id").
		Joins("LEFT JOIN branches ON branches.id = user_branches.branch_id").
		Where("user_branches.branch_id = ? AND user_branches.user_id = ?", branchID, userID).
		Scan(&profilStruct).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get userbranches failed", "Failed to fetch user branches with details")
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Otoritas : "+userRole, profilStruct)
}
