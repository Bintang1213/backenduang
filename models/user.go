package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Nama     string `gorm:"type:varchar(100);not null" json:"nama"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"type:user_role_enum" json:"role"` 
}