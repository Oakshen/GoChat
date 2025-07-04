import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, getUserInfo } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || '')
  
  // 计算属性
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  
  // 设置token
  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }
  
  // 设置用户信息
  const setUser = (userInfo) => {
    user.value = userInfo
  }
  
  // 登录
  const login = async (credentials) => {
    try {
      const response = await loginApi(credentials)
      const { token: newToken, user: userInfo } = response.data
      
      setToken(newToken)
      setUser(userInfo)
      
      return response
    } catch (error) {
      throw error
    }
  }
  
  // 获取用户信息
  const fetchUserInfo = async () => {
    try {
      const response = await getUserInfo()
      setUser(response.data)
      return response
    } catch (error) {
      // 如果获取用户信息失败，清除登录状态
      logout()
      throw error
    }
  }
  
  // 注销
  const logout = () => {
    user.value = null
    token.value = ''
    localStorage.removeItem('token')
  }
  
  // 初始化认证状态
  const initAuth = async () => {
    if (token.value) {
      try {
        await fetchUserInfo()
      } catch (error) {
        // 如果token失效，清除登录状态
        logout()
      }
    }
  }
  
  return {
    user,
    token,
    isAuthenticated,
    setToken,
    setUser,
    login,
    fetchUserInfo,
    logout,
    initAuth
  }
}) 