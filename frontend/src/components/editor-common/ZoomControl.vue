<template>
  <div class="zoom-control">
    <button class="zoom-btn" @click="handleZoomOut" title="缩小">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="5" y1="12" x2="19" y2="12"></line>
      </svg>
    </button>
    <span class="zoom-value">{{ formatZoom(editorStore.zoomLevel) }}</span>
    <button class="zoom-btn" @click="handleZoomIn" title="放大">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="12" y1="5" x2="12" y2="19"></line>
        <line x1="5" y1="12" x2="19" y2="12"></line>
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import { useGraphEditorStore } from '@/stores/graph/editor.store'
import { formatZoom } from '@/utils/graph-editor/helpers'

const editorStore = useGraphEditorStore()

const graphActions = inject<{
  zoomIn: () => void
  zoomOut: () => void
} | null>('graphActions', null)

function handleZoomIn() {
  graphActions?.zoomIn()
}

function handleZoomOut() {
  graphActions?.zoomOut()
}
</script>

<style scoped>
.zoom-control {
  display: flex;
  align-items: center;
  gap: 4px;
  background: var(--bg-panel);
  border-radius: var(--border-radius);
  padding: 2px;
}

.zoom-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: var(--border-radius-sm);
  color: #374151;
  background: transparent;
  transition: all var(--transition-fast);
}

.zoom-btn:hover {
  background: var(--bg-hover);
  color: #111827;
}

.zoom-value {
  font-size: var(--font-size-xs);
  color: #6b7280;
  min-width: 40px;
  text-align: center;
  font-variant-numeric: tabular-nums;
}
</style>
