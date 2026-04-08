<template>
  <div class="space-panel" :class="{ 'space-panel--visible': visible }">
    <div class="panel-header">
      <h2 class="panel-title">{{ space?.name }}</h2>
      <t-button theme="default" variant="text" @click="handleClose">
        <CloseIcon />
      </t-button>
    </div>

    <div class="panel-content">
      <div class="space-info">
        <img
          :src="space?.cover_url"
          :alt="space?.name"
          class="space-cover"
        />
        <div class="space-meta">
          <span class="space-location">
            <LocationIcon /> {{ space?.province }} {{ space?.city }}
          </span>
        </div>
        <p class="space-description">{{ space?.description }}</p>
      </div>

      <div class="topology-section">
        <h3 class="section-title">景点拓扑图</h3>
        <SceneTopology :scenes="space?.scenes || []" @scene-click="handleSceneClick" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { CloseIcon, LocationIcon } from 'tdesign-icons-vue-next'
import type { MockSpace, MockScene } from '@/utils/mockData'
import SceneTopology from './SceneTopology.vue'

defineProps<{
  visible: boolean
  space: MockSpace | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'sceneClick', scene: MockScene): void
}>()

const handleClose = () => {
  emit('close')
}

const handleSceneClick = (scene: MockScene) => {
  emit('sceneClick', scene)
}
</script>

<style scoped>
.space-panel {
  height: 100%;
  background: #ffffff;
  border-left: 1px solid #e5e6eb;
  display: flex;
  flex-direction: column;
  box-shadow: -10px 0 40px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
  overflow: hidden;
  width: 0;
  transition: width 0.3s ease;
}

.space-panel--visible {
  width: 420px;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e6eb;
}

.panel-title {
  color: #1d2129;
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  background: linear-gradient(135deg, #0052d9, #4080ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.space-info {
  margin-bottom: 24px;
}

.space-cover {
  width: 100%;
  height: 180px;
  object-fit: cover;
  border-radius: 12px;
  border: 1px solid #e5e6eb;
  margin-bottom: 16px;
}

.space-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
}

.space-location {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #86909c;
  font-size: 13px;
}

.space-description {
  color: #4e5969;
  font-size: 14px;
  line-height: 1.6;
  margin: 0;
}

.topology-section {
  margin-top: 24px;
}

.section-title {
  color: #1d2129;
  font-size: 14px;
  font-weight: 500;
  margin: 0 0 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-title::before {
  content: '';
  width: 3px;
  height: 14px;
  background: linear-gradient(180deg, #0052d9, #4080ff);
  border-radius: 2px;
}
</style>
