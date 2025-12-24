import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  // 登录状态
  const isLogin = ref(false)
  // 用户基础信息
  const userInfo = ref<{ id?: string; name?: string; email?: string }>({})
  // 认证令牌
  const accessToken = ref('')
  const refreshToken = ref('')

  function setLogin(status: boolean) {
    isLogin.value = status
  }

  function setUserInfo(info: { id?: string; name?: string; email?: string }) {
    userInfo.value = info
  }

  function setToken(accessTokenStr: string, refreshTokenStr: string) {
    accessToken.value = accessTokenStr
    refreshToken.value = refreshTokenStr
  }

  function logout() {
    isLogin.value = false
    userInfo.value = {}
    accessToken.value = ''
    refreshToken.value = ''
  }

  return { 
    isLogin, 
    userInfo, 
    accessToken, 
    refreshToken, 
    setLogin, 
    setUserInfo, 
    setToken, 
    logout 
  }
}, {
  persist: {
    // 显式指定需要持久化的字段
    key: 'user',
    storage: localStorage,
    // 使用localStorage作为存储
  },
})