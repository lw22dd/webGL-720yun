<template>
  <div class="map-toolbar">
    <div class="toolbar-section">
      <span class="toolbar-label">底图</span>
      <div class="toolbar-buttons">
        <button
          :class="['toolbar-btn', { active: basemapType === 'vector' }]"
          @click="$emit('switchBasemap', 'vector')"
          title="矢量图"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M1 6v16l7-4 8 4 7-4V2l-7 4-8-4-7 4z"/>
            <path d="M8 2v16"/>
            <path d="M16 6v16"/>
          </svg>
        </button>
        <button
          :class="['toolbar-btn', { active: basemapType === 'satellite' }]"
          @click="$emit('switchBasemap', 'satellite')"
          title="卫星图"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <path d="M2 12h20"/>
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
          </svg>
        </button>
      </div>
    </div>

    <div class="toolbar-section" v-if="showPOIToggle">
      <span class="toolbar-label">POI</span>
      <div class="toolbar-buttons">
        <button
          :class="['toolbar-btn', { active: showPOI }]"
          @click="$emit('togglePOI', !showPOI)"
          title="显示/隐藏POI"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
            <circle cx="12" cy="10" r="3"/>
          </svg>
        </button>
      </div>
    </div>

    <div class="toolbar-section" v-if="showMarkerToggle">
      <span class="toolbar-label">节点</span>
      <div class="toolbar-buttons">
        <button
          :class="['toolbar-btn', { active: showMarkers }]"
          @click="$emit('toggleMarkers', !showMarkers)"
          title="显示/隐藏场景节点"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="7" height="7"/>
            <rect x="14" y="3" width="7" height="7"/>
            <rect x="14" y="14" width="7" height="7"/>
            <rect x="3" y="14" width="7" height="7"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { BasemapType } from '@/composables/useL7Scene'

defineProps<{
  basemapType: BasemapType
  showPOI?: boolean
  showMarkers?: boolean
  showPOIToggle?: boolean
  showMarkerToggle?: boolean
}>()

defineEmits<{
  (e: 'switchBasemap', type: BasemapType): void
  (e: 'togglePOI', show: boolean): void
  (e: 'toggleMarkers', show: boolean): void
}>()
</script>

<style scoped>
.map-toolbar {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 100;
  display: flex;
  gap: 12px;
  padding: 10px 14px;
  background: rgba(26, 35, 50, 0.92);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.toolbar-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  font-weight: 500;
}

.toolbar-buttons {
  display: flex;
  gap: 4px;
}

.toolbar-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  transition: all 0.2s;
}

.toolbar-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.toolbar-btn.active {
  background: rgba(14, 165, 233, 0.2);
  border-color: rgba(14, 165, 233, 0.4);
  color: #0EA5E9;
}
</style>
