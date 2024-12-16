package models

import "time"

// Unit model
type Unit struct {
	ID       string `gorm:"type:varchar(15);primaryKey" json:"id" validate:"required"`
	Name     string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	BranchID string `gorm:"type:varchar(15);not null" json:"branch_id" validate:"required"`

	Branch Branch `gorm:"foreignKey:BranchID"`
}

// UnitConversion model
type UnitConversion struct {
	ID          string `gorm:"type:varchar(15);primaryKey" json:"id" validate:"required"`
	ProductId   string `gorm:"type:varchar(15);not null" json:"product_id" validate:"required"`
	UnitInitId  string `gorm:"type:varchar(15);not null" json:"unit_init_id" validate:"required"`
	UnitFinalId string `gorm:"type:varchar(15);not null" json:"unit_final_id" validate:"required"`
	ValueConv   int    `gorm:"type:int;not null;default:0" json:"value_conv" validate:"required"`
	BranchID    string `gorm:"type:varchar(15);not null" json:"branch_id" validate:"required"`

	Branch Branch `gorm:"foreignKey:BranchID"`
	// Product   Product `gorm:"foreignKey:ProductId"`
	UnitInit  Unit `gorm:"foreignKey:UnitInitId"`
	UnitFinal Unit `gorm:"foreignKey:UnitFinalId"`
}

// ProductCategory model
type ProductCategory struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	BranchID string `gorm:"type:varchar(15);not null" json:"branch_id" validate:"required"`

	Branch Branch `gorm:"foreignKey:BranchID"`
}

// Product model
type Product struct {
	ID                string    `gorm:"type:varchar(15);primaryKey" json:"id" validate:"required"`
	Name              string    `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Description       string    `gorm:"type:text;" json:"description"`
	UnitId            string    `gorm:"type:varchar(15);not null" json:"unit_id" validate:"required"`
	Stock             int       `gorm:"type:int;not null;default:0" json:"stock"`
	ExpiredDate       time.Time `gorm:"not null" json:"expired_date" validate:"required"`
	SalesPrice        int       `gorm:"type:int;not null;default:0" json:"sales_price" validate:"required"`
	AlternatePrice    int       `gorm:"type:int;not null;default:0" json:"alternate_price" validate:"required"`
	PurchasePrice     int       `gorm:"type:int;not null;default:0" json:"purchase_price" validate:"required"`
	ProductCategoryId uint      `gorm:"not null" json:"product_category_id" validate:"required"`
	BranchID          string    `gorm:"type:varchar(15);not null" json:"branch_id" validate:"required"`

	Branch          Branch          `gorm:"foreignKey:BranchID"`
	ProductCategory ProductCategory `gorm:"foreignKey:ProductCategoryId"`
	Unit            Unit            `gorm:"foreignKey:UnitId"`
}

// MemberCategory model
type MemberCategory struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	BranchID string `gorm:"type:varchar(15);not null" json:"branch_id" validate:"required"`

	Branch Branch `gorm:"foreignKey:BranchID"`
}

// Member model
type Member struct {
	ID               string `gorm:"type:varchar(15);primaryKey" json:"id" validate:"required"`
	Name             string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Phone            string `gorm:"type:varchar(100);" json:"phone"`
	Address          string `gorm:"type:text;" json:"address"`
	MemberCategoryId uint   `gorm:"not null" json:"member_category_id" validate:"required"`
	Saldo            int    `gorm:"type:int;not null;default:0" json:"saldo" validate:"required"`
	BranchID         string `gorm:"type:varchar(15);not null" json:"branch_id" validate:"required"`

	Branch         Branch         `gorm:"foreignKey:BranchID"`
	MemberCategory MemberCategory `gorm:"foreignKey:MemberCategoryId"`
}

// SupplierCategory model
type SupplierCategory struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	BranchID string `gorm:"type:varchar(15);not null" json:"branch_id" validate:"required"`

	Branch Branch `gorm:"foreignKey:BranchID"`
}

// Supplier model
type Supplier struct {
	ID                 string `gorm:"type:varchar(15);primaryKey" json:"id" validate:"required"`
	Name               string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Phone              string `gorm:"type:varchar(100);" json:"phone"`
	Address            string `gorm:"type:text;" json:"address"`
	SupplierCategoryId uint   `gorm:"not null" json:"supplier_category_id" validate:"required"`
	BranchID           string `gorm:"type:varchar(15);not null" json:"branch_id" validate:"required"`

	Branch           Branch           `gorm:"foreignKey:BranchID"`
	SupplierCategory SupplierCategory `gorm:"foreignKey:SupplierCategoryId"`
}

// SupplierProduct model
type SupplierProduct struct {
	ID         string `gorm:"type:varchar(15);primaryKey" json:"id" validate:"required"`
	SupplierId string `gorm:"type:varchar(15);not null" json:"supplier_id" validate:"required"`
	ProductId  string `gorm:"type:varchar(15);not null" json:"product_id" validate:"required"`
	BranchID   string `gorm:"type:varchar(15);not null" json:"branch_id" validate:"required"`

	Branch   Branch   `gorm:"foreignKey:BranchID"`
	Supplier Supplier `gorm:"foreignKey:SupplierId"`
	Product  Product  `gorm:"foreignKey:ProductId"`
}
