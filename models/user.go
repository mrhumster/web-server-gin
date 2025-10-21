package models

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrhumster/web-server-gin/dto/request"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login        *string `gorm:"uniqueIndex;not null" json:"login"`
	Email        *string `gorm:"uniqueIndex;not nul" json:"email"`
	PasswordHash *string `gorm:"not null" json:"-"`
	Name         *string `gorm:"" json:"name"`
	LastName     *string `gorm:"" json:"last_name"`
	Role         *string `gorm:"" json:"role"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) SetPassword(password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hash := string(hashedBytes)
	u.PasswordHash = &hash
	return nil
}

func (u *User) CheckPassword(password string) bool {
	hashedBytes := []byte(*u.PasswordHash)
	err := bcrypt.CompareHashAndPassword(hashedBytes, []byte(password))
	return err == nil
}

func (u *User) FillInTheRequest(r request.UserRequest) {
	u.Login = &r.Login
	u.Email = &r.Email
	u.Name = &r.Name
	u.LastName = &r.LastName
	u.SetPassword(r.Password)
}

func (u *User) FillInTheUpdateRequest(r request.UpdateUserRequest) {
	if r.Name != nil {
		u.Name = r.Name
	}

	if r.LastName != nil {
		u.LastName = r.LastName
	}

	if r.Email != nil {
		u.Email = r.Email
	}
}

func (u *User) Debug() {
	log.Printf("👤 Login: %s", *u.Login)
	log.Printf("\tEmail: %s", *u.Email)
	log.Printf("\tPasswordHash: %s", *u.PasswordHash)
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
