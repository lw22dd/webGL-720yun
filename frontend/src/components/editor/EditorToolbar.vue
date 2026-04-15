<template>
  <div class="toolbar">
    <div class="toolbar-left">
      <button class="toolbar-btn back-btn" @click="handleBack">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
          stroke-linecap="round" stroke-linejoin="round">
          <polyline points="15 18 9 12 15 6"></polyline>
        </svg>
        <span>返回</span>
      </button>
      <div class="toolbar-divider"></div>
      <h1 class="toolbar-title">{{ sceneStore.spaceInfo?.name || '图编辑器' }}</h1>
      <span class="toolbar-badge">
        已放置: {{ sceneStore.placedCount }} / {{ sceneStore.totalCount }}
      </span>
    </div>

    <div class="toolbar-center">
      <button class="toolbar-btn icon-btn" :disabled="!canUndo" title="撤销 (Ctrl+Z)" @click="handleUndo">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
          stroke-linecap="round" stroke-linejoin="round">
          <polyline points="1 4 1 10 7 10"></polyline>
          <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path>
        </svg>
      </button>
      <button class="toolbar-btn icon-btn" :disabled="!canRedo" title="重做 (Ctrl+Shift+Z)" @click="handleRedo">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
          stroke-linecap="round" stroke-linejoin="round">
          <polyline points="23 4 23 10 17 10"></polyline>
          <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path>
        </svg>
      </button>
      <div class="toolbar-divider"></div>
      <ZoomControl />
    </div>

    <div class="toolbar-right">
      <UploadButton />
      <button class="toolbar-btn icon-btn" title="适应画布" @click="handleFitCanvas">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
          stroke-linecap="round" stroke-linejoin="round">
          <path d="M8 3H5a2 2 0 0 0-2 2v3m18 0V5a2 2 0 0 0-2-2h-3m0 18h3a2 2 0 0 0 2-2v-3M3 16v3a2 2 0 0 0 2 2h3">
          </path>
        </svg>
      </button>
      <button class="toolbar-btn save-btn" :class="{ saving: sceneStore.saving }" :disabled="sceneStore.saving"
        @click="handleSave">
        <svg v-if="!sceneStore.saving" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
          stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
          <polyline points="17 21 17 13 7 13 7 21"></polyline>
          <polyline points="7 3 7 8 15 8"></polyline>
        </svg>
        <span v-if="sceneStore.saving" class="save-spinner"></span>
        <span>{{ sceneStore.saving ? '保存中...' : '保存' }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject, type Ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGraphSceneStore } from '@/stores/graph/scene.store'
import ZoomControl from '@/components/editor-common/ZoomControl.vue'
import UploadButton from '@/components/editor-common/UploadButton.vue'

const router = useRouter()
const sceneStore = useGraphSceneStore()

const historyState = inject<{
  canUndo: Ref<boolean>
  canRedo: Ref<boolean>
} | null>('historyState', null)

const graphActions = inject<{
  undo: () => void
  redo: () => void
  fitCanvas: () => void
} | null>('graphActions', null)

const canUndo = computed(() => historyState?.canUndo?.value ?? false)
const canRedo = computed(() => historyState?.canRedo?.value ?? false)

function handleBack() {
  router.push('/admin/spaces')
}

function handleUndo() {
  console.log('[Toolbar] handleUndo 被调用')
  console.log('[Toolbar] graphActions:', graphActions)
  console.log('[Toolbar] graphActions?.undo:', graphActions?.undo)
  if (graphActions?.undo) {
    console.log('[Toolbar] 调用 graphActions.undo()')
    const result = graphActions.undo()
    console.log('[Toolbar] undo 返回结果:', result)
  } else {
    console.log('[Toolbar] graphActions?.undo 不存在')
  }
}

function handleRedo() {
  console.log('[Toolbar] handleRedo 被调用')
  console.log('[Toolbar] graphActions:', graphActions)
  console.log('[Toolbar] graphActions?.redo:', graphActions?.redo)
  if (graphActions?.redo) {
    console.log('[Toolbar] 调用 graphActions.redo()')
    const result = graphActions.redo()
    console.log('[Toolbar] redo 返回结果:', result)
  } else {
    console.log('[Toolbar] graphActions?.redo 不存在')
  }
}

function handleFitCanvas() {
  graphActions?.fitCanvas()
}

async function handleSave() {
  const success = await sceneStore.saveScene()
  if (success) {
    console.log('保存成功')
  }
}
</script>

<style scoped>
.toolbar {
  height: var(--toolbar-height);
  background: var(--bg-toolbar);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--spacing-lg);
  flex-shrink: 0;
  z-index: 100;
}

.toolbar-left,
.toolbar-center,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.toolbar-left {
  flex: 1;
}

.toolbar-center {
  flex: 0;
}

.toolbar-right {
  flex: 1;
  justify-content: flex-end;
}

.toolbar-title {
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: #111827;
  white-space: nowrap;
}

.toolbar-badge {
  font-size: var(--font-size-xs);
  color: #6b7280;
  background: var(--bg-panel);
  padding: 2px 10px;
  border-radius: 12px;
  white-space: nowrap;
}

.toolbar-divider {
  width: 1px;
  height: 24px;
  background: var(--border-color);
  margin: 0 var(--spacing-sm);
}

.toolbar-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: var(--border-radius-sm);
  font-size: var(--font-size-sm);
  color: #374151;
  background: transparent;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.toolbar-btn:hover:not(:disabled) {
  background: var(--bg-hover);
  color: #111827;
}

.toolbar-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.back-btn {
  padding: 6px 10px;
}

.back-btn span {
  font-size: var(--font-size-sm);
}

.icon-btn {
  padding: 8px;
  border-radius: var(--border-radius);
}

.icon-btn svg {
  flex-shrink: 0;
}

.save-btn {
  background: var(--color-primary);
  color: #fff;
  padding: 6px 16px;
  border-radius: var(--border-radius);
  font-weight: 500;
}

.save-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.save-btn.saving {
  opacity: 0.8;
}

.save-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
