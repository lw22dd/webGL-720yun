<template>
  <div class="pending-list">
    <div class="list-header">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="8" y1="6" x2="21" y2="6"></line>
        <line x1="8" y1="12" x2="21" y2="12"></line>
        <line x1="8" y1="18" x2="21" y2="18"></line>
        <line x1="3" y1="6" x2="3.01" y2="6"></line>
        <line x1="3" y1="12" x2="3.01" y2="12"></line>
        <line x1="3" y1="18" x2="3.01" y2="18"></line>
      </svg>
      <span>待处理节点 ({{ sceneStore.pendingNodes.length }})</span>
    </div>

    <div class="list-body" v-if="sceneStore.pendingNodes.length > 0">
      <PendingNodeItem
        v-for="node in sceneStore.pendingNodes"
        :key="node.id"
        :node="node"
      />
    </div>

    <div class="list-empty" v-else>
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="#d1d5db" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
        <polyline points="22 4 12 14.01 9 11.01"></polyline>
      </svg>
      <p>所有节点已放置完成</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useGraphSceneStore } from '@/stores/graphSceneStore'
import PendingNodeItem from './PendingNodeItem.vue'

const sceneStore = useGraphSceneStore()
</script>

<style scoped>
.pending-list {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.list-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px;
  font-size: var(--font-size-base);
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}

.list-header svg {
  color: #6b7280;
}

.list-body {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.list-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px 20px;
  color: #9ca3af;
}

.list-empty p {
  font-size: var(--font-size-sm);
}
</style>
