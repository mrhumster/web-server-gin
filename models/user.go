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
	u.SetPassword(r.Password)
}
