package controllers

import (
	fmt "fmt"
	strconv "strconv"
	time "time"

	fiber "github.com/gofiber/fiber/v2"
	config "github.com/herumitra/ziidaapi/config"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
	gorm "gorm.io/gorm"
)

// CreateUser menangani pembuatan pengguna baru
func CreateUser(c *fiber.Ctx) error {
	// Buat instance baru untuk User
	var user models.User

	// Parse input JSON menjadi struct User
	if err := c.BodyParser(&user); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Generate ID untuk user
	userID, err := generateUserID(config.DB) // Pastikan generateUserID menerima DB dan mengembalikan ID yang valid
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Create user failed", "Failed to generate user ID")
	}

	// Set ID user yang sudah digenerate
	user.ID = userID

	// Cek panjang password
	if len(user.Password) < 8 {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Create user failed", "Password must be at least 8 characters long")
	}

	// Hash password sebelum disimpan
	if err := user.HashPassword(); err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Create user failed", "Failed to hash password")
	}

	// Simpan user ke database
	if err := config.DB.Create(&user).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to create user", err)
	}

	// Statemen terakhir jika tidak ditemukan error
	return helpers.JSONResponse(c, fiber.StatusCreated, "User created successfully", user)
}

// UpdateUser menangani pembaruan pengguna
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Cari user berdasarkan ID
	if err := config.DB.First(&user, id).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "User not found", err)
	}

	// Parsing request body untuk update data
	if err := c.BodyParser(&user); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Update data user
	if err := config.DB.Save(&user).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update user", err)
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "User updated successfully", user)
}

// DeleteUser menangani penghapusan pengguna
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "User not found", err)
	}

	// Hapus user
	if err := config.DB.Where("id = ?", id).Delete(&user).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete user", err)
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "User deleted successfully", id)
}

// GetUser menangani penampilan pengguna
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "User not found", err)
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "User found", user)
}

// GetAllUsers menangani penampilan semua pengguna
func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	// Mengambil semua data user dari database
	if err := config.DB.Find(&users).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get users failed", "Failed to fetch users")
	}

	// Mengembalikan response sukses dengan data pengguna
	return helpers.JSONResponse(c, fiber.StatusOK, "Users retrieved successfully", users)
}

// Fungsi untuk generate ID user
func generateUserID(db *gorm.DB) (string, error) {
	// Ambil tanggal saat ini dalam format DDMMYYYY
	now := time.Now()
	dateStr := now.Format("02012006") // Format DDMMYYYY

	var user models.User // Menggunakan model yang sudah ada

	// Ambil urutan terbesar untuk tanggal tersebut
	if err := db.Where("id LIKE ?", "USR"+dateStr+"%").Order("id DESC").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Jika tidak ada user sebelumnya, urutan awal adalah 1
			return fmt.Sprintf("USR%s001", dateStr), nil
		} else {
			return "", fmt.Errorf("error querying database: %v", err)
		}
	}

	// Jika ditemukan user sebelumnya, ambil urutan terakhir dan tambah 1
	lastID := user.ID                // Ambil ID user.ID
	seqStr := lastID[len(lastID)-3:] // Ambil 3 digit terakhir dari ID sebelumnya
	seq, err := strconv.Atoi(seqStr)
	if err != nil {
		return "", fmt.Errorf("error converting sequence: %v", err)
	}
	sequence := seq + 1

	// Format ID baru dengan urutan 3 digit
	sequenceStr := fmt.Sprintf("%03d", sequence)
	userID := fmt.Sprintf("USR%s%s", dateStr, sequenceStr)

	return userID, nil
}
