<template>
  <transition name="slide">
    <div v-if="visible" class="space-panel">
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
  </transition>
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
  position: fixed;
  top: 0;
  right: 0;
  width: 420px;
  height: 100vh;
  background: rgba(10, 14, 39, 0.95);
  backdrop-filter: blur(20px);
  border-left: 1px solid rgba(0, 240, 255, 0.2);
  z-index: 1000;
  display: flex;
  flex-direction: column;
  box-shadow: -10px 0 40px rgba(0, 0, 0, 0.5);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid rgba(0, 240, 255, 0.15);
}

.panel-title {
  color: #fff;
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  background: linear-gradient(135deg, #00f0ff, #7b2fff);
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
  border: 1px solid rgba(0, 240, 255, 0.2);
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
  color: rgba(255, 255, 255, 0.6);
  font-size: 13px;
}

.space-description {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  line-height: 1.6;
  margin: 0;
}

.topology-section {
  margin-top: 24px;
}

.section-title {
  color: #fff;
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
  background: linear-gradient(180deg, #00f0ff, #7b2fff);
  border-radius: 2px;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
}
</style>
