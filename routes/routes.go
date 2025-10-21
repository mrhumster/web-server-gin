package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/handler"
	"github.com/mrhumster/web-server-gin/middleware"
	"github.com/mrhumster/web-server-gin/repository"
	"github.com/mrhumster/web-server-gin/service"
	"gorm.io/gorm"
)

func AddPolicyIfNotExists(sub, obj, act string, enforcer *casbin.Enforcer) {
	if hasPolicy, _ := enforcer.HasPolicy(sub, obj, act); !hasPolicy {
		enforcer.AddPolicy(sub, obj, act)
	}
}

func SetupRoutes(db *gorm.DB, mode string) *gin.Engine {
	if mode == "test" {
		gin.SetMode(gin.TestMode)
	}
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	cfg, _ := config.LoadConfig()
	if mode == "test" || mode == "debug" {
		cfg, _ = config.TestConfig()
	}

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	enforcer, err := casbin.NewEnforcer(cfg.Server.CasbinModel, adapter)
	if err != nil {
		log.Printf("⚠️ Casbin Load Error, %s", err.Error())
		panic("⚠️ Error loading roles config")
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, enforcer)
	tokenService, err := service.NewTokenService(&cfg.JWT)
	if err != nil {
		fmt.Printf("⚠️ SetupRoutes: %v", err)
		panic("Error create new token service")
	}
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService, tokenService, cfg.JwtSecret)
	commonHandler := handler.NewCommonHandler(tokenService)

	AddPolicyIfNotExists("admin", "*", "*", enforcer)
	AddPolicyIfNotExists("member", "users", "read", enforcer)
	AddPolicyIfNotExists("*", "users", "read", enforcer)

	r.POST("/api/login", authHandler.Login)
	r.GET("/api/logout", authHandler.Logout)
	r.POST("/api/users", userHandler.CreateUser)
	r.GET("/api/auth/public-key", commonHandler.GetPublicKey)

	auth := r.Group("/api", middleware.AuthMiddleware(tokenService))
	{
		auth.GET("/users", middleware.Authorize("users", "read", enforcer), userHandler.ReadUsers)
		auth.GET("/users/:id", middleware.Authorize("users", "read", enforcer), userHandler.ReadUser)
		auth.PATCH("/users/:id", middleware.Authorize("users", "write", enforcer), userHandler.Update)
		auth.DELETE("/users/:id", middleware.Authorize("users", "delete", enforcer), userHandler.Delete)
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
