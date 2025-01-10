package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/herumitra/ziidaapi/controllers"
	"github.com/herumitra/ziidaapi/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Adding logger middleware of fiber
	app.Use(logger.New())

	// Auth Endpoints
	app.Post("/login", controllers.Login)
	app.Post("/logout", controllers.Logout)
	app.Get("/profile", middleware.JWTMiddleware, controllers.GetProfile)

	// SetBranch Endpoint
	app.Post("/set_branch", controllers.SetBranch)

	// Grouped API routes
	apiAdm := app.Group("/api", middleware.JWTMiddleware, middleware.RoleMiddleware("administrator"))
	apiAdmOp := app.Group("/api", middleware.JWTMiddleware, middleware.RoleMiddleware("administrator", "operator"))
	apiAdmOpCsFn := app.Group("/api", middleware.JWTMiddleware, middleware.RoleMiddleware("administrator", "operator", "cashier", "finance"))

	// User Endpoints
	apiAdm.Post("/users", controllers.CreateUser)
	apiAdm.Get("/users", controllers.GetAllUsers)
	apiAdm.Get("/users/:id", controllers.GetUser)
	apiAdm.Put("/users/:id", controllers.UpdateUser)
	apiAdm.Delete("/users/:id", controllers.DeleteUser)

	// Branch Endpoints
	apiAdmOp.Post("/branches", controllers.CreateBranch)
	apiAdmOp.Get("/branches", controllers.GetAllBranch)
	apiAdmOp.Get("/branches/:id", controllers.GetBranch)
	apiAdmOp.Put("/branches/:id", controllers.UpdateBranch)
	apiAdm.Delete("/branches/:id", controllers.DeleteBranch)

	// UserBranch Endpoints
	apiAdm.Post("/user_branches", controllers.CreateUserBranch)
	apiAdmOpCsFn.Get("/user_branches", controllers.GetAllUserBranch)
	apiAdmOpCsFn.Get("/user_branches/:userid", controllers.GetUserBranch)
	apiAdm.Put("/user_branches/:userid/:branchid", controllers.UpdateUserBranch)
	apiAdm.Delete("/user_branches/:userid/:branchid", controllers.DeleteUserBranch)

	// Unit Endpoints
	apiAdmOp.Post("/units", controllers.CreateUnit)
	apiAdmOpCsFn.Get("/units", controllers.GetAllUnit)
	apiAdmOpCsFn.Get("/units/:id", controllers.GetUnit)
	apiAdmOp.Put("/units/:id", controllers.UpdateUnit)
	apiAdm.Delete("/units/:id", controllers.DeleteUnit)

	// Add similar routes for ProductCategory, Product, UnitConversion, MemberCategory, Member, SupplierCategory, Supplier, and SupplierProduct
	// Repeat the pattern above for each resource
}
