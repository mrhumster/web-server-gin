package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetEnv(t *testing.T) {
	key := "NOT_EXISTS"
	v := getEnv(key, "exist")
	assert.Equal(t, v, "exist")
	v = getEnv("USER", "user")
	assert.NotEqual(t, v, "user")
}

func TestTestConfig(t *testing.T) {
	cfg := TestConfig()
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Database.Host)
	assert.NotEmpty(t, cfg.Database.Name)
	assert.NotEmpty(t, cfg.Server.ServerAddr)
	dsn := cfg.GetDsn()
	assert.NotEmpty(t, dsn)
}

func TestLoadConfig_WithEnvRestore(t *testing.T) {
	originalEnv := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASS":     os.Getenv("DB_PASS"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"SERVER_ADDR": os.Getenv("SERVER_ADDR"),
	}

	t.Cleanup(func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	})

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "myapp")
	os.Setenv("SERVER_ADDR", ":8080")

	cfg := LoadConfig()

	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "myapp", cfg.Database.Name)
	assert.Equal(t, ":8080", cfg.Server.ServerAddr)
	assert.Equal(t, "disable", cfg.Database.SslMode)
	assert.Equal(t, "UTC", cfg.Database.TimeZone)
}
