import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getRooms, getRoomMessages } from '@/api/room'

export const useChatStore = defineStore('chat', () => {
  // 状态
  const rooms = ref([])
  const currentRoomId = ref(null)
  const messages = ref({}) // roomId: [messages]
  const onlineUsers = ref([])
  const roomMembers = ref({}) // roomId: [members]
  
  // 计算属性
  const currentRoom = computed(() => {
    return rooms.value.find(room => room.id === currentRoomId.value)
  })
  
  const currentMessages = computed(() => {
    return messages.value[currentRoomId.value] || []
  })
  
  const currentRoomMembers = computed(() => {
    return roomMembers.value[currentRoomId.value] || []
  })
  
  // 设置聊天室列表
  const setRooms = (roomList) => {
    rooms.value = roomList
  }
  
  // 添加聊天室
  const addRoom = (room) => {
    const existingIndex = rooms.value.findIndex(r => r.id === room.id)
    if (existingIndex !== -1) {
      rooms.value[existingIndex] = room
    } else {
      rooms.value.push(room)
    }
  }
  
  // 设置当前聊天室
  const setCurrentRoom = (roomId) => {
    currentRoomId.value = roomId
  }
  
  // 添加消息
  const addMessage = (roomId, message) => {
    if (!messages.value[roomId]) {
      messages.value[roomId] = []
    }
    messages.value[roomId].push(message)
  }
  
  // 设置消息列表
  const setMessages = (roomId, messageList) => {
    messages.value[roomId] = messageList
  }
  
  // 设置在线用户
  const setOnlineUsers = (users) => {
    onlineUsers.value = users
  }
  
  // 设置聊天室成员
  const setRoomMembers = (roomId, members) => {
    roomMembers.value[roomId] = members
  }
  
  // 更新聊天室成员列表
  const updateRoomMembers = (roomId, members) => {
    roomMembers.value[roomId] = members
  }
  
  // 获取聊天室列表
  const fetchRooms = async () => {
    try {
      const response = await getRooms()
      setRooms(response.data || [])
      return response
    } catch (error) {
      throw error
    }
  }
  
  // 获取聊天记录
  const fetchMessages = async (roomId, page = 1, limit = 50) => {
    try {
      const response = await getRoomMessages(roomId, page, limit)
      if (page === 1) {
        setMessages(roomId, response.data || [])
      } else {
        // 加载更多消息，插入到前面
        const existingMessages = messages.value[roomId] || []
        const newMessages = response.data || []
        messages.value[roomId] = [...newMessages, ...existingMessages]
      }
      return response
    } catch (error) {
      throw error
    }
  }
  
  // 清除聊天数据
  const clearChatData = () => {
    rooms.value = []
    currentRoomId.value = null
    messages.value = {}
    onlineUsers.value = []
    roomMembers.value = {}
  }
  
  return {
    // 状态
    rooms,
    currentRoomId,
    messages,
    onlineUsers,
    roomMembers,
    
    // 计算属性
    currentRoom,
    currentMessages,
    currentRoomMembers,
    
    // 方法
    setRooms,
    addRoom,
    setCurrentRoom,
    addMessage,
    setMessages,
    setOnlineUsers,
    setRoomMembers,
    updateRoomMembers,
    fetchRooms,
    fetchMessages,
    clearChatData
  }
}) 