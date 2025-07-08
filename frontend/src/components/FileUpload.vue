<template>
  <div class="file-upload">
    <!-- 上传按钮 -->
    <el-dropdown 
      split-button 
      type="primary" 
      size="small"
      @click="handleButtonClick"
      @command="handleDropdownCommand"
    >
      <el-icon><Plus /></el-icon>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="image">
            <el-icon><Picture /></el-icon> 图片
          </el-dropdown-item>
          <el-dropdown-item command="file">
            <el-icon><Document /></el-icon> 文件
          </el-dropdown-item>
          <el-dropdown-item command="video">
            <el-icon><VideoPlay /></el-icon> 视频
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>

    <!-- 隐藏的文件输入 -->
    <input
      ref="fileInput"
      type="file"
      multiple
      style="display: none"
      :accept="currentAccept"
      @change="handleFileSelect"
    />

    <!-- 拖拽上传区域 -->
    <el-dialog
      v-model="showUploadDialog"
      title="文件上传"
      width="600px"
      :close-on-click-modal="false"
    >
      <div
        class="upload-area"
        :class="{ 'is-dragover': isDragging }"
        @dragover.prevent
        @dragenter.prevent="isDragging = true"
        @dragleave.prevent="isDragging = false"
        @drop.prevent="handleDrop"
      >
        <div class="upload-content">
          <el-icon class="upload-icon" size="48"><UploadFilled /></el-icon>
          <div class="upload-text">
            <p>拖拽文件到此处或 <el-button type="text" @click="selectFiles">点击选择文件</el-button></p>
            <p class="upload-tip">
              支持 {{ getTypeDescription() }}，单个文件最大 50MB
            </p>
          </div>
        </div>
      </div>

      <!-- 文件列表 -->
      <div v-if="fileList.length > 0" class="file-list">
        <div v-for="(file, index) in fileList" :key="file.id" class="file-item">
          <div class="file-info">
            <el-icon class="file-icon">
              <Picture v-if="file.category === 'image'" />
              <VideoPlay v-else-if="file.category === 'video'" />
              <Document v-else />
            </el-icon>
            <div class="file-details">
              <div class="file-name">{{ file.name }}</div>
              <div class="file-size">{{ formatFileSize(file.size) }}</div>
            </div>
          </div>
          
          <!-- 进度条 -->
          <el-progress 
            v-if="file.status === 'uploading'"
            :percentage="file.progress"
            :stroke-width="4"
            class="file-progress"
          />
          
          <!-- 状态显示 -->
          <div class="file-status">
            <el-icon v-if="file.status === 'success'" class="success-icon">
              <CircleCheck />
            </el-icon>
            <el-icon v-else-if="file.status === 'error'" class="error-icon">
              <CircleClose />
            </el-icon>
            <el-button 
              v-if="file.status !== 'uploading'"
              type="text" 
              size="small"
              @click="removeFile(index)"
            >
              <el-icon><Close /></el-icon>
            </el-button>
          </div>
        </div>
      </div>

      <!-- 消息输入 -->
      <el-input
        v-model="messageText"
        type="textarea"
        :rows="3"
        placeholder="添加消息描述（可选）..."
        class="message-input"
      />

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showUploadDialog = false">取消</el-button>
          <el-button 
            type="primary"
            :loading="isUploading"
            :disabled="fileList.length === 0"
            @click="handleUpload"
          >
            发送 ({{ fileList.filter(f => f.status !== 'error').length }})
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Plus, Picture, Document, VideoPlay, UploadFilled, 
  CircleCheck, CircleClose, Close 
} from '@element-plus/icons-vue'
import { uploadFile } from '@/api/file'

const emit = defineEmits(['file-uploaded'])

// 响应式数据
const showUploadDialog = ref(false)
const isDragging = ref(false)
const fileInput = ref()
const fileList = ref([])
const messageText = ref('')
const isUploading = ref(false)
const currentUploadType = ref('all')

// 计算属性
const currentAccept = computed(() => {
  const accepts = {
    image: 'image/*',
    video: 'video/*',
    file: '.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.md,.csv,.json,.xml,.zip,.rar,.7z,.tar,.gz',
    all: 'image/*,video/*,.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.md,.csv,.json,.xml,.zip,.rar,.7z,.tar,.gz'
  }
  return accepts[currentUploadType.value] || accepts.all
})

// 方法
const handleButtonClick = () => {
  currentUploadType.value = 'all'
  showUploadDialog.value = true
}

const handleDropdownCommand = (command) => {
  currentUploadType.value = command
  showUploadDialog.value = true
}

const getTypeDescription = () => {
  const descriptions = {
    image: '图片文件 (JPG, PNG, GIF 等)',
    video: '视频文件 (MP4, AVI, MOV 等)',
    file: '文档文件 (PDF, DOC, TXT 等)',
    all: '图片、视频、文档文件'
  }
  return descriptions[currentUploadType.value] || descriptions.all
}

const selectFiles = () => {
  fileInput.value.click()
}

const handleFileSelect = (event) => {
  const files = Array.from(event.target.files)
  addFiles(files)
  // 清空输入框，允许重复选择同一文件
  event.target.value = ''
}

