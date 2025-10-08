package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/handler"
	"github.com/mrhumster/web-server-gin/middleware"
	"github.com/mrhumster/web-server-gin/repository"
	"github.com/mrhumster/web-server-gin/service"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, mode string) *gin.Engine {
	if mode == "test" {
		gin.SetMode(gin.TestMode)
	}
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()

	cfg := config.LoadConfig()
	if mode == "test" || mode == "debug" {
		cfg = config.TestConfig()
	}
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService, cfg.JwtSecret)

	r.POST("/api/login", authHandler.Login)
	r.GET("/api/logout", authHandler.Logout)
	r.POST("/api/users", userHandler.CreateUser)

	auth := r.Group("/api")

	auth.Use(middleware.AuthMiddleware(cfg.JwtSecret))
	{
		auth.GET("/users", userHandler.ReadUsers)
		auth.GET("/users/:id", userHandler.ReadUser)
		auth.PATCH("/users/:id", userHandler.Update)
		auth.DELETE("/users/:id", userHandler.Delete)
	}

	r.GET("/health", func(c *gin.Context) {
		if _, err := db.DB(); err != nil {
			log.Println("⚠️ PG ERROR: ", err.Error())
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})
	return r
}
