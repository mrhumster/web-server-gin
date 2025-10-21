package config

import (
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
	cfg, _ := TestConfig()
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Database.Host)
	assert.NotEmpty(t, cfg.Database.Name)
	assert.NotEmpty(t, cfg.Server.ServerAddr)
	dsn := cfg.GetDsn()
	assert.NotEmpty(t, dsn)
}
