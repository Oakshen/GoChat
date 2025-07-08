package handlers

import (
	"context"
	"fmt"
	"gochat/internal/services"
	"gochat/pkg/response"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// FileHandler 文件处理器
type FileHandler struct {
	attachmentService *services.AttachmentService
}

// NewFileHandler 创建文件处理器实例
func NewFileHandler() *FileHandler {
	return &FileHandler{
		attachmentService: services.NewAttachmentService(),
	}
}

// UploadFile 上传文件
func (h *FileHandler) UploadFile(ctx context.Context, c *app.RequestContext) {
	// 获取消息ID参数（可选，如果为0则创建临时附件）
	messageIDStr := c.PostForm("message_id")
	var messageID uint = 0

	if messageIDStr != "" {
		id, err := strconv.ParseUint(messageIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.H{
				"error": "无效的消息ID",
			})
			return
		}
		messageID = uint(id)
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "未找到上传文件",
		})
		return
	}

	// 验证文件
	if err := h.attachmentService.ValidateFile(file); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}

	// 上传文件（messageID为0时创建临时附件）
	attachment, err := h.attachmentService.UploadFile(file, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"error": fmt.Sprintf("文件上传失败: %v", err),
		})
		return
	}

	response.Success(ctx, c, attachment)
}

// DownloadFile 下载文件
func (h *FileHandler) DownloadFile(ctx context.Context, c *app.RequestContext) {
	attachmentIDStr := c.Param("id")
	attachmentID, err := strconv.ParseUint(attachmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的附件ID",
		})
		return
	}

	// 获取附件信息
	attachment, err := h.attachmentService.GetAttachment(uint(attachmentID))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.H{
			"error": "附件不存在",
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", attachment.FileType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachment.FileName))
	c.Header("Content-Length", fmt.Sprintf("%d", attachment.FileSize))

	// 发送文件
	c.File(attachment.FilePath)
}

// PreviewFile 预览文件（主要用于图片和视频）
func (h *FileHandler) PreviewFile(ctx context.Context, c *app.RequestContext) {
	attachmentIDStr := c.Param("id")
	attachmentID, err := strconv.ParseUint(attachmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的附件ID",
		})
		return
	}

	// 获取附件信息
	attachment, err := h.attachmentService.GetAttachment(uint(attachmentID))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.H{
			"error": "附件不存在",
		})
		return
	}

	// 只允许预览图片和视频
	if attachment.Category != "image" && attachment.Category != "video" {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "该文件类型不支持预览",
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", attachment.FileType)
	c.Header("Cache-Control", "public, max-age=31536000") // 缓存一年

	// 发送文件
	c.File(attachment.FilePath)
}

// GetFileInfo 获取文件信息
func (h *FileHandler) GetFileInfo(ctx context.Context, c *app.RequestContext) {
	attachmentIDStr := c.Param("id")
	attachmentID, err := strconv.ParseUint(attachmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的附件ID",
		})
		return
	}

	// 获取附件信息
	attachment, err := h.attachmentService.GetAttachment(uint(attachmentID))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.H{
			"error": "附件不存在",
		})
		return
	}

	response.Success(ctx, c, attachment)
}

// DeleteFile 删除文件
func (h *FileHandler) DeleteFile(ctx context.Context, c *app.RequestContext) {
	attachmentIDStr := c.Param("id")
	attachmentID, err := strconv.ParseUint(attachmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"error": "无效的附件ID",
		})
		return
	}

	// 删除附件
	if err := h.attachmentService.DeleteAttachment(uint(attachmentID)); err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"error": fmt.Sprintf("删除文件失败: %v", err),
		})
		return
	}

	response.Success(ctx, c, utils.H{
		"message": "文件删除成功",
	})
}
