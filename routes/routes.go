package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/handler"
	"github.com/mrhumster/web-server-gin/repository"
	"github.com/mrhumster/web-server-gin/service"
)

func SetupRoutes(r *gin.Engine, cfg config.Config) {
	db, err := pgx.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		panic("🔴 Unable to connect to database: " + err.Error())
	}
	defer db.Close(context.Background())

	albumRepo := repository.NewAlbumRepository(db)
	albumService := service.NewAlbumService(albumRepo)
	albumHandler := handler.NewAlbumHandler(albumService)

	r.POST("/albums", albumHandler.CreateAlbum)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})
}
