<!--高德地图-->
<template>
  <div class="space-detail-map">
    
    <div ref="mapContainer" class="map-container"></div>

    <div v-if="loading" class="map-loading-overlay">
      <div class="loading-spinner"></div>
      <span class="loading-text">地图加载中...</span>
    </div>

    <MapToolbar
      :basemap-type="basemapType"
      :show-p-o-i="showPOI"
      :show-markers="showMarkers"
      :show-p-o-i-toggle="!isAdmin"
      :show-marker-toggle="true"
      :space-name="spaceData?.name"
      :space-province="spaceData?.province"
      :space-city="spaceData?.city"
      :scene-count="sceneList.length"
      :show-click-hint="isAdmin"
      @switch-basemap="handleBasemapSwitch"
      @toggle-p-o-i="handleTogglePOI"
      @toggle-markers="handleToggleMarkers"
      @back="goBack"
    />


  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Scene as L7Scene, Marker, Popup } from '@antv/l7'
import { GaodeMap } from '@antv/l7-maps'
import MapToolbar from './MapToolbar.vue'
import type { BasemapType } from '@/composables/useL7Scene'
import type { SceneMarkerData } from '@/composables/useSceneMarkers'
import SpaceApi from '@/apis/space.api'
import type { SpaceDetailResponse } from '@/models/space.model'

const props = defineProps<{
  isAdmin?: boolean
}>()

const emit = defineEmits<{
  (e: 'sceneClick', scene: SceneMarkerData): void
}>()

const router = useRouter()
const route = useRoute()

const mapContainer = ref<HTMLDivElement | null>(null)
const loading = ref(true)
const spaceData = ref<SpaceDetailResponse | null>(null)
const sceneList = ref<SceneMarkerData[]>([])
const basemapType = ref<BasemapType>('vector')
const showPOI = ref(false)
const showMarkers = ref(true)
const activeSceneId = ref<number | null>(null)

let l7Scene: L7Scene | null = null
let aMapInstance: any = null
let markerList: Marker[] = []
let satelliteLayer: any = null
let normalLayer: any = null

const AMAP_KEY = import.meta.env.VITE_AMAP_KEY || ''
const AMAP_SECURITY_KEY = import.meta.env.VITE_AMAP_SECURITY_KEY || ''

const getPreviewUrl = (sceneCode: string) => {
  return `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'}/api/v1/res/previews/${sceneCode}`
}

function goBack() {
  router.push({ name: 'home' })
}

const loadAmapScript = (): Promise<void> => {
  return new Promise((resolve, reject) => {
    if ((window as any).AMap) {
      resolve()
      return
    }

    // 配置安全密钥
    if (AMAP_SECURITY_KEY) {
      (window as any)._AMapSecurityConfig = {
        securityJsCode: AMAP_SECURITY_KEY,
      }
    }

    const script = document.createElement('script')
    script.src = `https://webapi.amap.com/maps?v=2.0&key=${AMAP_KEY}`
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('Failed to load AMap SDK'))
    document.head.appendChild(script)
  })
}

async function loadSpaceData() {
  const spaceId = Number(route.params.id)
  if (!spaceId) {
    console.error('No space ID provided')
    return false
  }

  try {
    loading.value = true
    const result = await SpaceApi.getSpaceDetail(spaceId)
    if (result.code === 200 && result.data) {
      spaceData.value = result.data
      sceneList.value = (result.data.scenes || []).map((s: any) => ({
        id: s.id,
        title: s.title,
        longitude: s.longitude || 0,
        latitude: s.latitude || 0,
        thumbnail_url: s.scene_code ? getPreviewUrl(s.scene_code) : s.thumbnail_url,
        scene_code: s.scene_code,
        view_count: s.view_count,
        hotspots: s.hotspots,
      }))
      console.log('Space data loaded:', spaceData.value?.name, 'with', sceneList.value.length, 'scenes')
      return true
    } else {
      console.error('Failed to load space data:', result)
      return false
    }
  } catch (error) {
    console.error('Failed to load space data:', error)
    return false
  }
}

const getMarkerColor = (index: number): string => {
  const colors = ['#0EA5E9', '#2563EB', '#10B981', '#F59E0B', '#8B5CF6', '#EC4899', '#DC2626', '#059669']
  return colors[index % colors.length]
}

