package models

import (
	"github.com/mrhumster/web-server-gin/dto/request"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login        string `gorm:"uniqueIndex;not null" json:"login"`
	Email        string `gorm:"uniqueIndex;not nul" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	Name         string `gorm:"" json:"name"`
	LastName     string `gorm:"" json:"last_name"`
}

func (u *User) SetPassword(password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedBytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) FillInTheRequest(r request.UserRequest) {
	u.Login = r.Login
	u.Email = r.Email
	u.Name = r.Name
	u.LastName = r.LastName
	u.SetPassword(r.Password)
}

func (u *User) FillInTheUpdateRequest(r request.UpdateUserRequest) {
	if r.Name != "" {
		u.Name = r.Name
	}

	if r.LastName != "" {
		u.LastName = r.LastName
	}

	if r.Email != "" {
		u.Email = r.Email
	}
}
