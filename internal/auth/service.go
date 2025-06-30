package auth

import (
	"errors"
	"gochat/internal/config"
	"gochat/internal/database"
	"gochat/pkg/utils"
	"strings"
	"time"
)

// Service 认证服务结构
type Service struct {
	jwtSecret   string
	expireHours time.Duration
}

// NewService 创建认证服务实例
func NewService(cfg *config.JWTConfig) *Service {
	return &Service{
		jwtSecret:   cfg.Secret,
		expireHours: cfg.ExpireHours,
	}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

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

// Register 用户注册
func (s *Service) Register(req *RegisterRequest) (*AuthResponse, error) {
	// 检查用户名是否已存在
	var existingUser database.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 密码哈希
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := database.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, errors.New("用户创建失败")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, s.jwtSecret, s.expireHours)
	if err != nil {
		return nil, errors.New("token生成失败")
	}

	return &AuthResponse{
		User: &UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
		},
		Token: token,
	}, nil
}

// Login 用户登录
func (s *Service) Login(req *LoginRequest) (*AuthResponse, error) {
	var user database.User

	// 支持用户名或邮箱登录
	query := database.DB
	if strings.Contains(req.Username, "@") {
		query = query.Where("email = ?", req.Username)
	} else {
		query = query.Where("username = ?", req.Username)
	}

	if err := query.First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
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

	return &AuthResponse{
		User: &UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
		},
		Token: token,
	}, nil
}

// GetUserByToken 通过token获取用户信息
func (s *Service) GetUserByToken(tokenString string) (*UserInfo, error) {
	claims, err := utils.ParseToken(tokenString, s.jwtSecret)
	if err != nil {
		return nil, errors.New("无效的token")
	}

	var user database.User
	if err := database.DB.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	return &UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}
 