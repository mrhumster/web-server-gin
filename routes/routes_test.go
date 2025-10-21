package routes

import (
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestRoutes(t *testing.T) {
	cfg, _ := config.TestConfig()
	db := database.SetupDatabase(cfg)
	SetupRoutes(db, "test")
	assert.IsType(t, &gorm.DB{}, db)
}
