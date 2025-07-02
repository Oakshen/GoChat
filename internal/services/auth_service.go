package services

import (
	"errors"
	"gochat/internal/config"
	"gochat/internal/dal"
	"gochat/internal/models/entities"
	"gochat/internal/models/requests"
	"gochat/internal/models/responses"
	"gochat/pkg/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

// AuthService 认证服务结构
type AuthService struct {
	userDAL     *dal.UserDAL
	jwtSecret   string
	expireHours time.Duration
}

// NewAuthService 创建认证服务实例
func NewAuthService(cfg *config.JWTConfig) *AuthService {
	return &AuthService{
		userDAL:     dal.NewUserDAL(),
		jwtSecret:   cfg.Secret,
		expireHours: cfg.ExpireHours,
	}
}

// Register 用户注册
func (s *AuthService) Register(req *requests.RegisterRequest) (*responses.AuthResponse, error) {
	// 检查用户名是否已存在
	exists, err := s.userDAL.CheckUsernameExists(req.Username)
	if err != nil {
		return nil, errors.New("系统错误")
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	exists, err = s.userDAL.CheckEmailExists(req.Email)
	if err != nil {
		return nil, errors.New("系统错误")
	}
	if exists {
		return nil, errors.New("邮箱已存在")
	}

	// 密码哈希
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &entities.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.userDAL.Create(user); err != nil {
		return nil, errors.New("用户创建失败")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, s.jwtSecret, s.expireHours)
	if err != nil {
		return nil, errors.New("token生成失败")
	}

	return &responses.AuthResponse{
		User: &responses.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
		},
		Token: token,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(req *requests.LoginRequest) (*responses.AuthResponse, error) {
	var user *entities.User
	var err error

	// 支持用户名或邮箱登录
	if strings.Contains(req.Username, "@") {
		user, err = s.userDAL.GetByEmail(req.Username)
	} else {
		user, err = s.userDAL.GetByUsername(req.Username)
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("系统错误")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("密码错误")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, s.jwtSecret, s.expireHours)
	if err != nil {
		return nil, errors.New("token生成失败")
	}

	return &responses.AuthResponse{
		User: &responses.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
		},
		Token: token,
	}, nil
}

// GetUserByToken 通过token获取用户信息
func (s *AuthService) GetUserByToken(tokenString string) (*responses.UserInfo, error) {
	claims, err := utils.ParseToken(tokenString, s.jwtSecret)
	if err != nil {
		return nil, errors.New("无效的token")
	}

	user, err := s.userDAL.GetByID(claims.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("系统错误")
	}

	return &responses.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}
