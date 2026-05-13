<template>
  <div class="admin-layout">
    <div class="admin-sidebar" :class="{ 'sidebar-hidden': sidebarHidden }">
      <AdminSideMenu />
    </div>

    <div class="admin-main" :class="{ 'main-full': sidebarHidden }">
      <div class="admin-header">
        <div class="header-left">
          <t-button
            variant="text"
            theme="default"
            size="large"
            @click="toggleSidebar"
          >
            <template #icon>
              <t-icon :name="sidebarHidden ? 'menu-fold' : 'menu-unfold'" />
            </template>
          </t-button>

          <t-breadcrumb>
            <t-breadcrumb-item @click="() => router.push('/')">首页</t-breadcrumb-item>
            <t-breadcrumb-item @click="() => router.push('/admin')">管理中心</t-breadcrumb-item>
            <t-breadcrumb-item v-if="route.path !== '/admin'">{{ currentPageName }}</t-breadcrumb-item>
          </t-breadcrumb>
        </div>

        <div class="header-right">
          <t-dropdown>
            <div class="user-info">
              <t-avatar size="small">
                <template #icon>
                  <t-icon name="user" />
                </template>
              </t-avatar>
              <span class="user-name">
                {{ userStore.userInfo.email || userStore.userInfo.name || '管理员' }}
              </span>
            </div>
            <template #dropdown>
              <t-dropdown-menu>
                <t-dropdown-item @click="navigateToHome">返回首页</t-dropdown-item>
                <t-dropdown-item @click="handleLogout">退出登录</t-dropdown-item>
              </t-dropdown-menu>
            </template>
          </t-dropdown>
        </div>
      </div>

      <div class="admin-content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import AdminSideMenu from '@/components/admin/AdminSideMenu.vue'
import { useUserStore } from '@/stores/user.store'
import UserApi from '@/apis/user.api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const sidebarHidden = ref(false)

const currentPageName = computed(() => {
  const pathMap: Record<string, string> = {
    '/admin': '管理中心',
    '/admin/users': '用户管理',
    '/admin/batch-register': '批量注册',
    '/admin/spaces': '空间管理',
    '/admin/scenes': '场景管理',
    '/admin/settings': '系统设置'
  }
  return pathMap[route.path] || route.path.split('/').pop()
})

const toggleSidebar = () => {
  sidebarHidden.value = !sidebarHidden.value
}



const navigateToHome = () => {
  router.push('/')
}

const handleLogout = async () => {
  try {
    const result = await UserApi.logout()
    if (result.code === 200) {
      userStore.logout()
      MessagePlugin.success('退出成功')
      router.push('/login')
    } else {
      userStore.logout()
      MessagePlugin.error(result.msg || '退出失败')
      router.push('/login')
    }
  } catch (error: any) {
    userStore.logout()
    MessagePlugin.success('已退出登录')
    router.push('/login')
  }
}

onMounted(() => {
  // 检查用户是否有权限访问管理后台 (role_id === 1 表示管理员)
  const isAdmin = userStore.userInfo.role_id === 1 || userStore.userInfo.is_super_admin === true
  if (!isAdmin) {
    console.log(userStore.userInfo)
    MessagePlugin.warning('您没有权限访问此页面')
    router.push('/')
  }
})
</script>

<style scoped>
.admin-layout {
  display: flex;
  width: 100%;
  min-height: 100vh;
  background-color: #f2f3f5;
}

.admin-sidebar {
  width: 256px;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  background-color: #ffffff;
  overflow-y: auto;
  overflow-x: hidden;
  flex-shrink: 0;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 100;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
}

.admin-sidebar.sidebar-hidden {
  transform: translateX(-100%);
}

.admin-main {
  flex: 1;
  margin-left: 256px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.admin-main.main-full {
  margin-left: 0;
}

.admin-header {
  background-color: #ffffff;
  border-bottom: 1px solid #e5e6eb;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
  height: 60px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.header-right {
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
  padding: 6px 12px;
  border-radius: 8px;
  transition: all 0.2s cubic-bezier(0.34, 0.69, 0.1, 1);
}

.user-info:hover {
  background: #f2f3f5;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #1d2129;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 120px;
}

/* 响应式适配 */
@media screen and (max-width: 768px) {
  .admin-sidebar {
    width: 200px;
  }

  .admin-main {
    margin-left: 200px;
  }

  .admin-header {
    padding: 0 16px;
  }

  .header-left {
    gap: 8px;
  }

  .user-name {
    max-width: 80px;
  }

  :deep(.t-breadcrumb__item) {
    font-size: 13px;
  }
}

@media screen and (max-width: 576px) {
  .admin-sidebar {
    width: 0;
    transform: translateX(-100%);
  }

  .admin-main {
    margin-left: 0;
  }

  .admin-header {
    padding: 0 12px;
  }

  .user-name {
    display: none;
  }
}

.admin-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background-color: #f2f3f5;
}

:deep(.t-breadcrumb) {
  height: 40px;
  display: flex;
  align-items: center;
}

:deep(.t-button) {
  border-radius: 6px;
}

/* 滚动条美化 */
.admin-sidebar::-webkit-scrollbar {
  width: 6px;
}

.admin-sidebar::-webkit-scrollbar-thumb {
  background-color: #c9cdd4;
  border-radius: 3px;
}

.admin-sidebar::-webkit-scrollbar-track {
  background-color: transparent;
}

.admin-content::-webkit-scrollbar {
  width: 8px;
}

.admin-content::-webkit-scrollbar-thumb {
  background-color: #c9cdd4;
  border-radius: 4px;
}

.admin-content::-webkit-scrollbar-track {
  background-color: transparent;
}
</style>
