<template>
  <header class="header">
    <div class="header-content">
<<<<<<< HEAD
      <div class="logo-section">
        <h1 class="logo">720云</h1>
      </div>

      <nav class="nav-menu">
        <a href="#" class="nav-item active">首页</a>
      <Dropdown>
        <span class="nav-item">产品服务</span>
        <template #dropdown>
          <DropdownMenu>
            <DropdownItem>VR全景创作</DropdownItem>
            <DropdownItem>数字孪生</DropdownItem>
            <DropdownItem>元宇宙平台</DropdownItem>
            <DropdownItem>VR直播</DropdownItem>
          </DropdownMenu>
        </template>
      </Dropdown>
      <Dropdown>
        <span class="nav-item">解决方案</span>
        <template #dropdown>
          <DropdownMenu>
            <DropdownItem>教育行业</DropdownItem>
            <DropdownItem>房地产</DropdownItem>
            <DropdownItem>旅游景区</DropdownItem>
            <DropdownItem>企业展示</DropdownItem>
          </DropdownMenu>
        </template>
      </Dropdown>
      <Dropdown>
        <span class="nav-item">内容社区</span>
        <template #dropdown>
          <DropdownMenu>
            <DropdownItem>作品展示</DropdownItem>
            <DropdownItem>教程中心</DropdownItem>
            <DropdownItem>社区论坛</DropdownItem>
          </DropdownMenu>
        </template>
      </Dropdown>
      <Dropdown>
        <span class="nav-item">定制服务</span>
        <template #dropdown>
          <DropdownMenu>
            <DropdownItem>VR内容定制</DropdownItem>
            <DropdownItem>数字孪生定制</DropdownItem>
            <DropdownItem>平台定制开发</DropdownItem>
          </DropdownMenu>
        </template>
      </Dropdown>
      </nav>

      <div class="search-section">
        <Input
          v-model="searchKeyword"
          placeholder="搜索"
          size="small"
          class="search-input"
          clearable
          @enter="handleSearch"
        />
      </div>

      <div class="auth-section">
        <template v-if="userStore.isLogin">
          <Dropdown>
            <span class="user-info">
              {{ userStore.userInfo.email || userStore.userInfo.name || '用户' }}
            </span>
            <template #dropdown>
              <DropdownMenu>
                <DropdownItem @click="navigateToUserCenter">用户中心</DropdownItem>
                <DropdownItem v-if="isAdmin" @click="navigateToAdmin">后台管理</DropdownItem>
                <DropdownItem @click="handleLogout">退出登录</DropdownItem>
              </DropdownMenu>
            </template>
          </Dropdown>
        </template>
        <template v-else>
          <Button theme="default" variant="outline" size="small" @click="showLoginDialog = true" class="login-btn">登录</Button>
          <Button theme="primary" size="small" @click="showRegisterDialog = true" class="register-btn">注册</Button>
        </template>
      </div>

      <Dialog v-model:visible="showLoginDialog" header="登录" width="420px" :footer="false" attach="body">
        <div class="dialog-content">
          <Form :data="loginForm" label-width="0">
            <FormItem name="username">
              <Input
                v-model="loginForm.username"
                placeholder="请输入用户名或邮箱"
                size="large"
              />
            </FormItem>
            <FormItem name="password">
              <Input
                v-model="loginForm.password"
                placeholder="请输入密码"
                type="password"
                size="large"
              />
            </FormItem>
            <FormItem>
              <Button theme="primary" block size="large" :loading="loginLoading" @click="handleLogin" class="submit-btn">
                登录
              </Button>
            </FormItem>
          </Form>
          <div class="dialog-footer">
            <span class="link-text" @click="showLoginDialog = false; showRegisterDialog = true">没有账号？去注册</span>
          </div>
        </div>
      </Dialog>

      <Dialog v-model:visible="showRegisterDialog" header="注册" width="420px" :footer="false" attach="body">
        <div class="dialog-content">
          <Form :data="registerForm" label-width="0">
            <FormItem name="username">
              <Input
                v-model="registerForm.username"
                placeholder="请输入用户名"
                size="large"
              />
            </FormItem>
            <FormItem name="email">
              <Input
                v-model="registerForm.email"
                placeholder="请输入邮箱"
                size="large"
              />
            </FormItem>
            <FormItem name="password">
              <Input
                v-model="registerForm.password"
                placeholder="请输入密码"
                type="password"
                size="large"
              />
            </FormItem>
            <FormItem name="confirmPassword">
              <Input
                v-model="registerForm.confirmPassword"
                placeholder="请确认密码"
                type="password"
                size="large"
              />
            </FormItem>
            <FormItem>
              <Button theme="primary" block size="large" :loading="registerLoading" @click="handleRegister" class="submit-btn">
                注册
              </Button>
            </FormItem>
          </Form>
          <div class="dialog-footer">
            <span class="link-text" @click="showRegisterDialog = false; showLoginDialog = true">已有账号？去登录</span>
          </div>
        </div>
      </Dialog>
    </div>
