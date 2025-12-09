<template>
  <header class="header">
    <div class="header-container">
      <!-- Logo区域 -->
      <div class="logo">
        <h1 class="logo-text">720云</h1>
      </div>
      
      <!-- 导航菜单 -->
      <nav class="nav-menu">
        <a href="#" class="nav-item">首页</a>
        <el-dropdown>
          <span class="nav-item dropdown">
            产品服务
            <el-icon class="icon-right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>VR全景创作</el-dropdown-item>
              <el-dropdown-item>数字孪生</el-dropdown-item>
              <el-dropdown-item>元宇宙平台</el-dropdown-item>
              <el-dropdown-item>VR直播</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-dropdown>
          <span class="nav-item dropdown">
            解决方案
            <el-icon class="icon-right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>教育行业</el-dropdown-item>
              <el-dropdown-item>房地产</el-dropdown-item>
              <el-dropdown-item>旅游景区</el-dropdown-item>
              <el-dropdown-item>企业展示</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-dropdown>
          <span class="nav-item dropdown">
            内容社区
            <el-icon class="icon-right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>作品展示</el-dropdown-item>
              <el-dropdown-item>教程中心</el-dropdown-item>
              <el-dropdown-item>社区论坛</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-dropdown>
          <span class="nav-item dropdown">
            定制服务
            <el-icon class="icon-right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>VR内容定制</el-dropdown-item>
              <el-dropdown-item>数字孪生定制</el-dropdown-item>
              <el-dropdown-item>平台定制开发</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-dropdown>
          <span class="nav-item dropdown">
            商城
            <el-icon class="icon-right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>VR设备</el-dropdown-item>
              <el-dropdown-item>素材资源</el-dropdown-item>
              <el-dropdown-item>服务套餐</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </nav>
      
      <!-- 搜索框 -->
      <div class="search-box">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索"
          size="small"
          prefix-icon="Search"
          class="search-input"
        >
          <template #append>
            <el-button @click="handleSearch" size="small" type="primary" icon="Search"></el-button>
          </template>
        </el-input>
      </div>
      
      <!-- 用户操作 -->
      <div class="user-actions">
        <el-button type="primary" size="small" @click="showLoginDialog = true" class="login-button mr-2">登录</el-button>
        <el-button type="default" size="small" @click="showRegisterDialog = true" class="register-button">注册</el-button>
      </div>
      
      <!-- 登录弹窗 -->
      <el-dialog v-model="showLoginDialog" title="登录" width="400px" center>
        <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" class="login-form">
          <el-form-item prop="email">
            <el-input 
              v-model="loginForm.email" 
              placeholder="请输入邮箱" 
              type="email"
              prefix-icon="Message"
            ></el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input 
              v-model="loginForm.password" 
              placeholder="请输入密码" 
              type="password" 
              show-password
              prefix-icon="Lock"
            ></el-input>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showLoginDialog = false">取消</el-button>
            <el-button type="primary" @click="handleLogin" :loading="loginLoading">
              登录
            </el-button>
          </span>
        </template>
      </el-dialog>
      
      <!-- 注册弹窗 -->
      <el-dialog v-model="showRegisterDialog" title="注册" width="400px" center>
        <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef" class="register-form">
          <el-form-item prop="name">
            <el-input 
              v-model="registerForm.name" 
              placeholder="请输入用户名"
              prefix-icon="User"
            ></el-input>
          </el-form-item>
          <el-form-item prop="email">
            <el-input 
              v-model="registerForm.email" 
              placeholder="请输入邮箱" 
              type="email"
              prefix-icon="Message"
            ></el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input 
              v-model="registerForm.password" 
              placeholder="请输入密码" 
              type="password" 
              show-password
              prefix-icon="Lock"
            ></el-input>
          </el-form-item>
          <el-form-item prop="confirmPassword">
            <el-input 
              v-model="registerForm.confirmPassword" 
              placeholder="请确认密码" 
              type="password" 
              show-password
              prefix-icon="Lock"
            ></el-input>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showRegisterDialog = false">取消</el-button>
            <el-button type="primary" @click="handleRegister" :loading="registerLoading">
              注册
            </el-button>
          </span>
        </template>
      </el-dialog>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowDown, Message, Lock, User } from '@element-plus/icons-vue'
import UserApi from '@/apis/userApi'
import { useUserStore } from '@/stores/userStore'

// 搜索功能
const searchKeyword = ref('')

// 搜索处理
const handleSearch = () => {
    if (searchKeyword.value.trim()) {
        ElMessage.info(`搜索内容：${searchKeyword.value}`)
    } else {
        ElMessage.warning('请输入搜索内容')
    }
}

// 用户状态
const userStore = useUserStore()

// 登录弹窗控制
const showLoginDialog = ref(false)
const showRegisterDialog = ref(false)

// 登录表单
const loginFormRef = ref()
const loginLoading = ref(false)
const loginForm = reactive({
    email: '',
    password: ''
})

// 登录表单验证规则
const loginRules = {
    email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
    ],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
    ]
}

// 注册表单
const registerFormRef = ref()
const registerLoading = ref(false)
const registerForm = reactive({
    name: '',
    email: '',
    password: '',
    confirmPassword: ''
})

