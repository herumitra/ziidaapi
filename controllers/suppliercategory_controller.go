package controllers

import (
	fiber "github.com/gofiber/fiber/v2"
	config "github.com/herumitra/ziidaapi/config"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
	services "github.com/herumitra/ziidaapi/services"
)

// GetAllSupplierCategory tampilkan semua supplier_category
func GetAllSupplierCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Inisialize variabel supplier_categories
	var supplier_categories []models.SupplierCategory

	// Mengambil semua data branch dari database dengan filter branch_id
	if err := config.DB.Where("branch_id = ?", branchID).Find(&supplier_categories).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get supplier_categories failed", "Failed to fetch supplier_categories")
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "SupplierCategories retrieved successfully", supplier_categories)
}

// GetSupplierCategory tampilkan supplier_category berdasarkan id
func GetSupplierCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id := c.Params("id")
	var supplier_category models.SupplierCategory

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&supplier_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "SupplierCategory not found", err)
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "SupplierCategory found", supplier_category)
}

// CreateSupplierCategory buat supplier_category
func CreateSupplierCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Buat instance baru untuk SupplierCategory
	var supplier_category models.SupplierCategory

	// Parse body request ke struct supplier_category
	if err := c.BodyParser(&supplier_category); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Tambahkan branchID ke supplier_category (relasi ke cabang)
	supplier_category.BranchID = branchID

	// Simpan supplier_category ke database menggunakan GORM
	if err := config.DB.Create(&supplier_category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create supplier_category",
		})
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "SupplierCategory created successfully", supplier_category)
}

// UpdateSupplierCategory update supplier_category
func UpdateSupplierCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil supplier_categoryID dari parameter
	id := c.Params("id")
	var supplier_category models.SupplierCategory

	// Cari supplier_category berdasarkan ID dan branchID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&supplier_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "SupplierCategory not found", err)
	}

	// Parse body request ke struct supplier_category
	if err := c.BodyParser(&supplier_category); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Simpan supplier_category ke database menggunakan GORM
	if err := config.DB.Model(&supplier_category).Updates(supplier_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update supplier_category", err)
	}

	// Mengembalikan response data supplier_category yang diupdate
	return helpers.JSONResponse(c, fiber.StatusOK, "SupplierCategory created successfully", supplier_category)

}

// DeleteSupplierCategory hapus supplier_category
func DeleteSupplierCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil supplier_categoryID dari parameter
	id := c.Params("id")
	var supplier_category models.SupplierCategory

	// Ambil data supplier_category
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&supplier_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete supplier_category", err)
	}

	// Hapus supplier_category
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).Delete(&supplier_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete supplier_category", err)
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "SupplierCategory deleted successfully", supplier_category)
}