=======
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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
  </header>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
<<<<<<< HEAD
import {
  Button,
  Dialog,
  Dropdown,
  DropdownMenu,
  DropdownItem,
  Form,
  FormItem,
  Input,
  Message
} from 'tdesign-vue-next'
import UserApi from '@/apis/userApi'
import { useUserStore } from '@/stores/userStore'

const router = useRouter()
const searchKeyword = ref('')

const handleSearch = () => {
    if (searchKeyword.value.trim()) {
        Message.info(`搜索内容：${searchKeyword.value}`)
    } else {
        Message.warning('请输入搜索内容')
    }
}

const userStore = useUserStore()

const isAdmin = computed(() => {
    return userStore.userInfo.is_super_admin === true
})

const navigateToUserCenter = () => {
    router.push('/user/center')
}

const navigateToAdmin = () => {
    router.push('/admin')
}

const showLoginDialog = ref(false)
const showRegisterDialog = ref(false)

const loginLoading = ref(false)
const loginForm = reactive({
    username: '',
    password: ''
=======
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
import UserApi from '@/services/api/user.api'
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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
})

const registerLoading = ref(false)
const registerForm = reactive({
<<<<<<< HEAD
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
})

const handleLogin = async () => {
    if (!loginForm.username || !loginForm.password) {
        Message.warning('请输入用户名和密码')
        return
    }

    loginLoading.value = true
    try {
        const result = await UserApi.login(loginForm.username, loginForm.password)
        if (result.code === 200) {
            userStore.setLogin(true)

            if (result.data?.user) {
                userStore.setUserInfo(result.data.user)
            } else {
                userStore.setUserInfo({
                    id: result.data?.user_id,
                    email: loginForm.username
                })
            }

            if (result.data && result.data?.access_token) {
                userStore.setToken(result.data.access_token, result.data.refresh_token)
            }

            showLoginDialog.value = false
            Message.success('登录成功')

            loginForm.username = ''
            loginForm.password = ''
        } else {
            Message.error(result.msg || '登录失败')
        }
    } catch (error: any) {
        Message.error(error.message || '登录失败，请稍后重试')
    } finally {
        loginLoading.value = false
    }
}

const handleRegister = async () => {
    if (!registerForm.username || !registerForm.email || !registerForm.password) {
        Message.warning('请填写完整信息')
        return
    }

    registerLoading.value = true
    try {
        const result = await UserApi.register({
            username: registerForm.username,
            email: registerForm.email,
            password: registerForm.password,
            role_id: 2,
            status: 1
        })
        if (result.code === 200) {
            Message.success('注册成功，请登录')
            showRegisterDialog.value = false
            showLoginDialog.value = true

            registerForm.username = ''
            registerForm.email = ''
            registerForm.password = ''
            registerForm.confirmPassword = ''
        } else {
            Message.error(result.msg || '注册失败')
        }
    } catch (error: any) {
        Message.error(error.message || '注册失败，请稍后重试')
    } finally {
        registerLoading.value = false
    }
}

const handleLogout = async () => {
    try {
        const result = await UserApi.logout()
        if (result.code === 200) {
            userStore.logout()
            Message.success('退出成功')
        } else {
            userStore.logout()
            Message.error(result.msg || '退出失败')
        }
    } catch (error: any) {
        userStore.logout()
        Message.success('已退出登录')
    }
=======
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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
}
</script>

<style scoped>
.header {
  position: relative;
  flex-shrink: 0;
<<<<<<< HEAD
  background: #ffffff;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  height: 64px;
  display: flex;
  align-items: center;
  backdrop-filter: saturate(180%) blur(20px);
=======
  background: rgba(15, 24, 38, 0.95);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  height: 64px;
  display: flex;
  align-items: center;
  backdrop-filter: blur(12px);
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
  z-index: 1000;
}

