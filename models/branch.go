package models

import "gorm.io/gorm"

type Branch struct {
	gorm.Model
	BrancheName   string `gorm:"unique;not null" json:"branche_name"`
	BrancheStatus int    `gorm:"not null" json:"branche_status" default:"0"`
}
