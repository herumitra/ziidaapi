package seeders

import (
	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/models"
)

// SeedUsers menambahkan data pengguna ke database
func SeedUsers() {
	users := []models.User{
		{Username: "john_doe", Password: "password123", UserStatus: 1},
		{Username: "jane_doe", Password: "password456", UserStatus: 1},
		{Username: "sam_smith", Password: "password789", UserStatus: 1},
		{Username: "alice_williams", Password: "password000", UserStatus: 1},
		{Username: "bob_jones", Password: "password111", UserStatus: 1},
	}

	// Hash password dan simpan ke database
	for _, user := range users {
		if err := user.HashPassword(); err != nil {
			continue
		}
		config.DB.Create(&user)
	}
}
