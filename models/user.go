package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login    string `gorm:"uniqueIndex;not null" json:"login"`
	Password string `gorm:"not null" json:"password"`
}
