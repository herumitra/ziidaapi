package main

import (
	log "log"
	os "os"

	fiber "github.com/gofiber/fiber/v2"
	logger "github.com/gofiber/fiber/v2/middleware/logger"
	config "github.com/herumitra/ziidaapi/config"
	controllers "github.com/herumitra/ziidaapi/controllers"
	middleware "github.com/herumitra/ziidaapi/middleware"
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
		seeders.BranchSeed()
		seeders.UserBranchSeed()
		seeders.UnitSeed()
		seeders.UnitConversionSeed()
		seeders.ProductCategorySeed()
		seeders.ProductSeed()
		os.Exit(0) // Exit after seeding
	}

	// Initialize app
	app := fiber.New()

	// Adding logger middleware of fiber
	app.Use(logger.New())

	// Endpoints for Auth
	app.Post("/login", controllers.Login)
	app.Post("/logout", controllers.Logout)

	// API routes with JWT middleware applied
	api := app.Group("/api", middleware.JWTMiddleware) // Perbaikan: panggil middleware tanpa tanda kurung

	// Endpoints for SetBranch
	api.Post("/set_branch", controllers.SetBranch)

	// API routes with JWT and role middleware applied
	api_role := app.Group("/api", middleware.JWTMiddleware, middleware.RoleMiddleware("administrator"))

	// Endpoints for User
	api_role.Post("/users", controllers.CreateUser)       //Create new User
	api_role.Get("/users", controllers.GetAllUsers)       //Menampilkan semua user
	api_role.Get("/users/:id", controllers.GetUser)       //Menampilkan user berdasarkan ID
	api_role.Put("/users/:id", controllers.UpdateUser)    //Update user berdasarkan ID
	api_role.Delete("/users/:id", controllers.DeleteUser) //Hapus user berdasarkan ID

	// Endpoints for Branch
	api_role.Post("/branches", controllers.CreateBranch) //Create new Branch

	// Start app
	app.Listen(":" + serverPort)

	// app.Get("/operational", JWTMiddleware, RoleMiddleware(models.Operator), OperationalHandler)
}
