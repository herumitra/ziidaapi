package models

import "time"

// Mendefinisikan custom type untuk ENUM UserRole
type JournalMethod string

const (
	Manual    JournalMethod = "manual"
	Automatic JournalMethod = "automatic"
)

type Branch struct {
	ID            string        `gorm:"type:varchar(15);primaryKey" json:"id" validate:"required"`
	BranchName    string        `gorm:"unique;not null" json:"branch_name"`
	Address       string        `gorm:"type:text;" json:"address"`
	Phone         string        `gorm:"type:varchar(100);" json:"phone"`
	Email         string        `gorm:"type:varchar(100);" json:"email"`
	SiaId         string        `gorm:"type:varchar(100);" json:"sia_id"`
	SiaName       string        `gorm:"type:varchar(255);" json:"sia_name"`
	PsaId         string        `gorm:"type:varchar(100);" json:"psa_id"`
	PsaName       string        `gorm:"type:varchar(255);" json:"psa_name"`
	Sipa          string        `gorm:"type:varchar(100);" json:"sipa"`
	SipaName      string        `gorm:"type:varchar(255);" json:"sipa_name"`
	ApingId       string        `gorm:"type:varchar(100);" json:"aping_id"`
	ApingName     string        `gorm:"type:varchar(255);" json:"aping_name"`
	BankName      string        `gorm:"type:varchar(255);" json:"bank_name"`
	AccountName   string        `gorm:"type:varchar(255);" json:"account_name"`
	AccountNumber string        `gorm:"type:varchar(100);" json:"account_number"`
	TaxPercentage int           `gorm:"type:int;default:0" json:"tax_percentage"`
	JournalMethod JournalMethod `gorm:"type:journal_method; default:'automatic'" json:"journal_method" validate:"required"`
	BranchStatus  DataStatus    `gorm:"type:data_status;default:'inactive'" json:"branch_status"`
	LicenseDate   time.Time     `gorm:"not null" json:"license_date" validate:"required"`
}

// Profile Struct
type ProfileStruct struct {
	UserID        string        `gorm:"type:varchar(15);primaryKey" json:"useri_id" validate:"required"`
	ProfileName   string        `gorm:"type:varchar(255);not null" json:"profile_name" validate:"required"`
	BranchID      string        `gorm:"type:varchar(15);primaryKey" json:"branch_id" validate:"required"`
	BranchName    string        `gorm:"unique;not null" json:"branch_name"`
	Address       string        `gorm:"type:text;" json:"address"`
	Phone         string        `gorm:"type:varchar(100);" json:"phone"`
	Email         string        `gorm:"type:varchar(100);" json:"email"`
	SiaId         string        `gorm:"type:varchar(100);" json:"sia_id"`
	SiaName       string        `gorm:"type:varchar(255);" json:"sia_name"`
	PsaId         string        `gorm:"type:varchar(100);" json:"psa_id"`
	PsaName       string        `gorm:"type:varchar(255);" json:"psa_name"`
	Sipa          string        `gorm:"type:varchar(100);" json:"sipa"`
	SipaName      string        `gorm:"type:varchar(255);" json:"sipa_name"`
	ApingId       string        `gorm:"type:varchar(100);" json:"aping_id"`
	ApingName     string        `gorm:"type:varchar(255);" json:"aping_name"`
	BankName      string        `gorm:"type:varchar(255);" json:"bank_name"`
	AccountName   string        `gorm:"type:varchar(255);" json:"account_name"`
	AccountNumber string        `gorm:"type:varchar(100);" json:"account_number"`
	TaxPercentage int           `gorm:"type:int;default:0" json:"tax_percentage"`
	JournalMethod JournalMethod `gorm:"type:journal_method; default:'automatic'" json:"journal_method" validate:"required"`
	BranchStatus  DataStatus    `gorm:"type:data_status;default:'inactive'" json:"branch_status"`
	LicenseDate   time.Time     `gorm:"not null" json:"license_date" validate:"required"`
}
