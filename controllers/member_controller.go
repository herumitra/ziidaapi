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

// GetAllMember tampilkan semua member
func GetAllMember(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Inisialize variabel members
	var members []models.Member

	// Mengambil semua data branch dari database dengan filter branch_id
	if err := config.DB.Where("branch_id = ?", branchID).Find(&members).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Get members failed", "Failed to fetch members")
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Members retrieved successfully", members)
}

// GetMember tampilkan member berdasarkan id
func GetMember(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id := c.Params("id")
	var member models.Member

	// Cari user berdasarkan ID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&member).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Member not found", err)
	}

	// Mengembalikan response data branch
	return helpers.JSONResponse(c, fiber.StatusOK, "Member found", member)
}

// CreateMember buat member
func CreateMember(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	// Buat instance baru untuk Member
	var member models.Member

	// Parse body request ke struct member
	if err := c.BodyParser(&member); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	memberID, err := generateMemberID(config.DB) // Pastikan generateMemberID menerima DB dan mengembalikan ID yang valid
	if err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Create member failed", "Failed to generate member ID")
	}

	member.ID = memberID
	// Tambahkan branchID ke member (relasi ke cabang)
	member.BranchID = branchID

	// Simpan member ke database menggunakan GORM
	if err := config.DB.Create(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create member",
		})
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "Member created successfully", member)
}

// UpdateMember update member
func UpdateMember(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil memberID dari parameter
	id := c.Params("id")
	var member models.Member

	// Cari member berdasarkan ID dan branchID
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&member).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusNotFound, "Member not found", err)
	}

	// Parse body request ke struct member
	if err := c.BodyParser(&member); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid input", err)
	}

	// Simpan member ke database menggunakan GORM
	if err := config.DB.Model(&member).Updates(member).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to update member", err)
	}

	// Mengembalikan response data member yang diupdate
	return helpers.JSONResponse(c, fiber.StatusOK, "Member created successfully", member)

}

// DeleteMember hapus member
func DeleteMember(c *fiber.Ctx) error {
	// Panggil fungsi GetBranchID
	branchID, err := services.GetBranchID(c)
	if err != nil {
		// Tangani error (misalnya kirim response dengan error)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil memberID dari parameter
	id := c.Params("id")
	var member models.Member

	// Ambil data member
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).First(&member).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete member", err)
	}

	// Hapus member
	if err := config.DB.Where("id = ? AND branch_id = ?", id, branchID).Delete(&member).Error; err != nil {
		return helpers.JSONResponse(c, fiber.StatusInternalServerError, "Failed to delete member", err)
	}

	// Berikan response sukses
	return helpers.JSONResponse(c, fiber.StatusOK, "Member deleted successfully", member)
}

// generateMemberID generate id member
func generateMemberID(db *gorm.DB) (string, error) {
	// Ambil tanggal saat ini dalam format DDMMYYYY
	now := time.Now()
	dateStr := now.Format("02012006") // Format DDMMYYYY

	var member models.Member // Menggunakan model yang sudah ada

	// Ambil urutan terbesar untuk tanggal tersebut
	if err := db.Where("id LIKE ?", "MBR"+dateStr+"%").Order("id DESC").First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Jika tidak ada user sebelumnya, urutan awal adalah 1
			return fmt.Sprintf("MBR%s001", dateStr), nil
		} else {
			return "", fmt.Errorf("error querying database: %v", err)
		}
	}

	// Jika ditemukan member sebelumnya, ambil urutan terakhir dan tambah 1
	lastID := member.ID              // Ambil ID member.ID
	seqStr := lastID[len(lastID)-3:] // Ambil 3 digit terakhir dari ID sebelumnya
	seq, err := strconv.Atoi(seqStr)
	if err != nil {
		return "", fmt.Errorf("error converting sequence: %v", err)
	}
	sequence := seq + 1

	// Format ID baru dengan urutan 3 digit
	sequenceStr := fmt.Sprintf("%03d", sequence)
	memberID := fmt.Sprintf("MBR%s%s", dateStr, sequenceStr)

	return memberID, nil
}
