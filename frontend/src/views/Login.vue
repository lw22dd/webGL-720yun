<template>
  <div class="login-page">
    <div class="login-left">
      <div class="login-left-content">
        <h1 class="login-display-title">探索世界</h1>
        <p class="login-display-subtitle">全景漫游实验教学平台</p>
        <div class="login-features">
          <div class="login-feature">
            <div class="login-feature-icon">🌍</div>
            <span class="login-feature-label">360°全景</span>
          </div>
          <div class="login-feature">
            <div class="login-feature-icon">🔬</div>
            <span class="login-feature-label">虚拟实验</span>
          </div>
          <div class="login-feature">
            <div class="login-feature-icon">📊</div>
            <span class="login-feature-label">智能评测</span>
          </div>
        </div>
      </div>
    </div>

    <div class="login-right">
      <div class="login-form-wrap">
        <h2 class="login-form-heading">欢迎回来</h2>
        <p class="login-form-subtitle">登录您的账户以继续</p>

        <div class="form-group">
          <label>用户名</label>
          <div class="form-input-wrap">
            <span class="input-icon">👤</span>
            <input
              v-model="loginForm.username"
              type="text"
              class="form-input"
              :class="{ 'error': errors.username }"
              placeholder="请输入用户名"
              @keyup.enter="handleLogin"
            />
          </div>
          <p v-if="errors.username" class="form-error-msg show">{{ errors.username }}</p>
        </div>

        <div class="form-group">
          <label>密码</label>
          <div class="form-input-wrap">
            <span class="input-icon">🔒</span>
            <input
              v-model="loginForm.password"
              :type="showPassword ? 'text' : 'password'"
              class="form-input"
              :class="{ 'error': errors.password }"
              placeholder="请输入密码"
              @keyup.enter="handleLogin"
            />
            <button class="toggle-pw" @click="showPassword = !showPassword">
              {{ showPassword ? '👁' : '👁‍🗨' }}
            </button>
          </div>
          <p v-if="errors.password" class="form-error-msg show">{{ errors.password }}</p>
        </div>

        <div class="form-group">
          <label>验证码</label>
          <div class="captcha-row">
            <div class="form-input-wrap">
              <span class="input-icon">🛡</span>
              <input
                v-model="loginForm.captcha"
                type="text"
                class="form-input"
                :class="{ 'error': errors.captcha }"
                placeholder="请输入验证码"
                @keyup.enter="handleLogin"
              />
            </div>
            <div class="captcha-img" @click="refreshCaptcha">{{ captchaText }}</div>
          </div>
          <p v-if="errors.captcha" class="form-error-msg show">{{ errors.captcha }}</p>
        </div>

        <div class="form-row-between">
          <label class="form-check">
            <input v-model="loginForm.remember" type="checkbox" />
            记住我
          </label>
          <a href="#" class="form-link" @click.prevent>忘记密码?</a>
        </div>

        <button
          class="btn btn-primary btn-block"
          :disabled="loading"
          @click="handleLogin"
        >
          <span v-if="loading" class="loading-spinner"></span>
          <span v-else>登 录</span>
        </button>

        <div class="login-footer">
          还没有账户？<a href="#" @click.prevent="showRegisterTip">立即注册</a>
        </div>
      </div>
    </div>

    <t-loading v-if="loading" :fullscreen="true" :text="loadingText" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import UserApi from '@/apis/userApi'
import { useUserStore } from '@/stores/userStore'

const router = useRouter()
const userStore = useUserStore()

const loginForm = reactive({
  username: '',
  password: '',
  captcha: '',
  remember: false
})

const errors = reactive({
  username: '',
  password: '',
  captcha: ''
})

const loading = ref(false)
const loadingText = ref('登录中...')
const showPassword = ref(false)

const captchaText = ref('')
const correctCaptcha = ref('')

