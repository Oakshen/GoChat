package entities

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email        string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	AvatarURL    string         `gorm:"size:255" json:"avatar_url"`
	IsOnline     bool           `gorm:"default:false" json:"is_online"`
	LastSeen     *time.Time     `json:"last_seen"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联关系
	CreatedRooms []Room       `gorm:"foreignKey:CreatedBy" json:"created_rooms,omitempty"`
	RoomMembers  []RoomMember `gorm:"foreignKey:UserID" json:"room_members,omitempty"`
	Messages     []Message    `gorm:"foreignKey:UserID" json:"messages,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Room 聊天室模型
type Room struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	IsPrivate   bool           `gorm:"default:false" json:"is_private"`
	CreatedBy   uint           `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联关系
	Creator     User         `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	RoomMembers []RoomMember `gorm:"foreignKey:RoomID" json:"room_members,omitempty"`
	Messages    []Message    `gorm:"foreignKey:RoomID" json:"messages,omitempty"`
}

// TableName 指定表名
func (Room) TableName() string {
	return "rooms"
}

// Attachment 附件模型
type Attachment struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	MessageID uint      `gorm:"not null;index" json:"message_id"`
	FileName  string    `gorm:"size:255;not null" json:"file_name"`  // 原始文件名
	FilePath  string    `gorm:"size:500;not null" json:"file_path"`  // 存储路径
	FileSize  int64     `gorm:"not null" json:"file_size"`           // 文件大小（字节）
	FileType  string    `gorm:"size:100;not null" json:"file_type"`  // MIME类型
	Category  string    `gorm:"size:20;not null" json:"category"`    // image, file, video
	Width     int       `gorm:"default:0" json:"width,omitempty"`    // 图片/视频宽度
	Height    int       `gorm:"default:0" json:"height,omitempty"`   // 图片/视频高度
	Duration  int       `gorm:"default:0" json:"duration,omitempty"` // 视频时长（秒）
	CreatedAt time.Time `json:"created_at"`

	// 关联关系
	Message Message `gorm:"foreignKey:MessageID" json:"message,omitempty"`
}

// TableName 指定表名
func (Attachment) TableName() string {
	return "attachments"
}

// Message 消息模型
type Message struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	RoomID      uint      `gorm:"not null;index" json:"room_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Content     string    `gorm:"type:text" json:"content"`                   // 文本内容，多媒体消息可为空
	MessageType string    `gorm:"size:20;default:'text'" json:"message_type"` // text, image, file, video
	CreatedAt   time.Time `json:"created_at"`

	// 关联关系
	Room        Room         `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	User        User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Attachments []Attachment `gorm:"foreignKey:MessageID" json:"attachments,omitempty"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}

// RoomMember 聊天室成员模型
type RoomMember struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	RoomID    uint           `gorm:"not null;index" json:"room_id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Role      string         `gorm:"size:20;default:'member'" json:"role"`
	JoinedAt  time.Time      `json:"joined_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联关系
	Room Room `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (RoomMember) TableName() string {
	return "room_members"
}

// BeforeCreate 为 RoomMember 设置唯一索引
func (rm *RoomMember) BeforeCreate(tx *gorm.DB) error {
	// 检查是否已存在相同的房间-用户组合（排除已删除的记录）
	var count int64
	tx.Model(&RoomMember{}).Where("room_id = ? AND user_id = ?", rm.RoomID, rm.UserID).Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	rm.JoinedAt = time.Now()
	return nil
}
