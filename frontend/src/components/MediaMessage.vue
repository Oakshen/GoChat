<template>
  <div class="media-message">
    <!-- 消息文本（如果有） -->
    <div v-if="message.content" class="message-text">
      {{ message.content }}
    </div>

    <!-- 附件列表 -->
    <div class="attachments">
      <div 
        v-for="attachment in message.attachments" 
        :key="attachment.id"
        class="attachment-item"
        :class="`attachment-${attachment.category}`"
      >
        <!-- 图片附件 -->
        <div v-if="attachment.category === 'image'" class="image-attachment">
          <el-image
            :src="getPreviewUrl(attachment.id)"
            :preview-src-list="[getPreviewUrl(attachment.id)]"
            class="attachment-image"
            fit="cover"
            :alt="attachment.file_name"
            @click="previewImage(attachment)"
          >
            <template #error>
              <div class="image-error">
                <el-icon><Picture /></el-icon>
                <span>图片加载失败</span>
              </div>
            </template>
          </el-image>
          <div class="attachment-info">
            <span class="file-name">{{ attachment.file_name }}</span>
            <span class="file-size">{{ formatFileSize(attachment.file_size) }}</span>
          </div>
        </div>

        <!-- 视频附件 -->
        <div v-else-if="attachment.category === 'video'" class="video-attachment">
          <div class="video-container">
            <video
              :src="getPreviewUrl(attachment.id)"
              controls
              preload="metadata"
              class="attachment-video"
              :poster="getVideoPoster(attachment)"
            >
              您的浏览器不支持视频播放。
            </video>
          </div>
          <div class="attachment-info">
            <span class="file-name">{{ attachment.file_name }}</span>
            <div class="video-meta">
              <span class="file-size">{{ formatFileSize(attachment.file_size) }}</span>
              <span v-if="attachment.duration" class="duration">
                {{ formatDuration(attachment.duration) }}
              </span>
            </div>
          </div>
        </div>

        <!-- 文件附件 -->
        <div v-else class="file-attachment" @click="downloadFile(attachment)">
          <div class="file-icon">
            <el-icon size="32">
              <Document v-if="isDocument(attachment.file_type)" />
              <FolderOpened v-else-if="isArchive(attachment.file_type)" />
              <Files v-else />
            </el-icon>
          </div>
          <div class="file-info">
            <div class="file-name">{{ attachment.file_name }}</div>
            <div class="file-meta">
              <span class="file-size">{{ formatFileSize(attachment.file_size) }}</span>
              <span class="file-type">{{ getFileTypeDisplay(attachment.file_type) }}</span>
            </div>
          </div>
          <div class="download-icon">
            <el-icon><Download /></el-icon>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片预览对话框 -->
    <el-dialog
      v-model="showImagePreview"
      title="图片预览"
      width="80%"
      :show-close="true"
      center
    >
      <div class="image-preview">
        <img 
          :src="previewImageUrl" 
          :alt="previewImageName"
          class="preview-image"
        />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="downloadCurrentImage">下载</el-button>
          <el-button type="primary" @click="showImagePreview = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Picture, Document, FolderOpened, Files, Download 
} from '@element-plus/icons-vue'
import { getPreviewUrl, getDownloadUrl } from '@/api/file'

const props = defineProps({
  message: {
    type: Object,
    required: true
  }
})

// 响应式数据
const showImagePreview = ref(false)
const previewImageUrl = ref('')
const previewImageName = ref('')
const currentAttachment = ref(null)

// 方法
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDuration = (seconds) => {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const isDocument = (fileType) => {
  const docTypes = [
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.ms-powerpoint',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    'text/plain',
    'text/markdown',
    'text/csv'
  ]
  return docTypes.includes(fileType)
}

const isArchive = (fileType) => {
  const archiveTypes = [
    'application/zip',
    'application/vnd.rar',
    'application/x-7z-compressed',
    'application/x-tar',
    'application/gzip'
  ]
  return archiveTypes.includes(fileType)
}

const getFileTypeDisplay = (fileType) => {
  const typeMap = {
    'application/pdf': 'PDF',
    'application/msword': 'DOC',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document': 'DOCX',
    'application/vnd.ms-excel': 'XLS',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': 'XLSX',
    'application/vnd.ms-powerpoint': 'PPT',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation': 'PPTX',
    'text/plain': 'TXT',
    'text/markdown': 'MD',
    'text/csv': 'CSV',
    'application/json': 'JSON',
    'application/xml': 'XML',
    'application/zip': 'ZIP',
    'application/vnd.rar': 'RAR',
    'application/x-7z-compressed': '7Z'
  }
  
  return typeMap[fileType] || fileType.split('/')[1]?.toUpperCase() || 'FILE'
}

const getVideoPoster = (attachment) => {
  // 可以返回视频缩略图，这里暂时返回空
  return ''
}

const previewImage = (attachment) => {
  currentAttachment.value = attachment
  previewImageUrl.value = getPreviewUrl(attachment.id)
  previewImageName.value = attachment.file_name
  showImagePreview.value = true
}

const downloadFile = (attachment) => {
  const downloadUrl = getDownloadUrl(attachment.id)
  
  // 创建隐藏的a标签进行下载
  const link = document.createElement('a')
  link.href = downloadUrl
  link.download = attachment.file_name
  link.target = '_blank'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  
  ElMessage.success('开始下载文件')
}

const downloadCurrentImage = () => {
  if (currentAttachment.value) {
    downloadFile(currentAttachment.value)
  }
}
</script>

<style scoped>
.media-message {
  max-width: 100%;
}

.message-text {
  margin-bottom: 8px;
  word-wrap: break-word;
  line-height: 1.4;
}

.attachments {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.attachment-item {
  border-radius: 8px;
  overflow: hidden;
  background-color: #f8f9fa;
}

/* 图片附件 */
.image-attachment {
  max-width: 300px;
}

.attachment-image {
  width: 100%;
  max-height: 200px;
  cursor: pointer;
  display: block;
}

.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  color: #909399;
  background-color: #f5f7fa;
}

.image-error span {
  margin-top: 8px;
  font-size: 12px;
}

/* 视频附件 */
.video-attachment {
  max-width: 400px;
}

.video-container {
  position: relative;
  background-color: #000;
}

.attachment-video {
  width: 100%;
  max-height: 300px;
  display: block;
}

.video-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.duration {
  color: #666;
  font-size: 12px;
}

/* 文件附件 */
.file-attachment {
  display: flex;
  align-items: center;
  padding: 12px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  max-width: 350px;
}

.file-attachment:hover {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.file-icon {
  color: #409eff;
  margin-right: 12px;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
  display: block;
  word-break: break-all;
}

.file-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #909399;
}

.download-icon {
  color: #409eff;
  margin-left: 8px;
}

/* 附件信息 */
.attachment-info {
  padding: 8px 12px;
  background-color: #fff;
  border-top: 1px solid #ebeef5;
}

.attachment-info .file-name {
  font-size: 13px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
  display: block;
  word-break: break-all;
}

.file-size {
  font-size: 12px;
  color: #909399;
}

/* 图片预览 */
.image-preview {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.preview-image {
  max-width: 100%;
  max-height: 70vh;
  object-fit: contain;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style> 