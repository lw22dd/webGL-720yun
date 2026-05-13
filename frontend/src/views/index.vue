<template>
  <div class="page-wrapper">
    <!-- 动态背景层 -->
    <div class="dynamic-background">
      <div class="bg-image"></div>
      <div class="bg-overlay"></div>
      
      <!-- 星星闪烁层 -->
      <div class="stars-container">
        <div 
          v-for="star in stars" 
          :key="star.id" 
          class="star"
          :style="{
            left: star.left,
            top: star.top,
            width: star.size + 'px',
            height: star.size + 'px',
            animationDelay: star.delay + 's',
            animationDuration: star.duration + 's'
          }"
        ></div>
      </div>

      <div class="light-spots">
        <div class="light-spot spot-1"></div>
        <div class="light-spot spot-2"></div>
        <div class="light-spot spot-3"></div>
      </div>
    </div>

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
                v-if="space.slug" 
                :src="getCoverUrl(space.slug)" 
                :alt="space.name"
                @error="console.error('封面加载失败:', space.name, getCoverUrl(space.slug))"
              />
              <img 
                v-else-if="space.cover_url" 
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

      <SpaceDetailDialog
        :visible="dialogVisible"
        :space="selectedSpace"
        @close="handleDialogClose"
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
import SpaceDetailDialog from '@/components/SpaceDetailDialog.vue'
import SpaceApi from '@/apis/space.api'
import type { SpaceListItem, SpaceWithScenes } from '@/models/space.model'

const router = useRouter()
const chinaMapRef = ref<InstanceType<typeof ChinaMap> | null>(null)
const panelVisible = ref(false)
const dialogVisible = ref(false)
const selectedSpace = ref<SpaceWithScenes | null>(null)
const viewMode = ref<'map' | 'card'>('map')
const spaceList = ref<SpaceListItem[]>([])
const loading = ref(false)
const stars = ref<any[]>([])

const initStars = () => {
  const count = 150 // 增加星星数量
  const result = []
  for (let i = 0; i < count; i++) {
    result.push({
      id: i,
      left: `${Math.random() * 100}%`,
      top: `${Math.random() * 100}%`,
      size: Math.random() * 1.5 + 0.5, // 细小的星星
      delay: Math.random() * 5,
      duration: Math.random() * 3 + 10
    })
  }
  stars.value = result
}

const getCoverUrl = (slug: string) => {
  return `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'}/api/v1/res/covers/${slug}`
}

const loadSpaceList = async () => {
  try {
    loading.value = true
    const result = await SpaceApi.getSpaceList({ page: 1, page_size: 100 })
    if (result.code === 200 && result.data) {
      console.log("全景景区列表:", result.data.spaces)
      // 打印每个空间的 slug 和生成的封面 URL
      result.data.spaces.forEach((space: any) => {
        console.log(`Space ${space.name}: slug=${space.slug}, coverUrl=${getCoverUrl(space.slug)}`)
      })
      spaceList.value = result.data.spaces
    }
  } catch (error) {
    console.error('Failed to load space list:', error)
  } finally {
    loading.value = false
  }
}

const handleSpaceClick = (space: SpaceWithScenes) => {
  router.push({ name: 'spaceDetail', params: { id: space.id } })
}

const handlePanelClose = () => {
  panelVisible.value = false
}

const handleDialogClose = () => {
  dialogVisible.value = false
}

const handleSceneClick = (scene: any) => {
  panelVisible.value = false
  dialogVisible.value = false
  sessionStorage.setItem('panoramaFrom', router.currentRoute.value.fullPath)
  router.push(`/panorama?scene=${scene.scene_code}`)
}

const handleViewModeChange = (mode: 'map' | 'card') => {
  viewMode.value = mode
}

const handleCardClick = async (space: SpaceListItem) => {
  router.push({ name: 'spaceDetail', params: { id: space.id } })
}

onMounted(() => {
  initStars()
  loadSpaceList()
})
</script>

