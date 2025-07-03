package main

import (
	"fmt"
	"gochat/internal/config"
	"gochat/internal/database"
	"gochat/internal/router"
	ws "gochat/internal/websocket"
	"gochat/pkg/logger"
	"log"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化日志
	logger.Init(&cfg.Log)
	logger.Info("Starting GoChat server...")

	// 连接数据库
	if err := database.Connect(&cfg.Database); err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}
	logger.Info("Database connected successfully")

	// 自动迁移数据库表结构
	if err := database.AutoMigrate(); err != nil {
		logger.Fatal("Failed to migrate database:", err)
	}
	logger.Info("Database migration completed")

	// 初始化WebSocket Hub
	wsHub := ws.NewHub()
	go wsHub.Run()
	logger.Info("WebSocket Hub started")

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Info("Server starting on", serverAddr)

	// 设置路由并启动服务器
	h := router.SetupRouter(serverAddr, wsHub)
	h.Spin()
}