const createMarkerElement = (scene: SceneMarkerData, index: number): HTMLDivElement => {
  const el = document.createElement('div')
  el.className = 'scene-marker-wrapper'
  el.innerHTML = `
    <div class="scene-marker" style="--marker-color: ${getMarkerColor(index)}">
      <div class="marker-pin">
        <svg width="28" height="38" viewBox="0 0 28 38">
          <path d="M14 0C6.268 0 0 6.268 0 14c0 10.5 14 24 14 24s14-13.5 14-24C28 6.268 21.732 0 14 0z"
                fill="var(--marker-color)" stroke="rgba(255,255,255,0.8)" stroke-width="1.5"/>
          <circle cx="14" cy="14" r="6" fill="rgba(255,255,255,0.95)"/>
          <circle cx="14" cy="14" r="3" fill="var(--marker-color)"/>
        </svg>
      </div>
      <div class="marker-pulse"></div>
    </div>
    <div class="marker-label">${scene.title}</div>
  `
  return el
}

const createPopupContent = (scene: SceneMarkerData): string => {
  return `
    <div class="scene-popup">
      <div class="popup-header">
        <h3 class="popup-title">${scene.title}</h3>
        <span class="popup-coords">${scene.longitude.toFixed(6)}, ${scene.latitude.toFixed(6)}</span>
      </div>
      ${scene.thumbnail_url ? `<div class="popup-thumbnail"><img src="${scene.thumbnail_url}" alt="${scene.title}" /></div>` : ''}
      <div class="popup-info">
        <div class="info-row">
          <span class="info-label">场景 ID</span>
          <span class="info-value">${scene.id}</span>
        </div>
        ${scene.hotspots && scene.hotspots.length > 0 ? `
        <div class="info-row">
          <span class="info-label">连接热点</span>
          <span class="info-value">${scene.hotspots.length} 个</span>
        </div>` : ''}
        <div class="info-row">
          <span class="info-label info-link" data-scene-id="${scene.id}">查看全景 →</span>
        </div>
      </div>
    </div>
  `
}

const loadSceneMarkers = () => {
  if (!l7Scene) return

  // 清除已有 Marker
  l7Scene!.removeAllMarkers()
  markerList = []

  sceneList.value.forEach((scene, index) => {
    if (!scene.longitude || !scene.latitude) return

    const el = createMarkerElement(scene, index)

    const marker = new Marker({
      element: el,
      offsets: [-14, -38],
      anchor: 'bottom-left' as any,
    }).setLnglat([scene.longitude, scene.latitude] as any)

    // 绑定 Popup
    const popup = new Popup({
      offsets: [0, -42],
      closeButton: true,
      maxWidth: '300px',
      title: '',
      html: createPopupContent(scene),
    })
    marker.setPopup(popup)

    // 点击 Marker
    el.addEventListener('click', (e) => {
      e.stopPropagation()
      activeSceneId.value = scene.id
      marker.togglePopup()
      if (!props.isAdmin) {
        emit('sceneClick', scene)
      }
    })

    if (l7Scene) {
      l7Scene.addMarker(marker)
      markerList.push(marker)
    }
  })
}

