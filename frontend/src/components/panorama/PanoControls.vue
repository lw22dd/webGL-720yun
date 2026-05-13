<template>
  <div class="pano-controls">
    <t-button theme="default" variant="outline" size="small" @click="$emit('reset')" title="重置视角">
      ⟲
    </t-button>
    <t-button
      :theme="autoRotate ? 'primary' : 'default'"
      variant="outline"
      size="small"
      @click="$emit('toggleRotate')"
      title="自动旋转"
    >
      ↻
    </t-button>
    <t-button theme="default" variant="outline" size="small" @click="$emit('zoomIn')" title="放大">
      +
    </t-button>
    <t-button theme="default" variant="outline" size="small" @click="$emit('zoomOut')" title="缩小">
      −
    </t-button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  autoRotate: boolean
}>()

defineEmits<{
  (e: 'reset'): void
  (e: 'toggleRotate'): void
  (e: 'zoomIn'): void
  (e: 'zoomOut'): void
}>()
</script>

<style scoped>
.pano-controls {
  position: absolute;
  right: 24px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 8px;
  z-index: 100;
}

.pano-controls :deep(.t-button) {
  width: 40px;
  height: 40px;
  padding: 0;
  background: rgba(255, 255, 255, 0.1) !important;
  border: 1px solid rgba(255, 255, 255, 0.2) !important;
  color: #fff !important;
  backdrop-filter: blur(10px);
  font-size: 18px;
}

.pano-controls :deep(.t-button:hover) {
  background: rgba(255, 255, 255, 0.2) !important;
}

.pano-controls :deep(.t-button.t-button--theme-primary) {
  background: #0EA5E9 !important;
  border-color: #0EA5E9 !important;
}
</style>