<template>
  <div class="china-map-container">
    <div ref="mapContainer" class="map-wrapper" v-loading="loading"></div>
    <div class="search-box">
      <t-input v-model="searchKeyword" placeholder="搜索景点或区域..." @enter="handleSearch" @change="handleSearchChange"
        clearable>
        <template #prefix-icon>
          <SearchIcon />
        </template>
      </t-input>

      <div v-if="showSearchResults" class="search-results">
        <div v-if="searchHistory.length > 0 && !searchKeyword.trim()" class="search-history">
          <div class="search-history-header">
            <span class="search-history-title">搜索历史</span>
            <button class="search-history-clear" @click="clearSearchHistory">清除</button>
          </div>
          <div class="search-history-tags">
            <span v-for="(term, index) in searchHistory.slice(0, 8)" :key="index" class="search-history-tag"
              @click="onHistoryTagClick(term)">
              {{ term }}
            </span>
          </div>
        </div>

        <div v-else-if="searchKeyword.trim() && searchResults.length === 0" class="search-empty">
          未找到「{{ searchKeyword }}」相关结果
        </div>

        <div v-else-if="searchResults.length > 0" class="search-result-list">
          <div v-for="result in searchResults" :key="result.id" class="search-result-item"
            @click="handleResultClick(result)">
            <span class="result-name">
              <template v-if="result.type === 'space'">🏛</template>
              <template v-else>🎥</template>
              {{ result.name }}
            </span>
            <span class="result-location">{{ result.location }}</span>
          </div>
        </div>
      </div>



      <div class="view-toggle">
        <t-button :variant="viewMode === 'map' ? 'base' : 'outline'" size="small"
          @click="$emit('viewModeChange', 'map')">
          🗺️
        </t-button>
        <t-button :variant="viewMode === 'card' ? 'base' : 'outline'" size="small"
          @click="$emit('viewModeChange', 'card')">
          📋
        </t-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { SearchIcon } from 'tdesign-icons-vue-next'
import * as echarts from 'echarts'
import { getChinaGeoJSON } from '@/apis/geo.api'
import SpaceApi from '@/apis/space.api'
import type { SpaceListItem, SpaceWithScenes } from '@/models/space.model'

const props = defineProps<{
  viewMode?: 'map' | 'card'
}>()

const emit = defineEmits<{
  (e: 'spaceClick', space: SpaceWithScenes): void
  (e: 'viewModeChange', mode: 'map' | 'card'): void
}>()

const mapContainer = ref<HTMLElement | null>(null)
const searchKeyword = ref('')
const searchResults = ref<any[]>([])
const showSearchResults = ref(false)
const loading = ref(true)
const spaceList = ref<SpaceListItem[]>([])

const searchHistory = ref<string[]>([])

let chartInstance: echarts.ECharts | null = null
let resizeObserver: ResizeObserver | null = null

const loadSearchHistory = () => {
  try {
    const history = localStorage.getItem('searchHistory')
    if (history) {
      searchHistory.value = JSON.parse(history)
    }
  } catch (e) {
    searchHistory.value = []
  }
}

const saveSearchHistory = (keyword: string) => {
  const newHistory = [keyword, ...searchHistory.value.filter(h => h !== keyword)].slice(0, 10)
  searchHistory.value = newHistory
  localStorage.setItem('searchHistory', JSON.stringify(newHistory))
}

const clearSearchHistory = () => {
  searchHistory.value = []
  localStorage.setItem('searchHistory', '[]')
}

const loadSpaceList = async () => {
  try {
    loading.value = true
    const result = await SpaceApi.getSpaceList({ page: 1, page_size: 100 })
    if (result.code === 200 && result.data) {
      console.log(result.data)
      spaceList.value = result.data.spaces
      updateMapData()
    }
  } catch (error) {
    console.error('Failed to load space list:', error)
  } finally {
    loading.value = false
  }
}

