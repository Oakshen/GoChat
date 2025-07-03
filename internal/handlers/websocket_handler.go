package handlers

import (
	"context"
	"gochat/pkg/logger"
	"gochat/pkg/utils"
	"net/http"
	"strconv"

	"gochat/internal/config"
	ws "gochat/internal/websocket"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
)

// WebSocketHandler WebSocket处理器
type WebSocketHandler struct {
	hub      *ws.Hub
	upgrader websocket.HertzUpgrader
}

// NewWebSocketHandler 创建WebSocket处理器实例
func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
		upgrader: websocket.HertzUpgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(c *app.RequestContext) bool {
				// 允许所有来源，简化测试
				return true
			},
		},
	}
}

// HandleWebSocket 处理WebSocket升级请求
func (h *WebSocketHandler) HandleWebSocket(ctx context.Context, c *app.RequestContext) {
	// 从查询参数或Header中获取用户信息
	token := string(c.Query("token"))
	if token == "" {
		token = string(c.GetHeader("Authorization"))
		if len(token) > 7 && string(token[:7]) == "Bearer " {
			token = token[7:]
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "缺少认证token",
		})
		return
	}

	// 从配置中获取JWT密钥
	cfg, err := config.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "服务器配置错误",
		})
		return
	}

	// 验证token并获取用户信息
	claims, err := utils.ParseToken(token, cfg.JWT.Secret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "无效的token",
		})
		return
	}

	// 升级为WebSocket连接
	logger.Info("正在为用户升级WebSocket连接:", claims.Username)

	err = h.upgrader.Upgrade(c, func(conn *websocket.Conn) {
		logger.Info("WebSocket升级成功，用户:", claims.Username)

		// 创建客户端
		client := ws.NewClient(claims.UserID, claims.Username, conn, h.hub)
		logger.Info("WebSocket客户端已创建:", claims.Username)

		// 注册客户端到Hub
		ws.Register <- client
		logger.Info("客户端已注册到Hub:", claims.Username)

		// 启动WritePump（发送消息循环）
		go client.WritePump()
		logger.Info("WritePump已启动:", claims.Username)

		// 启动ReadPump（读取消息循环）
		// 这个函数会阻塞直到连接关闭
		client.ReadPump()
		logger.Info("ReadPump已退出，连接关闭:", claims.Username)
	})

	if err != nil {
		logger.Error("WebSocket升级失败:", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "WebSocket升级失败",
		})
		return
	}
}

// GetStats 获取WebSocket统计信息
func (h *WebSocketHandler) GetStats(ctx context.Context, c *app.RequestContext) {
	stats := h.hub.GetStats()
	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// BroadcastToRoom 向指定聊天室广播消息
func (h *WebSocketHandler) BroadcastToRoom(ctx context.Context, c *app.RequestContext) {
	roomIDStr := c.Param("roomId")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "无效的聊天室ID",
		})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
		Type    string `json:"type"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "参数验证失败",
		})
		return
	}

	// 创建系统消息
	message := &ws.WSMessage{
		Type:    ws.MessageTypeSystem,
		RoomID:  uint(roomID),
		Content: req.Content,
	}

	// 广播消息
	ws.Broadcast <- &ws.BroadcastMessage{
		RoomID:  uint(roomID),
		Message: message,
		Exclude: nil,
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "消息已广播",
	})
}
