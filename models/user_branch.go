package models

import (
	"time"

	"gorm.io/gorm"
)

type UserBranch struct {
	UserID    int64 `gorm:"primaryKey"`
	BranchID  int64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete

	// Relasi ke User dan Branch (tidak perlu dikaitkan secara eksplisit)
	User   User   `gorm:"foreignKey:UserID"`
	Branch Branch `gorm:"foreignKey:BranchID"`
}
