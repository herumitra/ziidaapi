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

// GetAllUnitConversion tampilkan semua unit_conversion
func GetAllUnitConversion(c *fiber.Ctx) error {
	// Get branch id
	branch_id, _ := services.GetBranchID(c)
	var unit_conversions []models.UserBranchDetail

	// Melakukan LEFT OUTER JOIN menggunakan GORM
	if err := config.DB.
		Table("unit_conversions").
		Select("unit_conversions.id, unit_conversions.unit_init_id, uinit.name AS unit_init_name ,unit_conversions.unit_final_id, ufinal.name AS unit_final_name ,unit_conversions.value_conv, unit_conversions.branch_id").
		Joins("LEFT JOIN units uinit ON uinit.id = unit_conversions.unit_init_id").
		Joins("LEFT JOIN units ufinal ON ufinal.id = unit_conversions.unit_final_id").
		Where("unit_conversions.branch_id = ?", branch_id).
		Scan(&unit_conversions).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get userbranches failed", "Failed to fetch user branches with details")
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "UnitConversions retrieved successfully", unit_conversions)
}

// GetUnitConversion tampilkan unit_conversion berdasarkan id
func GetUnitConversion(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id := c.Params("id")
	var unit_conversion models.UnitConversion

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&unit_conversion).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "UnitConversion not found", err)
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "UnitConversion found", unit_conversion)
}

// CreateUnitConversion buat unit_conversion
func CreateUnitConversion(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Buat instance baru untuk UnitConversion
	var unit_conversion models.UnitConversion

	// Parse body request ke struct unit_conversion
	if err := c.BodyParser(&unit_conversion); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	unit_conversionID, err := generateUnitConversionID(config.DB) // Pastikan generateUnitConversionID menerima DB dan mengembalikan ID yang valid
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Create unit_conversion failed", "Failed to generate unit_conversion ID")
	}

	unit_conversion.ID = unit_conversionID
	// Tambahkan branchID ke unit_conversion (relasi ke cabang)
	unit_conversion.BranchID = branchID

	// Simpan unit_conversion ke database menggunakan GORM
	if err := config.DB.Create(&unit_conversion).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create unit_conversion",
		})
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "UnitConversion created successfully", unit_conversion)
}

// UpdateUnitConversion update unit_conversion
func UpdateUnitConversion(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil unit_conversionID dari parameter
	id := c.Params("id")
	var unit_conversion models.UnitConversion

	// Cari unit_conversion berdasarkan ID dan branchID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&unit_conversion).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "UnitConversion not found", err)
	}

	// Parse body request ke struct unit_conversion
	if err := c.BodyParser(&unit_conversion); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Simpan unit_conversion ke database menggunakan GORM
	if err := config.DB.Model(&unit_conversion).Updates(unit_conversion).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update unit_conversion", err)
	}

	// Mengembalikan response data unit_conversion yang diupdate
	return helpers.JSONResponse(c, fiber.StatusOK, "UnitConversion created successfully", unit_conversion)

}

// DeleteUnitConversion hapus unit_conversion
func DeleteUnitConversion(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil unit_conversionID dari parameter
	id := c.Params("id")
	var unit_conversion models.UnitConversion

	// Ambil data unit_conversion
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&unit_conversion).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete unit_conversion", err)
	}

	// Hapus unit_conversion
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).Delete(&unit_conversion).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete unit_conversion", err)
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "UnitConversion deleted successfully", unit_conversion)
}

// generateUnitConversionID generate id unit_conversion
func generateUnitConversionID(db *gorm.DB) (string, error) {
	// Ambil tanggal saat ini dalam format DDMMYYYY
	now := time.Now()
	dateStr := now.Format("02012006") // Format DDMMYYYY

	var unit_conversion models.UnitConversion // Menggunakan model yang sudah ada

	// Ambil urutan terbesar untuk tanggal tersebut
	if err := db.Where("id LIKE ?", "UNC"+dateStr+"%").Order("id DESC").First(&unit_conversion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Jika tidak ada user sebelumnya, urutan awal adalah 1
			return fmt.Sprintf("UNC%s001", dateStr), nil
		} else {
			return "", fmt.Errorf("error querying database: %v", err)
		}
	}

	// Jika ditemukan unit_conversion sebelumnya, ambil urutan terakhir dan tambah 1
	lastID := unit_conversion.ID     // Ambil ID unit_conversion.ID
	seqStr := lastID[len(lastID)-3:] // Ambil 3 digit terakhir dari ID sebelumnya
	seq, err := strconv.Atoi(seqStr)
	if err != nil {
		return "", fmt.Errorf("error converting sequence: %v", err)
	}
	sequence := seq + 1

	// Format ID baru dengan urutan 3 digit
	sequenceStr := fmt.Sprintf("%03d", sequence)
	unit_conversionID := fmt.Sprintf("UNC%s%s", dateStr, sequenceStr)

	return unit_conversionID, nil
}
