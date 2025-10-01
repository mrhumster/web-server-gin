package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login    string `gorm:"uniqueIndex;not null" json:"login"`
	Password string `gorm:"not null" json:"password"`
	Email    string `gotm:"uniqueIndex;not nul" json:"email"`
}
