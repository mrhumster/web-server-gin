package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/config"
	"github.com/mrhumster/web-server-gin/routes"
)

func main() {
	r := gin.Default()
	cfg := config.LoadConfig()
	routes.SetupRoutes(r, cfg)
	r.Run(cfg.ServerAddress)
}
