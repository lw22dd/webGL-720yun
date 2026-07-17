<!-- 全景漫游组件 -->
<template>
  <div class="panorama-container">
    <div ref="containerRef" class="panorama-canvas" :class="{ 'edit-mode': editMode }"></div>

    <PanoTopbar :title="sceneTitle" :is-admin="isAdmin" v-model:edit-mode="editMode" @back="goBack"
      @minimap="toggleMinimap" @fullscreen="toggleFullscreen" />

    <div v-if="engineLoading" class="pano-loading">
      <t-loading size="large" text="加载中..." />
    </div>

    <div v-if="engineError" class="pano-error">
      <t-icon name="error-circle" size="48" />
      <p>{{ engineError }}</p>
      <t-button theme="primary" @click="loadScene">重试</t-button>
    </div>

    <PanoHotspotDialog v-model:visible="infoDialogVisible" :hotspot="currentHotspot" />

    <PanoHotspotEditor v-model:visible="editorVisible" :form="editorForm" :mode="editorMode" :scene-list="sceneList"
      :space-id="currentScene?.space_id" :scene-code="currentScene?.scene_code" @submitted="onEditorSubmitted"
      @deleted="onEditorSubmitted" />

    <PanoMinimap v-model:visible="minimapVisible" :current-scene="currentScene" :get-viewer="getViewer"
      @scene-select="handleSceneSelect" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import SceneApi from '@/apis/scene.api'
import HotspotApi from '@/apis/hotspot.api'
import type { SceneDetailResponse, SceneListItem } from '@/models/scene.model'
import type { HotspotForViewer, HotspotEditorForm } from '@/models/hotspot.model'
import { useUserStore } from '@/stores/user.store'
import { usePanoramaEngine } from '@/composables/usePanoramaEngine'
import PanoTopbar from './panorama/PanoTopbar.vue'
import PanoHotspotDialog from './panorama/PanoHotspotDialog.vue'
import PanoHotspotEditor from './panorama/PanoHotspotEditor.vue'
import PanoMinimap from './panorama/PanoMinimap.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const containerRef = ref<HTMLDivElement | null>(null)
const {
  loading: engineLoading,
  error: engineError,
  initScene,
  destroy: destroyEngine,
  setOnHotspotClick,
  setOnViewerClick,
  reloadHotspots,
  getViewer,
} = usePanoramaEngine(containerRef)
const currentScene = ref<SceneDetailResponse | null>(null)
const currentSceneId = ref<number | null>(null)
const sceneTitle = computed(() => currentScene.value?.title || '全景漫游')
const isAdmin = computed(() => userStore.isAdmin)

const infoDialogVisible = ref(false)
const quizDialogVisible = ref(false)
const currentHotspot = ref<HotspotForViewer | null>(null)
const selectedQuizOption = ref<number | null>(null)
const minimapVisible = ref(false)

// 编辑模式状态
const editMode = ref(false)
const editorVisible = ref(false)
const editorMode = ref<'create' | 'edit'>('create')
const editorForm = ref<HotspotEditorForm>({
  scene_id: 0,
  type: 1,
  pitch: 0,
  yaw: 0,
  title: '',
  icon_source: 'preset',
  icon_preset_key: 'arrow-blue',
  transition_effect: 'fade',
})
const sceneList = ref<SceneListItem[]>([])

const loadScene = async () => {
  const sceneId = route.query.id ? Number(route.query.id) : null
  const sceneCode = route.query.scene as string

  if (!sceneId && !sceneCode) {
    engineError.value = '未指定场景'
    return
  }

  try {
    let targetSceneId: number | null = sceneId

    // 兼容旧链接：只有 sceneCode 时先查列表换 id
    if (!targetSceneId && sceneCode) {
      const result = await SceneApi.getSceneList({ keyword: sceneCode, page: 1, page_size: 10 })
      if (result.code === 200 && result.data?.scenes?.length) {
        const scene = result.data.scenes.find((s: SceneListItem) => s.scene_code === sceneCode)
        if (scene) {
          targetSceneId = scene.id
        }
      }
    }

    if (!targetSceneId) {
      engineError.value = '场景不存在'
      return
    }

    currentSceneId.value = targetSceneId
    const detailResult = await SceneApi.getSceneDetail(targetSceneId)
    if (detailResult.code === 200 && detailResult.data) {
      currentScene.value = detailResult.data
      await initScene(detailResult.data)
      await loadHotspotsForScene(targetSceneId)
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
        reloadHotspots(hotspots)
      }
    }
  } catch (e) {
    console.error('Failed to load hotspots:', e)
  }
}

