package services

import (
	fiber "github.com/gofiber/fiber/v2"
	helpers "github.com/herumitra/ziidaapi/helpers"
)

func GetUserID(c *fiber.Ctx) (string, error) {
	// Ambil user_id dari claims
	userID, _ := helpers.GetClaimsToken(c, "sub")

	return userID, nil
}

func GetBranchID(c *fiber.Ctx) (string, error) {
	// Ambil branch_id dari claims
	branchID, _ := helpers.GetClaimsToken(c, "branch_id")

	return branchID, nil
}

func GetUserRole(c *fiber.Ctx) (string, error) {
	// Ambil user_role dari claims
	userRole, _ := helpers.GetClaimsToken(c, "user_role")

	return userRole, nil
}
