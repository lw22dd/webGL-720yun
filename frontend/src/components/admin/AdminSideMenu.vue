<template>
  <div class="admin-side-menu">
    <div class="menu-header">
      <t-icon name="user-circle" size="24" />
      <span class="menu-title">管理中心</span>
    </div>

    <t-menu
      v-model:defaultValue="activePath"
      theme="light"
      @change="handleMenuChange"
    >
      <t-menu-item value="/admin/users">
        <template #icon>
          <t-icon name="user" />
        </template>
        <span>用户管理</span>
      </t-menu-item>

      <t-menu-item value="/admin/batch-register">
        <template #icon>
          <t-icon name="upload" />
        </template>
        <span>批量注册</span>
      </t-menu-item>

      <t-menu-item value="/admin/spaces">
        <template #icon>
          <t-icon name="file" />
        </template>
        <span>空间管理</span>
      </t-menu-item>

      <t-menu-item value="/admin/scenes">
        <template #icon>
          <t-icon name="photo" />
        </template>
        <span>场景管理</span>
      </t-menu-item>

      <t-menu-item value="/admin/settings">
        <template #icon>
          <t-icon name="setting" />
        </template>
        <span>系统设置</span>
      </t-menu-item>
    </t-menu>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const activePath = ref(route.path)

watch(
  () => route.path,
  (newPath) => {
    activePath.value = newPath
  }
)

const handleMenuChange = (value: string) => {
  router.push(value)
}
</script>

<style scoped>
.admin-side-menu {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #fff;
}

.menu-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px 20px;
  border-bottom: 1px solid #e5e6eb;
  color: #0052d9;
}

.menu-title {
  font-size: 18px;
  font-weight: 700;
  white-space: nowrap;
  letter-spacing: -0.02em;
}

:deep(.t-menu) {
  border-right: none;
  padding: 8px;
}

:deep(.t-menu__item) {
  border-radius: 8px;
  margin: 2px 0;
  width: 100%;
  padding: 10px 12px;
  font-size: 14px;
  font-weight: 500;
  color: #4e5969;
}

:deep(.t-menu__item .t-icon) {
  font-size: 18px;
}

:deep(.t-menu__item--hover) {
  background-color: #f2f3f5;
  color: #1d2129;
}

:deep(.t-menu__item--active) {
  background-color: #e6f0ff;
  color: #0052d9;
}

:deep(.t-menu__item--active::before) {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 24px;
  background-color: #0052d9;
  border-radius: 0 4px 4px 0;
}
</style>