const handleDrop = (event) => {
  isDragging.value = false
  const files = Array.from(event.dataTransfer.files)
  addFiles(files)
}

const addFiles = (files) => {
  files.forEach(file => {
    // 验证文件
    const error = validateFile(file)
    if (error) {
      ElMessage.error(error)
      return
    }

    const fileItem = {
      id: Date.now() + Math.random(),
      file: file,
      name: file.name,
      size: file.size,
      category: getFileCategory(file.name),
      status: 'ready', // ready, uploading, success, error
      progress: 0,
      error: null
    }

    fileList.value.push(fileItem)
  })
}

const validateFile = (file) => {
  // 文件大小限制 50MB
  const maxSize = 50 * 1024 * 1024
  if (file.size > maxSize) {
    return `文件 "${file.name}" 超过最大限制 50MB`
  }

  // 文件类型验证
  const allowedTypes = {
    image: ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp'],
    video: ['mp4', 'avi', 'mov', 'wmv', 'flv', 'webm'],
    file: ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'txt', 'md', 'csv', 'json', 'xml', 'zip', 'rar', '7z', 'tar', 'gz']
  }

  const extension = file.name.split('.').pop().toLowerCase()
  const category = getFileCategory(file.name)
  
  if (currentUploadType.value !== 'all' && currentUploadType.value !== category) {
    return `文件 "${file.name}" 不是 ${getTypeDescription()} 类型`
  }

  const allAllowed = [...allowedTypes.image, ...allowedTypes.video, ...allowedTypes.file]
  if (!allAllowed.includes(extension)) {
    return `不支持的文件类型: ${extension}`
  }

  return null
}

const getFileCategory = (filename) => {
  const ext = filename.split('.').pop().toLowerCase()
  
  const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp']
  const videoExts = ['mp4', 'avi', 'mov', 'wmv', 'flv', 'webm']
  
  if (imageExts.includes(ext)) return 'image'
  if (videoExts.includes(ext)) return 'video'
  return 'file'
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const removeFile = (index) => {
  fileList.value.splice(index, 1)
}

const handleUpload = async () => {
  const validFiles = fileList.value.filter(f => f.status !== 'error')
  if (validFiles.length === 0) {
    ElMessage.warning('没有可上传的文件')
    return
  }

  isUploading.value = true

  // 首先创建一个空消息以获取消息ID
  try {
    const uploadResults = []
    
    for (const fileItem of validFiles) {
      fileItem.status = 'uploading'
      fileItem.progress = 0

      try {
        // 创建 FormData
        const formData = new FormData()
        formData.append('file', fileItem.file)
        formData.append('message_id', '0') // 临时ID，后端需要处理

        // 上传文件
        const result = await uploadFile(formData, (progress) => {
          fileItem.progress = progress
        })

        fileItem.status = 'success'
        fileItem.progress = 100
        uploadResults.push({
          attachment: result.data,
          fileItem: fileItem
        })

      } catch (error) {
        fileItem.status = 'error'
        fileItem.error = error.response?.data?.error || '上传失败'
        ElMessage.error(`文件 "${fileItem.name}" 上传失败: ${fileItem.error}`)
      }
    }

    // 如果有成功上传的文件，触发文件上传完成事件
    if (uploadResults.length > 0) {
      emit('file-uploaded', {
        attachments: uploadResults.map(r => r.attachment),
        messageText: messageText.value,
        category: uploadResults[0].attachment.category
      })

      // 重置状态
      fileList.value = []
      messageText.value = ''
      showUploadDialog.value = false
      
      ElMessage.success(`成功上传 ${uploadResults.length} 个文件`)
    }

  } catch (error) {
    ElMessage.error('上传过程中出现错误')
  } finally {
    isUploading.value = false
  }
}
</script>

<style scoped>
.file-upload {
  display: inline-block;
}

.upload-area {
  border: 2px dashed #dcdfe6;
  border-radius: 8px;
  padding: 40px;
  text-align: center;
  background-color: #fafafa;
  transition: all 0.3s;
  cursor: pointer;
}

.upload-area:hover,
.upload-area.is-dragover {
  border-color: #409eff;
  background-color: #f0f9ff;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.upload-icon {
  color: #c0c4cc;
}

.upload-text {
  color: #606266;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}

.file-list {
  margin-top: 20px;
  max-height: 300px;
  overflow-y: auto;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  margin-bottom: 8px;
  background-color: #fafafa;
}

.file-info {
  display: flex;
  align-items: center;
  flex: 1;
  gap: 12px;
}

.file-icon {
  font-size: 24px;
  color: #409eff;
}

.file-details {
  flex: 1;
}

.file-name {
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.file-size {
  font-size: 12px;
  color: #909399;
}

.file-progress {
  flex: 1;
  margin: 0 16px;
}

.file-status {
  display: flex;
  align-items: center;
  gap: 8px;
}

.success-icon {
  color: #67c23a;
  font-size: 18px;
}

.error-icon {
  color: #f56c6c;
  font-size: 18px;
}

.message-input {
  margin-top: 16px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style> 