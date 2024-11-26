package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/helpers"
	"github.com/herumitra/ziidaapi/models"
)

// GetAllUsers mengembalikan daftar semua pengguna
func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	// Mengambil semua data user dari database
	if err := config.DB.Find(&users).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to fetch users", nil)
	}

	// Mengembalikan response sukses dengan data pengguna
	return helpers.JSONResponse(c, fiber.StatusOK, "Users retrieved successfully", users)
}

// CreateUser menangani pembuatan pengguna baru
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	// Hash password sebelum disimpan
	if err := user.HashPassword(); err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to hash password", nil)
	}

	// Simpan user ke database
	if err := config.DB.Create(&user).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to create user", nil)
	}

	return helpers.JSONResponse(c, fiber.StatusCreated, "User created successfully", user)
}

// GetUser mengembalikan data pengguna berdasarkan ID
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Cari user berdasarkan ID
	if err := config.DB.First(&user, id).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "User found", user)
}

// UpdateUser memperbarui data pengguna
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Cari user berdasarkan ID
	if err := config.DB.First(&user, id).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	// Parsing request body untuk update data
	if err := c.BodyParser(&user); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	// Update data user
	if err := config.DB.Save(&user).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update user", nil)
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "User updated successfully", user)
}

// DeleteUser menghapus pengguna berdasarkan ID
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Cari user berdasarkan ID
	if err := config.DB.First(&user, id).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	// Hapus user
	if err := config.DB.Delete(&user).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete user", nil)
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "User deleted successfully", nil)
}
