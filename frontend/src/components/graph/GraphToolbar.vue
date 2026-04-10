<template>
  <div class="graph-toolbar">
    <t-button-group>
      <t-button
        variant="outline"
        size="small"
        :disabled="!canUndo"
        @click="$emit('undo')"
      >
        <template #icon><t-icon name="undo" /></template>
      </t-button>
      <t-button
        variant="outline"
        size="small"
        :disabled="!canRedo"
        @click="$emit('redo')"
      >
        <template #icon><t-icon name="redo" /></template>
      </t-button>
    </t-button-group>

    <t-divider layout="vertical" />

    <t-button-group>
      <t-button variant="outline" size="small" @click="$emit('zoomOut')">
        <template #icon><t-icon name="remove" /></template>
      </t-button>
      <t-button variant="outline" size="small" class="zoom-display">
        {{ Math.round(zoom * 100) }}%
      </t-button>
      <t-button variant="outline" size="small" @click="$emit('zoomIn')">
        <template #icon><t-icon name="add" /></template>
      </t-button>
    </t-button-group>

    <t-divider layout="vertical" />

    <t-button
      variant="outline"
      size="small"
      @click="$emit('fitContent')"
    >
      <template #icon><t-icon name="fullscreen" /></template>
      适应画布
    </t-button>

    <t-divider layout="vertical" />

    <t-button
      theme="primary"
      size="small"
      @click="$emit('save')"
    >
      <template #icon><t-icon name="save" /></template>
      保存
    </t-button>
  </div>
</template>

<script setup lang="ts">
interface Props {
  canUndo: boolean
  canRedo: boolean
  zoom: number
}

defineProps<Props>()

defineEmits<{
  (e: 'undo'): void
  (e: 'redo'): void
  (e: 'zoomIn'): void
  (e: 'zoomOut'): void
  (e: 'zoomReset'): void
  (e: 'fitContent'): void
  (e: 'save'): void
}>()
</script>

<style scoped>
.graph-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
}

.zoom-display {
  min-width: 60px;
  cursor: default;
}

.zoom-display:hover {
  background-color: inherit;
}
</style>
