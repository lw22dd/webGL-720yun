<template>
  <div class="panorama-container">
    <div ref="containerRef" class="panorama-canvas"></div>

    <PanoTopbar
      :title="sceneTitle"
      @back="goBack"
      @share="copyShareLink"
      @fullscreen="toggleFullscreen"
    />

    <PanoControls
      :auto-rotate="autoRotate"
      @reset="resetCamera"
      @toggle-rotate="toggleAutoRotateAction"
      @zoom-in="zoomIn"
      @zoom-out="zoomOut"
    />

    <SceneStrip
      v-if="scenes.length > 0"
      :scenes="scenes"
      :current-scene-id="currentSceneId"
      @scene-select="handleSceneSelect"
    />

    <div v-if="engineLoading" class="pano-loading">
      <t-loading size="large" text="加载中..." />
    </div>

    <div v-if="engineError" class="pano-error">
      <t-icon name="error-circle" size="48" />
      <p>{{ engineError }}</p>
      <t-button theme="primary" @click="loadScene">重试</t-button>
    </div>

    <PanoHotspotDialog
      v-model:visible="infoDialogVisible"
      :hotspot="currentHotspot"
    />

    
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import SceneApi from '@/apis/scene.api'
import HotspotApi from '@/apis/hotspot.api'
import type { SceneDetailResponse, SceneListItem } from '@/models/scene.model'
import type { HotspotForViewer } from '@/models/hotspot.model'
import { usePanoramaEngine } from '@/composables/usePanoramaEngine'
import PanoTopbar from './panorama/PanoTopbar.vue'
import PanoControls from './panorama/PanoControls.vue'
import PanoHotspotDialog from './panorama/PanoHotspotDialog.vue'
import SceneStrip from './SceneStrip.vue'

const route = useRoute()
const router = useRouter()

const containerRef = ref<HTMLDivElement | null>(null)
const {
  loading: engineLoading,
  error: engineError,
  initScene,
  destroy: destroyEngine,
  zoomIn,
  zoomOut,
  toggleAutoRotate,
  resetCamera: resetCameraPose,
  setOnHotspotClick,
} = usePanoramaEngine(containerRef)

const autoRotate = ref(false)
const currentScene = ref<SceneDetailResponse | null>(null)
const scenes = ref<SceneListItem[]>([])
const currentSceneId = ref<number | null>(null)
const sceneTitle = computed(() => currentScene.value?.title || '全景漫游')

const infoDialogVisible = ref(false)
const quizDialogVisible = ref(false)
const currentHotspot = ref<HotspotForViewer | null>(null)
const selectedQuizOption = ref<number | null>(null)

const loadScene = async () => {
  const sceneCode = route.query.scene as string
  if (!sceneCode) {
    engineError.value = '未指定场景'
    return
  }

  try {
    const result = await SceneApi.getSceneList({ keyword: sceneCode, page: 1, page_size: 10 })

    if (result.code === 200 && result.data && result.data.scenes && result.data.scenes.length > 0) {
      const scene = result.data.scenes.find((s: SceneListItem) => s.scene_code === sceneCode)
      if (scene) {
        currentSceneId.value = scene.id
        scenes.value = result.data.scenes
        const detailResult = await SceneApi.getSceneDetail(scene.id)
        if (detailResult.code === 200 && detailResult.data) {
          currentScene.value = detailResult.data
          await initScene(detailResult.data)
          await loadHotspotsForScene(scene.id)
        }
      }
    } else {
      engineError.value = '场景不存在'
    }
  } catch (e: unknown) {
    engineError.value = e instanceof Error ? e.message : '加载失败'
  }
}

const loadHotspotsForScene = async (sceneId: number) => {
  try {
    const result = await HotspotApi.getHotspotList({
      scene_id: sceneId,
      page: 1,
      page_size: 100
    })
    if (result.code === 200 && result.data) {
      const hotspots = result.data.hotspots as unknown as HotspotForViewer[]
      if (currentScene.value) {
        currentScene.value.hotspots = hotspots as any
      }
    }
  } catch (e) {
    console.error('Failed to load hotspots:', e)
  }
}

const handleMarkerClick = (hotspot: HotspotForViewer) => {
  currentHotspot.value = hotspot

  switch (hotspot.type) {
    case 1:
      if (hotspot.target_scene_id) {
        handleSceneSelect(hotspot.target_scene_id)
      }
      break
    case 2:
      infoDialogVisible.value = true
      break
    case 3:
      selectedQuizOption.value = null
      quizDialogVisible.value = true
      break
  }
}

setOnHotspotClick(handleMarkerClick)

const handleSceneSelect = async (sceneId: number) => {
  try {
    const result = await SceneApi.getSceneDetail(sceneId)
    if (result.code === 200 && result.data) {
      currentScene.value = result.data
      currentSceneId.value = sceneId
      router.push({ query: { scene: result.data.scene_code } })
      await initScene(result.data)
      await loadHotspotsForScene(sceneId)
    }
  } catch (e: unknown) {
    MessagePlugin.error(e instanceof Error ? e.message : '切换场景失败')
  }
}

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    containerRef.value?.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}

const resetCamera = () => {
  if (!currentScene.value) return
  resetCameraPose(
    currentScene.value.initial_pitch || 0,
    currentScene.value.initial_yaw || 0,
    currentScene.value.initial_fov || 50
  )
}

const toggleAutoRotateAction = () => {
  autoRotate.value = !autoRotate.value
  toggleAutoRotate(autoRotate.value)
}

const copyShareLink = () => {
  const link = window.location.href
  navigator.clipboard.writeText(link).then(() => {
    MessagePlugin.success('链接已复制')
  })
}

const goBack = () => {
  if (window.history.state && window.history.length > 2) {
    router.back()
  } else {
    const from = sessionStorage.getItem('panoramaFrom')
    if (from) {
      sessionStorage.removeItem('panoramaFrom')
      router.push(from)
    } else {
      router.push('/index')
    }
  }
}

watch(() => route.query.scene, () => {
  if (route.query.scene) {
    loadScene()
  }
})

onMounted(() => {
  loadScene()
})

onUnmounted(() => {
  destroyEngine()
})
</script>

<style scoped>
.panorama-container {
  width: 100%;
  height: 100vh;
  position: relative;
  overflow: hidden;
  background-color: #000;
}

.panorama-canvas {
  width: 100%;
  height: 100%;
}

.pano-loading {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 200;
}

.pano-error {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: #fff;
  z-index: 200;
}

.pano-error p {
  margin: 0;
  font-size: 16px;
}
</style>

<style>
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
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.4);
}
</style>