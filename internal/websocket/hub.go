package websocket

import (
	"encoding/json"
	"gochat/internal/models/entities"
	"gochat/internal/services"
	"gochat/pkg/logger"
	"sync"
	"time"
)

// Hub WebSocket连接管理中心
type Hub struct {
	// 客户端连接管理
	clients    map[*Client]bool // 已注册的客户端
	register   chan *Client     // 客户端注册通道
	unregister chan *Client     // 客户端注销通道

	// 聊天室管理
	rooms     map[uint]map[*Client]bool // 聊天室中的客户端
	joinRoom  chan *JoinRoomMessage     // 加入聊天室通道
	leaveRoom chan *LeaveRoomMessage    // 离开聊天室通道

	// 消息处理
	handleMessage chan *ClientMessage    // 处理消息通道
	broadcast     chan *BroadcastMessage // 广播消息通道

	// 服务
	messageService *services.MessageService
	roomService    *services.RoomService
	userService    *services.UserService

	// 线程安全
	mutex sync.RWMutex
}

// JoinRoomMessage 加入聊天室消息
type JoinRoomMessage struct {
	Client *Client
	RoomID uint
}

// LeaveRoomMessage 离开聊天室消息
type LeaveRoomMessage struct {
	Client *Client
	RoomID uint
}

// ClientMessage 客户端消息
type ClientMessage struct {
	Client  *Client
	Message *WSMessage
}

// BroadcastMessage 广播消息
type BroadcastMessage struct {
	RoomID  uint
	Message *WSMessage
	Exclude *Client // 排除的客户端（通常是发送者）
}

// 公开的通道字段
var (
	Register      chan *Client
	Unregister    chan *Client
	JoinRoom      chan *JoinRoomMessage
	LeaveRoom     chan *LeaveRoomMessage
	HandleMessage chan *ClientMessage
	Broadcast     chan *BroadcastMessage
)

// NewHub 创建新的Hub实例
func NewHub() *Hub {
	hub := &Hub{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client, 10), // 添加缓冲
		unregister:     make(chan *Client, 10), // 添加缓冲
		rooms:          make(map[uint]map[*Client]bool),
		joinRoom:       make(chan *JoinRoomMessage, 100),   // 添加缓冲
		leaveRoom:      make(chan *LeaveRoomMessage, 100),  // 添加缓冲
		handleMessage:  make(chan *ClientMessage, 1000),    // 添加较大缓冲，处理消息
		broadcast:      make(chan *BroadcastMessage, 1000), // 添加较大缓冲，处理广播
		messageService: services.NewMessageService(),
		roomService:    services.NewRoomService(),
		userService:    services.NewUserService(),
	}

	// 设置公开通道
	Register = hub.register
	Unregister = hub.unregister
	JoinRoom = hub.joinRoom
	LeaveRoom = hub.leaveRoom
	HandleMessage = hub.handleMessage
	Broadcast = hub.broadcast

	return hub
}

// Run 启动Hub主循环
func (h *Hub) Run() {
	logger.Info("WebSocket Hub started")

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case joinMsg := <-h.joinRoom:
			logger.Info("Hub收到JoinRoomMessage:", joinMsg.Client.Username, "RoomID:", joinMsg.RoomID)
			h.handleJoinRoom(joinMsg)
			logger.Info("Hub处理JoinRoomMessage完成:", joinMsg.Client.Username, "RoomID:", joinMsg.RoomID)

		case leaveMsg := <-h.leaveRoom:
			h.handleLeaveRoom(leaveMsg)

		case clientMsg := <-h.handleMessage:
			logger.Info("Hub收到ClientMessage:", clientMsg.Client.Username, "消息类型:", clientMsg.Message.Type, "RoomID:", clientMsg.Message.RoomID)
			h.handleClientMessage(clientMsg)
			logger.Info("Hub处理ClientMessage完成:", clientMsg.Client.Username, "消息类型:", clientMsg.Message.Type)

		case broadcastMsg := <-h.broadcast:
			h.handleBroadcast(broadcastMsg)
		}
	}
}