.header-content {
  width: 100%;
  max-width: 1440px;
  margin: 0 auto;
  padding: 0 24px;
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
<<<<<<< HEAD
=======
  cursor: pointer;
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
}

.logo {
  font-size: 20px;
  font-weight: 700;
<<<<<<< HEAD
  color: #0052d9;
=======
  background: linear-gradient(135deg, #0EA5E9, #4080ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
  margin: 0;
  letter-spacing: -0.02em;
}

.nav-menu {
  display: flex;
  align-items: center;
  flex: 1;
  justify-content: flex-start;
<<<<<<< HEAD
  gap: 4px;
=======
  gap: 8px;
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
}

.nav-item {
  position: relative;
  font-size: 14px;
  font-weight: 500;
<<<<<<< HEAD
  color: #4e5969;
  cursor: pointer;
  padding: 8px 16px;
  border-radius: 6px;
  transition: all 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
  display: flex;
  align-items: center;
  user-select: none;
}

.nav-item:hover {
  color: #0052d9;
  background: rgba(0, 82, 217, 0.06);
}

.nav-item.active {
  color: #0052d9;
  font-weight: 600;
}

.search-section {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  width: 240px;
}

.search-input {
  width: 100%;
=======
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  padding: 8px 16px;
  border-radius: 8px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
}

.nav-item:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}

.nav-item.active {
  color: #fff;
  background: rgba(14, 165, 233, 0.2);
}

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  background: #0EA5E9;
  border-radius: 9px;
  font-size: 11px;
  font-weight: 600;
  color: #fff;
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
}

.auth-section {
  display: flex;
  align-items: center;
<<<<<<< HEAD
  gap: 8px;
  flex-shrink: 0;
}

.login-btn {
  font-weight: 500;
  font-size: 14px;
  border-radius: 6px;
  padding: 0 20px;
  border-color: #0052d9;
  color: #0052d9;
  background: transparent;
  transition: all 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
}

.login-btn:hover {
  background: #0052d9 !important;
  border-color: #0052d9 !important;
  color: #ffffff !important;
}

.register-btn {
  font-weight: 500;
  font-size: 14px;
  border-radius: 6px;
  padding: 0 20px;
  background: linear-gradient(135deg, #0052d9 0%, #4080ff 100%);
  border-color: #0052d9;
  color: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 82, 217, 0.2);
  transition: all 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
}

.register-btn:hover {
  background: linear-gradient(135deg, #4080ff 0%, #0052d9 100%);
  border-color: #4080ff;
  box-shadow: 0 4px 12px rgba(0, 82, 217, 0.3);
  transform: translateY(-1px);
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  color: #4e5969;
  transition: all 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
}

.user-info:hover {
  color: #0052d9;
  background: rgba(0, 82, 217, 0.06);
}

.dialog-content {
  padding: 8px 0;
=======
  gap: 12px;
  flex-shrink: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 6px 12px 6px 6px;
  border-radius: 24px;
  transition: all 0.2s;
}

.user-info:hover {
  background: rgba(255, 255, 255, 0.1);
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #0EA5E9, #4080ff);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
}

.user-name {
  color: #fff;
  font-size: 14px;
  font-weight: 500;
}

:deep(.t-button--variant-outline) {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.3);
  color: #fff;
}

:deep(.t-button--variant-outline:hover) {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.5);
}

:deep(.t-button--theme-primary) {
  background: linear-gradient(135deg, #0EA5E9, #4080ff);
  border: none;
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
}

.dialog-footer {
  margin-top: 16px;
  display: flex;
  justify-content: center;
}

<<<<<<< HEAD
.submit-btn {
  border-radius: 6px;
  font-weight: 500;
  margin-top: 8px;
}

.link-text {
  font-size: 14px;
  color: #0052d9;
  cursor: pointer;
  transition: color 0.2s;
  text-align: center;
=======
.link-text {
  font-size: 14px;
  color: #0EA5E9;
  cursor: pointer;
  transition: color 0.2s;
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
}

.link-text:hover {
  color: #4080ff;
  text-decoration: underline;
}
<<<<<<< HEAD
=======

:deep(.t-dropdown-menu) {
  background: rgba(26, 35, 50, 0.95);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

:deep(.t-dropdown-item) {
  color: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.t-dropdown-item:hover) {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
</style>
