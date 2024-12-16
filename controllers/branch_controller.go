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
	id := c.Params("id")
	var branch models.Branch

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ?", id).First(&branch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Branch not found", err)
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch found", branch)
}

// UpdateBranch menangani pembaruan branch
func UpdateBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	var branch models.Branch

	// Cari branch berdasarkan ID
	if err := config.DB.Where("id = ?", id).First(&branch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Branch not found", err)
	}

	// Parsing data body langsung ke struct `branch`
	// Namun, ini hanya akan mengupdate field-field tertentu.
	if err := c.BodyParser(&branch); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Pastikan hanya field yang ingin diperbarui yang diubah.
	// Gunakan `Model` untuk menghindari overwrite seluruh object.
	if err := config.DB.Model(&branch).Updates(branch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update branch", err)
	}

	// Mengembalikan response sukses dengan data branch yang diperbarui
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch updated successfully", branch)
}

// DeleteBranch menangani penghapusan branch
func DeleteBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	var branch models.Branch

	// Cari branch berdasarkan ID
	if err := config.DB.Where("id = ?", id).First(&branch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "branch not found", err)
	}

	// Hapus branch
	if err := config.DB.Where("id = ?", id).Delete(&branch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete branch", err)
	}

	// Mengembalikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch deleted successfully", branch)
}

// GetAllBranch menangani penampilan semua branch
func GetAllBranch(c *fiber.Ctx) error {
	var branches []models.Branch

	// Mengambil semua data branch dari database
	if err := config.DB.Find(&branches).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get branches failed", "Failed to fetch branches")
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Branches retrieved successfully", branches)
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

	// Jika ditemukan branch sebelumnya, ambil urutan terakhir dan tambah 1
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