async function initMap() {
  if (!mapContainer.value || !spaceData.value) {
    console.warn('Map container or space data not available')
    return
  }

  // 检查坐标是否有效
  const lng = spaceData.value.longitude
  const lat = spaceData.value.latitude
  if (!lng || !lat || lng === 0 || lat === 0) {
    console.warn('Invalid coordinates:', lng, lat)
    spaceData.value.longitude = 116.4074
    spaceData.value.latitude = 39.9042
  }

  try {
    console.log('Initializing map with center:', [spaceData.value.longitude, spaceData.value.latitude])

    // 加载高德地图 SDK
    await loadAmapScript()

    // 用 L7 Scene 直接创建地图
    l7Scene = new L7Scene({
      id: mapContainer.value,
      map: new GaodeMap({
        style: 'normal',
        center: [spaceData.value.longitude, spaceData.value.latitude],
        zoom: 15,
        pitch: 0,
      }),
      logoVisible: false,
    })

    // 尝试获取 AMap 实例的方法
    const tryGetAMapInstance = () => {
      const mapService = (l7Scene as any).mapService
      if (mapService?.map) return mapService.map
      if (mapService) {
        for (const key of Object.keys(mapService)) {
          const val = (mapService as any)[key]
          if (val && typeof val.setMapStyle === 'function') return val
        }
      }
      const container = mapContainer.value
      if (container) {
        const amapKey = Object.keys(container).find(k => k.startsWith('__amap__'))
        if (amapKey) return (container as any)[amapKey]
      }
      return null
    }

    // 定时尝试获取 AMap 实例
    const aMapReady = setInterval(() => {
      const instance = tryGetAMapInstance()
      if (instance) {
        clearInterval(aMapReady)
        aMapInstance = instance
        console.log('Got AMap instance:', !!aMapInstance)

        // 获取地图当前的图层数组，保存默认矢量图层
        if (!normalLayer) {
          const layers = aMapInstance.getLayers()
          console.log('Current layers:', layers)
          // 查找矢量底图
          const AMap = (window as any).AMap
          layers.forEach((layer: any) => {
            if (layer && typeof layer.setOpacity === 'function') {
              normalLayer = layer
            }
          })
          // 如果没找到，使用 createDefaultLayer 创建标准矢量图层
          if (!normalLayer && AMap.createDefaultLayer) {
            normalLayer = AMap.createDefaultLayer()
          }
        }

        if (typeof aMapInstance.setFeatures === 'function') {
          aMapInstance.setFeatures(['bg', 'road', 'building'])
        }

        loadSceneMarkers()
        if (props.isAdmin) {
          l7Scene?.on('click', handleMapClick)
        }
        loading.value = false
      }
    }, 500)

    // 监听 L7 loaded 事件
    l7Scene.on('loaded', () => {
      clearInterval(aMapReady)
      if (!aMapInstance) {
        aMapInstance = tryGetAMapInstance()
      }
      console.log('L7 loaded event fired, aMapInstance:', !!aMapInstance)

      // 获取并保存默认矢量图层
      if (!normalLayer && aMapInstance) {
        const layers = aMapInstance.getLayers()
        console.log('Current layers on loaded:', layers)
        const AMap = (window as any).AMap
        layers.forEach((layer: any) => {
          if (layer && typeof layer.setOpacity === 'function') {
            normalLayer = layer
          }
        })
        if (!normalLayer && AMap.createDefaultLayer) {
          normalLayer = AMap.createDefaultLayer()
        }
      }

      if (aMapInstance && typeof aMapInstance.setFeatures === 'function') {
        aMapInstance.setFeatures(['bg', 'road', 'building'])
      }

      loadSceneMarkers()
      if (props.isAdmin) {
        l7Scene?.on('click', handleMapClick)
      }
      loading.value = false
    })

    // 超时保护
    setTimeout(() => {
      if (loading.value) {
        clearInterval(aMapReady)
        console.warn('Timeout: forcing loading to false')
        loading.value = false
      }
    }, 10000)

  } catch (error) {
    console.error('Failed to init map:', error)
    loading.value = false
  }
}

function handleBasemapSwitch(type: BasemapType) {
  basemapType.value = type

  if (!aMapInstance || !(window as any).AMap) {
    console.warn('handleBasemapSwitch: aMapInstance or AMap not available')
    return
  }

  const AMap = (window as any).AMap
  console.log('Switching to:', type, 'normalLayer:', !!normalLayer, 'satelliteLayer:', !!satelliteLayer)

  if (type === 'satellite') {
    // 切换到卫星图：隐藏矢量图层，显示卫星图层
    if (normalLayer) {
      normalLayer.hide()
      console.log('Hid normalLayer')
    }
    if (!satelliteLayer) {
      satelliteLayer = new AMap.TileLayer.Satellite()
      aMapInstance.add(satelliteLayer)
      console.log('Created and added satelliteLayer')
    } else {
      satelliteLayer.show()
      console.log('Showed satelliteLayer')
    }
  } else {
    // 切换到矢量图：隐藏卫星图层，显示矢量图层
    if (satelliteLayer) {
      satelliteLayer.hide()
      console.log('Hid satelliteLayer')
    }
    if (!normalLayer) {
      normalLayer = AMap.createDefaultLayer()
      aMapInstance.add(normalLayer)
      console.log('Created and added normalLayer')
    } else {
      normalLayer.show()
      console.log('Showed normalLayer')
    }
    // 重新应用 setFeatures 确保样式一致
    if (typeof aMapInstance.setFeatures === 'function') {
      aMapInstance.setFeatures(showPOI.value ? ['bg', 'road', 'building', 'point'] : ['bg', 'road', 'building'])
    }
  }

  // 验证图层切换结果
  setTimeout(() => {
    const layers = aMapInstance.getLayers()
    console.log('Current layers after switch:', layers.length, 'layers')
  }, 100)
}

