package models

type Menu struct {
	ID           uint     `gorm:"primaryKey" json:"id"`
	Name         string   `gorm:"type:varchar(255);not null" json:"name"`
	Route        string   `gorm:"type:varchar(255);not null" json:"route"`
	AllowedRoles []string `gorm:"type:text[];not null" json:"allowed_roles"` // Menyimpan daftar role sebagai array
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}
