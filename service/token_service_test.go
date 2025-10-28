package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/stretchr/testify/assert"
)

func TestTokenService_GenerateAndValidateToken(t *testing.T) {
	cfg, _ := config.TestConfig()
	service, _ := NewTokenService(&cfg.JWT)

	login := "testuser"
	email := "testuser@test.local"
	role := "member"
	tokenVersion := "1"
	user := &models.User{
		Login:        &login,
		Email:        &email,
		Role:         &role,
		TokenVersion: &tokenVersion,
	}

	token, err := service.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := service.ValidateAccessToken(token.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%d", user.ID), claims.UserID)
	assert.Equal(t, *user.Login, claims.Username)
	assert.Equal(t, "auth-service", claims.Issuer)
}

func TestTokenService_ValidateToken_Invalid(t *testing.T) {
	cfg, _ := config.TestConfig()
	service, _ := NewTokenService(&cfg.JWT)

	_, err := service.ValidateAccessToken("invalid-token")
	assert.Error(t, err)

	claims := &models.AccessClaims{
		UserID:   "123",
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			Issuer:    "test-issuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	expiredToken, _ := token.SignedString([]byte("test-secret"))

	_, err = service.ValidateAccessToken(expiredToken)
	assert.Error(t, err)
}
