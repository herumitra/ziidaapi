package controllers

import (
	fiber "github.com/gofiber/fiber/v2"
	config "github.com/herumitra/ziidaapi/config"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
	"github.com/herumitra/ziidaapi/services"
)

// CreateUserBranch menangani penambahan userbranch
func CreateUserBranch(c *fiber.Ctx) error {
	// Buat instance baru untuk UserBranch
	var userbranch models.UserBranch

	// Parse input JSON menjadi struct UserBranch
	if err := c.BodyParser(&userbranch); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Simpan user ke database
	if err := config.DB.Create(&userbranch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to create user", err)
	}
	return helpers.JSONResponse(c, fiber.StatusOK, "UserBranch created successfully", userbranch)
}

// GetUserBranch menangani penampilan userbranch
func GetUserBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	var userbranch models.UserBranch

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ?", id).First(&userbranch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "UserBranch not found", err)
	}

	// Mengembalikan response data userbranch
	return helpers.JSONResponse(c, fiber.StatusOK, "UserBranch found", userbranch)
}

// UpdateUserBranch menangani pembaruan userbranch
func UpdateUserBranch(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	branch_id := c.Params("branch_id")

	var userbranch models.UserBranch

	// Cari userbranch berdasarkan ID
	if err := config.DB.Where("user_id	= ? AND branch_id = ?", user_id, branch_id).First(&userbranch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "UserBranch not found", err)
	}

	// Parsing data body langsung ke struct `userbranch`
	// Namun, ini hanya akan mengupdate field-field tertentu.
	if err := c.BodyParser(&userbranch); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Pastikan hanya field yang ingin diperbarui yang diubah.
	// Gunakan `Model` untuk menghindari overwrite seluruh object.
	if err := config.DB.Model(&userbranch).Where("user_id	= ? AND branch_id = ?", user_id, branch_id).Updates(userbranch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update userbranch", err)
	}

	// Mengembalikan response sukses dengan data userbranch yang diperbarui
	return helpers.JSONResponse(c, fiber.StatusOK, "UserBranch updated successfully", userbranch)
}

// DeleteUserBranch menangani penghapusan userbranch
func DeleteUserBranch(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	branch_id := c.Params("branch_id")
	var userbranch models.UserBranch

	// Cari userbranch berdasarkan ID
	if err := config.DB.Where("user_id	= ? AND branch_id = ?", user_id, branch_id).First(&userbranch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "userbranch not found", err)
	}

	// Hapus userbranch
	if err := config.DB.Where("user_id	= ? AND branch_id = ?", user_id, branch_id).Delete(&userbranch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete userbranch", err)
	}

	// Mengembalikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "UserBranch deleted successfully", userbranch)
}

// GetAllUserBranch menangani penampilan semua userbranch
func GetAllUserBranch(c *fiber.Ctx) error {
	branch_id, _ := services.GetBranchID(c)
	var userBranchDetails []models.UserBranchDetail

	// Melakukan LEFT OUTER JOIN menggunakan GORM
	if err := config.DB.
		Table("user_branches").
		Select("user_branches.user_id, users.name AS user_name, user_branches.branch_id, branches.name AS branch_name").
		Joins("LEFT JOIN users ON users.id = user_branches.user_id").
		Joins("LEFT JOIN branches ON branches.id = user_branches.branch_id").
		Where("user_branches.branch_id = ?", branch_id).
		Scan(&userBranchDetails).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get userbranches failed", "Failed to fetch user branches with details")
	}

	// Mengembalikan response dengan data hasil JOIN
	return helpers.JSONResponse(c, fiber.StatusOK, "UserBranches retrieved successfully", userBranchDetails)
}
