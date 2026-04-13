<template>
  <div class="user-center-page">
    <Header />
    
    <div class="page-container">
      <div class="profile-header">
        <div class="profile-avatar">
          <t-upload
            v-model="avatarFiles"
            theme="image"
            accept="image/*"
            :auto-upload="false"
            @change="handleAvatarChange"
          >
            <template #file-list-display>
              <div class="avatar-wrapper">
                <img v-if="userInfo.avatar" :src="userInfo.avatar" alt="头像" />
                <div v-else class="avatar-placeholder">
                  {{ userInfo.username?.charAt(0)?.toUpperCase() || 'U' }}
                </div>
                <div class="avatar-overlay">
                  <t-icon name="camera" />
                </div>
              </div>
            </template>
          </t-upload>
        </div>
        <div class="profile-info">
          <h1 class="profile-name">{{ userInfo.nickname || userInfo.username }}</h1>
          <p class="profile-role">
            <t-tag :theme="getRoleTheme(userInfo.role)">{{ getRoleName(userInfo.role) }}</t-tag>
          </p>
          <p class="profile-email" v-if="userInfo.email">{{ userInfo.email }}</p>
        </div>
      </div>

      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-value">{{ stats.spacesVisited }}</div>
          <div class="stat-label">浏览景区</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.favorites }}</div>
          <div class="stat-label">收藏数量</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.studyTime }}</div>
          <div class="stat-label">学习时长</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.loginDays }}</div>
          <div class="stat-label">累计登录</div>
        </div>
      </div>

      <div class="content-grid">
        <t-card title="个人信息" :bordered="false" class="info-card">
          <t-form :data="userInfo" @submit="handleUpdateProfile">
            <t-form-item label="用户名">
              <t-input v-model="userInfo.username" disabled />
            </t-form-item>
            <t-form-item label="昵称">
              <t-input v-model="userInfo.nickname" placeholder="请输入昵称" />
            </t-form-item>
            <t-form-item label="邮箱">
              <t-input v-model="userInfo.email" placeholder="请输入邮箱" />
            </t-form-item>
            <t-form-item label="手机号">
              <t-input v-model="userInfo.phone" placeholder="请输入手机号" />
            </t-form-item>
            <t-form-item label="班级" v-if="userInfo.role === 'student'">
              <t-input v-model="userInfo.class_name" placeholder="请输入班级" />
            </t-form-item>
            <t-form-item label="学号" v-if="userInfo.role === 'student'">
              <t-input v-model="userInfo.student_id" placeholder="请输入学号" />
            </t-form-item>
            <t-form-item>
              <t-button theme="primary" type="submit">保存修改</t-button>
            </t-form-item>
          </t-form>
        </t-card>

        <t-card title="快捷入口" :bordered="false" class="quick-card">
          <div class="quick-links">
            <div class="quick-link" @click="router.push('/favorites')">
              <t-icon name="star" size="24" />
              <span>我的收藏</span>
            </div>
            <div class="quick-link" @click="router.push('/history')">
              <t-icon name="history" size="24" />
              <span>浏览历史</span>
            </div>
            <div class="quick-link" @click="showPasswordDialog = true">
              <t-icon name="lock-on" size="24" />
              <span>修改密码</span>
            </div>
            <div class="quick-link" @click="handleLogout">
              <t-icon name="logout" size="24" />
              <span>退出登录</span>
            </div>
          </div>
        </t-card>
      </div>
    </div>

    <t-dialog
      v-model:visible="showPasswordDialog"
      header="修改密码"
      width="400px"
      @confirm="handleChangePassword"
    >
      <t-form :data="passwordForm" ref="passwordFormRef" :rules="passwordRules">
        <t-form-item label="当前密码" name="oldPassword">
          <t-input v-model="passwordForm.oldPassword" type="password" placeholder="请输入当前密码" />
        </t-form-item>
        <t-form-item label="新密码" name="newPassword">
          <t-input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码" />
        </t-form-item>
        <t-form-item label="确认密码" name="confirmPassword">
          <t-input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import Header from '@/components/Header.vue'
import { useUserStore } from '@/stores/userStore'
import UserApi from '@/apis/userApi'

const router = useRouter()
const userStore = useUserStore()

const avatarFiles = ref([])
const showPasswordDialog = ref(false)
const passwordFormRef = ref()

const userInfo = reactive({
  id: 0,
  username: '',
  nickname: '',
  email: '',
  phone: '',
  avatar: '',
  role: '',
  class_name: '',
  student_id: ''
})

const stats = reactive({
  spacesVisited: 0,
  favorites: 0,
  studyTime: '0h',
  loginDays: 0
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (val: string) => val === passwordForm.newPassword,
      message: '两次密码输入不一致',
      trigger: 'blur'
    }
  ]
}

