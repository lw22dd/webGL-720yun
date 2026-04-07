<template>
  <header class="header">
    <div class="header-content">
      <div class="logo-section">
        <h1 class="logo">720云</h1>
      </div>

      <nav class="nav-menu">
        <a href="#" class="nav-item active">首页</a>
        <t-dropdown>
          <span class="nav-item">产品服务</span>
          <template #dropdown>
            <t-dropdown-menu>
              <t-dropdown-item>VR全景创作</t-dropdown-item>
              <t-dropdown-item>数字孪生</t-dropdown-item>
              <t-dropdown-item>元宇宙平台</t-dropdown-item>
              <t-dropdown-item>VR直播</t-dropdown-item>
            </t-dropdown-menu>
          </template>
        </t-dropdown>
        <t-dropdown>
          <span class="nav-item">解决方案</span>
          <template #dropdown>
            <t-dropdown-menu>
              <t-dropdown-item>教育行业</t-dropdown-item>
              <t-dropdown-item>房地产</t-dropdown-item>
              <t-dropdown-item>旅游景区</t-dropdown-item>
              <t-dropdown-item>企业展示</t-dropdown-item>
            </t-dropdown-menu>
          </template>
        </t-dropdown>
        <t-dropdown>
          <span class="nav-item">内容社区</span>
          <template #dropdown>
            <t-dropdown-menu>
              <t-dropdown-item>作品展示</t-dropdown-item>
              <t-dropdown-item>教程中心</t-dropdown-item>
              <t-dropdown-item>社区论坛</t-dropdown-item>
            </t-dropdown-menu>
          </template>
        </t-dropdown>
        <t-dropdown>
          <span class="nav-item">定制服务</span>
          <template #dropdown>
            <t-dropdown-menu>
              <t-dropdown-item>VR内容定制</t-dropdown-item>
              <t-dropdown-item>数字孪生定制</t-dropdown-item>
              <t-dropdown-item>平台定制开发</t-dropdown-item>
            </t-dropdown-menu>
          </template>
        </t-dropdown>
      </nav>

      <div class="search-section">
        <t-input
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
          <t-dropdown>
            <span class="user-info">
              {{ userStore.userInfo.email || userStore.userInfo.name || '用户' }}
            </span>
            <template #dropdown>
              <t-dropdown-menu>
                <t-dropdown-item @click="handleLogout">退出登录</t-dropdown-item>
              </t-dropdown-menu>
            </template>
          </t-dropdown>
        </template>
        <template v-else>
          <t-button theme="default" variant="outline" size="small" @click="showLoginDialog = true" class="login-btn">登录</t-button>
          <t-button theme="primary" size="small" @click="showRegisterDialog = true" class="register-btn">注册</t-button>
        </template>
      </div>

      <t-dialog v-model="showLoginDialog" header="登录" width="420px" :footer="null">
        <div class="dialog-content">
          <t-form :data="loginForm" :rules="loginRules" ref="loginFormRef" label-width="0">
            <t-form-item name="username">
              <t-input
                v-model="loginForm.username"
                placeholder="请输入用户名或邮箱"
                size="large"
              />
            </t-form-item>
            <t-form-item name="password">
              <t-input
                v-model="loginForm.password"
                placeholder="请输入密码"
                type="password"
                size="large"
              />
            </t-form-item>
            <t-form-item>
              <t-button theme="primary" block size="large" :loading="loginLoading" @click="handleLogin" class="submit-btn">
                登录
              </t-button>
            </t-form-item>
          </t-form>
          <div class="dialog-footer">
            <span class="link-text" @click="showLoginDialog = false; showRegisterDialog = true">没有账号？去注册</span>
          </div>
        </div>
      </t-dialog>

      <t-dialog v-model="showRegisterDialog" header="注册" width="420px" :footer="null">
        <div class="dialog-content">
          <t-form :data="registerForm" :rules="registerRules" ref="registerFormRef" label-width="0">
            <t-form-item name="username">
              <t-input
                v-model="registerForm.username"
                placeholder="请输入用户名"
                size="large"
              />
            </t-form-item>
            <t-form-item name="email">
              <t-input
                v-model="registerForm.email"
                placeholder="请输入邮箱"
                size="large"
              />
            </t-form-item>
            <t-form-item name="password">
              <t-input
                v-model="registerForm.password"
                placeholder="请输入密码"
                type="password"
                size="large"
              />
            </t-form-item>
            <t-form-item name="confirmPassword">
              <t-input
                v-model="registerForm.confirmPassword"
                placeholder="请确认密码"
                type="password"
                size="large"
              />
            </t-form-item>
            <t-form-item>
              <t-button theme="primary" block size="large" :loading="registerLoading" @click="handleRegister" class="submit-btn">
                注册
              </t-button>
            </t-form-item>
          </t-form>
          <div class="dialog-footer">
            <span class="link-text" @click="showRegisterDialog = false; showLoginDialog = true">已有账号？去登录</span>
          </div>
        </div>
      </t-dialog>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Message } from 'tdesign-vue-next'
