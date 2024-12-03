package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Mendefinisikan custom type untuk ENUM StatusUser
type DataStatus string

const (
	Active   DataStatus = "active"
	Inactive DataStatus = "inactive"
)

// Mendefinisikan custom type untuk ENUM UserRole
type UserRole string

const (
	Operator      UserRole = "operator"
	Cashier       UserRole = "cashier"
	Finance       UserRole = "finance"
	Administrator UserRole = "administrator"
)

type User struct {
	gorm.Model
	Username   string     `gorm:"type:varchar(100);unique;not null" json:"username" validate:"required"`
	Password   string     `gorm:"type:text;not null" json:"password" validate:"required"`
	Name       string     `gorm:"type:varchar(255);not null" json:"name" validate:"required"`
	Address    string     `gorm:"type:text;" json:"address"`
	UserRole   UserRole   `gorm:"type:user_role;default:'operator'; not null" json:"user_role" validate:"required"`
	UserStatus DataStatus `gorm:"type:user_status;default:'inactive'" json:"user_status"`
}

// HashPassword function hashes the password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword function checks if the password is correct
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
