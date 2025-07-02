package router

import (
	"context"
	"gochat/internal/config"
	"gochat/internal/handlers"
	"gochat/internal/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// SetupRouter 设置路由
func SetupRouter(addr string) *server.Hertz {
	h := server.Default(server.WithHostPorts(addr))

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(&cfg.JWT)
	userHandler := handlers.NewUserHandler()

	// API 路由组
	api := h.Group("/api")

	// 公开路由
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

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

	// 健康检查
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{
			"status": "ok",
			"server": "GoChat Hertz",
		})
	})

	return h
}
