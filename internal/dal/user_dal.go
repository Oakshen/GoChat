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

// ===================== 软删除相关方法 =====================

// SoftDeleteByID 软删除用户（推荐使用）
func (d *UserDAL) SoftDeleteByID(userID uint) error {
	return d.db.Delete(&entities.User{}, userID).Error
}

// HardDeleteByID 硬删除用户（谨慎使用）
func (d *UserDAL) HardDeleteByID(userID uint) error {
	return d.db.Unscoped().Delete(&entities.User{}, userID).Error
}

// RestoreByID 恢复已软删除的用户
func (d *UserDAL) RestoreByID(userID uint) error {
	return d.db.Unscoped().Model(&entities.User{}).Where("id = ?", userID).Update("deleted_at", nil).Error
}

// GetDeletedUsers 获取已删除的用户列表
func (d *UserDAL) GetDeletedUsers() ([]entities.User, error) {
	var users []entities.User
	err := d.db.Unscoped().Where("deleted_at IS NOT NULL").Find(&users).Error
	return users, err
}

// GetByIDWithDeleted 根据ID查找用户（包括已删除的）
func (d *UserDAL) GetByIDWithDeleted(id uint) (*entities.User, error) {
	var user entities.User
	err := d.db.Unscoped().First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// IsUserDeleted 检查用户是否已被软删除
func (d *UserDAL) IsUserDeleted(userID uint) (bool, error) {
	var user entities.User
	err := d.db.Unscoped().Select("deleted_at").First(&user, userID).Error
	if err != nil {
		return false, err
	}
	return user.DeletedAt.Valid, nil
}

// ===================== 房间相关软删除方法 =====================

// RoomDAL 房间数据访问层
type RoomDAL struct {
	db *gorm.DB
}

// NewRoomDAL 创建房间数据访问层实例
func NewRoomDAL() *RoomDAL {
	return &RoomDAL{
		db: database.DB,
	}
}

// SoftDeleteRoom 软删除房间
func (d *RoomDAL) SoftDeleteRoom(roomID uint) error {
	return d.db.Delete(&entities.Room{}, roomID).Error
}

// RestoreRoom 恢复已软删除的房间
func (d *RoomDAL) RestoreRoom(roomID uint) error {
	return d.db.Unscoped().Model(&entities.Room{}).Where("id = ?", roomID).Update("deleted_at", nil).Error
}

// GetActiveRooms 获取活跃房间列表
func (d *RoomDAL) GetActiveRooms() ([]entities.Room, error) {
	var rooms []entities.Room
	err := d.db.Find(&rooms).Error // 自动排除已删除的
	return rooms, err
}

// ===================== 房间成员相关软删除方法 =====================

// RoomMemberDAL 房间成员数据访问层
type RoomMemberDAL struct {
	db *gorm.DB
}

// NewRoomMemberDAL 创建房间成员数据访问层实例
func NewRoomMemberDAL() *RoomMemberDAL {
	return &RoomMemberDAL{
		db: database.DB,
	}
}

// SoftDeleteMember 软删除房间成员（用户退出房间）
func (d *RoomMemberDAL) SoftDeleteMember(roomID, userID uint) error {
	return d.db.Where("room_id = ? AND user_id = ?", roomID, userID).Delete(&entities.RoomMember{}).Error
}

// RestoreMember 恢复房间成员
func (d *RoomMemberDAL) RestoreMember(roomID, userID uint) error {
	return d.db.Unscoped().Model(&entities.RoomMember{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("deleted_at", nil).Error
}

// GetRoomActiveMembers 获取房间活跃成员
func (d *RoomMemberDAL) GetRoomActiveMembers(roomID uint) ([]entities.RoomMember, error) {
	var members []entities.RoomMember
	err := d.db.Where("room_id = ?", roomID).Find(&members).Error // 自动排除已删除的
	return members, err
}

// GetMemberHistory 获取房间成员历史（包括已删除的）
func (d *RoomMemberDAL) GetMemberHistory(roomID uint) ([]entities.RoomMember, error) {
	var members []entities.RoomMember
	err := d.db.Unscoped().Where("room_id = ?", roomID).Find(&members).Error
	return members, err
}
