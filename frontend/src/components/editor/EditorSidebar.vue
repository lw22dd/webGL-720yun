<template>
  <div class="sidebar" :class="{ collapsed: editorStore.sidebarCollapsed }">
    <div class="sidebar-toggle" @click="editorStore.toggleSidebar">
      <svg
        width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
        stroke-linecap="round" stroke-linejoin="round"
        :style="{ transform: editorStore.sidebarCollapsed ? 'rotate(180deg)' : 'none' }"
      >
        <polyline points="15 18 9 12 15 6"></polyline>
      </svg>
    </div>
    <div class="sidebar-content" v-show="!editorStore.sidebarCollapsed">
      <PropertyPanel v-if="editorStore.hasSelection" />
      <PendingNodeList v-else />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useGraphEditorStore } from '@/stores/graphEditorStore'
import PropertyPanel from './PropertyPanel.vue'
import PendingNodeList from '@/components/editor-sidebar/PendingNodeList.vue'

const editorStore = useGraphEditorStore()
</script>

<style scoped>
.sidebar {
  width: var(--sidebar-width);
  height: 100%;
  background: var(--bg-sidebar);
  border-left: 1px solid var(--border-color);
  display: flex;
  flex-shrink: 0;
  position: relative;
  transition: width var(--transition-normal);
}

.sidebar.collapsed {
  width: 0;
  border-left: none;
}

.sidebar-toggle {
  position: absolute;
  left: -16px;
  top: 50%;
  transform: translateY(-50%);
  width: 16px;
  height: 48px;
  background: var(--bg-sidebar);
  border: 1px solid var(--border-color);
  border-right: none;
  border-radius: 6px 0 0 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #6b7280;
  z-index: 10;
  transition: all var(--transition-fast);
}

.sidebar-toggle:hover {
  background: var(--bg-hover);
  color: #374151;
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}
</style>
