<template>
  <t-dialog
    :visible="visible"
    header="全景预览"
    width="90vw"
    :footer="false"
    @close="handleClose"
    class="scene-preview-dialog"
  >
    <div class="preview-container" v-if="scene">
      <div class="preview-toolbar">
        <div class="toolbar-info">
          <span class="scene-title">{{ scene.title }}</span>
          <span class="scene-code">{{ scene.scene_code }}</span>
        </div>
        <div class="toolbar-actions">
          <t-button size="small" @click="resetCamera" title="重置视角">
            ⟲ 重置
          </t-button>
          <t-button
            :theme="autoRotate ? 'primary' : 'default'"
            size="small"
            @click="toggleAutoRotate"
          >
            ↻ 自动旋转
          </t-button>
        </div>
      </div>

      <div ref="viewerContainer" class="viewer-container"></div>

      <div v-if="loading" class="preview-loading">
        <t-loading size="large" text="加载中..." />
      </div>

      <div v-if="error" class="preview-error">
        <t-icon name="error-circle" size="48" />
        <p>{{ error }}</p>
        <t-button theme="primary" @click="loadPanorama">重试</t-button>
      </div>
    </div>
  </t-dialog>
</template>

<script setup lang="ts">
import { ref, watch, onUnmounted } from 'vue'
import { Viewer } from '@photo-sphere-viewer/core'
import { MarkersPlugin } from '@photo-sphere-viewer/markers-plugin'
import '@photo-sphere-viewer/core/index.css'
import '@photo-sphere-viewer/markers-plugin/index.css'
import SceneApi from '@/apis/scene.api'
import HotspotApi from '@/apis/hotspot.api'
import type { SceneListItem } from '@/models/scene.model'

interface Props {
  visible: boolean
  scene: SceneListItem | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const viewerContainer = ref<HTMLDivElement | null>(null)
const loading = ref(false)
const error = ref('')
const autoRotate = ref(false)

let viewer: Viewer | null = null
let markersPlugin: MarkersPlugin | null = null

const getPreviewUrl = (sceneCode: string) => {
  const base = import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'
  return `${base}/api/v1/res/previews/${sceneCode}`
}

const getSceneDetail = async () => {
  if (!props.scene) return null

  try {
    const result = await SceneApi.getSceneDetail(props.scene.id)
    if (result.code === 200 && result.data) {
      return result.data
    }
  } catch (e) {
    console.error('Failed to get scene detail:', e)
  }
  return null
}

const loadPanorama = async () => {
  if (!viewerContainer.value || !props.scene) return

  if (viewer) {
    viewer.destroy()
    viewer = null
  }

  loading.value = true
  error.value = ''

  try {
    const sceneDetail = await getSceneDetail()
    if (!sceneDetail) {
      error.value = '无法获取场景详情'
      return
    }

    const panoramaUrl = sceneDetail.tile_url || sceneDetail.source_url || getPreviewUrl(props.scene.scene_code)

    viewer = new Viewer({
      container: viewerContainer.value,
      panorama: panoramaUrl,
      defaultZoomLvl: 50,
      minFov: 30,
      maxFov: 90,
      navbar: false,
      plugins: [
        [MarkersPlugin, {}]
      ]
    })

    markersPlugin = viewer.getPlugin(MarkersPlugin) as MarkersPlugin

    viewer.addEventListener('ready', () => {
      loading.value = false
    })

    await loadHotspots(sceneDetail.id)

  } catch (e: any) {
    error.value = e.message || '加载失败'
    loading.value = false
  }
}

const loadHotspots = async (sceneId: number) => {
  if (!markersPlugin) return

  try {
    const result = await HotspotApi.getHotspotList({ scene_id: sceneId })
    if (result.code === 200 && result.data) {
      result.data.hotspots.forEach((hotspot: any) => {
        const markerConfig: any = {
          id: `hotspot-${hotspot.id}`,
          position: { pitch: hotspot.pitch, yaw: hotspot.yaw },
          tooltip: hotspot.title,
          data: hotspot
        }

        switch (hotspot.type) {
          case 1:
            markerConfig.html = `<div class="hotspot-marker hotspot-scene">→</div>`
            break
          case 2:
            markerConfig.html = `<div class="hotspot-marker hotspot-info">ℹ</div>`
            break
          case 3:
            markerConfig.html = `<div class="hotspot-marker hotspot-quiz">?</div>`
            break
        }

        markersPlugin?.addMarker(markerConfig)
      })

      markersPlugin?.addEventListener('select-marker', ({ marker }: any) => {
        const hotspotId = marker.config.id.replace('hotspot-', '')
        handleHotspotClick(parseInt(hotspotId))
      })
    }
  } catch (e) {
    console.error('Failed to load hotspots:', e)
  }
}

const handleHotspotClick = (hotspotId: number) => {
  console.log('Hotspot clicked:', hotspotId)
}

const resetCamera = () => {
  if (viewer) {
    viewer.animate({
      pitch: 0,
      yaw: 0,
      zoom: 50,
      speed: '2rpm'
    })
  }
}

const toggleAutoRotate = () => {
  autoRotate.value = !autoRotate.value
  if (viewer) {
    if (autoRotate.value) {
      (viewer as any).startAutoRotate({ speed: '2rpm' })
    } else {
      (viewer as any).stopAutoRotate()
    }
  }
}

const handleClose = () => {
  if (viewer) {
    viewer.destroy()
    viewer = null
  }
  emit('close')
}

watch(() => props.visible, (newVal) => {
  if (newVal && props.scene) {
    setTimeout(() => {
      loadPanorama()
    }, 100)
  }
})

onUnmounted(() => {
  if (viewer) {
    viewer.destroy()
    viewer = null
  }
})
</script>

<style scoped>
.preview-container {
  width: 100%;
  height: 70vh;
  min-height: 500px;
  display: flex;
  flex-direction: column;
  background: #000;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
}

.preview-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.8);
  z-index: 10;
}

.toolbar-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.scene-title {
  color: #fff;
  font-size: 16px;
  font-weight: 500;
}

.scene-code {
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  padding: 2px 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.toolbar-actions {
  display: flex;
  gap: 8px;
}

.viewer-container {
  flex: 1;
  width: 100%;
}

.preview-loading,
.preview-error {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: #fff;
  z-index: 20;
}

.preview-error p {
  margin: 0;
}
</style>

<style>
.scene-preview-dialog .t-dialog__body {
  padding: 0 !important;
}

.hotspot-marker {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.hotspot-scene {
  background: linear-gradient(135deg, #0EA5E9, #4080ff);
  color: #fff;
}

.hotspot-info {
  background: linear-gradient(135deg, #10B981, #059669);
  color: #fff;
}

.hotspot-quiz {
  background: linear-gradient(135deg, #F59E0B, #D97706);
  color: #fff;
}

.hotspot-marker:hover {
  transform: scale(1.2);
}
</style>
