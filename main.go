package main

import (
	log "log"
	os "os"

	fiber "github.com/gofiber/fiber/v2"
	logger "github.com/gofiber/fiber/v2/middleware/logger"
	config "github.com/herumitra/ziidaapi/config"
	controllers "github.com/herumitra/ziidaapi/controllers"
	seeders "github.com/herumitra/ziidaapi/seeders"
	godotenv "github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get port from environment
	serverPort := os.Getenv("SERVER_PORT")

	// Initialize database
	err := config.SetupDB()
	if err != nil {
		log.Fatal(err)
	}

	// Check for command line arguments
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		seeders.UserSeed()
		os.Exit(0) // Exit after seeding
	}

	// Initialize app
	app := fiber.New()

	// Adding logger middleware of fiber
	app.Use(logger.New())

	// Testing Route
	app.Get("/testing", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Post Create User
	app.Post("/users", controllers.CreateUser)

	// Start app
	app.Listen(":" + serverPort)

}
