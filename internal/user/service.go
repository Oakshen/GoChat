package user

import (
	"errors"
	"gochat/internal/database"
	"time"
)

// Service 用户服务结构
type Service struct{}

// NewService 创建用户服务实例
func NewService() *Service {
	return &Service{}
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
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

// GetProfile 获取用户资料
func (s *Service) GetProfile(userID uint) (*UserProfile, error) {
	var user database.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	return &UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		IsOnline:  user.IsOnline,
		LastSeen:  user.LastSeen,
		CreatedAt: user.CreatedAt,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *Service) UpdateProfile(userID uint, req *UpdateProfileRequest) (*UserProfile, error) {
	var user database.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新字段
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, errors.New("更新失败")
	}

	return &UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		IsOnline:  user.IsOnline,
		LastSeen:  user.LastSeen,
		CreatedAt: user.CreatedAt,
	}, nil
}

// SetOnlineStatus 设置用户在线状态
func (s *Service) SetOnlineStatus(userID uint, isOnline bool) error {
	now := time.Now()
	updates := map[string]interface{}{
		"is_online": isOnline,
	}

	if !isOnline {
		updates["last_seen"] = &now
	}

	return database.DB.Model(&database.User{}).Where("id = ?", userID).Updates(updates).Error
}

// GetOnlineUsers 获取在线用户列表
func (s *Service) GetOnlineUsers() ([]UserProfile, error) {
	var users []database.User
	if err := database.DB.Where("is_online = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}

	var profiles []UserProfile
	for _, user := range users {
		profiles = append(profiles, UserProfile{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			IsOnline:  user.IsOnline,
			LastSeen:  user.LastSeen,
			CreatedAt: user.CreatedAt,
		})
	}

	return profiles, nil
}

func (s *Service) CreateUsers(userName string)
