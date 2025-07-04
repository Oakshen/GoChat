import http from './http'

// 用户注册
export const register = (data) => {
  return http.post('/auth/register', data)
}

// 用户登录
export const login = (data) => {
  return http.post('/auth/login', data)
}

// 获取用户信息
export const getUserInfo = () => {
  return http.get('/auth/userinfo')
}

// 获取在线用户列表
export const getOnlineUsers = () => {
  return http.get('/users/online')
} 