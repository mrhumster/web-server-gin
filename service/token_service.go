package service

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/models"
)

type TokenService struct {
	accessPrivateKey  *rsa.PrivateKey
	accessPublicKey   *rsa.PublicKey
	refreshPrivateKey *rsa.PrivateKey
	refreshPublicKey  *rsa.PublicKey
	accessExpiry      time.Duration
	refreshExpiry     time.Duration
	issuer            string
}

func NewTokenService(cfg *config.JWT) (*TokenService, error) {
	accessPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.AccessPrivateKey))
	if err != nil {
		return nil, fmt.Errorf("parse access private key: %w", err)
	}
	accessPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.AccessPublicKey))
	if err != nil {
		return nil, fmt.Errorf("parse access public key: %w", err)
	}
	refreshPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.RefreshPrivateKey))
	if err != nil {
		return nil, fmt.Errorf("parse refresh private key: %w", err)
	}
	refreshPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.RefreshPublicKey))
	if err != nil {
		return nil, fmt.Errorf("parse refresh public key: %w", err)
	}
	return &TokenService{
		accessPrivateKey:  accessPrivateKey,
		accessPublicKey:   accessPublicKey,
		refreshPrivateKey: refreshPrivateKey,
		refreshPublicKey:  refreshPublicKey,
		accessExpiry:      cfg.AccessTokenExpiry,
		refreshExpiry:     cfg.RefreshTokenExpiry,
		issuer:            cfg.Issuer,
	}, nil
}

func (s *TokenService) GenerateToken(user *models.User) (*models.TokenPair, error) {
	accessExpiresAt := time.Now().Add(s.accessExpiry)
	accessClaims := &models.AccessClaims{
		UserID:   fmt.Sprintf("%d", user.ID),
		Username: *user.Login,
		Role:     *user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.accessPrivateKey)
	if err != nil {
		return nil, err
	}
	refreshExpiresAt := time.Now().Add(s.refreshExpiry)
	refreshClaims := &models.RefreshClaims{
		UserID:       fmt.Sprintf("%d", user.ID),
		TokenVersion: *user.TokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.refreshPrivateKey)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.accessExpiry.Seconds()),
		TokenType:    "bearer",
	}, nil
}

func (s *TokenService) ValidateAccessToken(tokenString string) (*models.AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.accessPublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*models.AccessClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (s *TokenService) ValidateRefreshToken(tokenString string) (*models.RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.refreshPublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*models.RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (s *TokenService) GetAccessPublicKey() *rsa.PublicKey {
	return s.accessPublicKey
}

func (s *TokenService) GetPublicKeyPEM() (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(s.refreshPublicKey)
	if err != nil {
		return "", err
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM), nil
}

func (s *TokenService) GetRefreshExpiry() time.Duration {
	return s.refreshExpiry
}

func (s *TokenService) GetAccessExpiry() time.Duration {
	return s.accessExpiry
}
