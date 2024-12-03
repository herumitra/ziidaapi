package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/controllers"
	"github.com/herumitra/ziidaapi/middleware"
	"github.com/herumitra/ziidaapi/seeders"
	godotenv "github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	app.Use(logger.New())
	config.SetupDatabase()

	// Check for command line arguments
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		seeders.SeedUsers()
		seeders.SeedBranch()
		seeders.SeedUserBranch()
		seeders.SeedUnit()
		seeders.SeedProductCategory()
		seeders.SeedProduct()
		seeders.SeedUnitConversion()
		os.Exit(0) // Exit after seeding
	}

	// Define routes
	app.Post("/login", controllers.Login)
	app.Get("/logout", controllers.Logout)
	app.Post("/set_branch", middleware.JWTMiddleware, controllers.SetBranch)
	app.Get("/profile", middleware.JWTMiddleware, controllers.GetBranchFromToken)

	// API routes with JWT middleware applied
	api := app.Group("/api", middleware.JWTMiddleware) // Perbaikan: panggil middleware tanpa tanda kurung

	// Endpoint Users
	api.Get("/users", controllers.GetAllUsers)
	api.Post("/users", controllers.CreateUser)
	api.Get("/users/:id", controllers.GetUser)
	api.Put("/users/:id", controllers.UpdateUser)
	api.Delete("/users/:id", controllers.DeleteUser)

	// Endpoint Branch
	api.Get("/branch", controllers.GetAllBranch)
	api.Post("/branch", controllers.CreateBranch)
	api.Get("/branch/:id", controllers.GetBranch)
	api.Put("/branch/:id", controllers.UpdateBranch)
	api.Delete("/branch/:id", controllers.DeleteBranch)
	api.Get("/get_branches", controllers.GetBranchByUserIdFromToken)

	log.Println("Starting server...")

	app.Listen(":" + os.Getenv("SERVER_PORT"))

	log.Println("Server stopped")
}
