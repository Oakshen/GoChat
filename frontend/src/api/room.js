import http from './http'

// 创建聊天室
export const createRoom = (data) => {
  return http.post('/rooms', data)
}

// 获取聊天室列表
export const getRooms = () => {
  return http.get('/rooms')
}

// 获取聊天室详情
export const getRoomDetail = (roomId) => {
  return http.get(`/rooms/${roomId}`)
}

// 加入聊天室
export const joinRoom = (roomId) => {
  return http.post(`/rooms/${roomId}/join`)
}

// 离开聊天室
export const leaveRoom = (roomId) => {
  return http.post(`/rooms/${roomId}/leave`)
}

// 获取聊天室成员
export const getRoomMembers = (roomId) => {
  return http.get(`/rooms/${roomId}/members`)
}

// 获取聊天记录
export const getRoomMessages = (roomId, page = 1, limit = 50) => {
  return http.get(`/rooms/${roomId}/messages`, {
    params: { page, limit }
  })
}

// 根据ID搜索聊天室（可搜索用户未加入的聊天室）
export const searchRoomById = (roomId) => {
  return http.get(`/rooms/${roomId}`)
}

// 根据名称模糊搜索聊天室
export const searchRoomsByName = (name) => {
  return http.get('/rooms/search', {
    params: { name }
  })
} 