package seeders

import (
	config "github.com/herumitra/ziidaapi/config"
	models "github.com/herumitra/ziidaapi/models"
)

func UserBranchSeed() {
	userBranch := []models.UserBranch{
		{UserID: "USR01072023001", BranchID: "BRC11122024001"},
		{UserID: "USR01072023002", BranchID: "BRC11122024001"},
		{UserID: "USR01072023001", BranchID: "BRC11122024002"},
	}
	config.DB.Create(&userBranch)
}
