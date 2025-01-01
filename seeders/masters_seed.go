package seeders

import (
	time "time"

	config "github.com/herumitra/ziidaapi/config"
	models "github.com/herumitra/ziidaapi/models"
)

// Function for unit seed
func UnitSeed() {
	unit := []models.Unit{
		{ID: "UNT12122024001", Name: "Pcs", BranchID: "BRC11122024001"},
		{ID: "UNT12122024002", Name: "Strip", BranchID: "BRC11122024001"},
		{ID: "UNT12122024003", Name: "Box", BranchID: "BRC11122024001"},
		{ID: "UNT12122024004", Name: "Box", BranchID: "BRC11122024002"},
	}
	config.DB.Create(&unit)
}

// Function for product category seed
func ProductCategorySeed() {
	productCategory := []models.ProductCategory{
		{ID: 1, Name: "Obat", BranchID: "BRC11122024001"},
		{ID: 2, Name: "Vitamin", BranchID: "BRC11122024001"},
		{ID: 3, Name: "Suplemen", BranchID: "BRC11122024001"},
		{ID: 4, Name: "Susu", BranchID: "BRC11122024001"},
	}
	config.DB.Create(&productCategory)
}

// Function for product seed
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

// Function for unit conversion seed
func UnitConversionSeed() {
	unitConversion := []models.UnitConversion{
		{ID: "UNC12122024001", ProductId: "PRD1212202400001", UnitInitId: "UNT12122024002", UnitFinalId: "UNT12122024001", ValueConv: 10, BranchID: "BRC11122024001"},
		{ID: "UNC12122024002", ProductId: "PRD1212202400001", UnitInitId: "UNT12122024003", UnitFinalId: "UNT12122024001", ValueConv: 100, BranchID: "BRC11122024001"},
		{ID: "UNC12122024003", ProductId: "PRD1212202400001", UnitInitId: "UNT12122024003", UnitFinalId: "UNT12122024002", ValueConv: 10, BranchID: "BRC11122024001"},
	}
	config.DB.Create(&unitConversion)
}

// Function for member category seed
func MemberCategorySeed() {
	memberCategory := []models.MemberCategory{
		{ID: 1, Name: "Reguler", BranchID: "BRC11122024001"},
		{ID: 2, Name: "Silver", BranchID: "BRC11122024001"},
		{ID: 3, Name: "Gold", BranchID: "BRC11122024001"},
	}
	config.DB.Create(&memberCategory)
}

// Function for member seed
func MemberSeed() {
	member := []models.Member{
		{ID: "MBR12122024001", Name: "Member 1", MemberCategoryId: 1, Phone: "08123456789", Address: "Jl. Member 1", Saldo: 100000, BranchID: "BRC11122024001"},
		{ID: "MBR12122024002", Name: "Member 2", MemberCategoryId: 2, Phone: "08523456789", Address: "Jl. Member 2", Saldo: 200000, BranchID: "BRC11122024001"},
		{ID: "MBR12122024003", Name: "Member 3", MemberCategoryId: 3, Phone: "08823456789", Address: "Jl. Member 3", Saldo: 300000, BranchID: "BRC11122024001"},
	}
	config.DB.Create(&member)
}

// Function for supplier category seed
func SupplierCategorySeed() {
	supplierCategory := []models.SupplierCategory{
		{ID: 1, Name: "PBF", BranchID: "BRC11122024001"},
		{ID: 2, Name: "Distributor", BranchID: "BRC11122024001"},
		{ID: 3, Name: "Sub Distributor", BranchID: "BRC11122024001"},
	}
	config.DB.Create(&supplierCategory)
}

// Function for supplier seed
func SupplierSeed() {
	supplier := []models.Supplier{
		{ID: "SPL12122024001", Name: "Supplier 1", SupplierCategoryId: 1, Phone: "08123456789", Address: "Jl. Supplier 1", BranchID: "BRC11122024001"},
		{ID: "SPL12122024002", Name: "Supplier 2", SupplierCategoryId: 2, Phone: "08523456789", Address: "Jl. Supplier 2", BranchID: "BRC11122024001"},
		{ID: "SPL12122024003", Name: "Supplier 3", SupplierCategoryId: 3, Phone: "08823456789", Address: "Jl. Supplier 3", BranchID: "BRC11122024001"},
	}
	config.DB.Create(&supplier)
}

// Function for supplier product seed
func SupplierProductSeed() {
	supplierProduct := []models.SupplierProduct{
		{ID: "SPP12122024001", SupplierId: "SPL12122024001", ProductId: "PRD1212202400001", BranchID: "BRC11122024001"},
		{ID: "SPP12122024002", SupplierId: "SPL12122024002", ProductId: "PRD1212202400002", BranchID: "BRC11122024001"},
		{ID: "SPP12122024003", SupplierId: "SPL12122024003", ProductId: "PRD1212202400003", BranchID: "BRC11122024001"},
	}
	config.DB.Create(&supplierProduct)
}
