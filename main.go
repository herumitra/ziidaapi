package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/controllers"
	"github.com/herumitra/ziidaapi/middleware"
	"github.com/herumitra/ziidaapi/seeders"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	config.SetupDatabase()

	// Check for command line arguments
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		seeders.SeedUsers()
		os.Exit(0) // Exit after seeding
	}

	// Define routes
	app.Post("/login", controllers.Login)
	app.Get("/logout", controllers.Logout)
	api := app.Group("/api", middleware.JWTMiddleware())
	api.Get("/users", controllers.GetAllUsers)
	api.Post("/users", controllers.CreateUser)
	api.Get("/users/:id", controllers.GetUser)
	api.Put("/users/:id", controllers.UpdateUser)
	api.Delete("/users/:id", controllers.DeleteUser)

	app.Listen(":3000")
}
