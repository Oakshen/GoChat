package router

import (
	"context"
	"gochat/internal/config"
	"gochat/internal/handlers"
	"gochat/internal/middleware"
	ws "gochat/internal/websocket"
	"io/ioutil"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// SetupRouter 设置路由
func SetupRouter(addr string, wsHub *ws.Hub) *server.Hertz {
	h := server.Default(server.WithHostPorts(addr))

	// 根据Hertz WebSocket官方文档设置
	h.NoHijackConnPool = true

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(&cfg.JWT)
	userHandler := handlers.NewUserHandler()
	roomHandler := handlers.NewRoomHandler()
	wsHandler := handlers.NewWebSocketHandler(wsHub)

	// API 路由组
	api := h.Group("/api")

	// 公开路由
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	// WebSocket 连接路由（需要token验证）
	h.GET("/ws", wsHandler.HandleWebSocket)

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

	// 认证相关路由
	protected.GET("/auth/userinfo", authHandler.GetUserInfo)

	// 用户相关路由
	protected.GET("/users/:id", userHandler.GetProfile)
	protected.PUT("/users/:id", userHandler.UpdateProfile)
	protected.GET("/users/online", userHandler.GetOnlineUsers)
	protected.GET("/users/delete/:id", userHandler.DeleteUser)

	// 聊天室相关路由
	protected.POST("/rooms", roomHandler.CreateRoom)
	protected.GET("/rooms", roomHandler.GetRooms)
	protected.GET("/rooms/:id", roomHandler.GetRoom)
	protected.POST("/rooms/:id/join", roomHandler.JoinRoom)
	protected.POST("/rooms/:id/leave", roomHandler.LeaveRoom)
	protected.GET("/rooms/:id/members", roomHandler.GetRoomMembers)
	protected.GET("/rooms/:id/messages", roomHandler.GetRoomMessages)

	// WebSocket 管理路由
	protected.GET("/ws/stats", wsHandler.GetStats)
	protected.POST("/ws/broadcast/:roomId", wsHandler.BroadcastToRoom)

	// 静态文件服务 - 手动处理
	h.GET("/web/*filepath", func(ctx context.Context, c *app.RequestContext) {
		file := c.Param("filepath")
		if file == "" || file == "/" {
			file = "/index.html"
		}

		filePath := filepath.Join("./web", file)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			c.String(404, "File not found")
			return
		}

		// 设置内容类型
		if filepath.Ext(file) == ".html" {
			c.Header("Content-Type", "text/html; charset=utf-8")
		} else if filepath.Ext(file) == ".css" {
			c.Header("Content-Type", "text/css")
		} else if filepath.Ext(file) == ".js" {
			c.Header("Content-Type", "application/javascript")
		}

		c.Data(200, "text/html; charset=utf-8", content)
	})

	// 根路径重定向到web页面
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.Redirect(302, []byte("/web/index.html"))
	})

	// 健康检查
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{
			"status": "ok",
			"server": "GoChat Hertz with WebSocket",
		})
	})

	return h
}