// registerClient 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	h.clients[client] = true
	h.mutex.Unlock()

	logger.Info("Client registered:", client.Username, "UserID:", client.UserID)

	// 暂时不发送欢迎消息，避免复杂性
	// 等连接稳定后再考虑添加欢迎消息
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)

		// 从所有聊天室中移除
		for roomID := range h.rooms {
			if clients, exists := h.rooms[roomID]; exists {
				if _, inRoom := clients[client]; inRoom {
					delete(clients, client)
					// 通知其他用户
					h.notifyUserLeft(roomID, client)
				}
			}
		}
	}
	h.mutex.Unlock()

	logger.Info("Client unregistered:", client.Username, "UserID:", client.UserID)
}

// handleJoinRoom 处理加入聊天室
func (h *Hub) handleJoinRoom(joinMsg *JoinRoomMessage) {
	client := joinMsg.Client
	roomID := joinMsg.RoomID

	logger.Info("处理加入聊天室请求:", client.Username, "RoomID:", roomID)

	// 在单独的goroutine中处理数据库验证，避免阻塞Hub主循环
	go func() {
		// 验证用户是否有权限加入聊天室
		canJoin, err := h.roomService.CanUserJoinRoom(client.UserID, roomID)
		if err != nil || !canJoin {
			logger.Error("用户无权限加入聊天室:", client.Username, "RoomID:", roomID, "Error:", err)
			client.SendError("无权限加入聊天室", 403)
			return
		}

		logger.Info("用户有权限加入聊天室:", client.Username, "RoomID:", roomID)

		h.mutex.Lock()
		// 初始化聊天室
		if h.rooms[roomID] == nil {
			h.rooms[roomID] = make(map[*Client]bool)
			logger.Info("初始化聊天室:", roomID)
		}
		h.rooms[roomID][client] = true
		h.mutex.Unlock()

		logger.Info("用户已添加到聊天室:", client.Username, "RoomID:", roomID)

		// 发送加入成功消息
		joinSuccessMsg := &WSMessage{
			Type:      MessageTypeSystem,
			RoomID:    roomID,
			Content:   "已加入聊天室",
			Timestamp: time.Now(),
		}

		logger.Info("准备发送加入成功消息给用户:", client.Username)
		err = client.SendWSMessage(joinSuccessMsg)
		if err != nil {
			logger.Error("发送加入成功消息失败:", client.Username, "Error:", err)
		} else {
			logger.Info("加入成功消息已发送给用户:", client.Username)
		}

		// 通知其他用户
		logger.Info("通知其他用户有新用户加入:", client.Username, "RoomID:", roomID)
		h.notifyUserJoined(roomID, client)

		// 发送当前在线用户列表
		logger.Info("发送用户列表:", "RoomID:", roomID)
		h.sendUserList(roomID)
	}()
}

// handleLeaveRoom 处理离开聊天室
func (h *Hub) handleLeaveRoom(leaveMsg *LeaveRoomMessage) {
	client := leaveMsg.Client
	roomID := leaveMsg.RoomID

	h.mutex.Lock()
	if clients, exists := h.rooms[roomID]; exists {
		if _, inRoom := clients[client]; inRoom {
			delete(clients, client)
		}
	}
	h.mutex.Unlock()

	logger.Info("Client left room:", client.Username, "RoomID:", roomID)

	// 通知其他用户
	h.notifyUserLeft(roomID, client)

	// 更新用户列表
	h.sendUserList(roomID)
}

