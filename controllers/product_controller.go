package controllers

import (
	fmt "fmt"
	strconv "strconv"
	time "time"

	fiber "github.com/gofiber/fiber/v2"
	config "github.com/herumitra/ziidaapi/config"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
	services "github.com/herumitra/ziidaapi/services"
	gorm "gorm.io/gorm"
)

// GetAllProduct tampilkan semua product
func GetAllProduct(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Inisialize variabel products
	var products []models.Product

	// Mengambil semua data branch dari database dengan filter branch_id
	if err := config.DB.Where("branch_id = ?", branchID).Find(&products).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get products failed", "Failed to fetch products")
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Products retrieved successfully", products)
}

// GetProduct tampilkan product berdasarkan id
func GetProduct(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id := c.Params("id")
	var product models.Product

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&product).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Product not found", err)
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Product found", product)
}

// CreateProduct buat product
func CreateProduct(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Buat instance baru untuk Product
	var product models.Product

	// Parse body request ke struct product
	if err := c.BodyParser(&product); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	productID, err := generateProductID(config.DB) // Pastikan generateProductID menerima DB dan mengembalikan ID yang valid
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Create product failed", "Failed to generate product ID")
	}

	product.ID = productID
	// Tambahkan branchID ke product (relasi ke cabang)
	product.BranchID = branchID

	// Simpan product ke database menggunakan GORM
	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "Product created successfully", product)
}

// UpdateProduct update product
func UpdateProduct(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil productID dari parameter
	id := c.Params("id")
	var product models.Product

	// Cari product berdasarkan ID dan branchID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&product).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Product not found", err)
	}

	// Parse body request ke struct product
	if err := c.BodyParser(&product); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Simpan product ke database menggunakan GORM
	if err := config.DB.Model(&product).Updates(product).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update product", err)
	}

	// Mengembalikan response data product yang diupdate
	return helpers.JSONResponse(c, fiber.StatusOK, "Product created successfully", product)

}

// DeleteProduct hapus product
func DeleteProduct(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil productID dari parameter
	id := c.Params("id")
	var product models.Product

	// Ambil data product
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&product).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete product", err)
	}

	// Hapus product
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).Delete(&product).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete product", err)
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "Product deleted successfully", product)
}

// generateProductID generate id product
func generateProductID(db *gorm.DB) (string, error) {
	// Ambil tanggal saat ini dalam format DDMMYYYY
	now := time.Now()
	dateStr := now.Format("02012006") // Format DDMMYYYY

	var product models.Product // Menggunakan model yang sudah ada

	// Ambil urutan terbesar untuk tanggal tersebut
	if err := db.Where("id LIKE ?", "PRD"+dateStr+"%").Order("id DESC").First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Jika tidak ada user sebelumnya, urutan awal adalah 1
			return fmt.Sprintf("PRD%s001", dateStr), nil
		} else {
			return "", fmt.Errorf("error querying database: %v", err)
		}
	}

	// Jika ditemukan product sebelumnya, ambil urutan terakhir dan tambah 1
	lastID := product.ID             // Ambil ID product.ID
	seqStr := lastID[len(lastID)-3:] // Ambil 3 digit terakhir dari ID sebelumnya
	seq, err := strconv.Atoi(seqStr)
	if err != nil {
		return "", fmt.Errorf("error converting sequence: %v", err)
	}
	sequence := seq + 1

	// Format ID baru dengan urutan 3 digit
	sequenceStr := fmt.Sprintf("%03d", sequence)
	productID := fmt.Sprintf("PRD%s%s", dateStr, sequenceStr)

	return productID, nil
}