const handleMarkerClick = (hotspot: HotspotForViewer) => {
  if (editMode.value) {
    editorMode.value = 'edit'
    editorForm.value = mapHotspotToForm(hotspot)
    editorVisible.value = true
    return
  }

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

function mapHotspotToForm(hotspot: HotspotForViewer): HotspotEditorForm {
  const iconSource: 'preset' | 'custom' =
    hotspot.icon_url && /^https?:|^data:/.test(hotspot.icon_url) ? 'custom' : 'preset'
  return {
    id: hotspot.id,
    scene_id: hotspot.scene_id,
    type: hotspot.type as 1 | 2 | 3,
    pitch: hotspot.pitch,
    yaw: hotspot.yaw,
    title: hotspot.title,
    icon_source: iconSource,
    icon_preset_key: iconSource === 'preset' ? (hotspot.style || 'arrow-blue') : undefined,
    icon_url: iconSource === 'custom' ? hotspot.icon_url : undefined,
    target_scene_id: hotspot.target_scene_id,
    content: hotspot.content,
    media_type: hotspot.media_type as 'image' | 'text' | undefined,
    media_url: hotspot.media_url,
    transition_effect: hotspot.transition_effect as 'fade' | 'zoom' | 'none' | undefined,
  }
}

setOnHotspotClick(handleMarkerClick)

setOnViewerClick(({ pitch, yaw }) => {
  if (!editMode.value || !currentSceneId.value) return
  editorMode.value = 'create'
  editorForm.value = {
    scene_id: currentSceneId.value,
    type: 1,
    pitch,
    yaw,
    title: '',
    icon_source: 'preset',
    icon_preset_key: 'arrow-blue',
    transition_effect: 'fade',
  }
  editorVisible.value = true
})

const onEditorSubmitted = async () => {
  editorVisible.value = false
  if (currentSceneId.value) {
    await loadHotspotsForScene(currentSceneId.value)
  }
}

const handleSceneSelect = async (sceneId: number) => {
  try {
    const result = await SceneApi.getSceneDetail(sceneId)
    if (result.code === 200 && result.data) {
      currentScene.value = result.data
      currentSceneId.value = sceneId
      router.push({ query: { id: sceneId, scene: result.data.scene_code } })
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

const toggleMinimap = () => {
  minimapVisible.value = !minimapVisible.value
}

const goBack = () => {
  const from = sessionStorage.getItem('panoramaFrom')
  if (from) {
    sessionStorage.removeItem('panoramaFrom')
    router.push(from)
  } else {
    router.push('/index')
  }
}

watch(() => route.query.scene, () => {
  if (route.query.scene) {
    loadScene()
  }
})

watch(editMode, async (val) => {
  if (val && currentScene.value?.space_id) {
    try {
      // 后端 page_size 上限 100，分页拉取全部场景
      const pageSize = 100
      const allScenes: SceneListItem[] = []
      let page = 1
      let hasMore = true
      while (hasMore) {
        const result = await SceneApi.getSceneList({
          space_id: currentScene.value.space_id,
          page,
          page_size: pageSize,
        })
        if (result.code === 200 && result.data) {
          const scenes = result.data.scenes || []
          allScenes.push(...scenes)
          const total = result.data.page_info?.total ?? 0
          if (allScenes.length >= total || scenes.length === 0) {
            hasMore = false
          } else {
            page++
          }
        } else {
          hasMore = false
        }
      }
      sceneList.value = allScenes
    } catch (e) {
      console.error('Failed to load scene list:', e)
    }
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

.panorama-canvas.edit-mode {
  cursor: crosshair;
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