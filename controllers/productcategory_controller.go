package controllers

import (
	fiber "github.com/gofiber/fiber/v2"
	config "github.com/herumitra/ziidaapi/config"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
	services "github.com/herumitra/ziidaapi/services"
)

// GetAllProductCategory tampilkan semua product_category
func GetAllProductCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Inisialize variabel product_categories
	var product_categories []models.ProductCategory

	// Mengambil semua data branch dari database dengan filter branch_id
	if err := config.DB.Where("branch_id = ?", branchID).Find(&product_categories).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get product_categories failed", "Failed to fetch product_categories")
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "ProductCategories retrieved successfully", product_categories)
}

// GetProductCategory tampilkan product_category berdasarkan id
func GetProductCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id := c.Params("id")
	var product_category models.ProductCategory

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&product_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "ProductCategory not found", err)
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "ProductCategory found", product_category)
}

// CreateProductCategory buat product_category
func CreateProductCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Buat instance baru untuk ProductCategory
	var product_category models.ProductCategory

	// Parse body request ke struct product_category
	if err := c.BodyParser(&product_category); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Tambahkan branchID ke product_category (relasi ke cabang)
	product_category.BranchID = branchID

	// Simpan product_category ke database menggunakan GORM
	if err := config.DB.Create(&product_category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product_category",
		})
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "ProductCategory created successfully", product_category)
}

// UpdateProductCategory update product_category
func UpdateProductCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil product_categoryID dari parameter
	id := c.Params("id")
	var product_category models.ProductCategory

	// Cari product_category berdasarkan ID dan branchID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&product_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "ProductCategory not found", err)
	}

	// Parse body request ke struct product_category
	if err := c.BodyParser(&product_category); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Simpan product_category ke database menggunakan GORM
	if err := config.DB.Model(&product_category).Updates(product_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update product_category", err)
	}

	// Mengembalikan response data product_category yang diupdate
	return helpers.JSONResponse(c, fiber.StatusOK, "ProductCategory created successfully", product_category)

}

// DeleteProductCategory hapus product_category
func DeleteProductCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil product_categoryID dari parameter
	id := c.Params("id")
	var product_category models.ProductCategory

	// Ambil data product_category
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&product_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete product_category", err)
	}

	// Hapus product_category
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).Delete(&product_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete product_category", err)
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "ProductCategory deleted successfully", product_category)
}
