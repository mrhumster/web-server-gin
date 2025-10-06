package database

import (
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase(cfg config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDsn()), &gorm.Config{})

	if err != nil {
		panic("⚠️ GORM not open DB")

	}
	db.AutoMigrate(&models.User{})
	return db
}