<style scoped>
.page-wrapper {
  position: relative;
  height: 100vh;
  background: #0a111d; /* 稍微深一点的底色 */
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* 动态背景样式 */
.dynamic-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 0;
  overflow: hidden;
}

.bg-image {
  position: absolute;
  top: -10%;
  left: -10%;
  width: 120%;
  height: 120%;
  background-image: url('/background.png'); /* 建议将 image.png 移动到 public/background.png */
  background-size: cover;
  background-position: center;
  filter: brightness(0.6) saturate(1.2) blur(2px);
  animation: bg-pan 60s linear infinite alternate;
}

.bg-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at center, transparent 0%, rgba(10, 17, 29, 0.8) 100%),
              linear-gradient(to bottom, rgba(10, 17, 29, 0.4) 0%, rgba(10, 17, 29, 0.9) 100%);
}

/* 星星样式 */
.stars-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.star {
  position: absolute;
  background: #fff;
  border-radius: 50%;
  box-shadow: 0 0 4px #fff, 0 0 8px rgba(255, 255, 255, 0.3);
  opacity: 0;
  animation: twinkle linear infinite;
}

@keyframes twinkle {
  0% { opacity: 0; transform: scale(0.5); }
  50% { opacity: 0.8; transform: scale(1.2); }
  100% { opacity: 0; transform: scale(0.5); }
}

.light-spots {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}

.light-spot {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  mix-blend-mode: screen;
  animation: float 20s ease-in-out infinite alternate;
}

.spot-1 {
  width: 400px;
  height: 400px;
  background: rgba(14, 165, 233, 0.3);
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.spot-2 {
  width: 500px;
  height: 500px;
  background: rgba(139, 92, 246, 0.2);
  bottom: 10%;
  right: 10%;
  animation-delay: -5s;
}

.spot-3 {
  width: 300px;
  height: 300px;
  background: rgba(16, 185, 129, 0.15);
  top: 40%;
  right: 20%;
  animation-delay: -10s;
}

@keyframes bg-pan {
  from { transform: scale(1) translate(0, 0); }
  to { transform: scale(1.1) translate(-2%, -2%); }
}

@keyframes float {
  0% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(5%, 5%) scale(1.1); }
  100% { transform: translate(-5%, 2%) scale(0.9); }
}

.main-content {
  position: relative;
  width: 100%;
  flex: 1;
  display: flex;
  overflow: hidden;
  z-index: 1; /* 确保在背景之上 */
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
  /* 使用玻璃拟态背景 */
  background: rgba(10, 17, 29, 0.6);
  backdrop-filter: blur(10px);
  border-left: 1px solid rgba(255, 255, 255, 0.05);
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
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
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
  background: rgba(30, 41, 59, 0.5);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.165, 0.84, 0.44, 1);
}

.space-card:hover {
  transform: translateY(-8px);
  background: rgba(30, 41, 59, 0.8);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
  border-color: rgba(14, 165, 233, 0.4);
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
  transition: transform 0.5s ease;
}

.space-card:hover .space-card-cover img {
  transform: scale(1.1);
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
  backdrop-filter: blur(4px);
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  color: #fff;
  box-shadow: 0 4px 12px rgba(14, 165, 233, 0.3);
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
  border-top: 1px solid rgba(255, 255, 255, 0.08);
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
  font-weight: 600;
  transition: all 0.3s;
}

.space-card:hover .space-card-enter {
  color: #38bdf8;
  transform: translateX(4px);
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
  background: rgba(30, 41, 59, 0.6);
  backdrop-filter: blur(4px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
}

:deep(.t-button:hover) {
  background: rgba(30, 41, 59, 0.8);
  border-color: rgba(14, 165, 233, 0.5);
  color: #fff;
}

:deep(.t-button.t-button--variant-base) {
  background: #0EA5E9;
  border-color: #0EA5E9;
  color: #fff;
  box-shadow: 0 4px 12px rgba(14, 165, 233, 0.3);
}

:deep(.t-loading) {
  color: #0EA5E9;
}
</style>