// 注册表单验证规则
const registerRules = {
    name: [
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

// 登录处理
const handleLogin = async () => {
    if (!loginFormRef.value) return
    
    // 表单验证
    await loginFormRef.value.validate(async (valid: boolean) => {
        if (valid) {
            loginLoading.value = true
            try {
                const result = await UserApi.login(loginForm.email, loginForm.password)
                if (result.code === 200) {
                    // 登录成功
                    userStore.setLogin(true)
                    userStore.setUserInfo({ email: loginForm.email })
                    showLoginDialog.value = false
                    ElMessage.success('登录成功')
                    
                    // 重置表单
                    loginForm.email = ''
                    loginForm.password = ''
                } else {
                    // 登录失败
                    ElMessage.error(result.msg || '登录失败')
                }
            } catch (error: any) {
                ElMessage.error(error.message || '登录失败，请稍后重试')
            } finally {
                loginLoading.value = false
            }
        }
    })
}

// 注册处理
const handleRegister = async () => {
    if (!registerFormRef.value) return
    
    // 表单验证
    await registerFormRef.value.validate(async (valid: boolean) => {
        if (valid) {
            registerLoading.value = true
            try {
                const result = await UserApi.register({
                    name: registerForm.name,
                    email: registerForm.email,
                    password: registerForm.password
                })
                if (result.code === 200) {
                    // 注册成功
                    ElMessage.success('注册成功，请登录')
                    showRegisterDialog.value = false
                    showLoginDialog.value = true
                    
                    // 重置表单
                    registerForm.name = ''
                    registerForm.email = ''
                    registerForm.password = ''
                    registerForm.confirmPassword = ''
                } else {
                    // 注册失败
                    ElMessage.error(result.msg || '注册失败')
                }
            } catch (error: any) {
                ElMessage.error(error.message || '注册失败，请稍后重试')
            } finally {
                registerLoading.value = false
            }
        }
    })
}
</script>

<style scoped lang="scss">
// 主题颜色变量
$primary-color: #3b82f6;
$primary-hover: #2563eb;
$red-color: #ef4444;
$orange-color: #f97316;
$orange-hover: #ea580c;
$gray-color: #6b7280;
$gray-hover: #374151;
$white-color: #ffffff;
$shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);

// 过渡效果
$transition: all 0.3s ease;

// Header 主容器
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 50;
  background-color: $white-color;
  box-shadow: $shadow-sm;
  transition: $transition;
  height: 64px;
  display: flex;
  align-items: center;
}

// Header 内容容器
.header-container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  gap: 12px;
}

// Logo 区域
.logo {
  display: flex;
  align-items: center;
  width: auto;
  min-width: 80px;
  padding-right: 4px;
  
  .logo-text {
    font-size: 20px;
    font-weight: bold;
    color: $primary-color;
    white-space: nowrap;
    margin: 0;
    line-height: 1;
  }
}

// 导航菜单
.nav-menu {
  display: flex;
  align-items: center;
  flex: 1;
  justify-content: flex-start;
  gap: 4px;
  
  .nav-item {
    position: relative;
    font-size: 14px;
    color: $gray-color;
    cursor: pointer;
    padding: 6px 14px;
    transition: $transition;
    white-space: nowrap;
    border-radius: 4px;
    
    &:hover {
      color: $primary-color;
      background-color: rgba(59, 130, 246, 0.05);
      
      &::after {
        width: 100%;
      }
    }
    
    &::after {
      content: '';
      position: absolute;
      bottom: -4px;
      left: 0;
      width: 0;
      height: 2px;
      background-color: $primary-color;
      transition: $transition;
    }
    
    &.dropdown {
      display: flex;
      align-items: center;
    }
  }
}

// 图标样式
.icon-right {
  margin-left: 4px;
  font-size: 12px;
}

.icon-left {
  margin-right: 4px;
  font-size: 14px;
}

// 用户操作按钮样式
.user-actions {
  display: flex;
  align-items: center;
  min-width: fit-content;
  
  .login-button,
  .register-button {
    font-size: 12px;
    padding: 6px 14px;
    transition: $transition;
    
    &:hover {
      transform: translateY(-1px);
    }
  }
  
  .mr-2 {
    margin-right: 8px;
  }
}

// 搜索框
.search-box {
  display: flex;
  align-items: center;
  flex: 0 1 200px;
  min-width: 150px;
  
  .search-input {
    width: 100%;
    min-width: 0;
  }
}

// 用户操作
.user-actions {
  display: flex;
  align-items: center;
  min-width: fit-content;
  
  .login-link {
    font-size: 13px;
    color: $primary-color;
    transition: $transition;
    padding: 6px 10px;
    border-radius: 4px;
    
    &:hover {
      color: $primary-hover;
      background-color: rgba(59, 130, 246, 0.05);
    }
  }
}



// 响应式设计
@media (max-width: 1200px) {
  .header-container {
    max-width: 1000px;
  }
  
  .special-tags {
    display: none;
  }
  
  .member-tags {
    display: none;
  }
  
  .search-box {
    flex: 0 1 180px;
  }
}

@media (max-width: 992px) {
  .nav-menu {
    display: none;
  }
  
  .search-box {
    flex: 1;
    min-width: 180px;
  }
}

@media (max-width: 576px) {
  .header-container {
    padding: 0 14px;
    gap: 8px;
  }
  
  .logo {
    min-width: 70px;
    
    .logo-text {
      font-size: 18px;
    }
  }
  
  .search-box {
    display: none;
  }
  
  .create-button {
    padding: 6px 12px;
    font-size: 13px;
  }
}
</style>