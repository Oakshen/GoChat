package middleware

import (
	"context"
	"gochat/pkg/response"
	"gochat/pkg/utils"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtSecret string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从请求头获取token
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			response.Unauthorized(ctx, c, "缺少认证token")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(ctx, c, "token格式错误")
			c.Abort()
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			response.Unauthorized(ctx, c, "token不能为空")
			c.Abort()
			return
		}

		// 验证token
		claims, err := utils.ParseToken(tokenString, jwtSecret)
		if err != nil {
			response.Unauthorized(ctx, c, "无效的token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("token", tokenString)

		c.Next(ctx)
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *app.RequestContext) uint {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(uint)
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *app.RequestContext) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}
