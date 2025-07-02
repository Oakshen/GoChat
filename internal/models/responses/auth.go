package responses

import "time"

// AuthResponse 认证响应结构
type AuthResponse struct {
	User  *UserInfo `json:"user"`
	Token string    `json:"token"`
}

// UserInfo 用户信息结构
type UserInfo struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// UserProfile 用户资料响应
type UserProfile struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	AvatarURL string     `json:"avatar_url"`
	IsOnline  bool       `json:"is_online"`
	LastSeen  *time.Time `json:"last_seen"`
	CreatedAt time.Time  `json:"created_at"`
}
