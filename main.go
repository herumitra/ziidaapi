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
		seeders.MemberCategorySeed()
		seeders.MemberSeed()
		seeders.SupplierCategorySeed()
		seeders.SupplierSeed()
		seeders.SupplierProductSeed()
		os.Exit(0) // Exit after seeding
	}

	// Initialize app
	app := fiber.New()

	// Adding logger middleware of fiber
	app.Use(logger.New())

	// Endpoints for Auth
	app.Post("/login", controllers.Login)
	app.Post("/logout", controllers.Logout)
	app.Get("/profile", middleware.JWTMiddleware, controllers.GetProfile)

	// API routes with JWT middleware applied
	api := app.Group("/api", middleware.JWTMiddleware) // Perbaikan: panggil middleware tanpa tanda kurung

	// Endpoints for SetBranch
	api.Post("/set_branch", controllers.SetBranch)

	// API routes with JWT and role middleware applied
	api_adm := app.Group("/api", middleware.JWTMiddleware, middleware.RoleMiddleware("administrator"))

	// API routes with JWT and role middleware applied
	api_adm_op := app.Group("/api", middleware.JWTMiddleware, middleware.RoleMiddleware("administrator", "operator"))

	// API routes with JWT and role middleware applied
	api_adm_op_cs_fn := app.Group("/api", middleware.JWTMiddleware, middleware.RoleMiddleware("administrator", "operator", "cashier", "finance"))

	// Endpoints for User
	api_adm.Post("/users", controllers.CreateUser)       //Create new User
	api_adm.Get("/users", controllers.GetAllUsers)       //Menampilkan semua user
	api_adm.Get("/users/:id", controllers.GetUser)       //Menampilkan user berdasarkan ID
	api_adm.Put("/users/:id", controllers.UpdateUser)    //Update user berdasarkan ID
	api_adm.Delete("/users/:id", controllers.DeleteUser) //Hapus user berdasarkan ID

	// Endpoints for Branch
	api_adm_op.Post("/branches", controllers.CreateBranch)    //Create new Branch
	api_adm_op.Get("/branches", controllers.GetAllBranch)     //Menampilkan semua branch
	api_adm_op.Get("/branches/:id", controllers.GetBranch)    //Menampilkan branch berdasarkan ID
	api_adm_op.Put("/branches/:id", controllers.UpdateBranch) //Update branch berdasarkan ID
	api_adm.Delete("/branches/:id", controllers.DeleteBranch) //Hapus branch berdasarkan ID

	// Endpoints for UserBranch
	api_adm.Post("/user_branches", controllers.CreateUserBranch)                     //Create new UserBranch
	api_adm_op_cs_fn.Get("/user_branches", controllers.GetAllUserBranch)             //Menampilkan semua user_branch
	api_adm_op_cs_fn.Get("/user_branches/:userid", controllers.GetUserBranch)        //Menampilkan user_branch berdasarkan ID
	api_adm.Put("/user_branches/:userid/:branchid", controllers.UpdateUserBranch)    //Update user_branch berdasarkan ID
	api_adm.Delete("/user_branches/:userid/:branchid", controllers.DeleteUserBranch) //Hapus user_branch berdasarkan ID

	// Endpoints for Unit
	api_adm_op.Post("/units", controllers.CreateUnit)       //Create new Unit
	api_adm_op_cs_fn.Get("/units", controllers.GetAllUnit)  //Menampilkan semua unit
	api_adm_op_cs_fn.Get("/units/:id", controllers.GetUnit) //Menampilkan unit berdasarkan ID
	api_adm_op.Put("/units/:id", controllers.UpdateUnit)    //Update unit berdasarkan ID
	api_adm.Delete("/units/:id", controllers.DeleteUnit)    //Hapus unit berdasarkan ID

	// Endpoints for ProductCategory
	api_adm_op.Post("/product_categories", controllers.CreateProductCategory)       //Create new ProductCategory
	api_adm_op_cs_fn.Get("/product_categories", controllers.GetAllProductCategory)  //Menampilkan semua product_category
	api_adm_op_cs_fn.Get("/product_categories/:id", controllers.GetProductCategory) //Menampilkan product_category berdasarkan ID
	api_adm_op.Put("/product_categories/:id", controllers.UpdateProductCategory)    //Update product_category berdasarkan ID
	api_adm.Delete("/product_categories/:id", controllers.DeleteProductCategory)    //Hapus product_category berdasarkan ID

	// Endpoints for Product
	api_adm_op.Post("/products", controllers.CreateProduct)       //Create new Product
	api_adm_op_cs_fn.Get("/products", controllers.GetAllProduct)  //Menampilkan semua product
	api_adm_op_cs_fn.Get("/products/:id", controllers.GetProduct) //Menampilkan product berdasarkan ID
	api_adm_op.Put("/products/:id", controllers.UpdateProduct)    //Update product berdasarkan ID
	api_adm.Delete("/products/:id", controllers.DeleteProduct)    //Hapus product berdasarkan ID

	// Endpoints for UnitConversion
	api_adm_op.Post("/unit_conversions", controllers.CreateUnitConversion)       //Create new UnitConversion
	api_adm_op_cs_fn.Get("/unit_conversions", controllers.GetAllUnitConversion)  //Menampilkan semua unit_conversion
	api_adm_op_cs_fn.Get("/unit_conversions/:id", controllers.GetUnitConversion) //Menampilkan unit_conversion berdasarkan ID
	api_adm_op.Put("/unit_conversions/:id", controllers.UpdateUnitConversion)    //Update unit_conversion berdasarkan ID
	api_adm.Delete("/unit_conversions/:id", controllers.DeleteUnitConversion)    //Hapus unit_conversion berdasarkan ID

	// Endpoints for MemberCategory
	api_adm_op.Post("/member_categories", controllers.CreateMemberCategory)       //Create new MemberCategory
	api_adm_op_cs_fn.Get("/member_categories", controllers.GetAllMemberCategory)  //Menampilkan semua member_category
	api_adm_op_cs_fn.Get("/member_categories/:id", controllers.GetMemberCategory) //Menampilkan member_category berdasarkan ID
	api_adm_op.Put("/member_categories/:id", controllers.UpdateMemberCategory)    //Update member_category berdasarkan ID
	api_adm_op.Delete("/member_categories/:id", controllers.DeleteMemberCategory) //Hapus member_category berdasarkan ID

	// Endpoints for Member
	api_adm_op.Post("/members", controllers.CreateMember)       //Create new Member
	api_adm_op_cs_fn.Get("/members", controllers.GetAllMember)  //Menampilkan semua member
	api_adm_op_cs_fn.Get("/members/:id", controllers.GetMember) //Menampilkan member berdasarkan ID
	api_adm_op.Put("/members/:id", controllers.UpdateMember)    //Update member berdasarkan ID
	api_adm_op.Delete("/members/:id", controllers.DeleteMember) //Hapus member berdasarkan ID

	// Endpoints for SupplierCategory
	api_adm_op.Post("/supplier_categories", controllers.CreateSupplierCategory)       //Create new SupplierCategory
	api_adm_op_cs_fn.Get("/supplier_categories", controllers.GetAllSupplierCategory)  //Menampilkan semua supplier_category
	api_adm_op_cs_fn.Get("/supplier_categories/:id", controllers.GetSupplierCategory) //Menampilkan supplier_category berdasarkan ID
	api_adm_op.Put("/supplier_categories/:id", controllers.UpdateSupplierCategory)    //Update supplier_category berdasarkan ID
	api_adm_op.Delete("/supplier_categories/:id", controllers.DeleteSupplierCategory) //Hapus supplier_category berdasarkan ID

	// Endpoints for Supplier
	api_adm_op.Post("/suppliers", controllers.CreateSupplier)       //Create new Supplier
	api_adm_op_cs_fn.Get("/suppliers", controllers.GetAllSupplier)  //Menampilkan semua supplier
	api_adm_op_cs_fn.Get("/suppliers/:id", controllers.GetSupplier) //Menampilkan supplier berdasarkan ID
	api_adm_op.Put("/suppliers/:id", controllers.UpdateSupplier)    //Update supplier berdasarkan ID
	api_adm_op.Delete("/suppliers/:id", controllers.DeleteSupplier) //Hapus supplier berdasarkan ID

	// Endpoints for SupplierProduct
	api_adm_op.Post("/supplier_products", controllers.CreateSupplierProduct)       //Create new SupplierProduct
	api_adm_op_cs_fn.Get("/supplier_products", controllers.GetAllSupplierProduct)  //Menampilkan semua supplier_product
	api_adm_op_cs_fn.Get("/supplier_products/:id", controllers.GetSupplierProduct) //Menampilkan supplier_product berdasarkan ID
	api_adm_op.Put("/supplier_products/:id", controllers.UpdateSupplierProduct)    //Update supplier_product berdasarkan ID
	api_adm_op.Delete("/supplier_products/:id", controllers.DeleteSupplierProduct) //Hapus supplier_product berdasarkan ID

	// Start app
	app.Listen(":" + serverPort)
}
