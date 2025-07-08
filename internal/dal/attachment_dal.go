package dal

import (
	"gochat/internal/database"
	"gochat/internal/models/entities"

	"gorm.io/gorm"
)

// AttachmentDAL 附件数据访问层
type AttachmentDAL struct {
	db *gorm.DB
}

// NewAttachmentDAL 创建附件DAL实例
func NewAttachmentDAL() *AttachmentDAL {
	return &AttachmentDAL{
		db: database.DB,
	}
}

// Create 创建附件
func (d *AttachmentDAL) Create(attachment *entities.Attachment) (*entities.Attachment, error) {
	if err := d.db.Create(attachment).Error; err != nil {
		return nil, err
	}
	return attachment, nil
}

// GetByID 根据ID获取附件
func (d *AttachmentDAL) GetByID(attachmentID uint) (*entities.Attachment, error) {
	var attachment entities.Attachment
	if err := d.db.First(&attachment, attachmentID).Error; err != nil {
		return nil, err
	}
	return &attachment, nil
}

// GetByMessageID 根据消息ID获取附件列表
func (d *AttachmentDAL) GetByMessageID(messageID uint) ([]*entities.Attachment, error) {
	var attachments []*entities.Attachment
	if err := d.db.Where("message_id = ?", messageID).Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

// Delete 删除附件
func (d *AttachmentDAL) Delete(attachmentID uint) error {
	return d.db.Delete(&entities.Attachment{}, attachmentID).Error
}

// DeleteByMessageID 根据消息ID删除所有相关附件
func (d *AttachmentDAL) DeleteByMessageID(messageID uint) error {
	return d.db.Where("message_id = ?", messageID).Delete(&entities.Attachment{}).Error
}

// UpdateMessageID 更新附件的消息ID
func (d *AttachmentDAL) UpdateMessageID(attachmentID, messageID uint) error {
	return d.db.Model(&entities.Attachment{}).Where("id = ?", attachmentID).Update("message_id", messageID).Error
}
