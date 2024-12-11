package models

import (
	time "time"

	gorm "gorm.io/gorm"
)

type UserBranch struct {
	UserID    string `gorm:"type:varchar(15);primaryKey" json:"user_id" validate:"required"`
	BranchID  string `gorm:"type:varchar(15);primaryKey" json:"branch_id" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