const getRoleName = (role: string) => {
  const roleMap: Record<string, string> = {
    admin: '管理员',
    teacher: '教师',
    student: '学生'
  }
  return roleMap[role] || '用户'
}

const getRoleTheme = (role: string) => {
  const themeMap: Record<string, string> = {
    admin: 'danger',
    teacher: 'warning',
    student: 'success'
  }
  return themeMap[role] || 'default'
}

const loadUserInfo = async () => {
  const storeUserInfo = userStore.userInfo
  if (storeUserInfo) {
    Object.assign(userInfo, storeUserInfo)
  }
  
  try {
    const result = await UserApi.getProfile()
    if (result.code === 200 && result.data) {
      Object.assign(userInfo, result.data)
    }
  } catch (e) {
    console.error('Failed to load user info:', e)
  }
}

const loadStats = () => {
  try {
    const favorites = JSON.parse(localStorage.getItem('favorites') || '[]')
    const history = JSON.parse(localStorage.getItem('browseHistory') || '[]')
    
    stats.favorites = favorites.length
    stats.spacesVisited = new Set(history.map((h: any) => h.id)).size
    
    const loginDays = localStorage.getItem('loginDays')
    stats.loginDays = loginDays ? parseInt(loginDays) : 1
    
    const studyMinutes = localStorage.getItem('studyMinutes')
    if (studyMinutes) {
      const hours = Math.floor(parseInt(studyMinutes) / 60)
      stats.studyTime = `${hours}h`
    }
  } catch (e) {
    console.error('Failed to load stats:', e)
  }
}

const handleAvatarChange = () => {
  if (avatarFiles.value.length > 0) {
    const file = (avatarFiles.value[0] as any)
    if (file.raw) {
      const reader = new FileReader()
      reader.onload = (e) => {
        userInfo.avatar = e.target?.result as string
        MessagePlugin.success('头像已更新')
      }
      reader.readAsDataURL(file.raw)
    }
  }
}

const handleUpdateProfile = async () => {
  try {
    const result = await UserApi.updateProfile({
      nickname: userInfo.nickname,
      email: userInfo.email,
      phone: userInfo.phone
    })
    if (result.code === 200) {
      MessagePlugin.success('个人信息已更新')
      userStore.setUserInfo({
        id: String(userInfo.id),
        username: userInfo.username,
        nickname: userInfo.nickname,
        email: userInfo.email,
        phone: userInfo.phone,
        avatar: userInfo.avatar
      })
    }
  } catch (e: any) {
    MessagePlugin.error(e.message || '更新失败')
  }
}

const handleChangePassword = async () => {
  const valid = await passwordFormRef.value?.validate()
  if (valid !== true) return
  
  try {
    const result = await UserApi.changePassword({
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword
    })
    if (result.code === 200) {
      MessagePlugin.success('密码修改成功')
      showPasswordDialog.value = false
      Object.assign(passwordForm, {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      })
    }
  } catch (e: any) {
    MessagePlugin.error(e.message || '修改失败')
  }
}

const handleLogout = () => {
  userStore.logout()
  router.push('/')
  MessagePlugin.success('已退出登录')
}

onMounted(() => {
  loadUserInfo()
  loadStats()
})
</script>

<style scoped>
.user-center-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0F1826 0%, #162032 50%, #1A2744 100%);
}

.page-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 24px;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 32px;
  padding: 32px;
  background: rgba(26, 35, 50, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
}

.profile-avatar {
  position: relative;
}

.avatar-wrapper {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
  position: relative;
}

.avatar-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0EA5E9, #4080ff);
  color: #fff;
  font-size: 40px;
  font-weight: 600;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  opacity: 0;
  transition: opacity 0.2s;
}

.avatar-wrapper:hover .avatar-overlay {
  opacity: 1;
}

.profile-info {
  flex: 1;
}

.profile-name {
  color: #fff;
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 8px;
}

.profile-role {
  margin: 0 0 8px;
}

.profile-email {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
  margin: 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: rgba(26, 35, 50, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 20px;
  text-align: center;
}

.stat-value {
  color: #fff;
  font-size: 32px;
  font-weight: 600;
  margin-bottom: 4px;
}

.stat-label {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
}

.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.info-card,
.quick-card {
  background: rgba(26, 35, 50, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

:deep(.t-card__header) {
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

:deep(.t-card__title) {
  color: #fff;
}

:deep(.t-form-item__label) {
  color: rgba(255, 255, 255, 0.8);
}

:deep(.t-input) {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
}

:deep(.t-input__inner) {
  color: #fff;
}

:deep(.t-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.4);
}

.quick-links {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.quick-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  color: rgba(255, 255, 255, 0.8);
}

.quick-link:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(14, 165, 233, 0.3);
  color: #fff;
}

.quick-link span {
  font-size: 14px;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .content-grid {
    grid-template-columns: 1fr;
  }
}
</style>