import UserApi from '@/apis/userApi'
import { useUserStore } from '@/stores/userStore'

const searchKeyword = ref('')

const handleSearch = () => {
    if (searchKeyword.value.trim()) {
        Message.info(`搜索内容：${searchKeyword.value}`)
    } else {
        Message.warning('请输入搜索内容')
    }
}

const userStore = useUserStore()

const showLoginDialog = ref(false)
const showRegisterDialog = ref(false)

const loginFormRef = ref()
const loginLoading = ref(false)
const loginForm = reactive({
    username: '',
    password: ''
})

const loginRules = {
    username: [
        { required: true, message: '请输入用户名或邮箱', trigger: 'blur' }
    ],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
    ]
}

const registerFormRef = ref()
const registerLoading = ref(false)
const registerForm = reactive({
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
})

const registerRules = {
    username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 2, max: 20, message: '用户名长度应在2-20个字符之间', trigger: 'blur' }
    ],
    email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
    ],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
    ],
    confirmPassword: [
        { required: true, message: '请确认密码', trigger: 'blur' },
        {
            validator: (rule: any, value: any, callback: any) => {
                if (value === '') {
                    callback(new Error('请确认密码'))
                } else if (value !== registerForm.password) {
                    callback(new Error('两次输入密码不一致'))
                } else {
                    callback()
                }
            },
            trigger: 'blur'
        }
    ]
}

const handleLogin = async () => {
    if (!loginFormRef.value) return

    const valid = await (loginFormRef.value as any).validate()
    if (valid) {
        loginLoading.value = true
        try {
            const result = await UserApi.login(loginForm.username, loginForm.password)
            console.log('login result', result)
            if (result.code === 200) {
                userStore.setLogin(true)

                userStore.setUserInfo({
                    id: result.data?.user_id,
                    email: loginForm.username
                })

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
}

const handleRegister = async () => {
    if (!registerFormRef.value) return

    const valid = await (registerFormRef.value as any).validate()
    if (valid) {
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
}
</script>

<style scoped>
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: #ffffff;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  height: 64px;
  display: flex;
  align-items: center;
  backdrop-filter: saturate(180%) blur(20px);
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
}

.logo {
  font-size: 20px;
  font-weight: 700;
  color: #0052d9;
  margin: 0;
  letter-spacing: -0.02em;
}

.nav-menu {
  display: flex;
  align-items: center;
  flex: 1;
  justify-content: flex-start;
  gap: 4px;
}

.nav-item {
  position: relative;
  font-size: 14px;
  font-weight: 500;
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
}

.auth-section {
  display: flex;
  align-items: center;
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
}

.dialog-footer {
  margin-top: 16px;
  display: flex;
  justify-content: center;
}

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
}

.link-text:hover {
  color: #4080ff;
  text-decoration: underline;
}
</style>
