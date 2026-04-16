<template>
  <div class="scene-panel">
    <div class="panel-header">
      <h3 class="panel-title">场景节点</h3>
      <span class="panel-count">{{ scenes.length }}</span>
    </div>
    <div class="scene-list">
      <div
        v-for="(scene, index) in scenes"
        :key="scene.id"
        :class="['scene-item', { active: activeSceneId === scene.id }]"
        @click="handleSceneClick(scene)"
      >
        <div class="scene-item-dot" :style="{ background: getMarkerColor(index) }"></div>
        <div class="scene-item-info">
          <span class="scene-item-name">{{ scene.title }}</span>
          <span class="scene-item-coords">
            {{ scene.longitude.toFixed(4) }}, {{ scene.latitude.toFixed(4) }}
          </span>
        </div>
        <span v-if="scene.hotspots && scene.hotspots.length > 0" class="scene-item-badge">
          {{ scene.hotspots.length }}
        </span>
      </div>
    </div>
    <div v-if="scenes.length === 0" class="scene-empty">
      <p>暂无场景节点</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SceneMarkerData } from '@/composables/useSceneMarkers'

const MARKER_COLORS = [
  '#0EA5E9', '#2563EB', '#10B981', '#F59E0B',
  '#8B5CF6', '#EC4899', '#DC2626', '#059669',
]

defineProps<{
  scenes: SceneMarkerData[]
  activeSceneId: number | null
}>()

const emit = defineEmits<{
  (e: 'sceneClick', scene: SceneMarkerData): void
}>()

function getMarkerColor(index: number): string {
  return MARKER_COLORS[index % MARKER_COLORS.length]
}

function handleSceneClick(scene: SceneMarkerData) {
  emit('sceneClick', scene)
}
</script>

<style scoped>
.scene-panel {
  position: absolute;
  bottom: 20px;
  left: 20px;
  z-index: 100;
  width: 280px;
  max-height: 400px;
  background: rgba(26, 35, 50, 0.92);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 14px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.panel-title {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.panel-count {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  background: rgba(14, 165, 233, 0.2);
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  color: #0EA5E9;
}

.scene-list {
  max-height: 340px;
  overflow-y: auto;
  padding: 6px;
}

.scene-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 10px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.scene-item:hover {
  background: rgba(255, 255, 255, 0.06);
}

.scene-item.active {
  background: rgba(14, 165, 233, 0.12);
}

.scene-item-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  box-shadow: 0 0 6px rgba(14, 165, 233, 0.4);
}

.scene-item-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.scene-item-name {
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.scene-item-coords {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
  font-family: 'SF Mono', 'Fira Code', monospace;
}

.scene-item-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  padding: 0 5px;
  background: rgba(245, 158, 11, 0.2);
  border-radius: 10px;
  font-size: 10px;
  font-weight: 600;
  color: #F59E0B;
}

.scene-empty {
  padding: 24px;
  text-align: center;
  color: rgba(255, 255, 255, 0.4);
  font-size: 13px;
}
</style>
