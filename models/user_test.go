package models

import (
	"testing"

	"github.com/mrhumster/web-server-gin/dto/request"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

func TestUser_SetPassword(t *testing.T) {
	user := User{}
	password := "mySecretPassword123"
	err := user.SetPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotEqual(t, user.PasswordHash, password)
}

func TestUser_CheckPassword(t *testing.T) {
	user := User{}
	password := "mySecretPassword123"
	err := user.SetPassword(password)
	assert.NoError(t, err)
	isValid := user.CheckPassword(password)
	assert.True(t, isValid)
	isValid = user.CheckPassword("myWrongPassword")
	assert.False(t, isValid)
}

func TestUser_FillTest(t *testing.T) {
	user := User{}
	req := request.UserRequest{
		Login:    "testuser",
		Password: "password",
	}
	user.FillInTheRequest(req)
	assert.NotEqual(t, *user.PasswordHash, req.Password)
	assert.Equal(t, user.Login, &req.Login)
	updReq := request.UpdateUserRequest{
		LastName: strPtr("Billy"),
	}
	user.FillInTheUpdateRequest(updReq)
	assert.Equal(t, *user.LastName, *updReq.LastName)
	updReq = request.UpdateUserRequest{
		Name:  strPtr("Coat"),
		Email: strPtr("username@test.local"),
	}
	user.FillInTheUpdateRequest(updReq)
	assert.Equal(t, *user.Name, *updReq.Name)
	assert.Equal(t, *user.Email, *updReq.Email)
}
