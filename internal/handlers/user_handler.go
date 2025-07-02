package handlers

import (
	"context"
	"gochat/internal/models/requests"
	"gochat/internal/services"
	"gochat/pkg/response"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的用户ID",
		})
		return
	}

	profile, err := h.userService.GetProfile(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, profile)
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的用户ID",
		})
		return
	}

	var req requests.UpdateProfileRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error":   "参数验证失败",
			"details": err.Error(),
		})
		return
	}

	profile, err := h.userService.UpdateProfile(uint(userID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, profile)
}

// GetOnlineUsers 获取在线用户列表
func (h *UserHandler) GetOnlineUsers(ctx context.Context, c *app.RequestContext) {
	profiles, err := h.userService.GetOnlineUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, profiles)
}

func (h *UserHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的用户ID",
		})
		return
	}

	err = h.userService.DeleteUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, "用户已删除")
}
