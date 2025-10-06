package database

import (
	"github.com/mrhumster/web-server-gin/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetupDatabase(t *testing.T) {
	t.Run("should create database connection with valid config", func(t *testing.T) {
		cfg := config.TestConfig()

		assert.NotPanics(t, func() {
			db := SetupDatabase(cfg)
			require.NotNil(t, db)

			sqlDB, err := db.DB()
			if assert.NoError(t, err) {
				defer sqlDB.Close()
				assert.NoError(t, sqlDB.Ping())
			}
		})
	})
}
