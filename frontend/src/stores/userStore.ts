import { defineStore } from 'pinia'
import { ref } from 'vue'

interface UserInfo {
  id?: string
  name?: string
  username?: string
  email?: string
  phone?: string
  nickname?: string
  role_id?: number
  is_super_admin?: boolean
  avatar?: string
  created_at?: string
  updated_at?: string
}

export const useUserStore = defineStore('user', () => {
  const isLogin = ref(false)
  const userInfo = ref<UserInfo>({})
  const accessToken = ref('')
  const refreshToken = ref('')

  function setLogin(status: boolean) {
    isLogin.value = status
  }

  function setUserInfo(info: UserInfo) {
    userInfo.value = { ...userInfo.value, ...info }
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
    key: 'user',
    storage: localStorage,
  },
})
