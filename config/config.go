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
	ServerAddr string
	JwtSecret  string
}

type Config struct {
	Database
	Server
}

func LoadConfig() Config {
	return Config{
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
			ServerAddr: os.Getenv("SERVER_ADDR"),
			JwtSecret:  os.Getenv("JWT_SECRET"),
		},
	}
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

func TestConfig() Config {
	return Config{
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
			ServerAddr: getEnv("TEST_SERVER_ADDR", ":8080"),
			JwtSecret:  getEnv("TEST_JWT_SECRET", "jwt-secret-jwt-secret"),
		},
	}

}
