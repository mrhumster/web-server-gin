package config

import (
	"fmt"
	"os"
	"time"
)

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SslMode  string
	TimeZone string
}

type Server struct {
	ServerAddr  string
	JwtSecret   string
	CasbinModel string
}

type JWT struct {
	AccessPrivateKey   string
	AccessPublicKey    string
	RefreshPrivateKey  string
	RefreshPublicKey   string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string `mapstructure:"jwt_issuer"`
}
type Config struct {
	Database
	Server `mapstructure:"server"`
	JWT    `mapstructure:"jwt"`
}

func LoadConfig() (*Config, error) {

	accessTokenExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_EXPIRY", "15m"))
	if err != nil {
		return nil, fmt.Errorf("Config error. Plase set ENV JWT_ACCESS_TOKEN_EXPIRY. %v", err)
	}
	refreshTokenExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_EXPIRY", "7d"))
	if err != nil {
		return nil, fmt.Errorf("Config error. Plase set ENV JWT_REFRESH_TOKEN_EXPIRY. %v", err)
	}

	cfg := &Config{
		Database: Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Name:     os.Getenv("DB_NAME"),
			SslMode:  "disable",
			TimeZone: "UTC",
		},
		Server: Server{
			ServerAddr:  os.Getenv("SERVER_ADDR"),
			JwtSecret:   os.Getenv("JWT_SECRET"),
			CasbinModel: os.Getenv("CASBIN_MODEL"),
		},
		JWT: JWT{
			AccessPrivateKey:   getEnv("JWT_ACCESS_PRIVATE_KEY", ""),
			AccessPublicKey:    getEnv("JWT_ACCESS_PUBLIC_KEY", ""),
			RefreshPrivateKey:  getEnv("JWT_REFRESH_PRIVATE_KEY", ""),
			RefreshPublicKey:   getEnv("JWT_REFRESH_PUBLIC_KEY", ""),
			AccessTokenExpiry:  accessTokenExpiry,
			RefreshTokenExpiry: refreshTokenExpiry,
			Issuer:             getEnv("JWT_ISSUER", "auth-service"),
		},
	}
	return cfg, nil
}

func (config *Config) GetDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Host,
		config.User,
		config.Password,
		config.Name,
		config.Port,
		config.SslMode,
		config.TimeZone)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func TestConfig() (*Config, error) {
	accessTokenExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_EXPIRY", "15m"))
	if err != nil {
		return nil, fmt.Errorf("Config error. Plase set ENV JWT_ACCESS_TOKEN_EXPIRY. %v", err)
	}
	refreshTokenExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_EXPIRY", "168h"))
	if err != nil {
		return nil, fmt.Errorf("Config error. Plase set ENV JWT_REFRESH_TOKEN_EXPIRY. %v", err)
	}

	accessPrivateKey, err := os.ReadFile("../config/keys/accessPrivate.pem")
	if err != nil {
		return nil, fmt.Errorf("Config error. AccessPrivateKey not read.")
	}

	accessPublicKey, err := os.ReadFile("../config/keys/accessPublic.pem")
	if err != nil {
		return nil, fmt.Errorf("Config error. AccessPublicKey not read.")
	}

	refreshPrivateKey, err := os.ReadFile("../config/keys/refreshPrivate.pem")
	if err != nil {
		return nil, fmt.Errorf("Config error. RefreshPrivateKey not read.")
	}

	refreshPublicKey, err := os.ReadFile("../config/keys/refreshPublic.pem")
	if err != nil {
		return nil, fmt.Errorf("Config error. RefreshPublicKey not read.")
	}

	return &Config{
		Database: Database{
			Host:     getEnv("TEST_DB_HOST", "localhost"),
			Port:     getEnv("TEST_DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASS", "Master1234"),
			Name:     getEnv("TEST_DB_NAME", "test_database1"),
			SslMode:  "disable",
			TimeZone: "UTC",
		},
		Server: Server{
			ServerAddr:  getEnv("TEST_SERVER_ADDR", ":8080"),
			JwtSecret:   getEnv("TEST_JWT_SECRET", "jwt-secret-jwt-secret"),
			CasbinModel: getEnv("TEST_CASBIN_MODEL", "../config/model.conf"),
		},
		JWT: JWT{
			AccessPrivateKey:   string(accessPrivateKey),
			AccessPublicKey:    string(accessPublicKey),
			RefreshPrivateKey:  string(refreshPrivateKey),
			RefreshPublicKey:   string(refreshPublicKey),
			AccessTokenExpiry:  accessTokenExpiry,
			RefreshTokenExpiry: refreshTokenExpiry,
			Issuer:             getEnv("JWT_ISSUER", "auth-service"),
		},
	}, nil
}
