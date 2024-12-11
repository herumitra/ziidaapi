package seeders

import (
	config "github.com/herumitra/ziidaapi/config"
	models "github.com/herumitra/ziidaapi/models"
)

func UserSeed() {
	users := []models.User{
		{ID: "USR01072023001", Username: "john_doe", Password: "password1", Name: "John Doe", UserRole: "operator", UserStatus: "active"},
		{ID: "USR01072023002", Username: "jane_smith", Password: "password2", Name: "Jane Smith", UserRole: "cashier", UserStatus: "active"},
		{ID: "USR01072023003", Username: "bob_jones", Password: "password3", Name: "Bob Jones", UserRole: "finance", UserStatus: "active"},
		{ID: "USR01072023004", Username: "sarah_wilson", Password: "password4", Name: "Sarah Wilson", UserRole: "administrator", UserStatus: "active"},
	}

	// Hash password for each user
	for _, user := range users {
		if err := user.HashPassword(); err != nil {
			continue
		}
		config.DB.Create(&user)
	}
}
