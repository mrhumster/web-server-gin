package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/database"
	"github.com/mrhumster/web-server-gin/routes"
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
