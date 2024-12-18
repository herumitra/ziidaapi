package controllers

import (
	fiber "github.com/gofiber/fiber/v2"
	config "github.com/herumitra/ziidaapi/config"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
	services "github.com/herumitra/ziidaapi/services"
)

// GetAllMemberCategory tampilkan semua member_category
func GetAllMemberCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Inisialize variabel member_categories
	var member_categories []models.MemberCategory

	// Mengambil semua data branch dari database dengan filter branch_id
	if err := config.DB.Where("branch_id = ?", branchID).Find(&member_categories).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get member_categories failed", "Failed to fetch member_categories")
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "MemberCategories retrieved successfully", member_categories)
}

// GetMemberCategory tampilkan member_category berdasarkan id
func GetMemberCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id := c.Params("id")
	var member_category models.MemberCategory

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&member_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "MemberCategory not found", err)
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "MemberCategory found", member_category)
}

// CreateMemberCategory buat member_category
func CreateMemberCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Buat instance baru untuk MemberCategory
	var member_category models.MemberCategory

	// Parse body request ke struct member_category
	if err := c.BodyParser(&member_category); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Tambahkan branchID ke member_category (relasi ke cabang)
	member_category.BranchID = branchID

	// Simpan member_category ke database menggunakan GORM
	if err := config.DB.Create(&member_category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create member_category",
		})
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "MemberCategory created successfully", member_category)
}

// UpdateMemberCategory update member_category
func UpdateMemberCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil member_categoryID dari parameter
	id := c.Params("id")
	var member_category models.MemberCategory

	// Cari member_category berdasarkan ID dan branchID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&member_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "MemberCategory not found", err)
	}

	// Parse body request ke struct member_category
	if err := c.BodyParser(&member_category); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Simpan member_category ke database menggunakan GORM
	if err := config.DB.Model(&member_category).Updates(member_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update member_category", err)
	}

	// Mengembalikan response data member_category yang diupdate
	return helpers.JSONResponse(c, fiber.StatusOK, "MemberCategory created successfully", member_category)

}

// DeleteMemberCategory hapus member_category
func DeleteMemberCategory(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil member_categoryID dari parameter
	id := c.Params("id")
	var member_category models.MemberCategory

	// Ambil data member_category
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&member_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete member_category", err)
	}

	// Hapus member_category
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).Delete(&member_category).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete member_category", err)
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "MemberCategory deleted successfully", member_category)
}
