package dal

import (
	"gochat/internal/database"
	"gochat/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

// UserDAL 用户数据访问层
type UserDAL struct {
	db *gorm.DB
}

// NewUserDAL 创建用户数据访问层实例
func NewUserDAL() *UserDAL {
	return &UserDAL{
		db: database.DB,
	}
}

// GetByUsername 根据用户名查找用户
func (d *UserDAL) GetByUsername(username string) (*entities.User, error) {
	var user entities.User
	err := d.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱查找用户
func (d *UserDAL) GetByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := d.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID 根据ID查找用户
func (d *UserDAL) GetByID(id uint) (*entities.User, error) {
	var user entities.User
	err := d.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (d *UserDAL) Create(user *entities.User) error {
	return d.db.Create(user).Error
}

// Update 更新用户
func (d *UserDAL) Update(user *entities.User) error {
	return d.db.Save(user).Error
}

// UpdateOnlineStatus 更新用户在线状态
func (d *UserDAL) UpdateOnlineStatus(userID uint, isOnline bool) error {
	now := time.Now()
	updates := map[string]interface{}{
		"is_online": isOnline,
	}

	if !isOnline {
		updates["last_seen"] = &now
	}

	return d.db.Model(&entities.User{}).Where("id = ?", userID).Updates(updates).Error
}

// GetOnlineUsers 获取在线用户列表
func (d *UserDAL) GetOnlineUsers() ([]entities.User, error) {
	var users []entities.User
	err := d.db.Where("is_online = ?", true).Find(&users).Error
	return users, err
}

// CheckUsernameExists 检查用户名是否存在
func (d *UserDAL) CheckUsernameExists(username string) (bool, error) {
	var count int64
	err := d.db.Model(&entities.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// CheckEmailExists 检查邮箱是否存在
func (d *UserDAL) CheckEmailExists(email string) (bool, error) {
	var count int64
	err := d.db.Model(&entities.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
