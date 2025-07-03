package handlers

import (
	"context"
	"gochat/internal/models/entities"
	"gochat/internal/models/requests"
	"gochat/internal/services"
	"gochat/pkg/response"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// RoomHandler 聊天室处理器
type RoomHandler struct {
	roomService *services.RoomService
}

// NewRoomHandler 创建聊天室处理器实例
func NewRoomHandler() *RoomHandler {
	return &RoomHandler{
		roomService: services.NewRoomService(),
	}
}

// CreateRoom 创建聊天室
func (h *RoomHandler) CreateRoom(ctx context.Context, c *app.RequestContext) {
	var req requests.CreateRoomRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error":   "参数验证失败",
			"details": err.Error(),
		})
		return
	}

	// 从JWT token中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.H{
			"error": "用户未认证",
		})
		return
	}

	room := &entities.Room{
		Name:        req.Name,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
		CreatedBy:   userID.(uint),
	}

	createdRoom, err := h.roomService.CreateRoom(room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, createdRoom)
}

// GetRooms 获取聊天室列表
func (h *RoomHandler) GetRooms(ctx context.Context, c *app.RequestContext) {
	// 从JWT token中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.H{
			"error": "用户未认证",
		})
		return
	}

	rooms, err := h.roomService.GetUserRooms(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, rooms)
}

// GetRoom 获取聊天室详情
func (h *RoomHandler) GetRoom(ctx context.Context, c *app.RequestContext) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的聊天室ID",
		})
		return
	}

	room, err := h.roomService.GetRoomByID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, room)
}

// JoinRoom 加入聊天室
func (h *RoomHandler) JoinRoom(ctx context.Context, c *app.RequestContext) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的聊天室ID",
		})
		return
	}

	// 从JWT token中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.H{
			"error": "用户未认证",
		})
		return
	}

	err = h.roomService.JoinRoom(userID.(uint), uint(roomID))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, "成功加入聊天室")
}

// LeaveRoom 离开聊天室
func (h *RoomHandler) LeaveRoom(ctx context.Context, c *app.RequestContext) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的聊天室ID",
		})
		return
	}

	// 从JWT token中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.H{
			"error": "用户未认证",
		})
		return
	}

	err = h.roomService.LeaveRoom(userID.(uint), uint(roomID))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, "成功离开聊天室")
}

// GetRoomMembers 获取聊天室成员
func (h *RoomHandler) GetRoomMembers(ctx context.Context, c *app.RequestContext) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的聊天室ID",
		})
		return
	}

	members, err := h.roomService.GetRoomMembers(uint(roomID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, members)
}

// GetRoomMessages 获取聊天室消息
func (h *RoomHandler) GetRoomMessages(ctx context.Context, c *app.RequestContext) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的聊天室ID",
		})
		return
	}

	// 分页参数
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	messageService := services.NewMessageService()
	messages, err := messageService.GetMessagesByRoom(uint(roomID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"error": err.Error(),
		})
		return
	}

	response.Success(ctx, c, messages)
}
