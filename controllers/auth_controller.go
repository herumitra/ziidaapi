package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/helpers"
	"github.com/herumitra/ziidaapi/models"
	"github.com/herumitra/ziidaapi/services"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login mengautentikasi pengguna dan menghasilkan JWT
func Login(c *fiber.Ctx) error {
	var users []models.User
	var req LoginRequest

	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	cek_status := config.DB.Where("username = ? AND user_status = 'active'", req.Username).First(&users)

	// Cek apakah user status aktif
	if err := cek_status.Error; err != nil {
		// Jika user status tidak aktif, kembalikan error
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "User not active, please call administrator", err)
	}

	// Verifikasi kredensial dan hasilkan token JWT
	_, resp, err := services.AuthenticateUser(config.DB, config.RDB, req.Username, req.Password)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Incorrect username or password", nil)
	}

	// Kembalikan token sebagai response
	return helpers.JSONResponse(c, fiber.StatusOK, "Login successful", resp.Data)
}

// Logout menangani proses logout pengguna dengan menghapus token di Redis
func Logout(c *fiber.Ctx) error {
	// Get token from header
	token := c.Get("Authorization")

	// Remove "Bearer " prefix
	if strings.HasPrefix(token, "Bearer ") {
		token = token[len("Bearer "):]
	}
	redisKey := fmt.Sprintf("auth:%s", token) // Menambahkan prefix "auth:" pada key Redis

	if err := config.RDB.Del(c.Context(), redisKey).Err(); err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to logout", nil)
	}
	return helpers.JSONResponse(c, fiber.StatusOK, "Logout successful", redisKey)
}

// Struktur untuk request body
type SetBranchRequest struct {
	BranchID int `json:"branch_id"`
}

func SetBranch(c *fiber.Ctx) error {
	// Membaca data JSON dari body request dan mengubahnya ke dalam struct SetBranchRequest
	var req SetBranchRequest
	if err := c.BodyParser(&req); err != nil {
		// Jika ada error dalam parsing, kembalikan response error
		return c.Status(fiber.StatusBadRequest).JSON(helpers.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	// Ambil branchID dari request body
	branchID := req.BranchID

	// Debug: Log data branchID yang diterima
	fmt.Println("Received branchID:", branchID)

	// Ambil data cabang dari database berdasarkan branchID
	var branch models.Branch
	if err := config.DB.First(&branch, branchID).Error; err != nil {
		// Jika gagal mengambil data cabang dari database, kembalikan error
		fmt.Println("Error fetching branch data from database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.Response{
			Status:  "error",
			Message: "Failed to fetch branch data from database",
		})
	}

	// Debug: Log data cabang yang ditemukan
	fmt.Println("Branch data found:", branch)

	// Ambil token dari header Authorization (Bearer token)
	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		authHeader = authHeader[len("Bearer "):] // Menghapus prefix "Bearer " dari token
	}

	// Jika token kosong, kembalikan response Unauthorized
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(helpers.Response{
			Status:  "error",
			Message: "Authorization token missing",
		})
	}

	// Debug: Log token yang diterima
	fmt.Println("Received token:", authHeader)

	// Memanggil fungsi JWTDecodeID untuk mendapatkan userID dari token
	userID := services.JWTDecodeID(authHeader)
	if userID == "" {
		// Debug: Log jika token tidak valid
		fmt.Println("Invalid token: userID not found")
		return c.Status(fiber.StatusUnauthorized).JSON(helpers.Response{
			Status:  "error",
			Message: "Invalid token",
		})
	}

	// Debug: Log userID yang berhasil didecode
	fmt.Println("Decoded userID:", userID)

	// Cek apakah user sudah terhubung dengan branch yang sama
	var userBranch models.UserBranch
	if err := config.DB.Where("user_id = ? AND branch_id = ?", userID, branchID).First(&userBranch).Error; err != nil {
		// Jika user tidak terhubung dengan branch tersebut, kembalikan response error
		fmt.Println("User not associated with this branch:", err)
		return c.Status(fiber.StatusForbidden).JSON(helpers.Response{
			Status:  "error",
			Message: "User not authorized for this branch",
		})
	}

	// Debug: Log jika user sudah terhubung dengan branch
	fmt.Println("User is authorized for this branch.")

	// Menyimpan data branch ke Redis dengan key profile:[branchID]
	branchRedisKey := fmt.Sprintf("profile:%d", branchID)
	branchData := map[string]interface{}{
		"branch_name": branch.BranchName,
		"address":     branch.Address,
		"sipa":        branch.Sipa,
		"sipa_name":   branch.SipaName,
		"sia_id":      branch.SiaId,
		"sia_name":    branch.SiaName,
	}

	// Debug: Log data yang akan disimpan ke Redis
	fmt.Println("Saving branch data to Redis:", branchData)

	err := config.RDB.HMSet(context.Background(), branchRedisKey, branchData).Err()
	if err != nil {
		// Jika gagal menyimpan data ke Redis, kembalikan response error
		fmt.Println("Error saving branch data to Redis:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.Response{
			Status:  "error",
			Message: "Failed to save branch data to Redis",
		})
	}

	// Menyimpan relasi token dengan profil branch di Redis (auth:[token] -> profile:[branchID])
	authRedisKey := fmt.Sprintf("auth:%s", authHeader)

	// Debug: Log relasi antara token dan branch profile yang disimpan
	fmt.Println("Associating token with branch profile in Redis:", authRedisKey, branchRedisKey)

	err = config.RDB.Set(context.Background(), authRedisKey, branchRedisKey, 8*time.Hour).Err()
	if err != nil {
		// Jika gagal menyimpan relasi ke Redis, kembalikan response error
		fmt.Println("Error associating token with branch profile in Redis:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.Response{
			Status:  "error",
			Message: "Failed to associate token with branch profile in Redis",
		})
	}

	// Mengembalikan response sukses setelah data disimpan
	return c.JSON(helpers.Response{
		Status:  "success",
		Message: "Branch set successfully and saved to Redis",
		Data:    branchID,
	})
}

