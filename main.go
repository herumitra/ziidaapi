package main

import (
	log "log"
	os "os"

	fiber "github.com/gofiber/fiber/v2"
	logger "github.com/gofiber/fiber/v2/middleware/logger"
	config "github.com/herumitra/ziidaapi/config"
	controllers "github.com/herumitra/ziidaapi/controllers"
	"github.com/herumitra/ziidaapi/middleware"
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

	// Endpoints for Auth
	app.Post("/login", controllers.Login)

	// API routes with JWT middleware applied
	api := app.Group("/api", middleware.JWTMiddleware) // Perbaikan: panggil middleware tanpa tanda kurung

	// Endpoints for User
	api.Post("/users", controllers.CreateUser)       //Create new User
	api.Get("/users", controllers.GetAllUsers)       //Menampilkan semua user
	api.Get("/users/:id", controllers.GetUser)       //Menampilkan user berdasarkan ID
	api.Put("/users/:id", controllers.UpdateUser)    //Update user berdasarkan ID
	api.Delete("/users/:id", controllers.DeleteUser) //Hapus user berdasarkan ID

	// Start app
	app.Listen(":" + serverPort)

}
