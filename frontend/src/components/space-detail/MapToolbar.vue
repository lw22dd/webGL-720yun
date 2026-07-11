<!--高德地图工具栏-->
<template>
  <div class="map-toolbar-container">
    <div class="map-toolbar">
      <button class="back-btn" @click="$emit('back')">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M19 12H5"/>
          <path d="M12 19l-7-7 7-7"/>
        </svg>
        返回总览
      </button>

      <div class="space-info-bar">
        <div class="space-info">
          <h1 class="space-name">{{ spaceName }}</h1>
          <span class="space-location">📍 {{ spaceProvince }} {{ spaceCity }}</span>
          <span class="space-divider">·</span>
          <span class="space-scene-count">🎥 {{ sceneCount }} 个场景节点</span>
        </div>
      </div>

      <div class="toolbar-actions">
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
    </div>

    <div v-if="showClickHint" class="click-hint">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/>
        <path d="M12 8v8M8 12h8"/>
      </svg>
      点击地图添加场景节点
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
  spaceName?: string
  spaceProvince?: string
  spaceCity?: string
  sceneCount?: number
  showClickHint?: boolean
}>()

defineEmits<{
  (e: 'switchBasemap', type: BasemapType): void
  (e: 'togglePOI', show: boolean): void
  (e: 'toggleMarkers', show: boolean): void
  (e: 'back'): void
}>()
</script>

<style scoped>
.map-toolbar-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 100;
  padding: 20px;
  pointer-events: none;
}

.map-toolbar-container > * {
  pointer-events: auto;
}

.map-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;

  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 12px;
  color: #333;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.85);
  color: #1a1a1a;
  border-color: rgba(0, 0, 0, 0.12);
}

.space-info-bar {
  display: flex;
  justify-content: center;
  padding: 10px 24px;

  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
}

.space-info {
  display: flex;
  align-items: center;
  gap: 12px;
  white-space: nowrap;
}

.space-name {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.space-location,
.space-scene-count {
  font-size: 13px;
  color: #666;
}

.space-divider {
  color: #ddd;
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
  padding: 10px 16px;

  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
}

.toolbar-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-label {
  font-size: 12px;
  color: #666;
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
  background: rgba(0, 0, 0, 0.04);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 6px;
  color: #555;
  cursor: pointer;
  transition: all 0.2s;
}

.toolbar-btn:hover {
  background: rgba(0, 0, 0, 0.08);
  color: #333;
}

.toolbar-btn.active {
  background: rgba(14, 165, 233, 0.1);
  border-color: rgba(14, 165, 233, 0.3);
  color: #0EA5E9;
}

.click-hint {
  position: absolute;
  bottom: 20px;
  right: 20px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: 10px;
  font-size: 13px;
  color: #555;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}
</style>
