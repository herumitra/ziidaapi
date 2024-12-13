package seeders

import (
	log "log"

	models "github.com/herumitra/ziidaapi/models"
	gorm "gorm.io/gorm"
)

func MenuSeed(db *gorm.DB) {
	menus := []models.Menu{
		{Name: "Dashboard", Route: "/dashboard", AllowedRoles: []string{"operator", "cashier", "finance", "administrator"}},
		{Name: "Users", Route: "/users", AllowedRoles: []string{"administrator"}},
		{Name: "User Branches", Route: "/user_branch", AllowedRoles: []string{"administrator"}},
		{Name: "Branches", Route: "/branches", AllowedRoles: []string{"operator", "cashier", "finance", "administrator"}},
		{Name: "Set Branch", Route: "/set_branch", AllowedRoles: []string{"operator", "cashier", "finance", "administrator"}},
		{Name: "Member Categories", Route: "/finance", AllowedRoles: []string{"operator", "administrator"}},
		{Name: "Members", Route: "/operational", AllowedRoles: []string{"operator", "finance", "administrator"}},
		{Name: "Product Categories", Route: "/admin", AllowedRoles: []string{"operator", "administrator"}},
		{Name: "Products", Route: "/admin", AllowedRoles: []string{"operator", "administrator"}},
		{Name: "Supplier Categories", Route: "/admin", AllowedRoles: []string{"operator", "administrator"}},
		{Name: "Suppliers", Route: "/admin", AllowedRoles: []string{"operator", "administrator"}},
		{Name: "Units", Route: "/admin", AllowedRoles: []string{"operator", "administrator"}},
		{Name: "Unit Conversions", Route: "/admin", AllowedRoles: []string{"operator", "administrator"}},
	}

	for _, menu := range menus {
		if err := db.Create(&menu).Error; err != nil {
			log.Printf("Failed to seed menu: %v", err)
		}
	}
}
