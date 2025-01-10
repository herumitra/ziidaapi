package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/router"
	"github.com/herumitra/ziidaapi/seeders"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get port from environment
	serverPort := os.Getenv("SERVER_PORT")

	// Initialize database
	if err := config.SetupDB(); err != nil {
		log.Fatal(err)
	}

	// Check for command line arguments
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		seeders.UserSeed()
		seeders.BranchSeed()
		seeders.UserBranchSeed()
		seeders.UnitSeed()
		seeders.UnitConversionSeed()
		seeders.ProductCategorySeed()
		seeders.ProductSeed()
		seeders.MemberCategorySeed()
		seeders.MemberSeed()
		seeders.SupplierCategorySeed()
		seeders.SupplierSeed()
		seeders.SupplierProductSeed()
		os.Exit(0) // Exit after seeding
	}

	// Initialize app
	app := fiber.New()

	// Setup routes
	router.SetupRoutes(app)

	// Start app
	log.Fatal(app.Listen(":" + serverPort))
}
