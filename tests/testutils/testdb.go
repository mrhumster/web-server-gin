package testutils

import (
	"fmt"
	"log"

	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var TestDB *gorm.DB

func InitTestDB() *gorm.DB {
	cfg := config.TestConfig()
	var err error

	TestDB, err := gorm.Open(postgres.Open(cfg.GetDsn()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Printf("%s", cfg.GetDsn())
		log.Fatalf("🔴 Failed to connect test database")
	}
	log.Printf("🟢 Connected to test database successfully")

	err = TestDB.AutoMigrate(
		&models.User{},
		&models.Album{},
	)
	if err != nil {
		log.Fatalf("🔴 Failed apply migrations")
	}
	log.Printf("🟢 Migrations apply")
	return TestDB
}

func SetupTestDB() *gorm.DB {
	if TestDB == nil {
		TestDB = InitTestDB()
	}
	CleanTestDatabase()
	return TestDB
}

func CleanTestDatabase() {
	if TestDB == nil {
		return
	}

	TestDB.Exec("SET session_replication_role = 'replica';")

	tables := []string{"users", "albums"}
	for _, table := range tables {
		TestDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table))
	}

	TestDB.Exec("SET session_replication_role = 'origin';")
}

func TeardownTestDatabase() {
	if TestDB != nil {
		sqlDB, err := TestDB.DB()
		if err != nil {
			log.Printf("Error getting SQL DB: %v", err)
			return
		}
		sqlDB.Close()
		TestDB = nil
	}
}

func GetTestDB() *gorm.DB {
	if TestDB == nil {
		return SetupTestDB()
	}
	return TestDB
}
