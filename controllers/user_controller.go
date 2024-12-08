package controllers

import (
	fmt "fmt"
	strconv "strconv"
	time "time"

	fiber "github.com/gofiber/fiber/v2"
	helpers "github.com/herumitra/ziidaapi/helpers"
	models "github.com/herumitra/ziidaapi/models"
	config "github.com/herumitra/ziidaapi/config"
)

var user models.User
var userCheck models.User

// Fungsi untuk generate ID user
func generateUserID() (string, error) {
	// Ambil tanggal saat ini dalam format DDMMYYYY
	now := time.Now()
	dateStr := now.Format("02012006") // Format DDMMYYYY

	// Ambil urutan terbesar untuk tanggal tersebut
	// var lastUser models.User
	if err := config.DB.Where("id LIKE ?", "USR"+now.Format("20060102")+"%").First(&user).Error; err != nil  && err !={
		// Jika tidak ada user sebelumnya, urutan awal adalah 1
		// sequence = 1
	}

	// Tentukan urutan, jika tidak ada user sebelumnya, mulai dari 1
	var sequence int
	if err ==  {
		sequence = 1
	} else {
		// Ambil angka urutan terakhir dan tambah 1
		lastID := user.ID                // Ambil ID user.ID
		seqStr := lastID[len(lastID)-3:] // Ambil 3 digit terakhir dari ID sebelumnya
		seq, _ := strconv.Atoi(seqStr)
		sequence = seq + 1
	}

	// Format ID baru dengan urutan 3 digit
	sequenceStr := fmt.Sprintf("%03d", sequence)
	userID := fmt.Sprintf("USR%s%s", dateStr, sequenceStr)

	return userID, nil
}

// Fungsi CreateUser
func CreateUser(c *fiber.Ctx) error {
	// user.ID, _ = generateUserID(db)
	// return db.Create(user).Error

	if err := c.BodyParser(&user); err != nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := db.Where("username = ?", user.Username).
	Order("id desc").
	First(&user).Error; err == nil {
		return helpers.JSONResponse(c, fiber.StatusBadRequest, "Username already exists", nil)
	}

	return helpers.JSONResponse(c, fiber.StatusOK, "User created successfully", user)
}
