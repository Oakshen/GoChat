package services

import (
	"fmt"
	"gochat/internal/dal"
	"gochat/internal/models/entities"
	"gochat/pkg/logger"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AttachmentService 附件服务
type AttachmentService struct {
	attachmentDAL *dal.AttachmentDAL
	uploadDir     string
}

// NewAttachmentService 创建附件服务实例
func NewAttachmentService() *AttachmentService {
	uploadDir := "./uploads"
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Error("Failed to create upload directory:", err)
	}

	return &AttachmentService{
		attachmentDAL: dal.NewAttachmentDAL(),
		uploadDir:     uploadDir,
	}
}

// UploadFile 上传文件
func (s *AttachmentService) UploadFile(fileHeader *multipart.FileHeader, messageID uint) (*entities.Attachment, error) {
	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("无法打开上传文件: %v", err)
	}
	defer file.Close()

	// 生成唯一的文件名
	ext := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.ReplaceAll(fileHeader.Filename, ext, ""), ext)

	// 根据文件类型创建子目录
	category := s.getCategoryFromFileName(fileHeader.Filename)
	subDir := filepath.Join(s.uploadDir, category)
	if err := os.MkdirAll(subDir, 0755); err != nil {
		return nil, fmt.Errorf("无法创建子目录: %v", err)
	}

	// 完整的文件路径
	filePath := filepath.Join(subDir, fileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法创建目标文件: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	size, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("文件复制失败: %v", err)
	}

	// 获取文件MIME类型
	fileType := s.getFileType(fileHeader.Filename)

	// 创建附件记录
	attachment := &entities.Attachment{
		MessageID: messageID, // 如果为0，表示临时附件，稍后会关联到实际消息
		FileName:  fileHeader.Filename,
		FilePath:  filePath,
		FileSize:  size,
		FileType:  fileType,
		Category:  category,
	}

	// 如果是图片或视频，可以在这里获取宽高信息
	// 这里先简化处理，后续可以添加图片处理库
	if category == "image" {
		// TODO: 获取图片宽高
	} else if category == "video" {
		// TODO: 获取视频宽高和时长
	}

	// 保存到数据库
	return s.attachmentDAL.Create(attachment)
}

// GetAttachment 获取附件
func (s *AttachmentService) GetAttachment(attachmentID uint) (*entities.Attachment, error) {
	return s.attachmentDAL.GetByID(attachmentID)
}

// GetMessageAttachments 获取消息的所有附件
func (s *AttachmentService) GetMessageAttachments(messageID uint) ([]*entities.Attachment, error) {
	return s.attachmentDAL.GetByMessageID(messageID)
}

// DeleteAttachment 删除附件
func (s *AttachmentService) DeleteAttachment(attachmentID uint) error {
	// 获取附件信息
	attachment, err := s.attachmentDAL.GetByID(attachmentID)
	if err != nil {
		return err
	}

	// 删除物理文件
	if err := os.Remove(attachment.FilePath); err != nil {
		logger.Error("Failed to delete file:", attachment.FilePath, "error:", err)
		// 即使物理文件删除失败，也继续删除数据库记录
	}

	// 删除数据库记录
	return s.attachmentDAL.Delete(attachmentID)
}

// UpdateAttachmentMessageID 更新附件的消息ID
func (s *AttachmentService) UpdateAttachmentMessageID(attachmentID, messageID uint) error {
	return s.attachmentDAL.UpdateMessageID(attachmentID, messageID)
}

// ValidateFile 验证文件
func (s *AttachmentService) ValidateFile(fileHeader *multipart.FileHeader) error {
	// 文件大小限制（50MB）
	const maxFileSize = 50 * 1024 * 1024
	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("文件大小超过限制，最大允许50MB")
	}

	// 允许的文件类型
	allowedTypes := map[string]bool{
		// 图片类型
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".webp": true,
		// 视频类型
		".mp4": true, ".avi": true, ".mov": true, ".wmv": true, ".flv": true, ".webm": true,
		// 文档类型
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true,
		".txt": true, ".md": true, ".csv": true,
		// 压缩文件
		".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
		// 其他常见类型
		".json": true, ".xml": true, ".log": true,
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedTypes[ext] {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}

	return nil
}

// getCategoryFromFileName 根据文件名获取文件分类
func (s *AttachmentService) getCategoryFromFileName(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".webp": true,
	}

	videoExts := map[string]bool{
		".mp4": true, ".avi": true, ".mov": true, ".wmv": true, ".flv": true, ".webm": true,
	}

	if imageExts[ext] {
		return "image"
	} else if videoExts[ext] {
		return "video"
	} else {
		return "file"
	}
}

// getFileType 获取文件MIME类型
func (s *AttachmentService) getFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".webp": "image/webp",
		".mp4":  "video/mp4",
		".avi":  "video/avi",
		".mov":  "video/quicktime",
		".wmv":  "video/x-ms-wmv",
		".flv":  "video/x-flv",
		".webm": "video/webm",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt":  "text/plain",
		".md":   "text/markdown",
		".csv":  "text/csv",
		".json": "application/json",
		".xml":  "application/xml",
		".zip":  "application/zip",
		".rar":  "application/vnd.rar",
		".7z":   "application/x-7z-compressed",
	}

	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}

	return "application/octet-stream"
}
