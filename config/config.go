package config

import "os"

type Config struct {
	ServerAddress string
	DatabaseURL   string
}

func LoadConfig() Config {
	return Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
	}
}
