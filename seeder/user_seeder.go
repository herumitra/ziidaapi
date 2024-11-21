package seeder

import (
	"log"

	"github.com/herumitra/ziidaapi/konfig"
	"github.com/herumitra/ziidaapi/models"
)

func SeedUsers() {
	users := []models.User{
		{Username: "admin", Password: "admin123", BranchID: 1},
		{Username: "user1", Password: "password1", BranchID: 2},
	}

	for _, user := range users {
		// Hash password sebelum menyimpan ke database
		if err := user.HashPassword(); err != nil {
			log.Fatal("Failed to hash password:", err)
		}

		// FirstOrCreate untuk menghindari duplikasi
		konfig.DB.FirstOrCreate(&user, models.User{Username: user.Username})
	}

	log.Println("Seeded initial users successfully!")
}