function handleTogglePOI(show: boolean) {
  showPOI.value = show
  if (!aMapInstance) return
  if (show) {
    aMapInstance.setFeatures(['bg', 'road', 'building', 'point'])
  } else {
    aMapInstance.setFeatures(['bg', 'road', 'building'])
  }
}

function handleToggleMarkers(show: boolean) {
  showMarkers.value = show
  markerList.forEach(marker => {
    if (show) marker.show()
    else marker.hide()
  })
}

function handleSceneClick(scene: SceneMarkerData) {
  activeSceneId.value = scene.id
  if (aMapInstance) {
    aMapInstance.setZoomAndCenter(17, [scene.longitude, scene.latitude], false, 500)
  }
  const index = sceneList.value.findIndex(s => s.id === scene.id)
  if (index >= 0 && markerList[index]) {
    markerList[index].openPopup()
  }
  if (!props.isAdmin) {
    emit('sceneClick', scene)
  }
}

async function handleMapClick(e: any) {
  if (e.target?.closest?.('.scene-marker-wrapper')) return

  const { lng, lat } = e.lngLat
  const newScene: SceneMarkerData = {
    id: Date.now(),
    title: `新场景节点`,
    longitude: lng,
    latitude: lat,
    scene_code: `scene-${Date.now()}`,
  }

  sceneList.value.push(newScene)

  const el = createMarkerElement(newScene, sceneList.value.length - 1)
  const marker = new Marker({
    element: el,
    offsets: [-14, -38],
    anchor: 'bottom-left' as any,
  }).setLnglat([lng, lat] as any)

  const popup = new Popup({
    offsets: [0, -42],
    closeButton: true,
    maxWidth: '300px',
    title: '',
    html: createPopupContent(newScene),
  })
  marker.setPopup(popup)

  el.addEventListener('click', (ev) => {
    ev.stopPropagation()
    activeSceneId.value = newScene.id
    marker.togglePopup()
  })

  if (l7Scene) {
    l7Scene.addMarker(marker)
    markerList.push(marker)
    marker.openPopup()
  }
}

onMounted(async () => {
  const loaded = await loadSpaceData()
  if (loaded && spaceData.value) {
    await initMap()
  } else {
    console.error('Cannot initialize map: space data not loaded')
    loading.value = false
  }
})

onUnmounted(() => {
  markerList = []
  if (l7Scene) {
    l7Scene.destroy()
    l7Scene = null
  }
  if (aMapInstance) {
    aMapInstance.destroy()
    aMapInstance = null
  }
})
</script>

<style scoped>
.space-detail-map {
  position: relative;
  width: 100%;
  height: 100%;
  background: #0F1826;
  overflow: hidden;
}

.map-container {
  width: 100%;
  height: 100%;
}

.map-loading-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(15, 24, 38, 0.7);
  backdrop-filter: blur(4px);
  z-index: 90;
  gap: 12px;
}

.loading-spinner {
  width: 36px;
  height: 36px;
  border: 3px solid rgba(14, 165, 233, 0.2);
  border-top-color: #0EA5E9;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-text {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}


</style>

<style>
.scene-marker-wrapper {
  position: relative;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.scene-marker {
  position: relative;
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.3));
  transition: transform 0.2s ease;
}

.scene-marker-wrapper:hover .scene-marker {
  transform: scale(1.1) translateY(-2px);
}

.marker-pulse {
  position: absolute;
  bottom: -4px;
  left: 50%;
  transform: translateX(-50%);
  width: 12px;
  height: 4px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 50%;
  filter: blur(2px);
}

.marker-label {
  margin-top: 2px;
  padding: 2px 8px;
  background: rgba(15, 24, 38, 0.85);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  font-size: 11px;
  font-family: 'Noto Sans SC', sans-serif;
  color: rgba(255, 255, 255, 0.85);
  white-space: nowrap;
  font-weight: 500;
}

.scene-popup {
  min-width: 240px;
  max-width: 300px;
  font-family: 'Noto Sans SC', sans-serif;
}

.popup-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
  padding-bottom: 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.popup-title {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.popup-coords {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.35);
  font-family: 'SF Mono', 'Fira Code', monospace;
  white-space: nowrap;
}

.popup-thumbnail {
  width: 100%;
  height: 110px;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 12px;
}

.popup-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.popup-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.45);
}

.info-value {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.85);
  font-weight: 500;
}

.info-link {
  color: #0EA5E9;
  cursor: pointer;
}

.info-link:hover {
  text-decoration: underline;
}
</style>
