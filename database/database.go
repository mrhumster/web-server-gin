package database

import (
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func SetupDatabase(cfg config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDsn()), &gorm.Config{})
	if err != nil {
		panic("⚠️ GORM not open DB")

	}
	sqlDb, err := db.DB()
	if err != nil {
		panic("⚠️ GORM not open DB")
	}
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(time.Hour)
	sqlDb.SetConnMaxIdleTime(30 * time.Minute)
	db.AutoMigrate(&models.User{})
	return db
}
