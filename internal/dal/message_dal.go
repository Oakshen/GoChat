package dal

import (
	"gochat/internal/database"
	"gochat/internal/models/entities"

	"gorm.io/gorm"
)

// MessageDAL 消息数据访问层
type MessageDAL struct {
	db *gorm.DB
}

// NewMessageDAL 创建消息DAL实例
func NewMessageDAL() *MessageDAL {
	return &MessageDAL{
		db: database.DB,
	}
}

// Create 创建消息
func (d *MessageDAL) Create(message *entities.Message) (*entities.Message, error) {
	if err := d.db.Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// GetByID 根据ID获取消息
func (d *MessageDAL) GetByID(messageID uint) (*entities.Message, error) {
	var message entities.Message
	if err := d.db.Preload("User").Preload("Room").First(&message, messageID).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

// GetByRoomID 根据聊天室ID获取消息列表
func (d *MessageDAL) GetByRoomID(roomID uint, limit, offset int) ([]*entities.Message, error) {
	var messages []*entities.Message
	query := d.db.Where("room_id = ?", roomID).
		Preload("User").
		Preload("Attachments").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// Update 更新消息内容
func (d *MessageDAL) Update(messageID uint, content string) (*entities.Message, error) {
	var message entities.Message
	if err := d.db.First(&message, messageID).Error; err != nil {
		return nil, err
	}

	message.Content = content
	if err := d.db.Save(&message).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

// Delete 删除消息
func (d *MessageDAL) Delete(messageID uint) error {
	return d.db.Delete(&entities.Message{}, messageID).Error
}

// GetUnreadCount 获取用户在指定聊天室的未读消息数量
func (d *MessageDAL) GetUnreadCount(userID, roomID uint) (int64, error) {
	var count int64
	// 这里简化处理，实际应该有读取状态表
	// 当前仅返回该用户在该聊天室的所有消息数量
	err := d.db.Model(&entities.Message{}).
		Where("room_id = ? AND user_id != ?", roomID, userID).
		Count(&count).Error
	return count, err
}

// GetRecentMessages 获取最近消息
func (d *MessageDAL) GetRecentMessages(roomID uint, limit int) ([]*entities.Message, error) {
	var messages []*entities.Message
	err := d.db.Where("room_id = ?", roomID).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}
