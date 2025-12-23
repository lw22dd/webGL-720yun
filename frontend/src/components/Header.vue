<template>
  <header class="fixed top-0 left-0 right-0 z-50 bg-white shadow-sm transition-all duration-300 h-16 flex items-center">
    <div class="w-full max-w-7xl mx-auto px-5 flex items-center justify-between h-full gap-3">
      <!-- Logo区域 -->
      <div class="flex items-center w-auto min-w-20 pr-1">
        <h1 class="text-xl font-bold text-blue-600 whitespace-nowrap m-0 leading-1">720云</h1>
      </div>
      
      <!-- 导航菜单 -->
      <nav class="flex items-center flex-1 justify-start gap-1">
        <a href="#" class="relative text-sm text-gray-600 cursor-pointer px-3.5 py-1.5 transition-all whitespace-nowrap rounded hover:text-blue-600 hover:bg-blue-50">首页</a>
        <el-dropdown>
          <span class="relative text-sm text-gray-600 cursor-pointer px-3.5 py-1.5 transition-all whitespace-nowrap rounded hover:text-blue-600 hover:bg-blue-50 flex items-center">
            产品服务
            <el-icon class="ml-1 text-xs"><arrow-down /></el-icon>
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
          <span class="relative text-sm text-gray-600 cursor-pointer px-3.5 py-1.5 transition-all whitespace-nowrap rounded hover:text-blue-600 hover:bg-blue-50 flex items-center">
            解决方案
            <el-icon class="ml-1 text-xs"><arrow-down /></el-icon>
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
          <span class="relative text-sm text-gray-600 cursor-pointer px-3.5 py-1.5 transition-all whitespace-nowrap rounded hover:text-blue-600 hover:bg-blue-50 flex items-center">
            内容社区
            <el-icon class="ml-1 text-xs"><arrow-down /></el-icon>
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
          <span class="relative text-sm text-gray-600 cursor-pointer px-3.5 py-1.5 transition-all whitespace-nowrap rounded hover:text-blue-600 hover:bg-blue-50 flex items-center">
            定制服务
            <el-icon class="ml-1 text-xs"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>VR内容定制</el-dropdown-item>
              <el-dropdown-item>数字孪生定制</el-dropdown-item>
              <el-dropdown-item>平台定制开发</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </nav>
      
      <!-- 搜索框 -->
      <div class="flex items-center flex-0 1 w-50 min-w-36">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索"
          size="small"
          prefix-icon="Search"
          class="w-full min-w-0"
        >
          <template #append>
            <el-button @click="handleSearch" size="small" type="primary" icon="Search"></el-button>
          </template>
        </el-input>
      </div>
      
      <!-- 用户操作 -->
      <div class="flex items-center min-w-fit">
        <el-button type="primary" size="small" @click="showLoginDialog = true" class="text-xs px-3.5 py-1.5 transition-all hover:-translate-y-0.5 mr-2">登录</el-button>
        <el-button type="default" size="small" @click="showRegisterDialog = true" class="text-xs px-3.5 py-1.5 transition-all hover:-translate-y-0.5">注册</el-button>
      </div>
      
      <!-- 登录弹窗 -->
      <el-dialog v-model="showLoginDialog" title="登录" width="400px" center>
        <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef">
          <el-form-item prop="username">
            <el-input 
              v-model="loginForm.username" 
              placeholder="请输入用户名或邮箱" 
              prefix-icon="User"
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
          <span>
            <el-button @click="showLoginDialog = false">取消</el-button>
            <el-button type="primary" @click="handleLogin" :loading="loginLoading">
              登录
            </el-button>
          </span>
        </template>
      </el-dialog>
      
      <!-- 注册弹窗 -->
      <el-dialog v-model="showRegisterDialog" title="注册" width="400px" center>
        <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef">
          <el-form-item prop="username">
            <el-input 
              v-model="registerForm.username" 
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
          <span>
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
    username: '',
    password: ''
})

// 登录表单验证规则
const loginRules = {
    username: [
        { required: true, message: '请输入用户名或邮箱', trigger: 'blur' }
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
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
})

// 注册表单验证规则
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

// 登录处理
const handleLogin = async () => {
    if (!loginFormRef.value) return
    
    // 表单验证
    await loginFormRef.value.validate(async (valid: boolean) => {
        if (valid) {
            loginLoading.value = true
            try {
                const result = await UserApi.login(loginForm.username, loginForm.password)
                if (result.code === 200) {
                    // 登录成功
                    userStore.setLogin(true)
                    userStore.setUserInfo({ email: loginForm.username })
                    showLoginDialog.value = false
                    ElMessage.success('登录成功')
                    
                    // 重置表单
                    loginForm.username = ''
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
                    username: registerForm.username,
                    email: registerForm.email,
                    password: registerForm.password,
                    role_id: 2 // 默认为普通用户角色，根据实际需求调整
                })
                if (result.code === 200) {
                    // 注册成功
                    ElMessage.success('注册成功，请登录')
                    showRegisterDialog.value = false
                    showLoginDialog.value = true
                    
                    // 重置表单
                    registerForm.username = ''
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

<style scoped>
/* Tailwind CSS 已应用，不再需要自定义 SCSS 样式 */
</style>