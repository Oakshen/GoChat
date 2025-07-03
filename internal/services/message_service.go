package services

import (
	"gochat/internal/dal"
	"gochat/internal/models/entities"
)

// MessageService 消息服务
type MessageService struct {
	messageDAL *dal.MessageDAL
}

// NewMessageService 创建消息服务实例
func NewMessageService() *MessageService {
	return &MessageService{
		messageDAL: dal.NewMessageDAL(),
	}
}

// CreateMessage 创建消息
func (s *MessageService) CreateMessage(message *entities.Message) (*entities.Message, error) {
	return s.messageDAL.Create(message)
}

// GetMessagesByRoom 获取聊天室消息
func (s *MessageService) GetMessagesByRoom(roomID uint, limit, offset int) ([]*entities.Message, error) {
	return s.messageDAL.GetByRoomID(roomID, limit, offset)
}

// GetMessageByID 根据ID获取消息
func (s *MessageService) GetMessageByID(messageID uint) (*entities.Message, error) {
	return s.messageDAL.GetByID(messageID)
}

// UpdateMessage 更新消息
func (s *MessageService) UpdateMessage(messageID uint, content string) (*entities.Message, error) {
	return s.messageDAL.Update(messageID, content)
}

// DeleteMessage 删除消息
func (s *MessageService) DeleteMessage(messageID uint) error {
	return s.messageDAL.Delete(messageID)
}

// GetUnreadCount 获取未读消息数量
func (s *MessageService) GetUnreadCount(userID, roomID uint) (int64, error) {
	return s.messageDAL.GetUnreadCount(userID, roomID)
}
