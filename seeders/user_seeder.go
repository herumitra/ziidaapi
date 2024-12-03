package seeders

import (
	"time"

	"github.com/herumitra/ziidaapi/config"
	"github.com/herumitra/ziidaapi/models"
)

// SeedUsers menambahkan data pengguna ke database
func SeedUsers() {
	users := []models.User{
		{Username: "john_doe", Password: "password123", UserStatus: "active", Address: "Jl. Contoh John, No. User jhon_doe, Kota Contoh John", UserRole: "operator"},
		{Username: "jane_doe", Password: "password456", UserStatus: "active", Address: "Jl. Contoh Jane, No. User jane_doe, Kota Contoh Jane", UserRole: "cashier"},
		{Username: "sam_smith", Password: "password789", UserStatus: "active", Address: "Jl. Contoh Sam, No. User sam_smith, Kota Contoh Sam", UserRole: "administrator"},
	}

	// Hash password dan simpan ke database
	for _, user := range users {
		if err := user.HashPassword(); err != nil {
			continue
		}
		config.DB.Create(&user)
	}
}

func SeedBranch() {
	t := time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC)

	branch := []models.Branch{
		{BranchName: "Branch 1", Address: "Jl. Raya Gudo, No. 101A, Gudo - Jombang", Phone: "085236990001", Email: "heru.oktafian@yahoo.com", SiaId: "xxxxxx", SiaName: "Mitra Farma", PsaId: "3517011710880001", PsaName: "Heru Oktafian", Sipa: "xxxxxx", SipaName: "Vita Fauzi. M", JournalMethod: "automatic", TaxPercentage: 11, BranchStatus: "active", LicenseDate: t},
		{BranchName: "Branch 2", Address: "Perum. Griya Nagari Singosari, Blok L, No. 1, Singosari - Malang", Phone: "0882009990001", Email: "oktafianheru@gmail.com", SiaId: "xxxxxx", SiaName: "Mitra Farma", PsaId: "3517011710880001", PsaName: "Heru Oktafian", Sipa: "xxxxxx", SipaName: "Vita Fauzi. M", JournalMethod: "automatic", TaxPercentage: 11, BranchStatus: "active", LicenseDate: t},
	}
	config.DB.Create(&branch)
}

func SeedUserBranch() {
	userBranch := []models.UserBranch{
		{UserID: 1, BranchID: 1},
		{UserID: 1, BranchID: 2},
		{UserID: 3, BranchID: 1},
	}
	config.DB.Create(&userBranch)
}

func SeedUnit() {
	unit := []models.Unit{
		{ID: "UN02012006150401", Name: "Pcs", BranchID: 1},
		{ID: "UN02012006150402", Name: "Strip", BranchID: 1},
		{ID: "UN02012006150403", Name: "Box", BranchID: 1},
		{ID: "UN02012006150404", Name: "Box", BranchID: 2},
	}
	config.DB.Create(&unit)
}

func SeedUnitConversion() {
	unitConversion := []models.UnitConversion{
		{ID: "UC02012006150401", ProductId: "PR02012006150401", UnitInitId: "UN02012006150402", UnitFinalId: "UN02012006150401", ValueConv: 10, BranchID: 1},
		{ID: "UC02012006150402", ProductId: "PR02012006150401", UnitInitId: "UN02012006150403", UnitFinalId: "UN02012006150401", ValueConv: 100, BranchID: 1},
		{ID: "UC02012006150403", ProductId: "PR02012006150401", UnitInitId: "UN02012006150403", UnitFinalId: "UN02012006150402", ValueConv: 10, BranchID: 1},
	}
	config.DB.Create(&unitConversion)
}

func SeedProductCategory() {
	productCategory := []models.ProductCategory{
		{ID: 1, Name: "Obat", BranchID: 1},
		{ID: 2, Name: "Vitamin", BranchID: 1},
		{ID: 3, Name: "Suplemen", BranchID: 1},
		{ID: 4, Name: "Susu", BranchID: 1},
	}
	config.DB.Create(&productCategory)
}

func SeedProduct() {
	layoutDate := "2006-01-02 15:04:05"
	strDate := "2025-12-31 00:00:00"
	t, _ := time.Parse(layoutDate, strDate)
	product := []models.Product{
		{ID: "PR02012006150401", Name: "Sanmol", Description: "Pracetamol", UnitId: "UN02012006150401", Stock: 100, ExpiredDate: t, SalesPrice: 10000, PurchasePrice: 5000, ProductCategoryId: 1, BranchID: 1},
		{ID: "PR02012006150402", Name: "Fenofibrate", Description: "Anti Kolesterol Tligiserida", UnitId: "UN02012006150401", Stock: 10, ExpiredDate: t, SalesPrice: 56900, PurchasePrice: 26000, ProductCategoryId: 1, BranchID: 1},
		{ID: "PR02012006150403", Name: "Simvastatin", Description: "Kolesterol Golongan Statin", UnitId: "UN02012006150401", Stock: 500, ExpiredDate: t, SalesPrice: 6500, PurchasePrice: 2000, ProductCategoryId: 1, BranchID: 1},
	}
	config.DB.Create(&product)
}
