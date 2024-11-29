package models

import "gorm.io/gorm"

type Branch struct {
	gorm.Model
	BranchName   string `gorm:"unique;not null" json:"branch_name"`
	BranchStatus int    `gorm:"not null" json:"branch_status" default:"0"`
}
