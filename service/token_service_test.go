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

	user := &models.User{
		Login: strPtr("testuser"),
		Email: strPtr("test@example.com"),
	}

	token, err := service.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%d", user.ID), claims.UserID)
	assert.Equal(t, *user.Login, claims.Username)
	assert.Equal(t, *user.Email, claims.Email)
	assert.Equal(t, "auth-service", claims.Issuer)
}

func TestTokenService_ValidateToken_Invalid(t *testing.T) {
	cfg, _ := config.TestConfig()
	service, _ := NewTokenService(&cfg.JWT)

	_, err := service.ValidateToken("invalid-token")
	assert.Error(t, err)

	claims := &models.Claims{
		UserID:   "123",
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			Issuer:    "test-issuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	expiredToken, _ := token.SignedString([]byte("test-secret"))

	_, err = service.ValidateToken(expiredToken)
	assert.Error(t, err)
}
