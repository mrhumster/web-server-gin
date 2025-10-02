package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
