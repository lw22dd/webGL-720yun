<template>
  <t-dialog
    v-model:visible="dialogVisible"
    :header="dialogTitle"
    :width="dialogWidth"
    :footer="false"
    top="50px"
    class="form-map-dialog"
    @opened="handleDialogOpen"
  >
    <div class="dialog-layout-container" :class="{ 'with-map': showMapPanel }">
      <div class="form-panel">
        <slot :formData="formData" :submitLoading="submitLoading" :closeDialog="closeDialog" :isEdit="isEdit" :editingId="editingId" />
      </div>

      <div v-if="showMapPanel" class="map-panel">
        <div class="map-picker-header">
          <div class="map-header-left">
            <t-icon name="map" size="18" />
            <span class="panel-title">地理位置选取</span>
          </div>
          <div class="map-search-compact">
            <t-input
              v-model="mapSearchKeyword"
              placeholder="搜索地点并定位..."
              size="small"
              @enter="doMapSearch"
            >
              <template #suffix-icon>
                <t-icon name="search" style="cursor: pointer" @click="doMapSearch" />
              </template>
            </t-input>
            <div v-if="mapSearchResults.length > 0" class="map-search-results">
              <div
                v-for="(item, index) in mapSearchResults"
                :key="index"
                class="map-search-item"
                @click="focusMapLocation(item)"
              >
                <div class="item-name">{{ item.name }}</div>
                <div class="item-address">{{ item.district }}</div>
              </div>
            </div>
          </div>
        </div>

        <div class="map-canvas-wrapper">
          <div id="map-picker-canvas" class="map-canvas"></div>
          <div v-if="mapLoading" class="map-loading-overlay">
            <t-loading size="large" text="地图加载中..." />
          </div>
        </div>

        <div class="map-status-bar">
          <div class="current-loc-tag" v-if="formData.longitude">
            <t-icon name="location-1" size="14" />
            <span>{{ formData.province }}{{ formData.city }} · {{ formData.name }}</span>
          </div>
          <div class="current-loc-hint" v-else>
            <t-icon name="info-circle" size="14" />
            <span>在地图上点击或搜索以拾取坐标</span>
          </div>
          <div class="map-coords" v-if="formData.longitude">
            <span class="coord-badge">{{ Number(formData.longitude).toFixed(4) }}°E</span>
            <span class="coord-badge">{{ Number(formData.latitude).toFixed(4) }}°N</span>
          </div>
        </div>
      </div>
    </div>
  </t-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, nextTick, computed } from 'vue'
import { initAMap } from '@/utils/amap'

interface LocationResult {
  name: string
  address: string
  district: string
  province: string
  city: string
  lng: number
  lat: number
}

const emit = defineEmits<{
  (e: 'mapLocationChange', data: Partial<LocationResult>): void
}>()

const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)
const isEdit = ref(false)
const editingId = ref<string | number | null>(null)

const formData = reactive<Record<string, any>>({
  name: '',
  longitude: '',
  latitude: '',
  province: '',
  city: ''
})

const mapSearchKeyword = ref('')
const mapSearchResults = ref<LocationResult[]>([])
const mapLoading = ref(false)
let pickerMap: any = null
let pickerMarker: any = null

const showMapPanel = computed(() => true)

const dialogWidth = computed(() => '1100px')

const handleDialogOpen = async () => {
  await nextTick()
  await new Promise(resolve => setTimeout(resolve, 200))

  const container = document.getElementById('map-picker-canvas')
  if (!container) return

  try {
    mapLoading.value = true

    if (!pickerMap) {
      pickerMap = await initAMap(container, {
        zoom: 12,
        center: formData.longitude ? [Number(formData.longitude), Number(formData.latitude)] : [116.397428, 39.90923]
      })

      pickerMap.on('click', (e: any) => {
        const { lng, lat } = e.lnglat
        updateLocationFromMap({ lng, lat })
      })
    } else {
      pickerMap.resize()
    }

    if (formData.longitude && formData.latitude) {
      const pos = [Number(formData.longitude), Number(formData.latitude)]
      if (window.AMap) {
        if (!pickerMarker) {
          pickerMarker = new window.AMap.Marker({ position: pos, map: pickerMap })
        } else {
          pickerMarker.setPosition(pos)
        }
        pickerMap.setCenter(pos)
      }
    }

    if (!pickerMap.isComplete()) {
      await new Promise<void>((resolve) => {
        pickerMap.on('complete', resolve)
      })
    }
  } catch (error) {
    console.error('地图加载失败:', error)
  } finally {
    mapLoading.value = false
  }
}

const updateLocationFromMap = (loc: Partial<LocationResult>) => {
  if (loc.lng) formData.longitude = String(loc.lng.toFixed(6))
  if (loc.lat) formData.latitude = String(loc.lat.toFixed(6))
  if (loc.name) formData.name = loc.name
  if (loc.province) formData.province = loc.province
  if (loc.city) formData.city = loc.city

  if (pickerMap && loc.lng && window.AMap) {
    const pos = [loc.lng, loc.lat]
    if (!pickerMarker) {
      pickerMarker = new window.AMap.Marker({ position: pos, map: pickerMap })
    } else {
      pickerMarker.setPosition(pos)
    }
  }

  emit('mapLocationChange', {
    lng: loc.lng,
    lat: loc.lat,
    name: loc.name,
    province: loc.province,
    city: loc.city
  })
}

