package router

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// SetupRouter 初始化和配置所有路由
func SetupRouter(addr string) *server.Hertz {
	// 创建路由器，并配置服务器地址
	h := server.Default(server.WithHostPorts(addr))

	// 添加中间件
	h.Use(recovery.Recovery())

	// 添加CORS中间件
	h.Use(func(ctx context.Context, c *app.RequestContext) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if string(c.Method()) == "OPTIONS" {
			c.AbortWithStatus(consts.StatusNoContent)
			return
		}

		c.Next(ctx)
	})

	// 设置基础路由
	setupBaseRoutes(h)

	// 设置API路由
	setupAPIRoutes(h)

	return h
}

// setupBaseRoutes 设置基础路由
func setupBaseRoutes(h *server.Hertz) {
	// 健康检查接口
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{
			"status":  "ok",
			"message": "GoChat server is running",
		})
	})
}

// setupAPIRoutes 设置API路由组
func setupAPIRoutes(h *server.Hertz) {
	// API路由组
	api := h.Group("/api")
	{
		// 公开接口
		api.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
			c.JSON(consts.StatusOK, utils.H{
				"message": "pong",
			})
		})
		api.POST("/test")

		// 认证相关路由
		auth := api.Group("/auth")
		{
			// TODO: 实现注册登录路由
			auth.POST("/register")
			// auth.POST("/login", handlers.Login)
			// auth.POST("/logout", middleware.AuthMiddleware(), handlers.Logout)
			// auth.POST("/refresh", handlers.RefreshToken)

			// 临时占位路由
			auth.GET("/info", func(ctx context.Context, c *app.RequestContext) {
				c.JSON(consts.StatusOK, utils.H{
					"message": "auth routes - coming soon",
				})
			})
		}

		// 用户相关路由
		users := api.Group("/users")
		{
			// TODO: 实现用户管理路由
			// users.GET("/profile", middleware.AuthMiddleware(), handlers.GetProfile)
			// users.PUT("/profile", middleware.AuthMiddleware(), handlers.UpdateProfile)
			// users.GET("/online", handlers.GetOnlineUsers)

			// 临时占位路由
			users.GET("/info", func(ctx context.Context, c *app.RequestContext) {
				c.JSON(consts.StatusOK, utils.H{
					"message": "user routes - coming soon",
				})
			})
		}

		// 聊天相关路由
		chat := api.Group("/chat")
		{
			// TODO: 实现聊天相关路由
			// chat.GET("/rooms", middleware.AuthMiddleware(), handlers.GetRooms)
			// chat.POST("/rooms", middleware.AuthMiddleware(), handlers.CreateRoom)
			// chat.GET("/rooms/:id/messages", middleware.AuthMiddleware(), handlers.GetMessages)
			// chat.GET("/ws", handlers.WebSocketHandler)

			// 临时占位路由
			chat.GET("/info", func(ctx context.Context, c *app.RequestContext) {
				c.JSON(consts.StatusOK, utils.H{
					"message": "chat routes - coming soon",
				})
			})
		}
	}
}
