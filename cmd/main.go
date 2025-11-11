package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/internal/database"
	"github.com/mrhumster/web-server-gin/internal/delivery/http/routes"
	"github.com/mrhumster/web-server-gin/internal/permission"
	"github.com/mrhumster/web-server-gin/internal/service"
	permissionpb "github.com/mrhumster/web-server-gin/proto/gen/go/permission"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("❌ Config: %s", err.Error()))
	}
	db := database.SetupDatabase(cfg)
	r := routes.SetupRoutes(db, "qa")

	defer func() {
		log.Println("🟡 Closing database pool...")
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("failed to get sql.DB: %s", err.Error())
		}
		if err := sqlDB.Close(); err != nil {
			log.Println("🟢 Database pool closed")
		}
	}()

	httpErr := make(chan error, 1)
	grpcErr := make(chan error, 1)

	srv := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server starting on %s\n", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("🔴 Server error: ", err)
			httpErr <- err
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("🔴 Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()

		adapter, err := gormadapter.NewAdapterByDB(db)
		if err != nil {
			panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
		}

		enforcer, err := casbin.NewEnforcer(cfg.Server.CasbinModel, adapter)
		if err != nil {
			log.Printf("⚠️ Casbin Load Error, %s", err.Error())
			panic("⚠️ Error loading roles config")
		}

		permissionService := service.NewPermissionService(enforcer)
		permissionServer := permission.NewPermissionGRPCServer(permissionService)
		permissionpb.RegisterPermissionServiceServer(grpcServer, permissionServer)
		log.Printf("🛰️ gRPC server listened at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("🔴 Failed to serve: %v", err)
			grpcErr <- err
		}

	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🟡 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("🔴 Server shutdown error: ", err)
	}

	log.Println("🟢 Server stoped")
}
