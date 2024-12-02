package seeders

import (
	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/models"
)

// SeedUsers menambahkan data pengguna ke database
func SeedUsers() {
	users := []models.User{
		{Username: "john_doe", Password: "password123", UserStatus: "active", Address: "Jl. Contoh John, No. User jhon_doe, Kota Contoh John", UserRole: "operator"},
		{Username: "jane_doe", Password: "password456", UserStatus: "active", Address: "Jl. Contoh Jane, No. User jane_doe, Kota Contoh Jane", UserRole: "cashier"},
		{Username: "sam_smith", Password: "password789", UserStatus: "active", Address: "Jl. Contoh Sam, No. User sam_smith, Kota Contoh Sam", UserRole: "administrator"},
	}

	// Hash password dan simpan ke database
	for _, user := range users {
		if err := user.HashPassword(); err != nil {
			continue
		}
		config.DB.Create(&user)
	}
}