// handleClientMessage 处理客户端消息
func (h *Hub) handleClientMessage(clientMsg *ClientMessage) {
	client := clientMsg.Client
	message := clientMsg.Message

	logger.Info("开始处理客户端消息:", client.Username, "消息类型:", message.Type, "RoomID:", message.RoomID)

	switch message.Type {
	case MessageTypeText:
		logger.Info("处理文本消息:", client.Username)
		h.handleTextMessage(client, message)
	case MessageTypeJoin:
		logger.Info("处理加入聊天室消息:", client.Username, "RoomID:", message.RoomID)
		if message.RoomID > 0 {
			logger.Info("直接处理加入聊天室逻辑:", client.Username, "RoomID:", message.RoomID)
			// 直接处理加入聊天室，避免循环依赖
			h.handleJoinRoomDirect(client, message.RoomID)
		} else {
			logger.Error("无效的RoomID:", message.RoomID, "用户:", client.Username)
		}
	case MessageTypeLeave:
		logger.Info("处理离开聊天室消息:", client.Username, "RoomID:", message.RoomID)
		if message.RoomID > 0 {
			// 直接处理离开聊天室
			h.handleLeaveRoomDirect(client, message.RoomID)
		}
	case MessageTypeTyping:
		logger.Info("处理正在输入消息:", client.Username)
		h.handleTypingMessage(client, message)
	default:
		logger.Error("不支持的消息类型:", message.Type, "用户:", client.Username)
		client.SendError("不支持的消息类型", 400)
	}

	logger.Info("客户端消息处理完成:", client.Username, "消息类型:", message.Type)
}

// handleTextMessage 处理文本消息
func (h *Hub) handleTextMessage(client *Client, message *WSMessage) {
	if message.RoomID == 0 {
		client.SendError("聊天室ID不能为空", 400)
		return
	}

	if !client.IsInRoom(message.RoomID) {
		client.SendError("您不在此聊天室中", 403)
		return
	}

	// 保存消息到数据库
	dbMessage := &entities.Message{
		RoomID:      message.RoomID,
		UserID:      client.UserID,
		Content:     message.Content,
		MessageType: string(MessageTypeText),
	}

	savedMessage, err := h.messageService.CreateMessage(dbMessage)
	if err != nil {
		logger.Error("Failed to save message:", err)
		client.SendError("消息发送失败", 500)
		return
	}

	// 设置消息ID
	message.MessageID = savedMessage.ID

	// 广播消息到聊天室
	h.broadcast <- &BroadcastMessage{
		RoomID:  message.RoomID,
		Message: message,
		Exclude: nil, // 不排除任何人，包括发送者
	}

	logger.Info("Message sent:", client.Username, "RoomID:", message.RoomID, "Content:", message.Content)
}

// handleTypingMessage 处理正在输入消息
func (h *Hub) handleTypingMessage(client *Client, message *WSMessage) {
	if message.RoomID == 0 || !client.IsInRoom(message.RoomID) {
		return
	}

	// 广播正在输入状态（排除发送者）
	h.broadcast <- &BroadcastMessage{
		RoomID:  message.RoomID,
		Message: message,
		Exclude: client,
	}
}

// handleBroadcast 处理广播消息
func (h *Hub) handleBroadcast(broadcastMsg *BroadcastMessage) {
	h.mutex.RLock()
	clients := h.rooms[broadcastMsg.RoomID]
	h.mutex.RUnlock()

	if clients == nil {
		return
	}

	messageData, err := json.Marshal(broadcastMsg.Message)
	if err != nil {
		logger.Error("Failed to marshal broadcast message:", err)
		return
	}

	for client := range clients {
		if broadcastMsg.Exclude != nil && client == broadcastMsg.Exclude {
			continue
		}

		select {
		case client.Send <- messageData:
		default:
			// 客户端发送缓冲区已满，移除客户端
			delete(clients, client)
			close(client.Send)
		}
	}
}

// notifyUserJoined 通知用户加入
func (h *Hub) notifyUserJoined(roomID uint, client *Client) {
	message := &WSMessage{
		Type:      MessageTypeSystem,
		RoomID:    roomID,
		Content:   client.Username + " 加入了聊天室",
		Timestamp: time.Now(),
	}

	h.broadcast <- &BroadcastMessage{
		RoomID:  roomID,
		Message: message,
		Exclude: client,
	}
}

// notifyUserLeft 通知用户离开
func (h *Hub) notifyUserLeft(roomID uint, client *Client) {
	message := &WSMessage{
		Type:      MessageTypeSystem,
		RoomID:    roomID,
		Content:   client.Username + " 离开了聊天室",
		Timestamp: time.Now(),
	}

	h.broadcast <- &BroadcastMessage{
		RoomID:  roomID,
		Message: message,
		Exclude: nil,
	}
}

