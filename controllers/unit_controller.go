package controllers

import (
	fmt "fmt"
	strconv "strconv"
	time "time"

	fiber "github.com/gofiber/fiber/v2"
	config "github.com/herumitra/ziidaapi/config"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
	services "github.com/herumitra/ziidaapi/services"
	gorm "gorm.io/gorm"
)

// GetAllUnit tampilkan semua unit
func GetAllUnit(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Inisialize variabel units
	var units []models.Unit

	// Mengambil semua data branch dari database dengan filter branch_id
	if err := config.DB.Where("branch_id = ?", branchID).Find(&units).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get units failed", "Failed to fetch units")
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Units retrieved successfully", units)
}

// GetUnit tampilkan unit berdasarkan id
func GetUnit(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id := c.Params("id")
	var unit models.Unit

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&unit).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Unit not found", err)
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Unit found", unit)
}

// CreateUnit buat unit
func CreateUnit(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Buat instance baru untuk Unit
	var unit models.Unit

	// Parse body request ke struct unit
	if err := c.BodyParser(&unit); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	unitID, err := generateUnitID(config.DB) // Pastikan generateUnitID menerima DB dan mengembalikan ID yang valid
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Create unit failed", "Failed to generate unit ID")
	}

	unit.ID = unitID
	// Tambahkan branchID ke unit (relasi ke cabang)
	unit.BranchID = branchID

	// Simpan unit ke database menggunakan GORM
	if err := config.DB.Create(&unit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create unit",
		})
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "Unit created successfully", unit)
}

// UpdateUnit update unit
func UpdateUnit(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil unitID dari parameter
	id := c.Params("id")
	var unit models.Unit

	// Cari unit berdasarkan ID dan branchID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&unit).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Unit not found", err)
	}

	// Parse body request ke struct unit
	if err := c.BodyParser(&unit); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Simpan unit ke database menggunakan GORM
	if err := config.DB.Model(&unit).Updates(unit).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update unit", err)
	}

	// Mengembalikan response data unit yang diupdate
	return helpers.JSONResponse(c, fiber.StatusOK, "Unit created successfully", unit)

}

// DeleteUnit hapus unit
func DeleteUnit(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil unitID dari parameter
	id := c.Params("id")
	var unit models.Unit

	// Ambil data unit
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&unit).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete unit", err)
	}

	// Hapus unit
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).Delete(&unit).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete unit", err)
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "Unit deleted successfully", unit)
}

// generateUnitID generate id unit
func generateUnitID(db *gorm.DB) (string, error) {
	// Ambil tanggal saat ini dalam format DDMMYYYY
	now := time.Now()
	dateStr := now.Format("02012006") // Format DDMMYYYY

	var unit models.Unit // Menggunakan model yang sudah ada

	// Ambil urutan terbesar untuk tanggal tersebut
	if err := db.Where("id LIKE ?", "UNT"+dateStr+"%").Order("id DESC").First(&unit).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Jika tidak ada user sebelumnya, urutan awal adalah 1
			return fmt.Sprintf("UNT%s001", dateStr), nil
		} else {
			return "", fmt.Errorf("error querying database: %v", err)
		}
	}

	// Jika ditemukan unit sebelumnya, ambil urutan terakhir dan tambah 1
	lastID := unit.ID                // Ambil ID unit.ID
	seqStr := lastID[len(lastID)-3:] // Ambil 3 digit terakhir dari ID sebelumnya
	seq, err := strconv.Atoi(seqStr)
	if err != nil {
		return "", fmt.Errorf("error converting sequence: %v", err)
	}
	sequence := seq + 1

	// Format ID baru dengan urutan 3 digit
	sequenceStr := fmt.Sprintf("%03d", sequence)
	unitID := fmt.Sprintf("UNT%s%s", dateStr, sequenceStr)

	return unitID, nil
}
