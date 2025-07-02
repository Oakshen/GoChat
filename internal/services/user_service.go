package services

import (
	"errors"
	"gochat/internal/dal"
	"gochat/internal/models/requests"
	"gochat/internal/models/responses"

	"gorm.io/gorm"
)

// UserService 用户服务结构
type UserService struct {
	userDAL *dal.UserDAL
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userDAL: dal.NewUserDAL(),
	}
}

// GetProfile 获取用户资料
func (s *UserService) GetProfile(userID uint) (*responses.UserProfile, error) {
	user, err := s.userDAL.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("系统错误")
	}

	return &responses.UserProfile{
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
func (s *UserService) UpdateProfile(userID uint, req *requests.UpdateProfileRequest) (*responses.UserProfile, error) {
	user, err := s.userDAL.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("系统错误")
	}

	// 更新字段
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userDAL.Update(user); err != nil {
		return nil, errors.New("更新失败")
	}

	return &responses.UserProfile{
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
func (s *UserService) SetOnlineStatus(userID uint, isOnline bool) error {
	return s.userDAL.UpdateOnlineStatus(userID, isOnline)
}

// GetOnlineUsers 获取在线用户列表
func (s *UserService) GetOnlineUsers() ([]responses.UserProfile, error) {
	users, err := s.userDAL.GetOnlineUsers()
	if err != nil {
		return nil, errors.New("获取在线用户失败")
	}

	var profiles []responses.UserProfile
	for _, user := range users {
		profiles = append(profiles, responses.UserProfile{
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

func (s *UserService) DeleteUser(userID uint) error {
	// todo 是否需要获取整个entities.User整个用户信息

	if err := s.userDAL.SoftDeleteByID(userID); err != nil {
		return errors.New("删除用户失败")
	}

	return nil
}
