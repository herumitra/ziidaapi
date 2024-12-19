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

// GetAllSupplier tampilkan semua supplier
func GetAllSupplier(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Inisialize variabel suppliers
	var suppliers []models.Supplier

	// Mengambil semua data branch dari database dengan filter branch_id
	if err := config.DB.Where("branch_id = ?", branchID).Find(&suppliers).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get suppliers failed", "Failed to fetch suppliers")
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Suppliers retrieved successfully", suppliers)
}

// GetSupplier tampilkan supplier berdasarkan id
func GetSupplier(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id := c.Params("id")
	var supplier models.Supplier

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&supplier).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Supplier not found", err)
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Supplier found", supplier)
}

// CreateSupplier buat supplier
func CreateSupplier(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Buat instance baru untuk Supplier
	var supplier models.Supplier

	// Parse body request ke struct supplier
	if err := c.BodyParser(&supplier); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	supplierID, err := generateSupplierID(config.DB) // Pastikan generateSupplierID menerima DB dan mengembalikan ID yang valid
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Create supplier failed", "Failed to generate supplier ID")
	}

	supplier.ID = supplierID
	// Tambahkan branchID ke supplier (relasi ke cabang)
	supplier.BranchID = branchID

	// Simpan supplier ke database menggunakan GORM
	if err := config.DB.Create(&supplier).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create supplier",
		})
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "Supplier created successfully", supplier)
}

// UpdateSupplier update supplier
func UpdateSupplier(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil supplierID dari parameter
	id := c.Params("id")
	var supplier models.Supplier

	// Cari supplier berdasarkan ID dan branchID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&supplier).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Supplier not found", err)
	}

	// Parse body request ke struct supplier
	if err := c.BodyParser(&supplier); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Simpan supplier ke database menggunakan GORM
	if err := config.DB.Model(&supplier).Updates(supplier).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update supplier", err)
	}

	// Mengembalikan response data supplier yang diupdate
	return helpers.JSONResponse(c, fiber.StatusOK, "Supplier created successfully", supplier)

}

// DeleteSupplier hapus supplier
func DeleteSupplier(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil supplierID dari parameter
	id := c.Params("id")
	var supplier models.Supplier

	// Ambil data supplier
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&supplier).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete supplier", err)
	}

	// Hapus supplier
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).Delete(&supplier).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete supplier", err)
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "Supplier deleted successfully", supplier)
}

// generateSupplierID generate id supplier
func generateSupplierID(db *gorm.DB) (string, error) {
	// Ambil tanggal saat ini dalam format DDMMYYYY
	now := time.Now()
	dateStr := now.Format("02012006") // Format DDMMYYYY

	var supplier models.Supplier // Menggunakan model yang sudah ada

	// Ambil urutan terbesar untuk tanggal tersebut
	if err := db.Where("id LIKE ?", "SPL"+dateStr+"%").Order("id DESC").First(&supplier).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Jika tidak ada user sebelumnya, urutan awal adalah 1
			return fmt.Sprintf("SPL%s001", dateStr), nil
		} else {
			return "", fmt.Errorf("error querying database: %v", err)
		}
	}

	// Jika ditemukan supplier sebelumnya, ambil urutan terakhir dan tambah 1
	lastID := supplier.ID            // Ambil ID supplier.ID
	seqStr := lastID[len(lastID)-3:] // Ambil 3 digit terakhir dari ID sebelumnya
	seq, err := strconv.Atoi(seqStr)
	if err != nil {
		return "", fmt.Errorf("error converting sequence: %v", err)
	}
	sequence := seq + 1

	// Format ID baru dengan urutan 3 digit
	sequenceStr := fmt.Sprintf("%03d", sequence)
	supplierID := fmt.Sprintf("SPL%s%s", dateStr, sequenceStr)

	return supplierID, nil
}
