package websocket

import (
	"encoding/json"
	"gochat/pkg/logger"
	"strconv"
	"sync"
	"time"

	"github.com/hertz-contrib/websocket"
)

// Client WebSocket 客户端连接
type Client struct {
	ID       string          // 客户端唯一标识
	UserID   uint            // 用户ID
	Username string          // 用户名
	Conn     *websocket.Conn // WebSocket 连接
	Send     chan []byte     // 发送消息通道
	Hub      *Hub            // 连接管理中心
	Rooms    map[uint]bool   // 用户加入的聊天室
	mutex    sync.RWMutex    // 读写锁
}

const (
	// 时间配置
	writeWait      = 10 * time.Second    // 写入超时时间
	pongWait       = 60 * time.Second    // pong响应等待时间
	pingPeriod     = (pongWait * 9) / 10 // ping发送周期
	maxMessageSize = 512                 // 最大消息大小
)

// NewClient 创建新的客户端连接
func NewClient(userID uint, username string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:       generateClientID(userID),
		UserID:   userID,
		Username: username,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Hub:      hub,
		Rooms:    make(map[uint]bool),
	}
}

// JoinRoom 加入聊天室
func (c *Client) JoinRoom(roomID uint) {
	logger.Info("客户端JoinRoom开始:", c.Username, "RoomID:", roomID)

	c.mutex.Lock()
	c.Rooms[roomID] = true
	c.mutex.Unlock()

	logger.Info("客户端本地Rooms已更新:", c.Username, "RoomID:", roomID)

	// 通知管理中心
	logger.Info("准备发送JoinRoomMessage到Hub:", c.Username, "RoomID:", roomID)
	JoinRoom <- &JoinRoomMessage{
		Client: c,
		RoomID: roomID,
	}
	logger.Info("JoinRoomMessage已发送到Hub:", c.Username, "RoomID:", roomID)
}

// LeaveRoom 离开聊天室
func (c *Client) LeaveRoom(roomID uint) {
	c.mutex.Lock()
	delete(c.Rooms, roomID)
	c.mutex.Unlock()

	// 通知管理中心
	LeaveRoom <- &LeaveRoomMessage{
		Client: c,
		RoomID: roomID,
	}
}

// IsInRoom 检查是否在聊天室中
func (c *Client) IsInRoom(roomID uint) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Rooms[roomID]
}

// GetRooms 获取用户加入的所有聊天室
func (c *Client) GetRooms() []uint {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	rooms := make([]uint, 0, len(c.Rooms))
	for roomID := range c.Rooms {
		rooms = append(rooms, roomID)
	}
	return rooms
}

// SendMessage 发送消息到客户端
func (c *Client) SendMessage(message []byte) {
	select {
	case c.Send <- message:
		// 消息发送成功
	default:
		// 发送缓冲区已满，记录警告但不立即关闭连接
		logger.Error("Client send buffer is full, dropping message for user:", c.Username)
		// 只有在无法发送重要消息时才关闭连接
		// close(c.Send)
		// Unregister <- c
	}
}

// SendWSMessage 发送WebSocket消息
func (c *Client) SendWSMessage(msg *WSMessage) error {
	logger.Info("准备发送WebSocket消息给客户端:", c.Username, "消息类型:", msg.Type, "内容:", msg.Content)

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("JSON编组失败:", c.Username, "Error:", err)
		return err
	}

	logger.Info("消息JSON编组成功:", c.Username, "数据:", string(data))

	// 使用defer recover防止panic
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic in SendWSMessage for client:", c.Username, "panic:", r)
		}
	}()

	// 检查channel是否已关闭
	select {
	case c.Send <- data:
		// 消息发送成功
		logger.Info("消息已放入发送通道:", c.Username)
		return nil
	default:
		// channel可能已关闭或缓冲区满
		logger.Error("Failed to send message to client:", c.Username, "- channel may be closed or full")
		return err
	}
}

// SendError 发送错误消息
func (c *Client) SendError(errorMsg string, code int) {
	errMsg := &ErrorMessage{
		Type:  MessageTypeError,
		Error: errorMsg,
		Code:  code,
	}

	data, err := json.Marshal(errMsg)
	if err != nil {
		logger.Error("Failed to marshal error message:", err)
		return
	}

	c.SendMessage(data)
}

// ReadPump 读取消息循环
func (c *Client) ReadPump() {
	defer func() {
		logger.Info("ReadPump exiting for client:", c.Username)
		// 离开所有聊天室
		for roomID := range c.Rooms {
			c.LeaveRoom(roomID)
		}
		Unregister <- c
		c.Conn.Close()
	}()

	logger.Info("ReadPump starting for client:", c.Username)

	// 最简化设置，只设置消息大小限制
	c.Conn.SetReadLimit(maxMessageSize)

	logger.Info("ReadPump ready for client:", c.Username)

	for {
		// 直接读取消息，不设置任何超时
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			logger.Info("ReadMessage error for client:", c.Username, "error:", err)
			break
		}

		logger.Info("Received message from client", c.Username, ":", string(message))

		// 解析消息
		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			logger.Error("Failed to unmarshal message:", err)
			c.SendError("Invalid message format", 400)
			continue
		}

		// 设置发送者信息
		wsMsg.UserID = c.UserID
		wsMsg.Username = c.Username
		wsMsg.Timestamp = time.Now()

		logger.Info("准备发送消息到HandleMessage通道:", c.Username, "消息类型:", wsMsg.Type, "RoomID:", wsMsg.RoomID)

		// 处理消息
		HandleMessage <- &ClientMessage{
			Client:  c,
			Message: &wsMsg,
		}

		logger.Info("消息已发送到HandleMessage通道:", c.Username, "消息类型:", wsMsg.Type)
	}
}

// WritePump 写入消息循环
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		logger.Info("WritePump exiting for client:", c.Username)
		ticker.Stop()
		c.Conn.Close()
	}()

	logger.Info("WritePump started for client:", c.Username)

	for {
		select {
		case message, ok := <-c.Send:
			logger.Info("WritePump收到消息:", c.Username, "数据长度:", len(message))

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了通道
				logger.Info("Send channel closed for client:", c.Username)
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			logger.Info("发送消息到WebSocket:", c.Username, "消息内容:", string(message))

			// 直接发送单个消息，不合并多个消息
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logger.Error("Failed to write message for client", c.Username, ":", err)
				return
			}

			logger.Info("消息发送完成:", c.Username)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Error("Failed to send ping to client", c.Username, ":", err)
				return
			}
		}
	}
}

// generateClientID 生成客户端ID
func generateClientID(userID uint) string {
	return time.Now().Format("20060102150405") + "_" + strconv.Itoa(int(userID))
}
