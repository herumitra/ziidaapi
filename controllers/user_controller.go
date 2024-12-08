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

var user models.User

// Fungsi untuk generate ID user
func generateUserID(db *gorm.DB) (string, error) {
	// Ambil tanggal saat ini dalam format YYYYMMDD
	now := time.Now()
	dateStr := now.Format("20061202") // Format YYYYMMDD

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

// CreateUser menangani pembuatan pengguna baru
func CreateUser(c *fiber.Ctx) error {
	var user models.User

	// Parse input JSON menjadi struct User
	if err := c.BodyParser(&user); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	// Generate ID untuk user
	userID, err := generateUserID(config.DB) // Pastikan generateUserID menerima DB dan mengembalikan ID yang valid
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to generate user ID", nil)
	}

	// Set ID user yang sudah digenerate
	user.ID = userID

	// Hash password sebelum disimpan
	if err := user.HashPassword(); err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to hash password", nil)
	}

	// Simpan user ke database
	if err := config.DB.Create(&user).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to create user", nil)
	}

	return helpers.JSONResponse(c, fiber.StatusCreated, "User created successfully", user)
}
