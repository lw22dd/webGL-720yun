<template>
  <div class="upload-btn-wrapper">
    <button class="upload-btn" @click="triggerUpload">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
        <circle cx="8.5" cy="8.5" r="1.5"></circle>
        <polyline points="21 15 16 10 5 21"></polyline>
      </svg>
      <span>底图</span>
    </button>
    <input
      ref="fileInput"
      type="file"
      accept="image/*"
      style="display: none"
      @change="handleFileChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, inject } from 'vue'

const fileInput = ref<HTMLInputElement | null>(null)

const graphActions = inject<{
  uploadBackground: (file: File) => void
} | null>('graphActions', null)

function triggerUpload() {
  fileInput.value?.click()
}

function handleFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file && graphActions) {
    graphActions.uploadBackground(file)
  }
  target.value = ''
}
</script>

<style scoped>
.upload-btn-wrapper {
  position: relative;
}

.upload-btn {
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

.upload-btn:hover {
  background: var(--bg-hover);
  color: #111827;
}
</style>
