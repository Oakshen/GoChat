import { ElMessage } from 'element-plus'

class WebSocketClient {
  constructor() {
    this.ws = null
    this.reconnectTimer = null
    this.heartbeatTimer = null
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectInterval = 3000
    this.heartbeatInterval = 30000
    this.listeners = {}
  }

  // 连接WebSocket
  connect(token) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected')
      return
    }

    const wsUrl = `ws://localhost:3000/ws?token=${token}`
    console.log('Connecting to WebSocket:', wsUrl)

    try {
      this.ws = new WebSocket(wsUrl)
      this.setupEventHandlers()
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
      ElMessage.error('无法连接到服务器')
    }
  }

  // 设置事件处理器
  setupEventHandlers() {
    this.ws.onopen = (event) => {
      console.log('WebSocket connected:', event)
      this.reconnectAttempts = 0
      this.startHeartbeat()
      this.emit('connected', event)
      ElMessage.success('连接成功')
    }

    this.ws.onmessage = (event) => {
      console.log('WebSocket message received:', event.data)
      try {
        const message = JSON.parse(event.data)
        this.handleMessage(message)
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    this.ws.onclose = (event) => {
      console.log('WebSocket disconnected:', event.code, event.reason)
      this.stopHeartbeat()
      this.emit('disconnected', event)
      
      // 如果不是主动关闭，尝试重连
      if (!event.wasClean && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.scheduleReconnect()
      }
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      this.emit('error', error)
    }
  }

  // 处理接收到的消息
  handleMessage(message) {
    const { type } = message
    
    switch (type) {
      case 'text':
        this.emit('message', message)
        break
      case 'system':
        this.emit('system', message)
        break
      case 'error':
        this.emit('error', message)
        ElMessage.error(message.error || '服务器错误')
        break
      case 'userlist':
        this.emit('userlist', message)
        break
      case 'pong':
        // 处理心跳pong响应
        console.log('Received pong from server')
        this.emit('pong', message)
        break
      default:
        console.log('Unknown message type:', type, message)
        this.emit('unknown', message)
    }
  }

  // 发送消息
  send(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const messageStr = JSON.stringify(message)
      console.log('Sending WebSocket message:', messageStr)
      this.ws.send(messageStr)
    } else {
      console.error('WebSocket is not connected')
      ElMessage.error('连接已断开，请重新连接')
    }
  }

  // 发送文本消息
  sendTextMessage(roomId, content) {
    this.send({
      type: 'text',
      room_id: roomId,
      content: content
    })
  }

  // 加入聊天室
  joinRoom(roomId) {
    this.send({
      type: 'join',
      room_id: roomId
    })
  }

  // 离开聊天室
  leaveRoom(roomId) {
    this.send({
      type: 'leave',
      room_id: roomId
    })
  }

  // 发送正在输入状态
  sendTyping(roomId) {
    this.send({
      type: 'typing',
      room_id: roomId
    })
  }

  // 开始心跳
  startHeartbeat() {
    this.stopHeartbeat()
    this.heartbeatTimer = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        // 发送ping消息
        this.ws.send(JSON.stringify({ type: 'ping' }))
      }
    }, this.heartbeatInterval)
  }

  // 停止心跳
  stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }

  // 计划重连
  scheduleReconnect() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
    }

    this.reconnectTimer = setTimeout(() => {
      this.reconnectAttempts++
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
      
      const token = localStorage.getItem('token')
      if (token) {
        this.connect(token)
      }
    }, this.reconnectInterval)
  }

  // 断开连接
  disconnect() {
    this.stopHeartbeat()
    
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  // 事件监听
  on(event, callback) {
    if (!this.listeners[event]) {
      this.listeners[event] = []
    }
    this.listeners[event].push(callback)
  }

  // 移除事件监听
  off(event, callback) {
    if (this.listeners[event]) {
      const index = this.listeners[event].indexOf(callback)
      if (index > -1) {
        this.listeners[event].splice(index, 1)
      }
    }
  }

  // 触发事件
  emit(event, data) {
    if (this.listeners[event]) {
      this.listeners[event].forEach(callback => {
        try {
          callback(data)
        } catch (error) {
          console.error('Error in event callback:', error)
        }
      })
    }
  }

  // 获取连接状态
  getReadyState() {
    return this.ws ? this.ws.readyState : WebSocket.CLOSED
  }

  // 是否已连接
  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN
  }
}

// 创建单例实例
const wsClient = new WebSocketClient()

export default wsClient 