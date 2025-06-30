package handlers

import (
	"context"
	"gochat/internal/config"
	"gochat/internal/models/requests"
	"gochat/internal/services"
	"gochat/pkg/response"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(cfg *config.JWTConfig) *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(cfg),
	}
}

// Register 用户注册
func (h *AuthHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req requests.RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error":   "参数验证失败",
			"details": err.Error(),
		})
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, resp)
}

// Login 用户登录
func (h *AuthHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req requests.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error":   "参数验证失败",
			"details": err.Error(),
		})
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, resp)
}

// GetUserInfo 获取当前用户信息
func (h *AuthHandler) GetUserInfo(ctx context.Context, c *app.RequestContext) {
	token := c.GetHeader("Authorization")
	if len(token) == 0 {
		c.JSON(http.StatusUnauthorized, utils.H{
			"error": "缺少Authorization头",
		})
		return
	}

	// 移除 "Bearer " 前缀
	if len(token) > 7 && string(token[:7]) == "Bearer " {
		token = token[7:]
	}

	userInfo, err := h.authService.GetUserByToken(string(token))
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, userInfo)
}
