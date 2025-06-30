package entities

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	Username     string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email        string     `gorm:"uniqueIndex;size:100;not null" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	AvatarURL    string     `gorm:"size:255" json:"avatar_url"`
	IsOnline     bool       `gorm:"default:false" json:"is_online"`
	LastSeen     *time.Time `json:"last_seen"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

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
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	IsPrivate   bool      `gorm:"default:false" json:"is_private"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	Creator     User         `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	RoomMembers []RoomMember `gorm:"foreignKey:RoomID" json:"room_members,omitempty"`
	Messages    []Message    `gorm:"foreignKey:RoomID" json:"messages,omitempty"`
}

// TableName 指定表名
func (Room) TableName() string {
	return "rooms"
}

// Message 消息模型
type Message struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	RoomID      uint      `gorm:"not null;index" json:"room_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	MessageType string    `gorm:"size:20;default:'text'" json:"message_type"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联关系
	Room Room `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}

// RoomMember 聊天室成员模型
type RoomMember struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	RoomID   uint      `gorm:"not null;index" json:"room_id"`
	UserID   uint      `gorm:"not null;index" json:"user_id"`
	Role     string    `gorm:"size:20;default:'member'" json:"role"`
	JoinedAt time.Time `json:"joined_at"`

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
	// 检查是否已存在相同的房间-用户组合
	var count int64
	tx.Model(&RoomMember{}).Where("room_id = ? AND user_id = ?", rm.RoomID, rm.UserID).Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	rm.JoinedAt = time.Now()
	return nil
}
