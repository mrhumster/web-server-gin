package service

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/models"
)

type TokenService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
}

func NewTokenService(jwt *config.JWT) (*TokenService, error) {
	privateKey, err := loadPrivateKey(jwt.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("NewTokenSrvice: error load private key: %w", err)
	}
	publicKey, err := loadPublicKey(jwt.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("NewTokenService: error load public key: %w", err)
	}
	return &TokenService{
		privateKey: privateKey,
		publicKey:  publicKey,
		issuer:     jwt.Issuer,
	}, nil
}

func (s *TokenService) GenerateToken(user *models.User) (string, error) {
	claims := &models.Claims{
		UserID:   fmt.Sprintf("%d", user.ID),
		Username: *user.Login,
		Email:    *user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	fmt.Printf("🔐 Generating token with method: %v\n", token.Method)
	fmt.Printf("🔐 Private key: %v\n", s.privateKey != nil)
	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		fmt.Printf("❌ Error signing token: %v\n", err)
		return "", err
	}

	fmt.Printf("✅ Token generated successfully\n")
	return tokenString, nil
}

func (s *TokenService) ValidateToken(tokenString string) (*models.Claims, error) {
	fmt.Printf("🔍 Validating token...\n")
	fmt.Printf("🔍 Public key available: %v\n", s.publicKey != nil)
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		fmt.Printf("🔍 Token signing method: %v\n", token.Method)
		fmt.Printf("🔍 Expected method: %v\n", jwt.SigningMethodRS256)
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		fmt.Printf("❌ Token validation error: %v\n", err)
		return nil, err
	}
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		fmt.Printf("✅ Token valid for user: %s\n", claims.Username)
		return claims, nil
	}
	fmt.Printf("❌ Invalid token claims\n")
	return nil, errors.New("invalid token")
}

func (s *TokenService) GetPublicKeyPEM() (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(s.publicKey)
	if err != nil {
		return "", err
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM), nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not RSA private key")
		}
		return rsaKey, nil
	}
	return privateKey, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return publicKey, nil
}
