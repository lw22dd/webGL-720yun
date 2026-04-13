<template>
  <div class="page-wrapper">
    <Header />
    <main class="main-content">
      <div class="map-section" v-show="viewMode === 'map'">
        <ChinaMap 
          ref="chinaMapRef" 
          :view-mode="viewMode"
          @space-click="handleSpaceClick"
          @view-mode-change="handleViewModeChange"
        />
      </div>
      
      <div class="card-section" v-show="viewMode === 'card'">
        <div class="card-view-header">
          <h2 class="card-view-title">全部景区</h2>
          <div class="card-view-toggle">
            <t-button 
              :variant="viewMode === 'map' ? 'base' : 'outline'" 
              size="small"
              @click="viewMode = 'map'"
            >
              🗺️ 地图
            </t-button>
            <t-button 
              :variant="viewMode === 'card' ? 'base' : 'outline'" 
              size="small"
              @click="viewMode = 'card'"
            >
              📋 卡片
            </t-button>
          </div>
        </div>
        
        <div class="card-view-grid" v-loading="loading">
          <div 
            v-for="space in spaceList" 
            :key="space.id" 
            class="space-card"
            @click="handleCardClick(space)"
          >
            <div class="space-card-cover">
              <img 
                v-if="space.cover_url" 
                :src="space.cover_url" 
                :alt="space.name"
              />
              <div v-else class="space-card-cover-placeholder">
                <t-icon name="image" size="48" />
              </div>
              <div class="space-card-badge">{{ space.scene_count || 0 }} 个景点</div>
            </div>
            <div class="space-card-body">
              <h3 class="space-card-name">{{ space.name }}</h3>
              <p class="space-card-location">
                <t-icon name="location" />
                {{ space.province }} {{ space.city }}
              </p>
              <p class="space-card-desc">{{ space.description || '暂无描述' }}</p>
              <div class="space-card-footer">
                <span class="space-card-scenes">
                  <t-icon name="view-module" />
                  {{ space.scene_count || 0 }} 个场景
                </span>
                <span class="space-card-enter">进入漫游 →</span>
              </div>
            </div>
          </div>
        </div>
        
        <div v-if="!loading && spaceList.length === 0" class="card-empty">
          <t-icon name="map-location" size="64" />
          <p>暂无景区数据</p>
        </div>
      </div>
      
      <SpaceDetailPanel
        :visible="panelVisible"
        :space="selectedSpace"
        @close="handlePanelClose"
        @scene-click="handleSceneClick"
      />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Header from '@/components/Header.vue'
import ChinaMap from '@/components/ChinaMap.vue'
import SpaceDetailPanel from '@/components/SpaceDetailPanel.vue'
import SpaceApi from '@/apis/spaceApi'
import type { SpaceListItem } from '@/models/SpaceModel'

interface SpaceWithScenes extends SpaceListItem {
  scenes?: any[]
}

const router = useRouter()
const chinaMapRef = ref<InstanceType<typeof ChinaMap> | null>(null)
const panelVisible = ref(false)
const selectedSpace = ref<SpaceWithScenes | null>(null)
const viewMode = ref<'map' | 'card'>('map')
const spaceList = ref<SpaceListItem[]>([])
const loading = ref(false)

const loadSpaceList = async () => {
  try {
    loading.value = true
    const result = await SpaceApi.getSpaceList({ page: 1, page_size: 100 })
    if (result.code === 200 && result.data) {
      spaceList.value = result.data.spaces
    }
  } catch (error) {
    console.error('Failed to load space list:', error)
  } finally {
    loading.value = false
  }
}

const handleSpaceClick = (space: SpaceWithScenes) => {
  selectedSpace.value = space
  panelVisible.value = true
}

const handlePanelClose = () => {
  panelVisible.value = false
}

const handleSceneClick = (scene: any) => {
  router.push(`/panorama?scene=${scene.scene_code}`)
}

const handleViewModeChange = (mode: 'map' | 'card') => {
  viewMode.value = mode
}

const handleCardClick = async (space: SpaceListItem) => {
  const detailResult = await SpaceApi.getSpaceDetail(space.id)
  if (detailResult.code === 200 && detailResult.data) {
    selectedSpace.value = {
      ...space,
      scenes: detailResult.data.scenes || []
    }
  } else {
    selectedSpace.value = space
  }
  panelVisible.value = true
}

onMounted(() => {
  loadSpaceList()
})
</script>

<style scoped>
.page-wrapper {
  height: 100vh;
  background: #0F1826;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.main-content {
  width: 100%;
  flex: 1;
  display: flex;
  overflow: hidden;
}

.map-section {
  position: relative;
  flex: 1;
  height: 100%;
  min-width: 0;
}

.card-section {
  flex: 1;
  height: 100%;
  overflow-y: auto;
  background: linear-gradient(135deg, #0F1826 0%, #162032 50%, #1A2744 100%);
  padding: 24px;
}

.card-view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 0 8px;
}

.card-view-title {
  font-size: 24px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.card-view-toggle {
  display: flex;
  gap: 8px;
}

.card-view-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.space-card {
  background: rgba(26, 35, 50, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
}

.space-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.3);
  border-color: rgba(14, 165, 233, 0.3);
}

.space-card-cover {
  position: relative;
  height: 180px;
  overflow: hidden;
}

.space-card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.space-card:hover .space-card-cover img {
  transform: scale(1.05);
}

.space-card-cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1E3A5F 0%, #0F1826 100%);
  color: rgba(255, 255, 255, 0.2);
}

.space-card-badge {
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

.space-card-body {
  padding: 20px;
}

.space-card-name {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 8px;
}

.space-card-location {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
  margin: 0 0 12px;
}

.space-card-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.5;
  margin: 0 0 16px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.space-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.space-card-scenes {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

.space-card-enter {
  font-size: 13px;
  color: #0EA5E9;
  font-weight: 500;
}

.card-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: rgba(255, 255, 255, 0.4);
}

.card-empty p {
  margin-top: 16px;
  font-size: 16px;
}

:deep(.t-button) {
  background: rgba(26, 35, 50, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
}

:deep(.t-button.t-button--variant-base) {
  background: #0EA5E9;
  border-color: #0EA5E9;
  color: #fff;
}

:deep(.t-loading) {
  color: #0EA5E9;
}
</style>
