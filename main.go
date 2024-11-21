package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/herumitra/ziidaapi/controllers"
	"github.com/herumitra/ziidaapi/konfig"
	"github.com/herumitra/ziidaapi/middleware"
	"github.com/herumitra/ziidaapi/models"
	"github.com/herumitra/ziidaapi/seeder"
)

func main() {
	// Define flags
	seed := flag.Bool("seed", false, "Run the database seeder")
	flag.Parse()

	// Connect to database
	konfig.ConnectDB()

	// AutoMigrate for models
	konfig.DB.AutoMigrate(&models.User{})

	// Jika flag `--seed` diaktifkan, jalankan seeder dan keluar
	if *seed {
		log.Println("Running database seeder...")
		seeder.SeedUsers()
		return
	}

	// Jika tidak, jalankan server
	log.Println("Starting server...")
	app := fiber.New()

	// Routes
	app.Post("/login", controllers.Login)
	app.Get("/protected", middleware.JWTProtected(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "You are authorized!"})
	})

	log.Fatal(app.Listen(":3000"))
}
