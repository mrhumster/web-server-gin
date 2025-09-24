package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/handler"
	"github.com/mrhumster/web-server-gin/repository"
	"github.com/mrhumster/web-server-gin/service"
	"log"
	"net/http"
)

func SetupRoutes(r *gin.Engine, cfg config.Config) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		panic("🔴 Unable to connect to database: " + err.Error())
	}

	albumRepo := repository.NewAlbumRepository(dbpool)
	albumService := service.NewAlbumService(albumRepo)
	albumHandler := handler.NewAlbumHandler(albumService)

	r.POST("/albums", albumHandler.CreateAlbum)
	r.GET("/albums/:id", albumHandler.GetAlbumByID)

	r.GET("/health", func(c *gin.Context) {
		if err := dbpool.Ping(context.Background()); err != nil {
			log.Println("⚠️ PG ERROR: ", err.Error())
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})

	return dbpool
}
