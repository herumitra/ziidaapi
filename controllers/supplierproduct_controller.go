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

// GetAllSupplierProduct tampilkan semua supplier_product
func GetAllSupplierProduct(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Inisialize variabel supplier_products
	var supplier_products []models.SupplierProduct

	// Mengambil semua data branch dari database dengan filter branch_id
	if err := config.DB.Where("branch_id = ?", branchID).Find(&supplier_products).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get supplier_products failed", "Failed to fetch supplier_products")
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "SupplierProducts retrieved successfully", supplier_products)
}

// GetSupplierProduct tampilkan supplier_product berdasarkan id
func GetSupplierProduct(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id := c.Params("id")
	var supplier_product models.SupplierProduct

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&supplier_product).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "SupplierProduct not found", err)
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "SupplierProduct found", supplier_product)
}

// CreateSupplierProduct buat supplier_product
func CreateSupplierProduct(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Buat instance baru untuk SupplierProduct
	var supplier_product models.SupplierProduct

	// Parse body request ke struct supplier_product
	if err := c.BodyParser(&supplier_product); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	supplier_productID, err := generateSupplierProductID(config.DB) // Pastikan generateSupplierProductID menerima DB dan mengembalikan ID yang valid
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Create supplier_product failed", "Failed to generate supplier_product ID")
	}

	supplier_product.ID = supplier_productID
	// Tambahkan branchID ke supplier_product (relasi ke cabang)
	supplier_product.BranchID = branchID

	// Simpan supplier_product ke database menggunakan GORM
	if err := config.DB.Create(&supplier_product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create supplier_product",
		})
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "SupplierProduct created successfully", supplier_product)
}

// UpdateSupplierProduct update supplier_product
func UpdateSupplierProduct(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil supplier_productID dari parameter
	id := c.Params("id")
	var supplier_product models.SupplierProduct

	// Cari supplier_product berdasarkan ID dan branchID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&supplier_product).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "SupplierProduct not found", err)
	}

	// Parse body request ke struct supplier_product
	if err := c.BodyParser(&supplier_product); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Simpan supplier_product ke database menggunakan GORM
	if err := config.DB.Model(&supplier_product).Updates(supplier_product).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update supplier_product", err)
	}

	// Mengembalikan response data supplier_product yang diupdate
	return helpers.JSONResponse(c, fiber.StatusOK, "SupplierProduct created successfully", supplier_product)

}

// DeleteSupplierProduct hapus supplier_product
func DeleteSupplierProduct(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil supplier_productID dari parameter
	id := c.Params("id")
	var supplier_product models.SupplierProduct

	// Ambil data supplier_product
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&supplier_product).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete supplier_product", err)
	}

	// Hapus supplier_product
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).Delete(&supplier_product).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete supplier_product", err)
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "SupplierProduct deleted successfully", supplier_product)
}

// generateSupplierProductID generate id supplier_product
func generateSupplierProductID(db *gorm.DB) (string, error) {
	// Ambil tanggal saat ini dalam format DDMMYYYY
	now := time.Now()
	dateStr := now.Format("02012006") // Format DDMMYYYY

	var supplier_product models.SupplierProduct // Menggunakan model yang sudah ada

	// Ambil urutan terbesar untuk tanggal tersebut
	if err := db.Where("id LIKE ?", "SPP"+dateStr+"%").Order("id DESC").First(&supplier_product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Jika tidak ada user sebelumnya, urutan awal adalah 1
			return fmt.Sprintf("SPP%s001", dateStr), nil
		} else {
			return "", fmt.Errorf("error querying database: %v", err)
		}
	}

	// Jika ditemukan supplier_product sebelumnya, ambil urutan terakhir dan tambah 1
	lastID := supplier_product.ID    // Ambil ID supplier_product.ID
	seqStr := lastID[len(lastID)-3:] // Ambil 3 digit terakhir dari ID sebelumnya
	seq, err := strconv.Atoi(seqStr)
	if err != nil {
		return "", fmt.Errorf("error converting sequence: %v", err)
	}
	sequence := seq + 1

	// Format ID baru dengan urutan 3 digit
	sequenceStr := fmt.Sprintf("%03d", sequence)
	supplier_productID := fmt.Sprintf("SPP%s%s", dateStr, sequenceStr)

	return supplier_productID, nil
}
