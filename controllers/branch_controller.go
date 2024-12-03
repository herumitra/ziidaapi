package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/helpers"
	"github.com/herumitra/ziidaapi/models"
	"github.com/herumitra/ziidaapi/services"
)

// GetAllBranchs mengembalikan daftar semua branch
func GetAllBranch(c *fiber.Ctx) error {
	var branch []models.Branch

	// Mengambil semua data branch dari database
	if err := config.DB.Find(&branch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to fetch branches", nil)
	}

	// Mengembalikan response sukses dengan data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch retrieved successfully", branch)
}

// CreateBranch menangani pembuatan branch baru
func CreateBranch(c *fiber.Ctx) error {
	var branch models.Branch
	if err := c.BodyParser(&branch); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	// Simpan branch ke database
	if err := config.DB.Create(&branch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to create branch", nil)
	}

	return helpers.JSONResponse(c, fiber.StatusCreated, "Branch created successfully", branch)
}

// GetBranch mengembalikan data branch berdasarkan ID
func GetBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	var branch models.Branch

	// Cari branch berdasarkan ID
	if err := config.DB.First(&branch, id).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Branch not found", nil)
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "Branch found", branch)
}

// UpdateBranch memperbarui data branch
func UpdateBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	var branch models.Branch

	// Cari branch berdasarkan ID
	if err := config.DB.First(&branch, id).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Branch not found", nil)
	}

	// Parsing request body untuk update data
	if err := c.BodyParser(&branch); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	// Update data branch
	if err := config.DB.Save(&branch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update branch", "Nama Branch yang ingin anda gunakan sudah ada.")
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "Branch updated successfully", branch)
}

// DeleteBranch menghapus branch berdasarkan ID
func DeleteBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	var branch models.Branch

	// Cari branch berdasarkan ID
	if err := config.DB.First(&branch, id).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Branch not found", nil)
	}

	// Hapus branch
	if err := config.DB.Delete(&branch).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete branch", nil)
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "Branch deleted successfully", nil)
}

func GetBranchByUserIdFromToken(c *fiber.Ctx) error {
	// Ambil token dari header Authorization (Bearer token)
	authHeader := c.Get("Authorization")

	// Memanggil fungsi JWTDecodeID dari package services
	userID := services.JWTDecodeID(authHeader)

	// Mendeklarasikan variabel untuk hasil query
	var result []struct {
		BranchID   int64  `json:"branch_id"`
		BranchName string `json:"branch_name"`
		Address    string `json:"address"`
		Sipa       string `json:"sipa"`
		SipaName   string `json:"sipa_name"`
		SiaID      string `json:"sia_id"`
		SiaName    string `json:"sia_name"`
	}

	// Menjalankan query dengan Join dan select yang tepat
	if err := config.DB.Table("user_branches").
		Select("user_branches.branch_id AS branch_id, branches.branch_name, branches.address, branches.sipa, branches.sipa_name, branches.sia_id, branches.sia_name").
		Joins("LEFT OUTER JOIN branches ON branches.ID = user_branches.branch_id").
		Where("user_branches.user_id = ?", userID).
		Find(&result).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to fetch branches", nil)
	}

	// Mengembalikan data yang ditemukan
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch retrieved successfully", result)

}
