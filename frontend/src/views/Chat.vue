<template>
  <div class="chat-container">
    <!-- 侧边栏 -->
    <div class="sidebar">
      <!-- 用户信息头部 -->
      <div class="user-header">
        <div class="user-info">
          <el-avatar :size="40" :src="userAvatar">
            {{ authStore.user?.username?.charAt(0)?.toUpperCase() }}
          </el-avatar>
          <div class="user-details">
            <div class="username">{{ authStore.user?.username }}</div>
            <div class="status">在线</div>
          </div>
        </div>
        <div class="header-actions">
          <el-button 
            :icon="Plus" 
            circle 
            size="small" 
            @click="showCreateRoomDialog = true"
            title="创建聊天室"
          />
          <el-button 
            :icon="Connection" 
            circle 
            size="small" 
            @click="showJoinRoomDialog = true"
            title="加入聊天室"
          />
          <el-dropdown @command="handleUserAction">
            <el-button :icon="Setting" circle size="small" />
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchQuery"
          placeholder="搜索联系人或聊天室"
          :prefix-icon="Search"
          clearable
        />
      </div>

      <!-- 标签切换 -->
      <el-tabs v-model="activeTab" class="sidebar-tabs">
        <!-- 聊天室列表 -->
        <el-tab-pane label="聊天室" name="rooms">
          <div class="room-list">
            <div
              v-for="room in filteredRooms"
              :key="room.id"
              class="room-item"
              :class="{ active: currentRoomId === room.id }"
              @click="selectRoom(room.id)"
            >
              <el-avatar :size="40" :src="room.avatar">
                {{ room.name?.charAt(0)?.toUpperCase() }}
              </el-avatar>
              <div class="room-info">
                <div class="room-name">{{ room.name }}</div>
                <div class="room-desc">{{ room.description }}</div>
              </div>
              <div class="room-meta">
                <div class="member-count">{{ room.member_count || 0 }}人</div>
              </div>
            </div>
          </div>
        </el-tab-pane>

        <!-- 在线用户列表 -->
        <el-tab-pane label="联系人" name="contacts">
          <div class="contact-list">
            <div
              v-for="user in filteredContacts"
              :key="user.user_id"
              class="contact-item"
            >
              <el-avatar :size="40" :src="user.avatar_url">
                {{ user.username?.charAt(0)?.toUpperCase() }}
              </el-avatar>
              <div class="contact-info">
                <div class="contact-name">{{ user.username }}</div>
                <div class="contact-status" :class="{ online: user.is_online }">
                  {{ user.is_online ? '在线' : '离线' }}
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 主聊天区域 -->
    <div class="chat-main">
      <div v-if="!currentRoomId" class="chat-placeholder">
        <div class="placeholder-content">
          <el-icon :size="80" color="#ccc"><ChatDotRound /></el-icon>
          <h3>选择一个聊天室开始聊天</h3>
          <p>从左侧选择一个聊天室，或创建新的聊天室</p>
        </div>
      </div>
      
      <!-- 聊天窗口 -->
      <div v-else class="chat-window">
        <!-- 聊天头部 -->
        <div class="chat-header">
          <div class="chat-info">
            <el-avatar :size="36" :src="currentRoom?.avatar">
              {{ currentRoom?.name?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <div class="chat-details">
              <div class="chat-name">{{ currentRoom?.name }}</div>
              <div class="chat-members">{{ currentRoomMembers.length }}人在线</div>
            </div>
          </div>
          <div class="chat-actions">
            <el-button :icon="User" size="small" @click="showMembersDialog = true">
              成员
            </el-button>
          </div>
        </div>

        <!-- 消息列表 -->
        <div ref="messagesContainer" class="messages-container">
          <div
            v-for="message in currentMessages"
            :key="message.id || message.timestamp"
            class="message-item"
            :class="{ 
              'own-message': message.user_id === authStore.user?.id,
              'system-message': message.type === 'system'
            }"
          >
            <div v-if="message.type === 'system'" class="system-content">
              {{ message.content }}
            </div>
            <div v-else class="message-content">
              <el-avatar 
                :size="32" 
                :src="message.avatar"
                class="message-avatar"
              >
                {{ message.username?.charAt(0)?.toUpperCase() }}
              </el-avatar>
              <div class="message-bubble">
                <div class="message-header">
                  <span class="message-username">{{ message.username }}</span>
                  <span class="message-time">{{ formatTime(message.timestamp) }}</span>
                </div>
                <!-- 多媒体消息或文本消息 -->
                <MediaMessage 
                  v-if="message.attachments && message.attachments.length > 0"
                  :message="message"
                />
                <div v-else class="message-text">{{ message.content }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="input-area">
          <div class="input-container">
            <el-input
              v-model="messageInput"
              type="textarea"
              :rows="2"
              placeholder="输入消息..."
              resize="none"
              @keydown.enter="handleInputKeydown"
            />
            <div class="input-actions">
              <FileUpload @file-uploaded="handleFileUploaded" />
              <el-button 
                type="primary" 
                :icon="Promotion"
                :loading="sendingMessage"
                @click="sendMessage"
              >
                发送
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建聊天室对话框 -->
    <el-dialog
      v-model="showCreateRoomDialog"
      title="创建聊天室"
      width="400px"
    >
      <el-form :model="createRoomForm" label-width="80px">
        <el-form-item label="聊天室名称">
          <el-input v-model="createRoomForm.name" placeholder="请输入聊天室名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input 
            v-model="createRoomForm.description" 
            type="textarea" 
            placeholder="请输入聊天室描述"
          />
        </el-form-item>
        <el-form-item label="类型">
          <el-radio-group v-model="createRoomForm.is_private">
            <el-radio :label="false">公开</el-radio>
            <el-radio :label="true">私密</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateRoomDialog = false">取消</el-button>
          <el-button type="primary" @click="createRoom" :loading="creatingRoom">
            创建
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 加入聊天室对话框 -->
    <el-dialog
      v-model="showJoinRoomDialog"
      title="加入聊天室"
      width="400px"
      @close="handleJoinDialogClose"
    >
      <el-form :model="joinRoomForm" label-width="80px">
        <el-form-item label="搜索方式">
          <el-radio-group v-model="joinRoomForm.searchType">
            <el-radio label="id">聊天室ID</el-radio>
            <el-radio label="name">聊天室名称</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="joinRoomForm.searchType === 'id' ? '聊天室ID' : '聊天室名称'">
          <el-input 
            v-model="joinRoomForm.searchValue" 
            :placeholder="joinRoomForm.searchType === 'id' ? '请输入聊天室ID' : '请输入聊天室名称'"
            @input="searchRooms"
          />
        </el-form-item>
        
        <!-- 搜索结果 -->
        <el-form-item label="搜索结果" v-if="searchResults.length > 0">
          <div class="search-results">
            <div
              v-for="room in searchResults"
              :key="room.id"
              class="search-result-item"
              :class="{ selected: joinRoomForm.selectedRoomId === room.id }"
              @click="selectSearchResult(room)"
            >
              <el-avatar :size="32" :src="room.avatar">
                {{ room.name?.charAt(0)?.toUpperCase() }}
              </el-avatar>
              <div class="result-info">
                <div class="result-name">{{ room.name }}</div>
                <div class="result-desc">ID: {{ room.id }} | {{ room.description || '暂无描述' }}</div>
                <div class="result-members">{{ room.member_count || 0 }}人 | {{ room.is_private ? '私密' : '公开' }}</div>
              </div>
            </div>
          </div>
        </el-form-item>
        
        <!-- 无搜索结果提示 -->
        <el-form-item v-if="joinRoomForm.searchValue && searchResults.length === 0 && !searchLoading">
          <el-alert
            title="未找到匹配的聊天室"
            type="info"
            description="请检查输入的ID或名称是否正确"
            :closable="false"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showJoinRoomDialog = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="joinSelectedRoom" 
            :loading="joiningRoom"
            :disabled="!joinRoomForm.selectedRoomId"
          >
            加入聊天室
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 成员列表对话框 -->
    <el-dialog
      v-model="showMembersDialog"
      title="聊天室成员"
      width="400px"
    >
      <div class="members-list">
        <div
          v-for="member in currentRoomMembers"
          :key="member.user_id"
          class="member-item"
        >
          <el-avatar :size="36" :src="member.avatar_url">
            {{ member.username?.charAt(0)?.toUpperCase() }}
          </el-avatar>
          <div class="member-info">
            <div class="member-name">{{ member.username }}</div>
            <div class="member-status" :class="{ online: member.is_online }">
              {{ member.is_online ? '在线' : '离线' }}
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Plus, Setting, Search, User, ChatDotRound, Promotion, Connection 
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useChatStore } from '@/stores/chat'
import { createRoom as createRoomApi, joinRoom, searchRoomById, searchRoomsByName } from '@/api/room'
import wsClient from '@/utils/websocket'
import FileUpload from '@/components/FileUpload.vue'
import MediaMessage from '@/components/MediaMessage.vue'

const router = useRouter()
const authStore = useAuthStore()
const chatStore = useChatStore()

// 响应式数据
const searchQuery = ref('')
const activeTab = ref('rooms')
const messageInput = ref('')
const sendingMessage = ref(false)
const showCreateRoomDialog = ref(false)
const showJoinRoomDialog = ref(false)
const showMembersDialog = ref(false)
const creatingRoom = ref(false)
const joiningRoom = ref(false)
const searchLoading = ref(false)
const messagesContainer = ref()
const searchResults = ref([])

// 用户头像（暂时使用默认值）
const userAvatar = ref('')

// 创建聊天室表单
const createRoomForm = reactive({
  name: '',
  description: '',
  is_private: false
})

// 加入聊天室表单
const joinRoomForm = reactive({
  searchType: 'name', // 'id' 或 'name'
  searchValue: '',
  selectedRoomId: null
})

// 计算属性
const currentRoomId = computed(() => chatStore.currentRoomId)
const currentRoom = computed(() => chatStore.currentRoom)
const currentMessages = computed(() => chatStore.currentMessages)
const currentRoomMembers = computed(() => chatStore.currentRoomMembers)

const filteredRooms = computed(() => {
  if (!searchQuery.value) return chatStore.rooms
  return chatStore.rooms.filter(room => 
    room.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const filteredContacts = computed(() => {
  if (!searchQuery.value) return chatStore.onlineUsers
  return chatStore.onlineUsers.filter(user => 
    user.username.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// 方法
const selectRoom = async (roomId) => {
  try {
    // 检查用户是否已经在该聊天室中
    const isAlreadyInRoom = chatStore.rooms.some(room => room.id === roomId)
    
    if (!isAlreadyInRoom) {
      // 如果没有在聊天室中，发送HTTP请求加入
      console.log('User not in room, sending join request:', roomId)
      await joinRoom(roomId)
      
      // 重新获取聊天室列表以更新状态
      await chatStore.fetchRooms()
    } else {
      console.log('User already in room, skipping join request:', roomId)
    }
    
    // 设置当前聊天室
    chatStore.setCurrentRoom(roomId)
    
    // 获取聊天记录
    await chatStore.fetchMessages(roomId)
    
    // WebSocket加入聊天室（总是需要，因为WebSocket连接可能断开过）
    wsClient.joinRoom(roomId)
    
    // 滚动到底部
    nextTick(() => {
      scrollToBottom()
    })
  } catch (error) {
    console.error('Failed to select room:', error)
    ElMessage.error('进入聊天室失败')
  }
}

const sendMessage = async () => {
  if (!messageInput.value.trim() || !currentRoomId.value) return
  
  sendingMessage.value = true
  try {
    wsClient.sendTextMessage(currentRoomId.value, messageInput.value.trim())
    messageInput.value = ''
  } catch (error) {
    console.error('Failed to send message:', error)
    ElMessage.error('发送消息失败')
  } finally {
    sendingMessage.value = false
  }
}

// 处理文件上传
const handleFileUploaded = (uploadData) => {
  if (!currentRoomId.value) {
    ElMessage.error('请先选择聊天室')
    return
  }

  // 构建附件信息
  const attachments = uploadData.attachments.map(attachment => ({
    id: attachment.id,
    file_name: attachment.file_name,
    file_size: attachment.file_size,
    file_type: attachment.file_type,
    category: attachment.category,
    width: attachment.width || 0,
    height: attachment.height || 0,
    duration: attachment.duration || 0,
    url: `/api/files/${attachment.id}/preview`
  }))

  // 发送多媒体消息
  const messageType = uploadData.category === 'image' ? 'image' : 
                     uploadData.category === 'video' ? 'video' : 'file'
  
  wsClient.sendMediaMessage(currentRoomId.value, {
    type: messageType,
    content: uploadData.messageText || '',
    attachments: attachments
  })

  ElMessage.success('文件发送成功')
}

const handleInputKeydown = (event) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    sendMessage()
  }
}

const createRoom = async () => {
  if (!createRoomForm.name.trim()) {
    ElMessage.error('请输入聊天室名称')
    return
  }
  
  creatingRoom.value = true
  try {
    const response = await createRoomApi(createRoomForm)
    chatStore.addRoom(response.data)
    ElMessage.success('聊天室创建成功')
    showCreateRoomDialog.value = false
    
    // 清空表单
    Object.assign(createRoomForm, {
      name: '',
      description: '',
      is_private: false
    })
  } catch (error) {
    console.error('Failed to create room:', error)
  } finally {
    creatingRoom.value = false
  }
}

// 搜索聊天室（防抖处理）
let searchTimer = null
const searchRooms = () => {
  // 清除之前的定时器
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  
  // 设置新的定时器
  searchTimer = setTimeout(async () => {
    if (!joinRoomForm.searchValue.trim()) {
      searchResults.value = []
      return
    }
    
    searchLoading.value = true
    try {
      if (joinRoomForm.searchType === 'id') {
        // 按ID搜索 - 调用GetRoom接口
        const roomId = parseInt(joinRoomForm.searchValue)
        if (!isNaN(roomId)) {
          try {
            const response = await searchRoomById(roomId)
            searchResults.value = response.data ? [response.data] : []
          } catch (error) {
            // 如果聊天室不存在或无权访问，返回空结果
            console.log('Room not found or access denied:', error)
            searchResults.value = []
          }
        } else {
          searchResults.value = []
        }
      } else {
        // 按名称搜索 - 调用GetRoomsByBlurName接口
        try {
          const response = await searchRoomsByName(joinRoomForm.searchValue)
          searchResults.value = response.data || []
        } catch (error) {
          console.error('Failed to search rooms by name:', error)
          searchResults.value = []
        }
      }
    } catch (error) {
      console.error('Failed to search rooms:', error)
      searchResults.value = []
    } finally {
      searchLoading.value = false
    }
  }, 300) // 300ms 防抖延迟
}

// 选择搜索结果
const selectSearchResult = (room) => {
  joinRoomForm.selectedRoomId = room.id
}

// 加入选中的聊天室
const joinSelectedRoom = async () => {
  if (!joinRoomForm.selectedRoomId) {
    ElMessage.error('请选择一个聊天室')
    return
  }
  
  joiningRoom.value = true
  try {
    await joinRoom(joinRoomForm.selectedRoomId)
    
    // 更新聊天室列表
    await chatStore.fetchRooms()
    
    // 自动选择并进入刚加入的聊天室
    await selectRoom(joinRoomForm.selectedRoomId)
    
    ElMessage.success('成功加入聊天室')
    showJoinRoomDialog.value = false
    
    // 清空表单和搜索结果
    handleJoinDialogClose()
  } catch (error) {
    console.error('Failed to join room:', error)
  } finally {
    joiningRoom.value = false
  }
}

// 处理加入聊天室对话框关闭
const handleJoinDialogClose = () => {
  // 清除搜索定时器
  if (searchTimer) {
    clearTimeout(searchTimer)
    searchTimer = null
  }
  
  // 重置表单和搜索结果
  Object.assign(joinRoomForm, {
    searchType: 'name',
    searchValue: '',
    selectedRoomId: null
  })
  searchResults.value = []
  searchLoading.value = false
}

const handleUserAction = (command) => {
  if (command === 'logout') {
    logout()
  }
}

const logout = () => {
  // 断开WebSocket连接
  wsClient.disconnect()
  
  // 清除登录状态
  authStore.logout()
  chatStore.clearChatData()
  
  // 跳转到登录页
  router.push('/login')
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date
  
  if (diff < 60000) { // 1分钟内
    return '刚刚'
  } else if (diff < 3600000) { // 1小时内
    return `${Math.floor(diff / 60000)}分钟前`
  } else if (date.toDateString() === now.toDateString()) { // 今天
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  } else {
    return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
}

// WebSocket事件处理
const setupWebSocketEvents = () => {
  // 连接成功
  wsClient.on('connected', () => {
    console.log('WebSocket connected')
  })
  
  // 接收消息
  wsClient.on('message', (message) => {
    chatStore.addMessage(message.room_id, message)
    nextTick(() => {
      scrollToBottom()
    })
  })
  
  // 系统消息
  wsClient.on('system', (message) => {
    chatStore.addMessage(message.room_id, message)
    nextTick(() => {
      scrollToBottom()
    })
  })
  
  // 用户列表更新
  wsClient.on('userlist', (message) => {
    chatStore.updateRoomMembers(message.room_id, message.users)
  })
  
  // 连接断开
  wsClient.on('disconnected', () => {
    console.log('WebSocket disconnected')
  })
}

// 生命周期
onMounted(async () => {
  try {
    // 获取聊天室列表
    await chatStore.fetchRooms()
    
    // 连接WebSocket
    const token = authStore.token
    if (token) {
      wsClient.connect(token)
      setupWebSocketEvents()
    }
  } catch (error) {
    console.error('Failed to initialize chat:', error)
  }
})

onUnmounted(() => {
  wsClient.disconnect()
})

// 监听消息变化，自动滚动到底部
watch(currentMessages, () => {
  nextTick(() => {
    scrollToBottom()
  })
}, { deep: true })
</script>

<style scoped>
.chat-container {
  height: 100vh;
  display: flex;
  background: #f5f5f5;
}

/* 侧边栏 */
.sidebar {
  width: 300px;
  background: white;
  border-right: 1px solid #e6e6e6;
  display: flex;
  flex-direction: column;
}

.user-header {
  padding: 15px;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-details {
  flex: 1;
}

.username {
  font-weight: 600;
  font-size: 14px;
  color: #333;
}

.status {
  font-size: 12px;
  color: #67c23a;
}

.header-actions {
  display: flex;
  gap: 5px;
}

.search-bar {
  padding: 15px;
  border-bottom: 1px solid #e6e6e6;
}

.sidebar-tabs {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

:deep(.el-tabs__content) {
  flex: 1;
  overflow: hidden;
}

:deep(.el-tab-pane) {
  height: 100%;
  overflow-y: auto;
}

.room-list, .contact-list {
  padding: 0;
}

.room-item, .contact-item {
  display: flex;
  align-items: center;
  padding: 12px 15px;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.3s;
}

.room-item:hover, .contact-item:hover {
  background: #f5f5f5;
}

.room-item.active {
  background: #e6f7ff;
  border-right: 3px solid #1890ff;
}

.room-info, .contact-info {
  flex: 1;
  margin-left: 10px;
  min-width: 0;
}

.room-name, .contact-name {
  font-weight: 500;
  font-size: 14px;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.room-desc {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.contact-status {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
}

.contact-status.online {
  color: #67c23a;
}

.room-meta {
  text-align: right;
}

.member-count {
  font-size: 12px;
  color: #999;
}

/* 主聊天区域 */
.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.chat-placeholder {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.placeholder-content {
  text-align: center;
  color: #999;
}

.placeholder-content h3 {
  margin: 20px 0 10px;
  color: #666;
}

.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
}

.chat-header {
  padding: 15px 20px;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chat-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.chat-details {
  flex: 1;
}

.chat-name {
  font-weight: 600;
  font-size: 16px;
  color: #333;
}

.chat-members {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #fafafa;
}

.message-item {
  margin-bottom: 15px;
}

.message-item.system-message {
  text-align: center;
}

.system-content {
  display: inline-block;
  padding: 5px 10px;
  background: #e6e6e6;
  border-radius: 12px;
  font-size: 12px;
  color: #666;
}

.message-content {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.message-item.own-message .message-content {
  flex-direction: row-reverse;
}

.message-bubble {
  max-width: 60%;
  background: white;
  border-radius: 12px;
  padding: 10px 15px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.message-item.own-message .message-bubble {
  background: #1890ff;
  color: white;
}

.message-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 5px;
}

.message-username {
  font-size: 12px;
  font-weight: 500;
  color: #666;
}

.message-item.own-message .message-username {
  color: rgba(255,255,255,0.8);
}

.message-time {
  font-size: 11px;
  color: #999;
}

.message-item.own-message .message-time {
  color: rgba(255,255,255,0.6);
}

.message-text {
  font-size: 14px;
  line-height: 1.4;
  word-wrap: break-word;
}

.input-area {
  padding: 20px;
  border-top: 1px solid #e6e6e6;
  background: white;
}

.input-container {
  display: flex;
  gap: 10px;
  align-items: flex-end;
}

.input-container :deep(.el-textarea) {
  flex: 1;
}

.input-actions {
  display: flex;
  gap: 10px;
}

/* 对话框样式 */
.members-list {
  max-height: 400px;
  overflow-y: auto;
}

.member-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}

.member-item:last-child {
  border-bottom: none;
}

.member-info {
  flex: 1;
  margin-left: 10px;
}

.member-name {
  font-weight: 500;
  font-size: 14px;
  color: #333;
}

.member-status {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
}

.member-status.online {
  color: #67c23a;
}

/* 加入聊天室对话框样式 */
.search-results {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #e6e6e6;
  border-radius: 6px;
}

.search-result-item {
  display: flex;
  align-items: center;
  padding: 10px;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.3s;
}

.search-result-item:last-child {
  border-bottom: none;
}

.search-result-item:hover {
  background: #f5f5f5;
}

.search-result-item.selected {
  background: #e6f7ff;
  border-color: #1890ff;
}

.result-info {
  flex: 1;
  margin-left: 10px;
  min-width: 0;
}

.result-name {
  font-weight: 500;
  font-size: 14px;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.result-desc {
  font-size: 12px;
  color: #666;
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.result-members {
  font-size: 11px;
  color: #999;
  margin-top: 2px;
}
</style> 