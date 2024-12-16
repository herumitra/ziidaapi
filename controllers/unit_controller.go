package controllers

import (
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
	return c.SendString("Hello, World!")
}

// UpdateUnit update unit
func UpdateUnit(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// DeleteUnit hapus unit
func DeleteUnit(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// generateUnitID generate id unit
func generateUnitID(db *gorm.DB) (string, error) {
	return "", nil
}
