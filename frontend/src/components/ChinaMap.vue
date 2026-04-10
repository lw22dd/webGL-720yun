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
let resizeObserver: ResizeObserver | null = null

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
      geo: {// 地图配置
        map: 'china',
        roam: false,
        zoom: 1.1,
        center: [105, 36],
        layoutCenter: ['50%', '50%'],
        layoutSize: '95%',
        label: {
          color: '#000000'
        },
        itemStyle: { // 地图区域样式
          areaColor: '#e6f0ff',
          borderColor: '#4080ff',
          borderWidth: 1,
          shadowColor: 'rgba(0, 82, 217, 0.15)',
          shadowBlur: 8
        },
        emphasis: { // 地图区域选中样式
          label: {
            color: '#000000'
          },
          itemStyle: {
            areaColor: '#c8dfff',
          }
        }
      },
      series: [ // 散点图配置
        {
          type: 'scatter',
          coordinateSystem: 'geo',
          data: scatterData,
          symbolSize: 20,
          symbol: 'circle',
          itemStyle: {
            color: '#0052d9',
            shadowBlur: 15,
            shadowColor: 'rgba(0, 82, 217, 0.4)'
          },
          emphasis: { // 散点图选中样式
            scale: 1.5,
            itemStyle: {
              shadowBlur: 25,
              shadowColor: 'rgba(0, 82, 217, 0.5)',
              color: '#4080ff'
            }
          },
          label: { // 散点图标签样式
            show: true,
            formatter: '{b}',
            position: 'top',
            color: '#000000',
            fontSize: 12,
            fontWeight: 'bold',
            backgroundColor: '#ffffff',
            padding: [4, 8],
            borderRadius: 4,
            borderColor: '#e5e6eb',
            borderWidth: 1,
            shadowBlur: 4,
            shadowColor: 'rgba(0, 0, 0, 0.1)'
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

    resizeObserver = new ResizeObserver(() => {
      chartInstance?.resize()
    })
    if (mapContainer.value) {
      resizeObserver.observe(mapContainer.value)
    }

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
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
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

  background: #213e6a; /** 地图背景颜色 */
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
  background: #ffffff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  margin-top: 8px;
  max-height: 300px;
  overflow-y: auto;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.search-result-item {
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.2s;
  border-bottom: 1px solid #f2f3f5;
}

.search-result-item:last-child {
  border-bottom: none;
}

.search-result-item:hover {
  background: #f2f3f5;
}

.result-name {
  display: block;
  color: #1d2129;
  font-size: 14px;
  font-weight: 500;
}

.result-location {
  display: block;
  color: #86909c;
  font-size: 12px;
  margin-top: 4px;
}
</style>
