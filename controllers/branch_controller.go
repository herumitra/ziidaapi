package controllers

import (
	fiber "github.com/gofiber/fiber/v2"
	helpers "github.com/herumitra/ziidaapi/helpers"
)

// CreateBranch menangani penambahan branch
func CreateBranch(c *fiber.Ctx) error {
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch created successfully", nil)
}

// GetBranch menangani penampilan branch
func GetBranch(c *fiber.Ctx) error {
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch found", nil)
}

// UpdateBranch menangani pembaruan branch
func UpdateBranch(c *fiber.Ctx) error {
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch updated successfully", nil)
}

// DeleteBranch menangani penghapusan branch
func DeleteBranch(c *fiber.Ctx) error {
	return helpers.JSONResponse(c, fiber.StatusOK, "Branch deleted successfully", nil)
}

// GetAllBranch menangani penampilan semua branch
func GetAllBranch(c *fiber.Ctx) error {
	return helpers.JSONResponse(c, fiber.StatusOK, "Branches retrieved successfully", nil)
}
