package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/internal/delivery/http/handler"
	"github.com/mrhumster/web-server-gin/internal/delivery/http/middleware"
	"github.com/mrhumster/web-server-gin/internal/repository"
	"github.com/mrhumster/web-server-gin/internal/service"
	"github.com/mrhumster/web-server-gin/pkg/auth"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, mode string, permissionClient auth.PermissionClient) *gin.Engine {

	// MODE
	if mode == "test" {
		gin.SetMode(gin.TestMode)
	}
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	// GIN ROUTE
	r := gin.Default()

	// CONFIGURATION
	cfg, _ := config.LoadConfig()
	if mode == "test" || mode == "debug" {
		cfg, _ = config.TestConfig()
	}

	// DATABSE
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	// CASBIN RULES
	enforcer, err := casbin.NewEnforcer(cfg.Server.CasbinModel, adapter)
	if err != nil {
		log.Printf("⚠️ Casbin Load Error, %s", err.Error())
		panic("⚠️ Error loading roles config")
	}

	// REPOSITORIES
	userRepo := repository.NewUserRepository(db)
	permissionService := service.NewPermissionService(enforcer)

	// SERVICES
	userService := service.NewUserService(userRepo, permissionService)
	tokenService, err := service.NewTokenService(&cfg.JWT)
	if err != nil {
		fmt.Printf("⚠️ SetupRoutes: %v", err)
		panic("Error create new token service")
	}

	// HANDLERS
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService, tokenService, cfg.JwtSecret, cfg.Server.Domain)
	commonHandler := handler.NewCommonHandler(tokenService)

	// PERMISSIONS

	type permission struct {
		sub string
		obj string
		act string
	}

	policies := []struct {
		params permission
		desc   string
	}{
		{permission{"*", "users", "read"}, "users read policy"},
		{permission{"*", "stream", "read"}, "stream read policy"},
		{permission{"*", "stream", "write"}, "stream write policy"},
	}

	for _, policy := range policies {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		_, err := permissionClient.AddPolicyIfNotExists(ctx,
			policy.params.sub,
			policy.params.obj,
			policy.params.act,
		)

		cancel()

		if err != nil {
			log.Printf("⚠️ Failed to add %s: %v", policy.desc, err)
		} else {
			log.Printf("✅ Successfully added %s", policy.desc)
		}
	}

	// ROUTE
	r.POST("/api/login", authHandler.Login)
	r.GET("/api/logout", authHandler.Logout)
	r.POST("/api/users", userHandler.CreateUser)
	r.POST("/api/refresh", authHandler.Refresh)

	auth := r.Group("/api", middleware.AuthMiddleware(tokenService))
	{
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/logout-all", authHandler.LogoutAll)
		auth.GET("/users", middleware.Authorize("users", "read", permissionClient), userHandler.ReadUsers)
		auth.GET("/users/:id", middleware.Authorize("users", "read", permissionClient), userHandler.ReadUser)
		auth.PATCH("/users/:id", middleware.Authorize("users", "write", permissionClient), userHandler.Update)
		auth.DELETE("/users/:id", middleware.Authorize("users", "delete", permissionClient), userHandler.Delete)
	}

	r.GET("/api/auth/public-key", commonHandler.GetPublicKey)
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
