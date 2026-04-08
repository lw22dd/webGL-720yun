<template>
  <div class="china-map-container">
    <div ref="mapContainer" class="map-wrapper"></div>
    <div class="search-box">
      <t-input
        v-model="searchKeyword"
        placeholder="搜索景点或区域..."
        @enter="handleSearch"
        @change="handleSearchChange"
        clearable
      >
        <template #prefix-icon>
          <SearchIcon />
        </template>
      </t-input>
      <div v-if="searchResults.length > 0" class="search-results">
        <div
          v-for="result in searchResults"
          :key="result.id"
          class="search-result-item"
          @click="handleResultClick(result)"
        >
          <span class="result-name">{{ result.name }}</span>
          <span class="result-location">{{ result.province }} {{ result.city }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { SearchIcon } from 'tdesign-icons-vue-next'
import * as echarts from 'echarts'
import { getChinaGeoJSON } from '@/api/geoApi'
import { mockSpaces, type MockSpace } from '@/utils/mockData'

const emit = defineEmits<{
  (e: 'spaceClick', space: MockSpace): void
}>()

const mapContainer = ref<HTMLElement | null>(null)
const searchKeyword = ref('')
const searchResults = ref<MockSpace[]>([])

let chartInstance: echarts.ECharts | null = null

const initMap = async () => {
  if (!mapContainer.value) return

  try {
    const chinaGeoJSON = await getChinaGeoJSON()
    echarts.registerMap('china', chinaGeoJSON)

    chartInstance = echarts.init(mapContainer.value, 'dark', {
      renderer: 'canvas'
    })

    const scatterData = mockSpaces.map(space => ({
      name: space.name,
      value: [space.longitude, space.latitude],
      space
    }))

    const option: echarts.EChartsOption = {
      backgroundColor: 'transparent',
      geo: {
        map: 'china',
        roam: false,
        zoom: 1.2,
        center: [105, 36],
        label: {
          show: false
        },
        itemStyle: {
          areaColor: 'rgba(48, 74, 216, 0.1)',
          borderColor: '#304AD8',
          borderWidth: 1,
          shadowColor: 'rgba(48, 74, 216, 0.3)',
          shadowBlur: 10
        },
        emphasis: {
          itemStyle: {
            areaColor: 'rgba(48, 74, 216, 0.2)'
          }
        }
      },
      series: [
        {
          type: 'scatter',
          coordinateSystem: 'geo',
          data: scatterData,
          symbolSize: 20,
          symbol: 'circle',
          itemStyle: {
            color: '#304AD8',
            shadowBlur: 20,
            shadowColor: '#304AD8'
          },
          emphasis: {
            scale: 1.5,
            itemStyle: {
              shadowBlur: 30,
              shadowColor: '#304AD8'
            }
          },
          label: {
            show: true,
            formatter: '{b}',
            position: 'top',
            color: '#fff',
            fontSize: 12,
            fontWeight: 'bold',
            backgroundColor: 'rgba(48, 74, 216, 0.8)',
            padding: [4, 8],
            borderRadius: 4,
            borderColor: 'rgba(48, 74, 216, 0.5)',
            borderWidth: 1
          }
        }
      ]
    }

    chartInstance.setOption(option)

    chartInstance.on('click', (params: any) => {
      if (params.data && params.data.space) {
        emit('spaceClick', params.data.space)
      }
    })

    const handleResize = () => {
      chartInstance?.resize()
    }
    window.addEventListener('resize', handleResize)

  } catch (error) {
    console.error('Failed to load map:', error)
  }
}

const handleSearchChange = (value: string) => {
  if (!value.trim()) {
    searchResults.value = []
    return
  }

  const keyword = value.toLowerCase()
  searchResults.value = mockSpaces.filter(
    space =>
      space.name.toLowerCase().includes(keyword) ||
      space.province.toLowerCase().includes(keyword) ||
      space.city.toLowerCase().includes(keyword)
  )
}

const handleSearch = () => {
  if (searchResults.value.length > 0) {
    handleResultClick(searchResults.value[0])
  }
}

const handleResultClick = (space: MockSpace) => {
  searchKeyword.value = space.name
  searchResults.value = []
  emit('spaceClick', space)
}

onMounted(() => {
  initMap()
})

onUnmounted(() => {
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
})
</script>

<style scoped>
.china-map-container {
  position: relative;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #1a2a6c 0%, #304AD8 50%, #1a2a6c 100%);
}

.map-wrapper {
  width: 100%;
  height: 100%;
}

.search-box {
  position: absolute;
  top: 20px;
  left: 20px;
  width: 320px;
  z-index: 100;
}

.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(48, 74, 216, 0.3);
  border-radius: 8px;
  margin-top: 8px;
  max-height: 300px;
  overflow-y: auto;
  backdrop-filter: blur(10px);
}

.search-result-item {
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.2s;
  border-bottom: 1px solid rgba(48, 74, 216, 0.1);
}

.search-result-item:last-child {
  border-bottom: none;
}

.search-result-item:hover {
  background: rgba(48, 74, 216, 0.1);
}

.result-name {
  display: block;
  color: #1a2a6c;
  font-size: 14px;
  font-weight: 500;
}

.result-location {
  display: block;
  color: rgba(26, 42, 108, 0.6);
  font-size: 12px;
  margin-top: 4px;
}
</style>
