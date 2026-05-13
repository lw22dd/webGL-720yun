<template>
  <div class="admin-layout">
    <!-- 全宽 Header -->
    <header class="admin-header">
      <div class="header-left">
        <div class="logo-area" @click="router.push('/admin')">
          <img src="/logo-cuit.jpg" alt="成都信息工程大学" class="logo-image" />
        </div>
        
        <t-button
          variant="text"
          shape="circle"
          class="sidebar-toggle"
          @click="toggleSidebar"
        >
          <template #icon>
            <t-icon :name="sidebarHidden ? 'indent-right' : 'indent-left'" size="20" />
          </template>
        </t-button>

        <nav class="breadcrumb-container">
          <t-breadcrumb class="admin-breadcrumb">
            <t-breadcrumb-item @click="() => router.push('/')">首页</t-breadcrumb-item>
            <t-breadcrumb-item @click="() => router.push('/admin')">控制台</t-breadcrumb-item>
            <t-breadcrumb-item v-if="route.path !== '/admin'">{{ currentPageName }}</t-breadcrumb-item>
          </t-breadcrumb>
        </nav>
      </div>

      <div class="header-right">
        <div class="action-items">
          <t-tooltip content="刷新页面">
            <t-button variant="text" shape="circle" @click="refreshPage">
              <t-icon name="refresh" />
            </t-button>
          </t-tooltip>
        </div>

        <t-dropdown trigger="click">
          <div class="user-info-card">
            <t-avatar size="32px" class="user-avatar">
              {{ (userStore.userInfo.username || 'A').charAt(0).toUpperCase() }}
            </t-avatar>
            <div class="user-meta">
              <span class="user-name">{{ userStore.userInfo.username || '管理员' }}</span>
              <span class="user-role">{{ isAdmin ? '超级管理员' : '运营人员' }}</span>
            </div>
            <t-icon name="chevron-down" size="14" class="dropdown-arrow" />
          </div>
          <template #dropdown>
            <t-dropdown-menu class="user-dropdown">
              <t-dropdown-item @click="navigateToHome">
                <t-icon name="home" /> <span>返回首页</span>
              </t-dropdown-item>
              <t-dropdown-item @click="handleLogout" class="logout-item">
                <t-icon name="logout" /> <span>安全退出</span>
              </t-dropdown-item>
            </t-dropdown-menu>
          </template>
        </t-dropdown>
      </div>
    </header>

    <div class="admin-body">
      <!-- 侧边栏 -->
      <aside class="admin-sidebar" :class="{ 'sidebar-hidden': sidebarHidden }">
        <AdminSideMenu />
      </aside>

      <!-- 主体内容区 -->
      <main class="admin-main">
        <div class="admin-content-wrapper">
          <router-view v-slot="{ Component }">
            <transition name="page-fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
        
        <footer class="admin-footer">
          <p>© {{ new Date().getFullYear() }} 720全景漫游系统 · 数字化管理后台</p>
        </footer>
      </main>
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

const isAdmin = computed(() => userStore.userInfo.role_id === 1 || userStore.userInfo.is_super_admin)

const currentPageName = computed(() => {
  const pathMap: Record<string, string> = {
    '/admin': '数据概览',
    '/admin/users': '用户管理',
    '/admin/batch-register': '批量导入',
    '/admin/spaces': '全景空间',
    '/admin/settings': '全局配置',
    '/admin/logs': '操作审计'
  }
  return pathMap[route.path] || route.path.split('/').pop()
})

const toggleSidebar = () => {
  sidebarHidden.value = !sidebarHidden.value
}

const refreshPage = () => {
  window.location.reload()
}

const navigateToHome = () => {
  router.push('/')
}

const handleLogout = async () => {
  try {
    await UserApi.logout()
    userStore.logout()
    MessagePlugin.success('已安全退出')
    router.push('/login')
  } catch (error: any) {
    userStore.logout()
    router.push('/login')
  }
}

onMounted(() => {
  if (!userStore.isLogin || !(userStore.userInfo.role_id === 1 || userStore.userInfo.is_super_admin)) {
    MessagePlugin.warning('无权访问管理后台')
    router.push('/')
  }
})
</script>

<style scoped>
.admin-layout {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100vh;
  background-color: #f6f8f9;
  color: #1d2129;
  overflow: hidden;
}

/* Header: 现代磨砂玻璃效果 */
.admin-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(229, 230, 235, 0.8);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  z-index: 1001;
  transition: all 0.3s;
}

.header-left {
  display: flex;
  align-items: center;
}

.logo-area {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding-right: 24px;
  margin-right: 12px;
  border-right: 1.5px solid #f2f3f5;
  transition: opacity 0.2s;
  height: 64px;
}

.logo-area:hover {
  opacity: 0.8;
}

.logo-image {
  height: 64px;
  width: auto;
  border-radius: 8px;
  object-fit: contain;
}

.logo-text {
  font-size: 17px;
  font-weight: 800;
  color: #1d2129;
  letter-spacing: -0.01em;
}

.sidebar-toggle {
  margin-right: 16px;
  color: #4e5969;
}

.sidebar-toggle:hover {
  background-color: #f2f3f5;
  color: #0052d9;
}

.breadcrumb-container {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.action-items {
  display: flex;
  gap: 8px;
  padding-right: 16px;
  border-right: 1.5px solid #f2f3f5;
}

.user-info-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 4px 12px 4px 8px;
  border-radius: 12px;
  cursor: pointer;
  transition: background-color 0.3s;
  background: #f8f9fa;
  border: 1px solid #f0f1f2;
}

.user-info-card:hover {
  background: #f2f3f5;
  border-color: #e5e6eb;
}

.user-avatar {
  background: linear-gradient(135deg, #0052d9, #0ea5e9);
  color: #fff;
  font-weight: 700;
  box-shadow: 0 2px 6px rgba(0, 82, 217, 0.15);
}

.user-meta {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.user-name {
  font-size: 13px;
  font-weight: 700;
  color: #1d2129;
}

.user-role {
  font-size: 11px;
  color: #86909c;
  font-weight: 500;
}

.dropdown-arrow {
  color: #c9cdd4;
  margin-left: 4px;
}

/* Body 结构 */
.admin-body {
  display: flex;
  margin-top: 64px;
  height: calc(100vh - 64px);
  width: 100%;
}

.admin-sidebar {
  width: 240px;
  height: 100%;
  background-color: #ffffff;
  border-right: 1px solid #f2f3f5;
  transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
  z-index: 1000;
}

.admin-sidebar.sidebar-hidden {
  width: 0;
  transform: translateX(-100%);
  opacity: 0;
  pointer-events: none;
}

.admin-main {
  flex: 1;
  height: 100%;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  position: relative;
}

.admin-content-wrapper {
  flex: 1;
  padding: 24px 32px;
}

.admin-footer {
  padding: 20px 32px;
  text-align: center;
  color: #86909c;
  font-size: 12px;
  border-top: 1px solid rgba(229, 230, 235, 0.5);
  background: #ffffff;
}

/* 页面切换动画 */
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.25s, transform 0.25s;
}

.page-fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.page-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* 针对深色模式或特定 TDesign 组件的微调 */
:deep(.user-dropdown) {
  min-width: 160px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1) !important;
}

:deep(.logout-item) {
  color: #d54941 !important;
}

:deep(.logout-item:hover) {
  background-color: #fff1f0 !important;
}
</style>
