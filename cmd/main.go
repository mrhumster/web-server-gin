package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/routes"
)

func main() {
	r := gin.Default()
	cfg := config.LoadConfig()
	dbpools := routes.SetupRoutes(r, cfg)

	defer func() {
		log.Println("🟡 Closing database pool...")
		dbpools.Close()
		log.Println("🟢 Database pool closed")
	}()

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	go func() {
		log.Printf("🚀 Serever starting on %s\n", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("🔴 Server error: ", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🟡 Shutting down server...")

	// Gracefull shutdown - даем завершить текущие запросы
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("🔴 Server shutdown error: ", err)
	}

	log.Println("🟢 Server stoped")
}
