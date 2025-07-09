package websocket

import "time"

// MessageType 消息类型枚举
type MessageType string

const (
	MessageTypeText     MessageType = "text"     // 文本消息
	MessageTypeImage    MessageType = "image"    // 图片消息
	MessageTypeFile     MessageType = "file"     // 文件消息
	MessageTypeVideo    MessageType = "video"    // 视频消息
	MessageTypeJoin     MessageType = "join"     // 加入聊天室
	MessageTypeLeave    MessageType = "leave"    // 离开聊天室
	MessageTypeTyping   MessageType = "typing"   // 正在输入
	MessageTypeSystem   MessageType = "system"   // 系统消息
	MessageTypeUserList MessageType = "userlist" // 用户列表更新
	MessageTypeError    MessageType = "error"    // 错误消息
	MessageTypePing     MessageType = "ping"     // 心跳ping消息
	MessageTypePong     MessageType = "pong"     // 心跳pong响应
)

// WSMessage WebSocket 消息结构
type WSMessage struct {
	Type        MessageType      `json:"type"`                  // 消息类型
	RoomID      uint             `json:"room_id,omitempty"`     // 聊天室ID
	UserID      uint             `json:"user_id,omitempty"`     // 用户ID
	Username    string           `json:"username,omitempty"`    // 用户名
	Content     string           `json:"content,omitempty"`     // 消息内容
	Timestamp   time.Time        `json:"timestamp"`             // 时间戳
	MessageID   uint             `json:"message_id,omitempty"`  // 消息ID（用于持久化）
	Attachments []AttachmentInfo `json:"attachments,omitempty"` // 附件信息
}

// AttachmentInfo 附件信息结构
type AttachmentInfo struct {
	ID       uint   `json:"id"`                 // 附件ID
	FileName string `json:"file_name"`          // 文件名
	FileSize int64  `json:"file_size"`          // 文件大小
	FileType string `json:"file_type"`          // MIME类型
	Category string `json:"category"`           // 文件分类：image, file, video
	Width    int    `json:"width,omitempty"`    // 图片/视频宽度
	Height   int    `json:"height,omitempty"`   // 图片/视频高度
	Duration int    `json:"duration,omitempty"` // 视频时长（秒）
	URL      string `json:"url"`                // 访问URL
}

// JoinRoomRequest 加入聊天室请求
type JoinRoomRequest struct {
	RoomID uint   `json:"room_id"` // 聊天室ID
	Token  string `json:"token"`   // JWT Token
}

// ErrorMessage 错误消息
type ErrorMessage struct {
	Type    MessageType `json:"type"`
	Error   string      `json:"error"`
	Code    int         `json:"code,omitempty"`
	Details string      `json:"details,omitempty"`
}

// SystemMessage 系统消息
type SystemMessage struct {
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
	RoomID  uint        `json:"room_id,omitempty"`
}

// UserListMessage 用户列表消息
type UserListMessage struct {
	Type   MessageType `json:"type"`
	RoomID uint        `json:"room_id"`
	Users  []UserInfo  `json:"users"`
}

// UserInfo 用户信息
type UserInfo struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url,omitempty"`
	IsOnline  bool   `json:"is_online"`
}