const initMap = async () => {
  if (!mapContainer.value) return

  try {
    const chinaGeoJSON = await getChinaGeoJSON()
    echarts.registerMap('china', chinaGeoJSON)

    chartInstance = echarts.init(mapContainer.value, 'dark', {
      renderer: 'canvas'
    })

    const option: echarts.EChartsOption = {
      backgroundColor: 'transparent',
      tooltip: {
        trigger: 'item',
        backgroundColor: 'rgba(26,35,50,0.92)',
        borderColor: 'rgba(255,255,255,0.1)',
        borderWidth: 1,
        padding: [12, 16],
        textStyle: {
          color: '#fff',
          fontSize: 13,
          fontFamily: 'Noto Sans SC, sans-serif'
        },
        extraCssText: 'backdrop-filter:blur(12px);border-radius:10px;box-shadow:0 8px 32px rgba(0,0,0,0.3);',
        formatter: (params: any) => {
          if (params.seriesType === 'effectScatter') {
            const space = spaceList.value.find(s => s.name === params.name)
            if (space) {
              return `<b style="font-size:14px;">${params.name}</b><br/>
                      <span style="color:rgba(255,255,255,0.6);">📍 ${space.province} ${space.city}</span><br/>
                      <span style="color:rgba(255,255,255,0.6);">🎥 ${space.scene_count} 个景点</span>`
            }
          }
          return params.name
        }
      },
      geo: {
        map: 'china',
        roam: true,
        zoom: 1.2,
        center: [104.5, 36],
        scaleLimit: { min: 0.8, max: 6 },
        label: { show: false },
        itemStyle: {
          areaColor: {
            type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: '#162032' },
              { offset: 1, color: '#0F1826' }
            ]
          },
          borderColor: '#1E3A5F',
          borderWidth: 0.6,
          shadowColor: 'rgba(14,165,233,0.15)',
          shadowBlur: 15
        },
        emphasis: {
          itemStyle: {
            areaColor: '#1E4A6E',
            borderColor: '#0EA5E9',
            borderWidth: 1.2,
            shadowColor: 'rgba(14,165,233,0.3)',
            shadowBlur: 20
          },
          label: { show: false }
        },
        select: { disabled: true }
      },
      series: [
        {
          name: '景区',
          type: 'effectScatter',
          coordinateSystem: 'geo',
          data: [],
          symbolSize: 16,
          showEffectOn: 'render',
          rippleEffect: { brushType: 'stroke', scale: 3.5, period: 4 },
          label: {
            show: true,
            position: 'top',
            formatter: '{b}',
            fontSize: 12,
            fontWeight: 500,
            fontFamily: 'Noto Sans SC, sans-serif',
            color: '#fff',
            textBorderColor: 'rgba(0,0,0,0.7)',
            textBorderWidth: 2.5,
            distance: 10
          },
          itemStyle: {
            shadowBlur: 12,
            shadowColor: 'rgba(0,0,0,0.4)'
          },
          zlevel: 1
        }
      ]
    }

    chartInstance.setOption(option)

    chartInstance.on('click', (params: any) => {
      if (params.seriesType === 'effectScatter') {
        const space = spaceList.value.find(s => s.name === params.name)
        if (space) {
          SpaceApi.getSpaceDetail(space.id).then(detailResult => {
            if (detailResult.code === 200 && detailResult.data) {
              emit('spaceClick', {
                ...space,
                scenes: detailResult.data.scenes || []
              })
            } else {
              emit('spaceClick', space)
            }
          })
        }
      }
    })

    const handleResize = () => {
      chartInstance?.resize()
    }
    window.addEventListener('resize', handleResize)

    resizeObserver = new ResizeObserver(() => {
      chartInstance?.resize()
    })
    if (mapContainer.value) {
      resizeObserver.observe(mapContainer.value)
    }

    // 如果初始化时已经有数据了（例如 loadSpaceList 先完成了），主动触发一次更新
    if (spaceList.value.length > 0) {
      updateMapData()
    }

  } catch (error) {
    console.error('Failed to load map:', error)
  }
}

const updateMapData = () => {
  if (!chartInstance) return

  const colors = ['#2563EB', '#DC2626', '#0EA5E9', '#8B5CF6', '#10B981', '#F59E0B', '#059669', '#D97706']

  const scatterData = spaceList.value.map((space, index) => ({
    name: space.name,
    value: [space.longitude, space.latitude, space.scene_count || 1],
    itemStyle: { color: colors[index % colors.length] },
    space
  }))

  chartInstance.setOption({
    series: [{
      data: scatterData
    }]
  })
}

