package seeders

import (
	time "time"

	config "github.com/herumitra/ziidaapi/config"
	models "github.com/herumitra/ziidaapi/models"
)

func UnitSeed() {
	unit := []models.Unit{
		{ID: "UNT12122024001", Name: "Pcs", BranchID: "BRC11122024001"},
		{ID: "UNT12122024002", Name: "Strip", BranchID: "BRC11122024001"},
		{ID: "UNT12122024003", Name: "Box", BranchID: "BRC11122024001"},
		{ID: "UNT12122024004", Name: "Box", BranchID: "BRC11122024002"},
	}
	config.DB.Create(&unit)
}

func ProductCategorySeed() {
	productCategory := []models.ProductCategory{
		{ID: 1, Name: "Obat", BranchID: "BRC11122024001"},
		{ID: 2, Name: "Vitamin", BranchID: "BRC11122024001"},
		{ID: 3, Name: "Suplemen", BranchID: "BRC11122024001"},
		{ID: 4, Name: "Susu", BranchID: "BRC11122024001"},
	}
	config.DB.Create(&productCategory)
}

func ProductSeed() {
	layoutDate := "2006-01-02 15:04:05"
	strDate := "2025-12-31 00:00:00"
	t, _ := time.Parse(layoutDate, strDate)
	product := []models.Product{
		{ID: "PRD1212202400001", Name: "Sanmol", Description: "Pracetamol", UnitId: "UNT12122024001", Stock: 100, ExpiredDate: t, SalesPrice: 10000, AlternatePrice: 10000, PurchasePrice: 5000, ProductCategoryId: 1, BranchID: "BRC11122024001"},
		{ID: "PRD1212202400002", Name: "Fenofibrate", Description: "Anti Kolesterol Tligiserida", UnitId: "UNT12122024001", Stock: 10, ExpiredDate: t, SalesPrice: 56900, AlternatePrice: 56900, PurchasePrice: 26000, ProductCategoryId: 1, BranchID: "BRC11122024001"},
		{ID: "PRD1212202400003", Name: "Simvastatin", Description: "Kolesterol Golongan Statin", UnitId: "UNT12122024001", Stock: 500, ExpiredDate: t, SalesPrice: 6500, AlternatePrice: 6500, PurchasePrice: 2000, ProductCategoryId: 1, BranchID: "BRC11122024001"},
	}
	config.DB.Create(&product)
}

func UnitConversionSeed() {
	unitConversion := []models.UnitConversion{
		{ID: "UNC12122024001", ProductId: "PRD1212202400001", UnitInitId: "UNT12122024002", UnitFinalId: "UNT12122024001", ValueConv: 10, BranchID: "BRC11122024001"},
		{ID: "UNC12122024002", ProductId: "PRD1212202400001", UnitInitId: "UNT12122024003", UnitFinalId: "UNT12122024001", ValueConv: 100, BranchID: "BRC11122024001"},
		{ID: "UNC12122024003", ProductId: "PRD1212202400001", UnitInitId: "UNT12122024003", UnitFinalId: "UNT12122024002", ValueConv: 10, BranchID: "BRC11122024001"},
	}
	config.DB.Create(&unitConversion)
}

func MemberCategorySeed() {
	memberCategory := []models.MemberCategory{
		{ID: 1, Name: "Reguler", BranchID: "BRC11122024001"},
		{ID: 2, Name: "Silver", BranchID: "BRC11122024001"},
		{ID: 2, Name: "Gold", BranchID: "BRC11122024001"},
	}
	config.DB.Create(&memberCategory)
}
