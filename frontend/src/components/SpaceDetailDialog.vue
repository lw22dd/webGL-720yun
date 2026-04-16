<template>
  <t-dialog
    :visible="visible"
    :header="space?.name"
    :width="720"
    :footer="false"
    @close="handleClose"
    class="space-detail-dialog"
  >
    <div class="dialog-content" v-if="space">
      <div class="space-cover">
        <img
          v-if="space.slug"
          :src="getCoverUrl(space.slug)"
          :alt="space.name"
        />
        <div v-else class="space-cover-placeholder">
          <t-icon name="image" size="48" />
        </div>
        <div class="space-cover-badge">
          {{ space.scene_count || 0 }} 个景点
        </div>
      </div>

      <div class="space-info">
        <div class="space-meta">
          <span class="space-location">
            <t-icon name="location" /> {{ space.province }} {{ space.city }}
          </span>
        </div>
        <p class="space-description">{{ space.description || '暂无描述' }}</p>
      </div>

      <div class="scene-section" v-if="space.scenes && space.scenes.length > 0">
        <h3 class="section-title">
          <t-icon name="view-module" />
          景点列表（{{ space.scenes.length }}个）
        </h3>
        <div class="scene-list">
          <div
            v-for="scene in space.scenes"
            :key="scene.id"
            class="scene-item"
            @click="handleSceneClick(scene)"
          >
            <div class="scene-thumb">
              <img
                v-if="scene.scene_code"
                :src="getPreviewUrl(scene.scene_code)"
                :alt="scene.title"
              />
              <div v-else class="scene-thumb-placeholder">
                <t-icon name="image" />
              </div>
            </div>
            <div class="scene-info">
              <div class="scene-name">{{ scene.title }}</div>
              <div class="scene-code">{{ scene.scene_code }}</div>
            </div>
            <t-button theme="primary" variant="text" size="small">
              进入 →
            </t-button>
          </div>
        </div>
      </div>

      <div class="dialog-actions">
        <t-button theme="primary" size="large" block @click="handleEnterPanorama">
          进入全景漫游
        </t-button>
      </div>
    </div>
  </t-dialog>
</template>

<script setup lang="ts">
import type { SpaceListItem } from '@/models/space.model'

interface SceneSimple {
  id: number
  title: string
  scene_code: string
  thumbnail_url: string
  view_count: number
  sort_order: number
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

const getCoverUrl = (slug: string) => {
  return `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'}/api/v1/res/covers/${slug}`
}

const getPreviewUrl = (sceneCode: string) => {
  return `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'}/api/v1/res/previews/${sceneCode}`
}

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
</script>

<style scoped>
.dialog-content {
  padding: 0;
}

.space-cover {
  position: relative;
  width: 100%;
  height: 240px;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 20px;
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
  color: rgba(255, 255, 255, 0.3);
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

.space-info {
  margin-bottom: 24px;
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
  font-size: 14px;
}

.space-description {
  color: rgba(255, 255, 255, 0.8);
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

.scene-thumb {
  width: 60px;
  height: 40px;
  border-radius: 6px;
  overflow: hidden;
  flex-shrink: 0;
  margin-right: 12px;
}

.scene-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.scene-thumb-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.3);
}

.scene-info {
  flex: 1;
  min-width: 0;
}

.scene-name {
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.scene-code {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  padding: 2px 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  display: inline-block;
}

.dialog-actions {
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}
</style>
