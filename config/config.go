package config

import (
	"fmt"
	"os"
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
	PrivateKeyPath string `mapstructure:"jwt_private_key_path"`
	PublicKeyPath  string `mapstructure:"jwt_public_key_path"`
	Issuer         string `mapstructure:"jwt_issuer"`
}
type Config struct {
	Database
	Server `mapstructure:"server"`
	JWT    `mapstructure:"jwt"`
}

func LoadConfig() (*Config, error) {
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
			PrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "/app/config/keys/private.pem"),
			PublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH", "/app/config/keys/public.pem"),
			Issuer:         getEnv("JWT_ISSUER", "auth-service"),
		},
	}
	if _, err := os.Stat(cfg.JWT.PrivateKeyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("private key file not found: %s", cfg.JWT.PrivateKeyPath)
	}
	if _, err := os.Stat(cfg.JWT.PublicKeyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("public key file not found: %s", cfg.JWT.PublicKeyPath)
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
			PrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "../config/keys/private.pem"),
			PublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH", "../config/keys/public.pem"),
			Issuer:         getEnv("JWT_ISSUER", "auth-service"),
		},
	}, nil
}