const handleSearchChange = (value: string) => {
  if (!value.trim()) {
    searchResults.value = []
    showSearchResults.value = searchHistory.value.length > 0
    return
  }

  showSearchResults.value = true
  const keyword = value.toLowerCase()
  const results: any[] = []

  spaceList.value.forEach(space => {
    if (space.name.toLowerCase().includes(keyword) ||
      space.province.toLowerCase().includes(keyword) ||
      space.city.toLowerCase().includes(keyword)) {
      results.push({
        id: space.id,
        name: space.name,
        location: `${space.province} ${space.city} · ${space.scene_count}个景点`,
        type: 'space',
        data: space
      })
    }
  })

  searchResults.value = results.slice(0, 8)
}

const handleSearch = () => {
  if (searchResults.value.length > 0) {
    handleResultClick(searchResults.value[0])
  }
}

const handleResultClick = async (result: any) => {
  searchKeyword.value = result.name
  showSearchResults.value = false

  saveSearchHistory(result.name)

  if (result.type === 'space') {
    const detailResult = await SpaceApi.getSpaceDetail(result.id)
    if (detailResult.code === 200 && detailResult.data) {
      emit('spaceClick', {
        ...result.data,
        scenes: detailResult.data.scenes || []
      })
    } else {
      emit('spaceClick', result.data)
    }
  }
}

const onHistoryTagClick = (term: string) => {
  searchKeyword.value = term
  handleSearchChange(term)
}

watch(() => props.viewMode, () => {
  setTimeout(() => {
    chartInstance?.resize()
  }, 100)
})

onMounted(async () => {
  loadSearchHistory()
  await initMap()
  await loadSpaceList()
})

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
})

defineExpose({
  focusSpace: (spaceId: number) => {
    const space = spaceList.value.find(s => s.id === spaceId)
    if (space && chartInstance) {
      chartInstance.setOption({
        geo: { center: [space.longitude, space.latitude], zoom: 4 }
      })
    }
  }
})
</script>

<style scoped>
.china-map-container {
  position: relative;
  width: 100%;
  height: 100%;
  background: transparent; /* 改为透明以显示父级动态背景 */
}

.map-wrapper {
  width: 100%;
  height: 100%;
}

.search-box {
  position: absolute;
  top: 30px;
  left: 30px;
  width: 380px;
  z-index: 100;
  filter: drop-shadow(0 8px 24px rgba(0, 0, 0, 0.4));
}

.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: rgba(30, 41, 59, 0.85);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  margin-top: 12px;
  max-height: 400px;
  overflow-y: auto;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5);
}

.search-history {
  padding: 16px;
}

.search-history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.search-history-title {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.search-history-clear {
  background: none;
  border: none;
  color: #0EA5E9;
  font-size: 12px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s;
}

.search-history-clear:hover {
  background: rgba(14, 165, 233, 0.1);
}

.search-history-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.search-history-tag {
  display: inline-block;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.2s;
}

.search-history-tag:hover {
  background: rgba(14, 165, 233, 0.2);
  color: #fff;
}

.search-empty {
  padding: 24px;
  text-align: center;
  color: rgba(255, 255, 255, 0.5);
  font-size: 13px;
}

.search-result-list {
  padding: 8px 0;
}

.search-result-item {
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.2s;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.search-result-item:last-child {
  border-bottom: none;
}

.search-result-item:hover {
  background: rgba(14, 165, 233, 0.1);
}

.result-name {
  display: block;
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.result-location {
  display: block;
  color: rgba(255, 255, 255, 0.4);
  font-size: 12px;
}

.view-toggle {
  position: absolute;
  top: 0;
  right: -56px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

:deep(.t-input) {
  background: rgba(30, 41, 59, 0.7) !important;
  backdrop-filter: blur(8px) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 12px !important;
  padding: 8px 16px !important;
  transition: all 0.3s !important;
}

:deep(.t-input:hover) {
  border-color: rgba(14, 165, 233, 0.4) !important;
  background: rgba(30, 41, 59, 0.85) !important;
}

:deep(.t-input.t-is-focused) {
  border-color: #0EA5E9 !important;
  box-shadow: 0 0 0 3px rgba(14, 165, 233, 0.2) !important;
}

:deep(.t-input__inner) {
  color: #fff !important;
  font-size: 14px !important;
}

:deep(.t-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.4) !important;
}

:deep(.t-button) {
  background: rgba(30, 41, 59, 0.7);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  width: 44px;
  height: 44px;
  border-radius: 12px;
}

:deep(.t-button.t-button--variant-base) {
  background: #0EA5E9;
  border-color: #0EA5E9;
  box-shadow: 0 4px 12px rgba(14, 165, 233, 0.3);
}
</style>
