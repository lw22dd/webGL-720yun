<template>
  <div class="admin-side-menu">
    <t-menu
      v-model:value="activePath"
      theme="light"
      class="custom-side-menu"
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

      <t-menu-item value="/admin/settings">
        <template #icon>
          <t-icon name="setting" />
        </template>
        <span>系统设置</span>
      </t-menu-item>

      <t-menu-item value="/admin/logs">
        <template #icon>
          <t-icon name="history" />
        </template>
        <span>日志管理</span>
      </t-menu-item>
    </t-menu>
    
    <div class="menu-footer">
      <div class="version-tag">v1.0.0 Stable</div>
    </div>
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
  padding: 12px 0;
}

.custom-side-menu {
  flex: 1;
  border-right: none;
}

:deep(.t-menu) {
  width: 100% !important;
}

:deep(.t-menu__item) {
  height: 48px;
  line-height: 48px;
  margin: 4px 12px;
  border-radius: 10px;
  width: auto !important;
  transition: all 0.3s cubic-bezier(0.34, 0.69, 0.1, 1);
  color: #4e5969;
}

:deep(.t-menu__item:hover) {
  background-color: #f2f3f5;
  color: #1d2129;
}

:deep(.t-menu__item.t-is-active) {
  background: linear-gradient(135deg, rgba(0, 82, 217, 0.1), rgba(0, 82, 217, 0.05));
  color: #0052d9;
  font-weight: 600;
}

:deep(.t-menu__item.t-is-active::after) {
  display: none;
}

:deep(.t-menu__item .t-icon) {
  font-size: 18px;
  margin-right: 8px;
  transition: transform 0.3s;
}

:deep(.t-menu__item:hover .t-icon) {
  transform: scale(1.1);
}

.menu-footer {
  padding: 16px 24px;
  border-top: 1px solid #f2f3f5;
  display: flex;
  justify-content: center;
}

.version-tag {
  font-size: 12px;
  color: #c9cdd4;
  font-family: 'JetBrains Mono', monospace;
}
</style>
