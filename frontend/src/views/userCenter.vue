<template>
  <div class="user-center-page">
    <Header />

    <div class="content-wrapper">
      <div class="user-center-container">
        <t-card class="info-card">
          <template #header>
            <div class="card-header">
              <t-icon name="user" class="header-icon" />
              <span class="header-title">个人信息</span>
            </div>
          </template>

          <div class="user-info-section">
            <div class="info-item">
              <label class="info-label">用户名</label>
              <div class="info-value">{{ userStore.userInfo.username || '-' }}</div>
            </div>

            <div class="info-item">
              <label class="info-label">邮箱</label>
              <div class="info-value">{{ userStore.userInfo.email || '-' }}</div>
            </div>

            <div class="info-item">
              <label class="info-label">昵称</label>
              <div class="info-value">{{ userStore.userInfo.nickname || '-' }}</div>
            </div>

            <div class="info-item">
              <label class="info-label">手机号</label>
              <div class="info-value">{{ userStore.userInfo.phone || '-' }}</div>
            </div>

            <div class="info-item">
              <label class="info-label">角色</label>
              <div class="info-value">
                <t-tag :theme="roleTheme">{{ roleLabel }}</t-tag>
              </div>
            </div>
          </div>
        </t-card>

        <t-card class="password-card">
          <template #header>
            <div class="card-header">
              <t-icon name="lock-on" class="header-icon" />
              <span class="header-title">修改密码</span>
            </div>
          </template>

          <t-form
            ref="passwordFormRef"
            :data="passwordForm"
            :rules="passwordRules"
            label-width="100px"
            @submit="handlePasswordSubmit"
          >
            <t-form-item label="当前密码" name="oldPassword">
              <t-input
                v-model="passwordForm.oldPassword"
                type="password"
                placeholder="请输入当前密码"
              />
            </t-form-item>

            <t-form-item label="新密码" name="newPassword">
              <t-input
                v-model="passwordForm.newPassword"
                type="password"
                placeholder="请输入新密码"
              />
            </t-form-item>

            <t-form-item label="确认密码" name="confirmPassword">
              <t-input
                v-model="passwordForm.confirmPassword"
                type="password"
                placeholder="请再次输入新密码"
              />
            </t-form-item>

            <t-form-item>
              <t-space>
                <t-button theme="primary" type="submit" :loading="passwordLoading">
                  修改密码
                </t-button>
                <t-button theme="default" @click="resetPasswordForm">
                  重置
                </t-button>
              </t-space>
            </t-form-item>
          </t-form>
        </t-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import Header from '@/components/Header.vue'
import { useUserStore } from '@/stores/userStore'

const userStore = useUserStore()

const roleLabel = computed(() => {
  if (userStore.userInfo.is_super_admin) {
    return '超级管理员'
  }
  const roleId = userStore.userInfo.role_id
  switch (roleId) {
    case 2:
      return '教师'
    case 3:
      return '学生'
    default:
      return '未知角色'
  }
})

const roleTheme = computed(() => {
  if (userStore.userInfo.is_super_admin) {
    return 'danger'
  }
  const roleId = userStore.userInfo.role_id
  switch (roleId) {
    case 1:
      return 'danger'
    case 2:
      return 'warning'
    case 3:
      return 'success'
    default:
      return 'default'
  }
})

const passwordFormRef = ref()
const passwordLoading = ref(false)

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (val: string) => {
  if (val !== passwordForm.newPassword) {
    return { result: false, message: '两次输入的密码不一致' }
  }
  return { result: true }
}

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handlePasswordSubmit = async () => {
  if (!passwordFormRef.value) return

  const valid = await (passwordFormRef.value as any).validate()
  if (valid) {
    passwordLoading.value = true
    try {
      MessagePlugin.success('密码修改功能开发中...')
    } finally {
      passwordLoading.value = false
    }
  }
}

const resetPasswordForm = () => {
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}
</script>

<style scoped>
.user-center-page {
  width: 100%;
  min-height: 100vh;
  background-color: #f5f7fa;
}

.content-wrapper {
  padding: 120px 24px 80px;
  display: flex;
  justify-content: center;
}

.user-center-container {
  width: 100%;
  max-width: 600px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.info-card,
.password-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-icon {
  font-size: 20px;
  color: #0052d9;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}

.user-info-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f2f3f5;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  width: 100px;
  font-size: 14px;
  color: #86909c;
  flex-shrink: 0;
}

.info-value {
  font-size: 14px;
  color: #1d2129;
  font-weight: 500;
}
</style>
