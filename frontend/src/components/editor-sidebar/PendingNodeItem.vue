<template>
  <div
    class="pending-item"
    :class="{ 'has-pano': node.thumbnail_url }"
    draggable="true"
    @mousedown="handleDragStart"
  >
    <div class="item-thumb">
      <div class="thumb-placeholder" :class="{ 'has-pano': node.thumbnail_url }">
        <svg v-if="node.thumbnail_url" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <circle cx="12" cy="12" r="3"></circle>
        </svg>
        <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
          <circle cx="8.5" cy="8.5" r="1.5"></circle>
          <polyline points="21 15 16 10 5 21"></polyline>
        </svg>
      </div>
    </div>
    <div class="item-info">
      <div class="item-name">{{ node.title }}</div>
      <div class="item-meta">
        <span class="item-id">{{ node.id }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { SceneNodeData } from '@/models/graphEditor/node'
import { DND_ACTIONS_KEY } from '@/composables/useGraphDnd'

const props = defineProps<{
  node: SceneNodeData
}>()

const dndActions = inject(DND_ACTIONS_KEY, null)

function handleDragStart(e: MouseEvent) {
  e.preventDefault()
  dndActions?.startDrag(props.node, e)
}
</script>

<style scoped>
.pending-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border-radius: var(--border-radius);
  cursor: grab;
  transition: all var(--transition-fast);
  border: 1px solid transparent;
}

.pending-item:hover {
  background: var(--bg-hover);
  border-color: var(--border-color);
}

.pending-item:active {
  cursor: grabbing;
  opacity: 0.7;
}

.item-thumb {
  flex-shrink: 0;
}

.thumb-placeholder {
  width: 40px;
  height: 40px;
  border-radius: var(--border-radius-sm);
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
}

.thumb-placeholder.has-pano {
  background: #d1fae5;
  color: #059669;
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: var(--font-size-sm);
  font-weight: 500;
  color: #1f2937;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 2px;
}

.item-id {
  font-size: 11px;
  color: #9ca3af;
  font-variant-numeric: tabular-nums;
}
</style>