const generateCaptcha = () => {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789'
  let result = ''
  for (let i = 0; i < 4; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  captchaText.value = result
  correctCaptcha.value = result.toLowerCase()
}

const refreshCaptcha = () => {
  generateCaptcha()
}

const validateForm = (): boolean => {
  let isValid = true
  errors.username = ''
  errors.password = ''
  errors.captcha = ''

  if (!loginForm.username.trim()) {
    errors.username = '请输入用户名'
    isValid = false
  }

  if (!loginForm.password) {
    errors.password = '请输入密码'
    isValid = false
  }

  if (!loginForm.captcha.trim()) {
    errors.captcha = '请输入验证码'
    isValid = false
  } else if (loginForm.captcha.toLowerCase() !== correctCaptcha.value) {
    errors.captcha = '验证码错误'
    isValid = false
    refreshCaptcha()
  }

  return isValid
}

const handleLogin = async () => {
  if (!validateForm()) return

  try {
    loading.value = true
    loadingText.value = '登录中...'

    const result = await UserApi.login(loginForm.username, loginForm.password)

    if (result.code === 200 && result.data) {
      const { access_token, refresh_token, user } = result.data

      userStore.setLogin(true)
      userStore.setUserInfo({
        id: user.id,
        username: user.username,
        email: user.email,
        phone: user.phone,
        nickname: user.nickname,
        role_id: user.role_id,
        is_super_admin: user.is_super_admin
      })
      userStore.setToken(access_token, refresh_token)

      MessagePlugin.success('登录成功')
      router.push('/index')
    } else {
      MessagePlugin.error(result.message || '登录失败')
      refreshCaptcha()
    }
  } catch (error: any) {
    MessagePlugin.error(error?.message || '网络错误，请重试')
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

const showRegisterTip = () => {
  MessagePlugin.info('注册功能开发中')
}

generateCaptcha()
</script>

<style scoped>
.login-page {
  display: flex;
  align-items: stretch;
  min-height: 100vh;
  background: #0F1826;
}

.login-left {
  flex: 0 0 55%;
  background: #0F1826;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 60px;
  position: relative;
  overflow: hidden;
}

.login-left::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background:
    radial-gradient(ellipse 600px 600px at 20% 30%, rgba(14,165,233,0.2) 0%, transparent 70%),
    radial-gradient(ellipse 500px 500px at 80% 60%, rgba(245,158,11,0.12) 0%, transparent 70%),
    radial-gradient(ellipse 400px 400px at 50% 80%, rgba(16,185,129,0.1) 0%, transparent 70%);
  animation: gradientMesh 8s ease-in-out infinite alternate;
}

@keyframes gradientMesh {
  0% {
    background:
      radial-gradient(ellipse 600px 600px at 20% 30%, rgba(14,165,233,0.2) 0%, transparent 70%),
      radial-gradient(ellipse 500px 500px at 80% 60%, rgba(245,158,11,0.12) 0%, transparent 70%),
      radial-gradient(ellipse 400px 400px at 50% 80%, rgba(16,185,129,0.1) 0%, transparent 70%);
  }
  33% {
    background:
      radial-gradient(ellipse 600px 600px at 40% 50%, rgba(14,165,233,0.18) 0%, transparent 70%),
      radial-gradient(ellipse 500px 500px at 60% 30%, rgba(245,158,11,0.15) 0%, transparent 70%),
      radial-gradient(ellipse 400px 400px at 30% 70%, rgba(16,185,129,0.12) 0%, transparent 70%);
  }
  66% {
    background:
      radial-gradient(ellipse 600px 600px at 60% 40%, rgba(14,165,233,0.22) 0%, transparent 70%),
      radial-gradient(ellipse 500px 500px at 30% 70%, rgba(245,158,11,0.1) 0%, transparent 70%),
      radial-gradient(ellipse 400px 400px at 70% 50%, rgba(16,185,129,0.14) 0%, transparent 70%);
  }
  100% {
    background:
      radial-gradient(ellipse 600px 600px at 50% 60%, rgba(14,165,233,0.16) 0%, transparent 70%),
      radial-gradient(ellipse 500px 500px at 70% 40%, rgba(245,158,11,0.13) 0%, transparent 70%),
      radial-gradient(ellipse 400px 400px at 20% 50%, rgba(16,185,129,0.11) 0%, transparent 70%);
  }
}

.login-left-content {
  position: relative;
  z-index: 2;
  text-align: center;
}

.login-display-title {
  font-size: 48px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 16px;
  letter-spacing: 2px;
}

.login-display-subtitle {
  font-size: 18px;
  font-weight: 300;
  color: rgba(255,255,255,0.6);
  margin-bottom: 60px;
  letter-spacing: 4px;
}

.login-features {
  display: flex;
  gap: 40px;
  justify-content: center;
}

.login-feature {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.login-feature-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  background: rgba(255,255,255,0.06);
  border: 1px solid rgba(255,255,255,0.1);
  backdrop-filter: blur(8px);
}

.login-feature-label {
  font-size: 13px;
  color: rgba(255,255,255,0.5);
  font-weight: 400;
}

.login-right {
  flex: 0 0 45%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px;
  background: #162032;
}

.login-form-wrap {
  width: 100%;
  max-width: 380px;
}

.login-form-heading {
  font-size: 24px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 8px;
}

.login-form-subtitle {
  font-size: 14px;
  color: rgba(255,255,255,0.6);
  margin-bottom: 36px;
}

.form-group {
  margin-bottom: 22px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: rgba(255,255,255,0.6);
  margin-bottom: 8px;
}

.form-input-wrap {
  position: relative;
}

.form-input-wrap .input-icon {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: rgba(255,255,255,0.3);
  font-size: 16px;
  pointer-events: none;
}

.form-input {
  width: 100%;
  padding: 12px 14px 12px 42px;
  border: 1.5px solid rgba(255,255,255,0.1);
  border-radius: 8px;
  font-size: 14px;
  background: #1E2A3A;
  color: #fff;
  transition: all 0.25s ease;
}

.form-input:focus {
  outline: none;
  border-color: #0EA5E9;
  box-shadow: 0 0 0 3px rgba(14,165,233,0.15);
}

.form-input::placeholder {
  color: rgba(255,255,255,0.3);
}

.form-input.error {
  border-color: #F87171;
  box-shadow: 0 0 0 3px rgba(239,68,68,0.1);
}

.toggle-pw {
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  font-size: 16px;
  color: rgba(255,255,255,0.3);
  cursor: pointer;
  padding: 4px;
  transition: color 0.2s;
  border: none;
}

.toggle-pw:hover {
  color: rgba(255,255,255,0.6);
}

.captcha-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.captcha-row .form-input-wrap {
  flex: 1;
}

.captcha-img {
  width: 120px;
  height: 44px;
  border-radius: 8px;
  background: #1E2A3A;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255,255,255,0.8);
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 6px;
  cursor: pointer;
  user-select: none;
  font-family: 'Courier New', monospace;
  transition: all 0.2s;
  border: 1.5px dashed rgba(255,255,255,0.2);
}

.captcha-img:hover {
  border-color: #0EA5E9;
  color: #0EA5E9;
  background: rgba(14,165,233,0.1);
}

.form-error-msg {
  font-size: 12px;
  color: #F87171;
  margin-top: 6px;
  display: none;
}

.form-error-msg.show {
  display: block;
}

.form-check {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: rgba(255,255,255,0.6);
}

.form-check input[type="checkbox"] {
  width: 16px;
  height: 16px;
  accent-color: #0EA5E9;
}

.form-row-between {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
}

.form-link {
  font-size: 13px;
  color: #0EA5E9;
  transition: color 0.2s;
  text-decoration: none;
}

.form-link:hover {
  color: #38BDF8;
}

.btn {
  width: 100%;
  height: 48px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.25s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn-primary {
  background: #0EA5E9;
  border: none;
  color: #fff;
}

.btn-primary:hover {
  background: #38BDF8;
}

.btn-primary:disabled {
  background: rgba(14,165,233,0.5);
  cursor: not-allowed;
}

.btn-block {
  width: 100%;
}

.login-footer {
  text-align: center;
  margin-top: 32px;
  font-size: 13px;
  color: rgba(255,255,255,0.4);
}

.login-footer a {
  color: #0EA5E9;
  font-weight: 500;
  text-decoration: none;
}

.login-footer a:hover {
  text-decoration: underline;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
