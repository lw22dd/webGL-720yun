<template>
  <div class="space-panel" :class="{ 'space-panel--visible': visible }">
    <div class="panel-header">
      <h2 class="panel-title">{{ space?.name }}</h2>
      <t-button theme="default" variant="text" @click="handleClose">
        <template #icon><CloseIcon /></template>
      </t-button>
    </div>

    <div class="panel-content">
      <div class="space-info">
        <div class="space-cover">
          <img
            v-if="space?.cover_url"
            :src="space.cover_url"
            :alt="space?.name"
          />
          <div v-else class="space-cover-placeholder">
            <t-icon name="image" size="48" />
          </div>
          <div class="space-cover-badge">
            {{ space?.scenes?.length || space?.scene_count || 0 }} 个景点
          </div>
        </div>
        <div class="space-meta">
          <span class="space-location">
            <LocationIcon /> {{ space?.province }} {{ space?.city }}
          </span>
        </div>
        <p class="space-description">{{ space?.description || '暂无描述' }}</p>
      </div>

      <div class="scene-section">
        <h3 class="section-title">
          <span class="section-title-icon">🎥</span>
          景点列表
        </h3>
        <div class="scene-list" v-if="space?.scenes && space.scenes.length > 0">
          <div 
            v-for="scene in space.scenes" 
            :key="scene.id" 
            class="scene-item"
            @click="handleSceneClick(scene)"
          >
            <div class="scene-item-thumb">
              <img 
                v-if="scene.thumbnail_url" 
                :src="scene.thumbnail_url" 
                :alt="scene.title"
              />
              <div v-else class="scene-item-thumb-placeholder">
                <t-icon name="image" />
              </div>
            </div>
            <div class="scene-item-info">
              <div class="scene-item-name">{{ scene.title }}</div>
              <div class="scene-item-meta">
                <span class="scene-item-code">{{ scene.scene_code }}</span>
                <span v-if="scene.view_count" class="scene-item-views">
                  {{ scene.view_count }} 次浏览
                </span>
              </div>
            </div>
            <div class="scene-item-arrow">→</div>
          </div>
        </div>
        <div v-else class="scene-empty">
          <t-icon name="view-module" size="32" />
          <p>暂无景点数据</p>
        </div>
      </div>

      <div class="topology-section" v-if="space?.scenes && space.scenes.length > 0">
        <h3 class="section-title">
          <span class="section-title-icon">🗺</span>
          景点拓扑图
        </h3>
        <SceneTopology :scenes="space.scenes" @scene-click="handleSceneClick" />
      </div>

      <div class="panel-actions">
        <t-button theme="primary" size="large" block @click="handleEnterPanorama">
          进入漫游
        </t-button>
        <t-button 
          theme="default" 
          size="large" 
          block 
          @click="handleToggleFavorite"
        >
          <template #icon>
            <StarIcon :class="{ 'is-favorite': isFavorite }" />
          </template>
          {{ isFavorite ? '取消收藏' : '收藏' }}
        </t-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CloseIcon, LocationIcon, StarIcon } from 'tdesign-icons-vue-next'
import SceneTopology from './SceneTopology.vue'
import type { SpaceListItem } from '@/models/SpaceModel'

interface SceneSimple {
  id: number
  title: string
  scene_code: string
  thumbnail_url: string
  view_count: number
  sort_order: number
  longitude: number
  latitude: number
}

interface SpaceWithScenes extends SpaceListItem {
  scenes?: SceneSimple[]
}

const props = defineProps<{
  visible: boolean
  space: SpaceWithScenes | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'sceneClick', scene: SceneSimple): void
}>()

const isFavorite = computed(() => {
  if (!props.space) return false
  try {
    const favorites = JSON.parse(localStorage.getItem('favorites') || '[]')
    return favorites.some((f: any) => f.id === props.space?.id)
  } catch {
    return false
  }
})

const handleClose = () => {
  emit('close')
}

const handleSceneClick = (scene: SceneSimple) => {
  emit('sceneClick', scene)
}

const handleEnterPanorama = () => {
  if (props.space?.scenes && props.space.scenes.length > 0) {
    const entryScene = props.space.scenes.find(s => s.sort_order === 1) || props.space.scenes[0]
    handleSceneClick(entryScene)
  }
}

const handleToggleFavorite = () => {
  if (!props.space) return
  
  try {
    let favorites = JSON.parse(localStorage.getItem('favorites') || '[]')
    const index = favorites.findIndex((f: any) => f.id === props.space?.id)
    
    if (index > -1) {
      favorites.splice(index, 1)
    } else {
      favorites.push({
        id: props.space.id,
        name: props.space.name,
        cover_url: props.space.cover_url,
        province: props.space.province,
        city: props.space.city,
        description: props.space.description,
        scene_count: props.space.scene_count,
        addedAt: new Date().toISOString()
      })
    }
    
    localStorage.setItem('favorites', JSON.stringify(favorites))
  } catch (e) {
    console.error('Failed to toggle favorite:', e)
  }
}
</script>

<style scoped>
.space-panel {
  height: 100%;
  background: rgba(26, 35, 50, 0.95);
  backdrop-filter: blur(12px);
  border-left: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  flex-direction: column;
  box-shadow: -10px 0 40px rgba(0, 0, 0, 0.3);
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
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  flex-shrink: 0;
}

.panel-title {
  color: #fff;
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  background: linear-gradient(135deg, #0EA5E9, #4080ff);
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
  position: relative;
  width: 100%;
  height: 180px;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 16px;
}

.space-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.space-cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1E3A5F 0%, #0F1826 100%);
  color: rgba(255, 255, 255, 0.2);
}

.space-cover-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 6px 12px;
  background: rgba(14, 165, 233, 0.9);
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  color: #fff;
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
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  line-height: 1.6;
  margin: 0;
}

.scene-section {
  margin-bottom: 24px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 16px;
}

.section-title-icon {
  font-size: 18px;
}

.scene-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
}

.scene-item {
  display: flex;
  align-items: center;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
}

.scene-item:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(14, 165, 233, 0.3);
}

.scene-item-thumb {
  width: 60px;
  height: 40px;
  border-radius: 6px;
  overflow: hidden;
  flex-shrink: 0;
  margin-right: 12px;
}

.scene-item-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.scene-item-thumb-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.3);
}

.scene-item-info {
  flex: 1;
  min-width: 0;
}

.scene-item-name {
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.scene-item-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.scene-item-code {
  padding: 2px 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  font-size: 11px;
}

.scene-item-arrow {
  color: #0EA5E9;
  font-size: 16px;
  margin-left: 8px;
}

.scene-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: rgba(255, 255, 255, 0.3);
}

.scene-empty p {
  margin-top: 12px;
  font-size: 14px;
}

.topology-section {
  margin-bottom: 24px;
}

.topology-section :deep(.topology-container) {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
}

.panel-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.panel-actions :deep(.t-button) {
  border-radius: 10px;
}

.is-favorite {
  color: #F59E0B;
}

:deep(.t-button--variant-text) {
  color: rgba(255, 255, 255, 0.6);
}

:deep(.t-button--variant-text:hover) {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}

:deep(.t-button--theme-primary) {
  background: linear-gradient(135deg, #0EA5E9 0%, #4080ff 100%);
  border: none;
}

:deep(.t-button--theme-default) {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: rgba(255, 255, 255, 0.8);
}
</style>