// sendUserList 发送用户列表
func (h *Hub) sendUserList(roomID uint) {
	h.mutex.RLock()
	clients := h.rooms[roomID]
	h.mutex.RUnlock()

	if clients == nil {
		return
	}

	var users []UserInfo
	for client := range clients {
		users = append(users, UserInfo{
			UserID:   client.UserID,
			Username: client.Username,
			IsOnline: true,
		})
	}

	userListMsg := &UserListMessage{
		Type:   MessageTypeUserList,
		RoomID: roomID,
		Users:  users,
	}

	messageData, err := json.Marshal(userListMsg)
	if err != nil {
		logger.Error("Failed to marshal user list:", err)
		return
	}

	for client := range clients {
		client.SendMessage(messageData)
	}
}

// GetStats 获取Hub统计信息
func (h *Hub) GetStats() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	stats := make(map[string]interface{})
	stats["total_clients"] = len(h.clients)
	stats["total_rooms"] = len(h.rooms)

	roomStats := make(map[uint]int)
	for roomID, clients := range h.rooms {
		roomStats[roomID] = len(clients)
	}
	stats["room_clients"] = roomStats

	return stats
}

// handleJoinRoomDirect 直接处理加入聊天室（避免循环依赖）
func (h *Hub) handleJoinRoomDirect(client *Client, roomID uint) {
	logger.Info("直接处理加入聊天室请求:", client.Username, "RoomID:", roomID)

	// 在单独的goroutine中处理数据库验证，避免阻塞Hub主循环
	go func() {
		// 验证用户是否有权限加入聊天室
		canJoin, err := h.roomService.CanUserJoinRoom(client.UserID, roomID)
		if err != nil || !canJoin {
			logger.Error("用户无权限加入聊天室:", client.Username, "RoomID:", roomID, "Error:", err)
			client.SendError("无权限加入聊天室", 403)
			return
		}

		logger.Info("用户有权限加入聊天室:", client.Username, "RoomID:", roomID)

		// 更新客户端本地房间列表
		client.mutex.Lock()
		client.Rooms[roomID] = true
		client.mutex.Unlock()

		h.mutex.Lock()
		// 初始化聊天室
		if h.rooms[roomID] == nil {
			h.rooms[roomID] = make(map[*Client]bool)
			logger.Info("初始化聊天室:", roomID)
		}
		h.rooms[roomID][client] = true
		h.mutex.Unlock()

		logger.Info("用户已添加到聊天室:", client.Username, "RoomID:", roomID)

		// 发送加入成功消息
		joinSuccessMsg := &WSMessage{
			Type:      MessageTypeSystem,
			RoomID:    roomID,
			Content:   "已加入聊天室",
			Timestamp: time.Now(),
		}

		logger.Info("准备发送加入成功消息给用户:", client.Username)
		err = client.SendWSMessage(joinSuccessMsg)
		if err != nil {
			logger.Error("发送加入成功消息失败:", client.Username, "Error:", err)
		} else {
			logger.Info("加入成功消息已发送给用户:", client.Username)
		}

		// 通知其他用户
		logger.Info("通知其他用户有新用户加入:", client.Username, "RoomID:", roomID)
		h.notifyUserJoined(roomID, client)

		// 发送当前在线用户列表
		logger.Info("发送用户列表:", "RoomID:", roomID)
		h.sendUserList(roomID)
	}()
}

// handleLeaveRoomDirect 直接处理离开聊天室（避免循环依赖）
func (h *Hub) handleLeaveRoomDirect(client *Client, roomID uint) {
	logger.Info("直接处理离开聊天室请求:", client.Username, "RoomID:", roomID)

	// 更新客户端本地房间列表
	client.mutex.Lock()
	delete(client.Rooms, roomID)
	client.mutex.Unlock()

	h.mutex.Lock()
	if clients, exists := h.rooms[roomID]; exists {
		if _, inRoom := clients[client]; inRoom {
			delete(clients, client)
		}
	}
	h.mutex.Unlock()

	logger.Info("用户已从聊天室移除:", client.Username, "RoomID:", roomID)

	// 通知其他用户
	h.notifyUserLeft(roomID, client)

	// 更新用户列表
	h.sendUserList(roomID)
}
