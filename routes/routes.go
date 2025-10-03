package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/handler"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/repository"
	"github.com/mrhumster/web-server-gin/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, cfg config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDsn()), &gorm.Config{})

	if err != nil {
		panic("⚠️ GORM not open DB")

	}

	db.AutoMigrate(&models.Album{})
	db.AutoMigrate(&models.User{})

	albumRepo := repository.NewAlbumRepository(db)
	albumService := service.NewAlbumService(albumRepo)
	albumHandler := handler.NewAlbumHandler(albumService)

	r.POST("/albums", albumHandler.CreateAlbum)
	r.GET("/albums/:id", albumHandler.GetAlbumByID)
	r.DELETE("/albums/:id", albumHandler.DeleteAlbumByID)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r.GET("/users", userHandler.ReadUsers)
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/:id", userHandler.ReadUser)
	r.PATCH("/users/:id", userHandler.Update)
	r.DELETE("/users/:id", userHandler.Delete)

	r.GET("/health", func(c *gin.Context) {
		if _, err := db.DB(); err != nil {
			log.Println("⚠️ PG ERROR: ", err.Error())
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})

	return db
}
