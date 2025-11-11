package routes

import (
	"testing"

	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/internal/database"
	"github.com/mrhumster/web-server-gin/pkg/auth"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRoutes(t *testing.T) {
	cfg, _ := config.TestConfig()
	db := database.SetupDatabase(cfg)
	permissionGRPCClient, _ := auth.NewPermissionClient(cfg.Server.AuthServiceAddr)
	defer permissionGRPCClient.Close()
	SetupRoutes(db, "test", permissionGRPCClient)
	assert.IsType(t, &gorm.DB{}, db)
}
