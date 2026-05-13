<template>
  <header class="header">
    <div class="header-content">
      <div class="logo-section" @click="router.push('/')">
        <h1 class="logo">全景漫游平台</h1>
      </div>

      <nav class="nav-menu">
        <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }">首页</router-link>
        <router-link to="/favorites" class="nav-item" :class="{ active: $route.path === '/favorites' }">
          <StarIcon />
          收藏
          <span v-if="favoriteStore.count > 0" class="badge">{{ favoriteStore.count }}</span>
        </router-link>
        <router-link to="/history" class="nav-item" :class="{ active: $route.path === '/history' }">
          <HistoryIcon />
          历史
        </router-link>
      </nav>

      <div class="auth-section">
        <template v-if="userStore.isLogin">
          <t-dropdown>
            <div class="user-info">
              <div class="user-avatar">
                {{ (userStore.userInfo.nickname || userStore.userInfo.username || 'U').charAt(0).toUpperCase() }}
              </div>
              <span class="user-name">{{ userStore.userInfo.nickname || userStore.userInfo.username || '用户' }}</span>
            </div>
            <template #dropdown>
              <t-dropdown-menu>
                <t-dropdown-item @click="router.push('/user/center')">
                  <UserIcon /> 个人中心
                </t-dropdown-item>
                <t-dropdown-item v-if="isAdmin" @click="router.push('/admin')">
                  <SettingIcon /> 后台管理
                </t-dropdown-item>
                <t-dropdown-item @click="handleLogout">
                  <LogoutIcon /> 退出登录
                </t-dropdown-item>
              </t-dropdown-menu>
            </template>
          </t-dropdown>
        </template>
        <template v-else>
          <t-button theme="default" variant="outline" size="small" @click="showLoginDialog = true">
            登录
          </t-button>
          <t-button theme="primary" size="small" @click="showRegisterDialog = true">
            注册
          </t-button>
        </template>
      </div>
    </div>

    <t-dialog
      v-model:visible="showLoginDialog"
      header="登录"
      width="400px"
      :footer="false"
    >
      <t-form :data="loginForm" @submit="handleLogin">
        <t-form-item>
          <t-input
            v-model="loginForm.username"
            placeholder="请输入用户名或邮箱"
            size="large"
            clearable
          >
            <template #prefix-icon><UserIcon /></template>
          </t-input>
        </t-form-item>
        <t-form-item>
          <t-input
            v-model="loginForm.password"
            placeholder="请输入密码"
            type="password"
            size="large"
            clearable
          >
            <template #prefix-icon><LockOnIcon /></template>
          </t-input>
        </t-form-item>
        <t-form-item>
          <t-button theme="primary" type="submit" block size="large" :loading="loginLoading">
            登录
          </t-button>
        </t-form-item>
      </t-form>
      <div class="dialog-footer">
        <span class="link-text" @click="showLoginDialog = false; showRegisterDialog = true">
          没有账号？去注册
        </span>
      </div>
    </t-dialog>

    <t-dialog
      v-model:visible="showRegisterDialog"
      header="注册"
      width="400px"
      :footer="false"
    >
      <t-form :data="registerForm" :rules="registerRules" ref="registerFormRef" @submit="handleRegister">
        <t-form-item name="username">
          <t-input
            v-model="registerForm.username"
            placeholder="请输入用户名"
            size="large"
            clearable
          >
            <template #prefix-icon><UserIcon /></template>
          </t-input>
        </t-form-item>
        <t-form-item name="email">
          <t-input
            v-model="registerForm.email"
            placeholder="请输入邮箱"
            size="large"
            clearable
          >
            <template #prefix-icon><MailIcon /></template>
          </t-input>
        </t-form-item>
        <t-form-item name="password">
          <t-input
            v-model="registerForm.password"
            placeholder="请输入密码"
            type="password"
            size="large"
            clearable
          >
            <template #prefix-icon><LockOnIcon /></template>
          </t-input>
        </t-form-item>
        <t-form-item name="confirmPassword">
          <t-input
            v-model="registerForm.confirmPassword"
            placeholder="请确认密码"
            type="password"
            size="large"
            clearable
          >
            <template #prefix-icon><LockOnIcon /></template>
          </t-input>
        </t-form-item>
        <t-form-item>
          <t-button theme="primary" type="submit" block size="large" :loading="registerLoading">
            注册
          </t-button>
        </t-form-item>
      </t-form>
      <div class="dialog-footer">
        <span class="link-text" @click="showRegisterDialog = false; showLoginDialog = true">
          已有账号？去登录
        </span>
      </div>
    </t-dialog>
  </header>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { 
  UserIcon, 
  LockOnIcon, 
  MailIcon, 
  SettingIcon, 
  LogoutIcon,
  StarIcon,
  HistoryIcon
} from 'tdesign-icons-vue-next'
import UserApi from '@/apis/user.api'
import { useUserStore } from '@/stores/user.store'
import { useFavoriteStore } from '@/stores/favorite.store'

const router = useRouter()
const userStore = useUserStore()
const favoriteStore = useFavoriteStore()

const isAdmin = computed(() => {
  return userStore.userInfo.role?.name === 'admin' || userStore.userInfo.is_super_admin === true
})

const showLoginDialog = ref(false)
const showRegisterDialog = ref(false)
const registerFormRef = ref()

