<template>
  <div class="admin-layout">
    <Header />

    <div class="admin-content">
      <div class="admin-container">
        <t-card class="admin-card">
          <t-tabs v-model="activeTab" @change="handleTabChange">
            <t-tab-panel value="users" label="用户管理">
              <UserManagement />
            </t-tab-panel>
            <t-tab-panel value="batch-register" label="批量注册">
              <ExcelUpload />
            </t-tab-panel>
          </t-tabs>
        </t-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Message } from 'tdesign-vue-next'
import Header from '@/components/Header.vue'
import UserManagement from '@/views/userManagement.vue'
import ExcelUpload from '@/components/ExcelUpload.vue'
import { useUserStore } from '@/stores/userStore'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeTab = ref('users')

const handleTabChange = (value: string) => {
  router.push(`/admin/${value}`)
}

onMounted(() => {
  const roleId = userStore.userInfo.role_id
  if (roleId !== 1 && roleId !== 2) {
    Message.warning('您没有权限访问此页面')
    router.push('/')
    return
  }

  const path = route.path
  if (path.includes('batch-register')) {
    activeTab.value = 'batch-register'
  } else {
    activeTab.value = 'users'
  }
})
</script>

<style scoped>
.admin-layout {
  width: 100%;
  min-height: 100vh;
  background-color: #f5f7fa;
}

.admin-content {
  padding: 120px 24px 80px;
  min-height: calc(100vh - 64px);
}

.admin-container {
  max-width: 1400px;
  margin: 0 auto;
}

.admin-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

:deep(.t-card__body) {
  padding: 0;
}

:deep(.t-tabs__header) {
  padding: 0 24px;
  background-color: #fafafa;
  border-bottom: 1px solid #e7e7e7;
}

:deep(.t-tab-panel) {
  padding: 24px;
}
</style>