const doMapSearch = async () => {
  if (!mapSearchKeyword.value) return
  try {
    const { autoComplete } = await import('@/utils/amap')
    const tips = await autoComplete(mapSearchKeyword.value)
    mapSearchResults.value = tips.map((tip: any) => ({
      name: tip.name,
      district: tip.district,
      address: tip.address,
      lng: tip.location.lng,
      lat: tip.location.lat,
      province: extractProvince(tip.district),
      city: extractCity(tip.district)
    }))
  } catch (e) {
    console.error(e)
  }
}

const focusMapLocation = (item: LocationResult) => {
  updateLocationFromMap(item)
  pickerMap.setCenter([item.lng, item.lat])
  pickerMap.setZoom(15)
  mapSearchResults.value = []
  mapSearchKeyword.value = item.name
}

const extractProvince = (district: string): string => {
  if (!district) return ''
  const parts = district.split(/省|自治区|直辖市|特别行政区/)
  if (parts.length > 1) {
    const prefix = parts[0]
    if (district.includes('省')) return prefix + '省'
    if (district.includes('自治区')) return prefix + '自治区'
    if (district.includes('直辖市')) return prefix + '直辖市'
    if (district.includes('特别行政区')) return prefix + '特别行政区'
  }
  const match = district.match(/^(.+?)(省|市|自治区)/)
  return match ? match[0] : district.split(/市/)[0] + '省'
}

const extractCity = (district: string): string => {
  if (!district) return ''
  const parts = district.split(/省|自治区/)
  if (parts.length > 1) {
    const afterProvince = parts[1]
    const cityMatch = afterProvince.match(/^(.+?市)/)
    return cityMatch ? cityMatch[1] : afterProvince.split(/地区|州|盟/)[0]
  }
  const match = district.match(/市(.+?区|.+?县|.+?市)/)
  if (match) {
    const cityPart = district.substring(0, district.indexOf(match[1]) + match[1].length)
    const cityMatch2 = cityPart.match(/(.+?市)/)
    return cityMatch2 ? cityMatch2[1] : ''
  }
  return ''
}

const openAddDialog = (title: string = '新增') => {
  isEdit.value = false
  editingId.value = null
  dialogTitle.value = title
  Object.keys(formData).forEach(key => {
    formData[key] = ''
  })
  dialogVisible.value = true
}

const openEditDialog = (data: Record<string, any>, id: string | number, title: string = '编辑') => {
  isEdit.value = true
  editingId.value = id
  dialogTitle.value = title
  Object.keys(formData).forEach(key => {
    delete formData[key]
  })
  Object.assign(formData, { id, ...data })
  dialogVisible.value = true
}

const closeDialog = () => {
  dialogVisible.value = false
}

const setSubmitLoading = (loading: boolean) => {
  submitLoading.value = loading
}

defineExpose({ openAddDialog, openEditDialog, closeDialog, setSubmitLoading, formData })
</script>

<style scoped>
.dialog-layout-container {
  display: flex;
  gap: 24px;
  min-height: 520px;
}

.dialog-layout-container.with-map .form-panel {
  width: 380px;
  flex-shrink: 0;
  padding-right: 24px;
  border-right: 1px solid #f0f1f2;
}

.form-panel {
  flex: 1;
  padding: 4px 0;
}

:deep(.t-form__item) {
  margin-bottom: 14px;
}

:deep(.t-form__label) {
  font-weight: 500;
  color: #4e5969;
}

:deep(.t-input--disabled .t-input__inner) {
  color: #1d2129;
  -webkit-text-fill-color: #1d2129;
}

:deep(.t-input--disabled) {
  background: #f7f8fa;
  border-color: #e5e6eb;
}

:deep(.t-textarea--disabled) {
  background: #f7f8fa;
  border-color: #e5e6eb;
}

:deep(.t-textarea--disabled textarea) {
  color: #86909c;
}

/* ===== 地图面板 ===== */
.map-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #ffffff;
  border-radius: 12px;
  padding: 16px;
  border: 1px solid #eef0f2;
  min-width: 0;
}

.map-picker-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.map-header-left {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.map-header-left :deep(.t-icon) {
  color: #0052d9;
}

.panel-title {
  font-weight: 600;
  font-size: 15px;
  color: #1d2129;
}

.map-search-compact {
  flex: 1;
  margin-left: 16px;
  position: relative;
  max-width: 320px;
}

.map-canvas-wrapper {
  flex: 1;
  position: relative;
  min-height: 400px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e5e6eb;
  background: #f5f7fa;
}

.map-canvas {
  width: 100%;
  height: 100%;
  min-height: 400px;
}

.map-loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.85);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
  backdrop-filter: blur(2px);
}

.map-status-bar {
  margin-top: 10px;
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #eef0f2;
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 36px;
}

.current-loc-tag {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #0052d9;
  font-weight: 500;
  font-size: 13px;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.current-loc-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #86909c;
  font-size: 13px;
  flex: 1;
}

.map-coords {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.coord-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  background: #e8f3ff;
  color: #0052d9;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  font-family: 'SF Mono', 'Cascadia Code', monospace;
}

/* 地图内搜索结果 */
.map-search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 6px;
  margin-top: 4px;
  max-height: 180px;
  overflow-y: auto;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  z-index: 100;
}

.map-search-item {
  padding: 8px 12px;
  cursor: pointer;
  border-bottom: 1px solid #f2f3f5;
}

.map-search-item:hover { background: #f5f7fa; }
.map-search-item:last-child { border-bottom: none; }

.item-name { font-size: 13px; font-weight: 500; }
.item-address { font-size: 11px; color: #86909c; }
</style>