const loginLoading = ref(false)
const loginForm = reactive({
  username: '',
  password: ''
})

const registerLoading = ref(false)
const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const registerRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { email: true, message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (val: string) => val === registerForm.password,
      message: '两次密码输入不一致',
      trigger: 'blur'
    }
  ]
}

const handleLogin = async () => {
  if (!loginForm.username || !loginForm.password) {
    MessagePlugin.warning('请输入用户名和密码')
    return
  }

  loginLoading.value = true
  try {
    const result = await UserApi.login(loginForm.username, loginForm.password)
    if (result.code === 200) {
      userStore.setLogin(true)
      if (result.data?.user) {
        userStore.setUserInfo({
          id: String(result.data.user.id || ''),
          username: result.data.user.username,
          email: result.data.user.email,
          phone: result.data.user.phone,
          nickname: result.data.user.nickname,
          role_id: result.data.user.role_id,
          is_super_admin: result.data.user.is_super_admin
        })
      }
      if (result.data?.access_token) {
        userStore.setToken(result.data.access_token, result.data.refresh_token)
      }
      showLoginDialog.value = false
      MessagePlugin.success('登录成功')
      loginForm.username = ''
      loginForm.password = ''
    } else {
      MessagePlugin.error(result.msg || '登录失败')
    }
  } catch (error: any) {
    MessagePlugin.error(error.message || '登录失败，请稍后重试')
  } finally {
    loginLoading.value = false
  }
}

const handleRegister = async () => {
  const valid = await registerFormRef.value?.validate()
  if (valid !== true) return

  registerLoading.value = true
  try {
    const result = await UserApi.register({
      username: registerForm.username,
      email: registerForm.email,
      password: registerForm.password,
      role_id: 2,
    })
    if (result.code === 200) {
      MessagePlugin.success('注册成功，请登录')
      showRegisterDialog.value = false
      showLoginDialog.value = true
      registerForm.username = ''
      registerForm.email = ''
      registerForm.password = ''
      registerForm.confirmPassword = ''
    } else {
      MessagePlugin.error(result.msg || '注册失败')
    }
  } catch (error: any) {
    MessagePlugin.error(error.message || '注册失败，请稍后重试')
  } finally {
    registerLoading.value = false
  }
}

const handleLogout = async () => {
  try {
    await UserApi.logout()
  } catch (e) {
    // ignore
  }
  userStore.logout()
  MessagePlugin.success('已退出登录')
  router.push('/login')
}
</script>

<style scoped>
.header {
  position: relative;
  flex-shrink: 0;
  background: rgba(15, 24, 38, 0.4); /* 降低透明度以显示背景动画 */
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  height: 64px;
  display: flex;
  align-items: center;
  backdrop-filter: blur(12px) saturate(1.5); /* 增加饱和度 */
  z-index: 1000;
}

.header-content {
  width: 100%;
  max-width: 1440px;
  margin: 0 auto;
  padding: 0 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  gap: 32px;
}

.logo-section {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.logo-section:hover {
  transform: scale(1.02);
}

.logo {
  font-size: 22px;
  font-weight: 800;
  background: linear-gradient(135deg, #38bdf8, #0EA5E9, #4080ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
  letter-spacing: -0.02em;
  text-shadow: 0 4px 12px rgba(14, 165, 233, 0.3);
}

.nav-menu {
  display: flex;
  align-items: center;
  flex: 1;
  justify-content: flex-start;
  gap: 12px;
}

.nav-item {
  position: relative;
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  padding: 8px 20px;
  border-radius: 20px;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
}

.nav-item:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.08);
}

.nav-item.active {
  color: #fff;
  background: rgba(14, 165, 233, 0.15);
  box-shadow: inset 0 0 0 1px rgba(14, 165, 233, 0.3);
}

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  background: #ef4444; /* 改为醒目的红色 */
  border-radius: 9px;
  font-size: 10px;
  font-weight: 700;
  color: #fff;
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.5);
}

.auth-section {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 4px 12px 4px 4px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  transition: all 0.3s;
}

.user-info:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(14, 165, 233, 0.4);
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #38bdf8, #4080ff);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 14px;
  font-weight: 700;
  box-shadow: 0 2px 8px rgba(14, 165, 233, 0.4);
}

.user-name {
  color: rgba(255, 255, 255, 0.9);
  font-size: 14px;
  font-weight: 600;
}

:deep(.t-button--variant-outline) {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.15);
  color: #fff;
  border-radius: 20px;
}

:deep(.t-button--variant-outline:hover) {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.3);
}

:deep(.t-button--theme-primary) {
  background: linear-gradient(135deg, #0EA5E9, #4080ff);
  border: none;
  border-radius: 20px;
  box-shadow: 0 4px 12px rgba(14, 165, 233, 0.3);
}

.dialog-footer {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.link-text {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: all 0.2s;
}

.link-text:hover {
  color: #0EA5E9;
  text-decoration: none;
}

:deep(.t-dropdown-menu) {
  background: rgba(30, 41, 59, 0.85);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 6px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.4);
}

:deep(.t-dropdown-item) {
  color: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  gap: 8px;
  border-radius: 8px;
  margin: 2px 0;
}

:deep(.t-dropdown-item:hover) {
  background: rgba(14, 165, 233, 0.15);
  color: #fff;
}
</style>
