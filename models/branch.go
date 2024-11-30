package models

import (
	"time"

	"gorm.io/gorm"
)

type Branch struct {
	gorm.Model
	BranchName    string    `gorm:"unique;not null" json:"branch_name"`
	Address       string    `gorm:"type:text;" json:"address"`
	Phone         string    `gorm:"type:varchar(100);" json:"phone"`
	Email         string    `gorm:"type:varchar(100);" json:"email"`
	SiaId         string    `gorm:"type:varchar(100);" json:"sia_id"`
	SiaName       string    `gorm:"type:varchar(255);" json:"sia_name"`
	PsaId         string    `gorm:"type:varchar(100);" json:"psa_id"`
	PsaName       string    `gorm:"type:varchar(255);" json:"psa_name"`
	Sipa          string    `gorm:"type:varchar(100);" json:"sipa"`
	SipaName      string    `gorm:"type:varchar(255);" json:"sipa_name"`
	ApingId       string    `gorm:"type:varchar(100);" json:"aping_id"`
	ApingName     string    `gorm:"type:varchar(255);" json:"aping_name"`
	BankName      string    `gorm:"type:varchar(255);" json:"bank_name"`
	AccountName   string    `gorm:"type:varchar(255);" json:"account_name"`
	AccountNumber string    `gorm:"type:varchar(100);" json:"account_number"`
	TaxPercentage int       `gorm:"type:int(3);default:0" json:"tax_percentage"`
	JournalMethod string    `gorm:"type:enum('manual','automatic'); default:'automatic'" json:"journal_method" validate:"required"`
	BranchStatus  int       `gorm:"not null" json:"branch_status" default:"0"`
	LicenseDate   time.Time `gorm:"not null" json:"license_date" validate:"required"`
}
