package seeders

import (
	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/models"
)

// SeedUsers menambahkan data pengguna ke database
func SeedUsers() {
	users := []models.User{
		{Username: "john_doe", Password: "password123"},
		{Username: "jane_doe", Password: "password456"},
		{Username: "sam_smith", Password: "password789"},
		{Username: "alice_williams", Password: "password000"},
		{Username: "bob_jones", Password: "password111"},
	}

	// Hash password dan simpan ke database
	for _, user := range users {
		if err := user.HashPassword(); err != nil {
			continue
		}
		config.DB.Create(&user)
	}
}
