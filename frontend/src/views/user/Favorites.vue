<template>
  <div class="favorites-page">
    <Header />
    
    <div class="page-container">
      <div class="page-header">
        <h1 class="page-title">我的收藏</h1>
        <div class="page-stats">共 {{ favorites.length }} 个收藏</div>
      </div>

      <div class="filter-bar">
        <t-radio-group v-model="filterType" variant="default-filled">
          <t-radio-button value="all">全部</t-radio-button>
          <t-radio-button value="space">景区</t-radio-button>
          <t-radio-button value="scene">景点</t-radio-button>
        </t-radio-group>
      </div>

      <div class="favorites-grid" v-if="filteredFavorites.length > 0">
        <div 
          v-for="item in filteredFavorites" 
          :key="item.id" 
          class="favorite-card"
          @click="handleItemClick(item)"
        >
          <div class="card-cover">
            <img v-if="item.slug" :src="getCoverUrl(item.slug)" :alt="item.name" />
            <div v-else class="card-cover-placeholder">
              <t-icon name="image" size="48" />
            </div>
            <div class="card-badge">{{ item.type === 'space' ? '景区' : '景点' }}</div>
            <button class="card-remove" @click.stop="removeFavorite(item)">
              <t-icon name="close" />
            </button>
          </div>
          <div class="card-body">
            <h3 class="card-name">{{ item.name }}</h3>
            <p class="card-location">
              <t-icon name="location" />
              {{ item.province }} {{ item.city }}
            </p>
            <p class="card-desc">{{ item.description || '暂无描述' }}</p>
            <div class="card-footer">
              <span class="card-time">收藏于 {{ formatDate(item.addedAt) }}</span>
              <span class="card-enter">进入漫游 →</span>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <t-icon name="star" size="64" />
        <p>暂无收藏内容</p>
        <t-button theme="primary" @click="router.push('/')">去探索</t-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import Header from '@/components/Header.vue'

interface FavoriteItem {
  id: number
  name: string
  type: 'space' | 'scene'
  cover_url?: string
  slug?: string
  province?: string
  city?: string
  description?: string
  scene_count?: number
  scene_code?: string
  addedAt: string
}

const router = useRouter()
const favorites = ref<FavoriteItem[]>([])
const filterType = ref<'all' | 'space' | 'scene'>('all')

const filteredFavorites = computed(() => {
  if (filterType.value === 'all') return favorites.value
  return favorites.value.filter(item => item.type === filterType.value)
})

const loadFavorites = () => {
  try {
    const data = localStorage.getItem('favorites')
    if (data) {
      favorites.value = JSON.parse(data)
    }
  } catch (e) {
    favorites.value = []
  }
}

const removeFavorite = (item: FavoriteItem) => {
  const index = favorites.value.findIndex(f => f.id === item.id && f.type === item.type)
  if (index > -1) {
    favorites.value.splice(index, 1)
    localStorage.setItem('favorites', JSON.stringify(favorites.value))
    MessagePlugin.success('已取消收藏')
  }
}

const handleItemClick = (item: FavoriteItem) => {
  if (item.type === 'space') {
    router.push(`/?space=${item.id}`)
  } else if (item.scene_code) {
    router.push(`/panorama?scene=${item.scene_code}`)
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}月${date.getDate()}日`
}

const getCoverUrl = (slug?: string) => {
  if (!slug) return ''
  return `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'}/api/v1/res/covers/${slug}`
}

onMounted(() => {
  loadFavorites()
})
</script>

<style scoped>
.favorites-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0F1826 0%, #162032 50%, #1A2744 100%);
}

.page-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  color: #fff;
  font-size: 28px;
  font-weight: 600;
  margin: 0;
}

.page-stats {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
}

.filter-bar {
  margin-bottom: 24px;
}

.filter-bar :deep(.t-radio-group) {
  background: rgba(255, 255, 255, 0.1);
}

.filter-bar :deep(.t-radio-button) {
  color: rgba(255, 255, 255, 0.8);
}

.filter-bar :deep(.t-radio-button.t-is-checked) {
  background: #0EA5E9;
  color: #fff;
}

.favorites-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

.favorite-card {
  background: rgba(26, 35, 50, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
}

.favorite-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.3);
  border-color: rgba(14, 165, 233, 0.3);
}

.card-cover {
  position: relative;
  height: 160px;
  overflow: hidden;
}

.card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.favorite-card:hover .card-cover img {
  transform: scale(1.05);
}

.card-cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1E3A5F 0%, #0F1826 100%);
  color: rgba(255, 255, 255, 0.2);
}

.card-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  padding: 4px 10px;
  background: rgba(14, 165, 233, 0.9);
  border-radius: 12px;
  font-size: 12px;
  color: #fff;
}

.card-remove {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 28px;
  height: 28px;
  border: none;
  background: rgba(0, 0, 0, 0.5);
  color: #fff;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: all 0.2s;
}

.favorite-card:hover .card-remove {
  opacity: 1;
}

.card-remove:hover {
  background: #EF4444;
}

.card-body {
  padding: 16px;
}

.card-name {
  color: #fff;
  font-size: 16px;
  font-weight: 500;
  margin: 0 0 8px;
}

.card-location {
  display: flex;
  align-items: center;
  gap: 4px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 13px;
  margin: 0 0 8px;
}

.card-desc {
  color: rgba(255, 255, 255, 0.7);
  font-size: 13px;
  line-height: 1.5;
  margin: 0 0 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.card-time {
  color: rgba(255, 255, 255, 0.5);
  font-size: 12px;
}

.card-enter {
  color: #0EA5E9;
  font-size: 13px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: rgba(255, 255, 255, 0.4);
}

.empty-state p {
  margin: 16px 0;
  font-size: 16px;
}
</style>