func GetBranchFromToken(c *fiber.Ctx) error {
	// Ambil token dari header Authorization (Bearer token)
	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		authHeader = authHeader[len("Bearer "):] // Menghapus prefix "Bearer " dari token
	}

	// Jika token kosong, kembalikan response Unauthorized
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(helpers.Response{
			Status:  "error",
			Message: "Authorization token missing",
		})
	}

	// Debug: Log token yang diterima
	fmt.Println("Received token:", authHeader)

	// Memanggil fungsi JWTDecodeID untuk mendapatkan userID dari token
	userID := services.JWTDecodeID(authHeader)
	if userID == "" {
		// Debug: Log jika token tidak valid
		fmt.Println("Invalid token: userID not found")
		return c.Status(fiber.StatusUnauthorized).JSON(helpers.Response{
			Status:  "error",
			Message: "Invalid token",
		})
	}

	// Debug: Log userID yang berhasil didecode
	fmt.Println("Decoded userID:", userID)

	// Ambil key Redis untuk auth:[token]
	authRedisKey := fmt.Sprintf("auth:%s", authHeader)

	// Debug: Log key Redis auth:[token]
	fmt.Println("Checking Redis for key:", authRedisKey)

	// Periksa apakah key auth:[token] ada di Redis
	exists, err := config.RDB.Exists(context.Background(), authRedisKey).Result()
	if err != nil {
		// Debug: Log error saat memeriksa eksistensi key
		fmt.Println("Error checking token existence in Redis:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.Response{
			Status:  "error",
			Message: "Failed to check token existence in Redis",
		})
	}

	if exists == 0 {
		// Debug: Log jika key token tidak ada di Redis
		fmt.Println("Token not found in Redis")
		return c.Status(fiber.StatusUnauthorized).JSON(helpers.Response{
			Status:  "error",
			Message: "Token expired or not found",
		})
	}

	// Debug: Log jika key token ditemukan
	fmt.Println("Token found in Redis, fetching related branch data...")

	// Ambil key profile:[branchID] yang terhubung dengan auth:[token]
	branchRedisKey, err := config.RDB.Get(context.Background(), authRedisKey).Result()
	if err != nil {
		// Debug: Log error saat mengambil data profil cabang dari Redis
		fmt.Println("Error retrieving branch key from Redis:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.Response{
			Status:  "error",
			Message: "Failed to retrieve branch key from Redis",
		})
	}

	// Debug: Log key profile yang terhubung dengan token
	fmt.Println("Branch key from Redis:", branchRedisKey)

	// Ambil data profil cabang dari Redis menggunakan key profile:[branchID]
	branchData, err := config.RDB.HGetAll(context.Background(), branchRedisKey).Result()
	if err != nil {
		// Debug: Log error saat mengambil data profil cabang
		fmt.Println("Error retrieving branch data from Redis:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.Response{
			Status:  "error",
			Message: "Failed to retrieve branch data from Redis",
		})
	}

	// Debug: Log data profil cabang yang ditemukan
	fmt.Println("Branch data retrieved from Redis:", branchData)

	// Kembalikan data cabang sebagai response
	return c.JSON(helpers.Response{
		Status:  "success",
		Message: "Branch data retrieved successfully",
		Data:    branchData,
	})
}
