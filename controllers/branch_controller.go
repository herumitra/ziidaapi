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

// CreateBranch menangani penambahan branch
func CreateBranch(c *fiber.Ctx) error {
	// Buat instance baru untuk Branch
	var branch models.Branch

	// Parse input JSON menjadi struct Branch
	if err := c.BodyParser(&branch); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Generate ID untuk branch
	branchID, err := generateBranchID(config.DB) // Pastikan generateBranchID menerima DB dan mengembalikan ID yang valid
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Create branch failed", "Failed to generate branch ID")
	}

	// Set ID branch yang sudah digenerate
	branch.ID = branchID

	// Simpan user ke database
	if err := config.DB.Create(&branch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to create user", err)
	}
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch created successfully", branch)
}

// GetBranch menangani penampilan branch
func GetBranch(c *fiber.Ctx) error {
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch found", nil)
}

// UpdateBranch menangani pembaruan branch
func UpdateBranch(c *fiber.Ctx) error {
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch updated successfully", nil)
}

// DeleteBranch menangani penghapusan branch
func DeleteBranch(c *fiber.Ctx) error {
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch deleted successfully", nil)
}

// GetAllBranch menangani penampilan semua branch
func GetAllBranch(c *fiber.Ctx) error {
	return helpers.JSONResponse(c, fiber.StatusOK, "Branches retrieved successfully", nil)
}

// Fungsi untuk generate ID user
func generateBranchID(db *gorm.DB) (string, error) {
	// Ambil tanggal saat ini dalam format DDMMYYYY
	now := time.Now()
	dateStr := now.Format("02012006") // Format DDMMYYYY

	var branch models.Branch // Menggunakan model yang sudah ada

	// Ambil urutan terbesar untuk tanggal tersebut
	if err := db.Where("id LIKE ?", "BRC"+dateStr+"%").Order("id DESC").First(&branch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Jika tidak ada user sebelumnya, urutan awal adalah 1
			return fmt.Sprintf("BRC%s001", dateStr), nil
		} else {
			return "", fmt.Errorf("error querying database: %v", err)
		}
	}

	// Jika ditemukan user sebelumnya, ambil urutan terakhir dan tambah 1
	lastID := branch.ID              // Ambil ID branch.ID
	seqStr := lastID[len(lastID)-3:] // Ambil 3 digit terakhir dari ID sebelumnya
	seq, err := strconv.Atoi(seqStr)
	if err != nil {
		return "", fmt.Errorf("error converting sequence: %v", err)
	}
	sequence := seq + 1

	// Format ID baru dengan urutan 3 digit
	sequenceStr := fmt.Sprintf("%03d", sequence)
	branchID := fmt.Sprintf("BRC%s%s", dateStr, sequenceStr)

	return branchID, nil
}
