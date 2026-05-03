<template>
  <div class="location-picker">
    <div class="search-row">
      <t-input
        v-model="searchKeyword"
        placeholder="输入地点名称搜索"
        clearable
        @enter="handleSearch"
      >
        <template #suffix-icon>
          <t-icon name="search" />
        </template>
      </t-input>
      <t-button theme="primary" @click="handleSearch">搜索</t-button>
    </div>

    <div v-if="searchResults.length > 0" class="results-list">
      <div
        v-for="(item, index) in searchResults"
        :key="index"
        class="result-item"
        :class="{ active: selectedIndex === index }"
        @click="selectLocation(item, index)"
      >
        <div class="result-name">{{ item.name }}</div>
        <div class="result-address">{{ item.district }}{{ item.address }}</div>
      </div>
    </div>

    <div v-else-if="hasSearched" class="no-results">
      <t-icon name="info-circle" size="20" />
      <span>未找到匹配地点，请尝试其他关键词</span>
    </div>

    <div class="map-container" ref="mapContainerRef"></div>

    <div v-if="selectedLocation" class="selected-info">
      <div class="info-row">
        <span class="info-label">已选地点：</span>
        <span class="info-value">{{ selectedLocation.name }}</span>
      </div>
      <div class="info-row">
        <span class="info-label">经纬度：</span>
        <span class="info-value">{{ selectedLocation.lng.toFixed(6) }}, {{ selectedLocation.lat.toFixed(6) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { initAMap, setMapCenter, createMarker } from '@/utils/amap'

interface LocationResult {
  name: string
  address: string
  district: string
  lng: number
  lat: number
}

const emit = defineEmits<{
  (e: 'select', location: { name: string; lng: number; lat: number }): void
}>()

const searchKeyword = ref('')
const searchResults = ref<LocationResult[]>([])
const hasSearched = ref(false)
const selectedIndex = ref<number>(-1)
const selectedLocation = ref<LocationResult | null>(null)

const mapContainerRef = ref<HTMLElement>()
let mapInstance: any = null
let currentMarker: any = null

const handleSearch = async () => {
  if (!searchKeyword.value.trim()) {
    MessagePlugin.warning('请输入搜索关键词')
    return
  }

  try {
    const AMap = (window as any).AMap
    if (!AMap) {
      MessagePlugin.error('地图加载中，请稍后再试')
      return
    }

    const placeSearch = new AMap.PlaceSearch({
      pageSize: 10,
      pageIndex: 1
    })

    placeSearch.search(searchKeyword.value, (status: string, result: any) => {
      hasSearched.value = true
      if (status === 'complete' && result.poiList && result.poiList.pois) {
        searchResults.value = result.poiList.pois.map((poi: any) => ({
          name: poi.name,
          address: poi.address || '',
          district: poi.adname || poi.district || '',
          lng: parseFloat(poi.location.lng),
          lat: parseFloat(poi.location.lat)
        }))
      } else {
        searchResults.value = []
      }
    })
  } catch (error) {
    console.error('搜索失败:', error)
    MessagePlugin.error('搜索失败')
    searchResults.value = []
  }
}

const selectLocation = (item: LocationResult, index: number) => {
  selectedIndex.value = index
  selectedLocation.value = item

  if (mapInstance) {
    setMapCenter(mapInstance, item.lng, item.lat, 15)

    if (currentMarker) {
      currentMarker.setMap(null)
    }
    currentMarker = createMarker(mapInstance, [item.lng, item.lat], {
      title: item.name
    })
  }

  emit('select', {
    name: item.name,
    lng: item.lng,
    lat: item.lat
  })
}

onMounted(async () => {
  if (mapContainerRef.value) {
    try {
      mapInstance = await initAMap(mapContainerRef.value, {
        zoom: 4,
        center: [105.397428, 38.04623]
      })
    } catch (error) {
      console.error('地图初始化失败:', error)
    }
  }
})

onUnmounted(() => {
  if (currentMarker) {
    currentMarker.setMap(null)
    currentMarker = null
  }
})
</script>

<style scoped>
.location-picker {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.search-row {
  display: flex;
  gap: 8px;
}

.search-row :deep(.t-input) {
  flex: 1;
}

.results-list {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid var(--td-component-border);
  border-radius: 6px;
  background: var(--td-bg-color-container);
}

.result-item {
  padding: 10px 12px;
  cursor: pointer;
  border-bottom: 1px solid var(--td-component-border);
  transition: background 0.2s;
}

.result-item:last-child {
  border-bottom: none;
}

.result-item:hover,
.result-item.active {
  background: var(--td-bg-color-container-hover);
}

.result-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--td-text-color-primary);
  margin-bottom: 2px;
}

.result-address {
  font-size: 12px;
  color: var(--td-text-color-secondary);
}

.no-results {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 20px;
  color: var(--td-text-color-secondary);
  font-size: 13px;
}

.map-container {
  width: 100%;
  height: 280px;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid var(--td-component-border);
}

.selected-info {
  padding: 10px 12px;
  background: var(--td-brand-color-light);
  border-radius: 6px;
  border: 1px solid var(--td-brand-color-focus);
}

.info-row {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  line-height: 1.6;
}

.info-label {
  color: var(--td-text-color-secondary);
  white-space: nowrap;
}

.info-value {
  color: var(--td-text-color-primary);
  font-weight: 500;
}
</style>
