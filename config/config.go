package config

import (
	"fmt"
	"os"
)

// 	dsn := "host=postgresql user=postgres password=Master1234 dbname=database1 port=5432 sslmode=disable TimeZone=Asia/Omsk"

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
