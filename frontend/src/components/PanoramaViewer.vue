<template>
  <div class="panorama-container">
    <div ref="containerRef" class="panorama-canvas"></div>
<<<<<<< HEAD
    <div class="panorama-controls">
      <t-button theme="primary" @click="toggleFullscreen">全屏</t-button>
      <t-button @click="resetCamera">重置视角</t-button>
    </div>
=======
    
    <div class="pano-topbar">
      <div class="pano-topbar-left">
        <t-button theme="default" variant="outline" @click="goBack" class="back-btn">
          ← 返回
        </t-button>
        <span class="pano-title">{{ sceneTitle }}</span>
      </div>
      <div class="pano-topbar-right">
        <t-button theme="default" variant="outline" @click="copyShareLink" class="action-btn">
          🔗 分享
        </t-button>
        <t-button theme="default" variant="outline" @click="toggleFullscreen" class="action-btn">
          ⛶ 全屏
        </t-button>
      </div>
    </div>

    <div class="pano-controls">
      <t-button theme="default" variant="outline" size="small" @click="resetCamera" title="重置视角">
        ⟲
      </t-button>
      <t-button 
        :theme="autoRotate ? 'primary' : 'default'" 
        variant="outline" 
        size="small" 
        @click="toggleAutoRotate"
        title="自动旋转"
      >
        ↻
      </t-button>
      <t-button theme="default" variant="outline" size="small" @click="zoomIn" title="放大">
        +
      </t-button>
      <t-button theme="default" variant="outline" size="small" @click="zoomOut" title="缩小">
        −
      </t-button>
    </div>

    <SceneStrip 
      v-if="scenes.length > 0"
      :scenes="scenes"
      :current-scene-id="currentSceneId"
      @scene-select="handleSceneSelect"
    />

    <div v-if="loading" class="pano-loading">
      <t-loading size="large" text="加载中..." />
    </div>

    <div v-if="error" class="pano-error">
      <t-icon name="error-circle" size="48" />
      <p>{{ error }}</p>
      <t-button theme="primary" @click="loadScene">重试</t-button>
    </div>

    <t-dialog
      v-model:visible="infoDialogVisible"
      :header="currentHotspot?.title"
      width="400px"
      :footer="false"
      attach=".panorama-container"
    >
      <div class="hotspot-info-content">
        <p>{{ currentHotspot?.content }}</p>
      </div>
    </t-dialog>

    <t-dialog
      v-model:visible="quizDialogVisible"
      :header="currentHotspot?.title || '知识问答'"
      width="500px"
      :footer="false"
      attach=".panorama-container"
    >
      <div class="hotspot-quiz-content">
        <p class="quiz-question">{{ currentHotspot?.question }}</p>
        <div class="quiz-options">
          <div 
            v-for="(option, index) in quizOptions" 
            :key="index"
            class="quiz-option"
            :class="{ 'selected': selectedQuizOption === index }"
            @click="selectedQuizOption = Number(index)"
          >
            {{ String.fromCharCode(65 + Number(index)) }}. {{ String(option) }}
          </div>
        </div>
        <t-button 
          theme="primary" 
          block 
          :disabled="selectedQuizOption === null"
          @click="submitQuiz"
        >
          提交答案
        </t-button>
      </div>
    </t-dialog>
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
  </div>
</template>

<script setup lang="ts">
<<<<<<< HEAD
import { ref, onMounted, onUnmounted } from 'vue';
import { Viewer } from '@photo-sphere-viewer/core';
import '@photo-sphere-viewer/core/index.css';

const containerRef = ref<HTMLDivElement | null>(null);
let viewer: Viewer | null = null;

const panoramaPath = new URL('@/assets/宝瓶口1.jpg', import.meta.url).href;

const initViewer = () => {
  if (!containerRef.value) return;

  viewer = new Viewer({
    container: containerRef.value,
    panorama: panoramaPath,
    defaultZoomLvl: 0,
    minFov: 30,
    maxFov: 90,
    navbar: false,
  });
};

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    containerRef.value?.requestFullscreen();
  } else {
    document.exitFullscreen();
  }
};

const resetCamera = () => {
  viewer?.rotate({ yaw: 0, pitch: 0 });
  viewer?.zoom(0);
};

onMounted(() => {
  initViewer();
});

onUnmounted(() => {
  viewer?.destroy();
});
=======
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { Viewer } from '@photo-sphere-viewer/core'
import { CubemapAdapter } from '@photo-sphere-viewer/cubemap-adapter'
import { CubemapTilesAdapter } from '@photo-sphere-viewer/cubemap-tiles-adapter'
import { MarkersPlugin } from '@photo-sphere-viewer/markers-plugin'
import '@photo-sphere-viewer/core/index.css'
import '@photo-sphere-viewer/markers-plugin/index.css'
import SceneApi from '@/services/api/scene.api'
import HotspotApi from '@/services/api/hotspot.api'
import SceneStrip from './SceneStrip.vue'
import type { SceneDetailResponse } from '@/models/scene.model'
import { TilePreloader, getCubemapUrls, getTileUrl } from '@/utils/tileLoader'


const baseApiUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'


const route = useRoute()
const router = useRouter()

const containerRef = ref<HTMLDivElement | null>(null)
const loading = ref(true)
const error = ref('')
const autoRotate = ref(false)

const currentScene = ref<SceneDetailResponse | null>(null)
const scenes = ref<any[]>([])
const currentSceneId = ref<number | null>(null)
const sceneTitle = computed(() => currentScene.value?.title || '全景漫游')

const infoDialogVisible = ref(false)
const quizDialogVisible = ref(false)
const currentHotspot = ref<any>(null)
const selectedQuizOption = ref<number | null>(null)
const quizOptions = computed(() => {
  if (!currentHotspot.value?.options) return []
  try {
    return JSON.parse(currentHotspot.value.options)
  } catch {
    return []
  }
})

let viewer: any = null
let markersPlugin: any = null
let tilePreloader: TilePreloader | null = null

const loadScene = async () => {
  const sceneCode = route.query.scene as string
  if (!sceneCode) {
    error.value = '未指定场景'
    loading.value = false
    return
  }

  try {
    loading.value = true
    error.value = ''

    const result = await SceneApi.getSceneList({ keyword: sceneCode, page: 1, page_size: 10 })
      console.log(result)

    if (result.code === 200 && result.data && result.data.scenes && result.data.scenes.length > 0) {
      const scene = result.data.scenes.find((s: any) => s.scene_code === sceneCode)
      if (scene) {
        currentSceneId.value = scene.id
        const detailResult = await SceneApi.getSceneDetail(scene.id)
        if (detailResult.code === 200 && detailResult.data) {
          currentScene.value = detailResult.data
          await initViewer()
          loadHotspots()
        }
      }
    } else {
      error.value = '场景不存在'
    }
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

const initViewer = async () => {
  if (!containerRef.value || !currentScene.value) return

  if (viewer) {
    viewer.destroy()
    viewer = null
  }

  if (tilePreloader) {
    tilePreloader.clear()
    tilePreloader = null
  }

  const scene = currentScene.value
  const sceneCode = scene.scene_code

  const useCubemap = (scene.slice_status === 'ready' || scene.slice_status === 'completed') && sceneCode

  try {
    let panoramaConfig: any
    let adapterConfig: any = null

    if (useCubemap) {
      const cubemapUrls = getCubemapUrls(sceneCode)
      const faceMap: Record<string, string> = {
        left: 'nx',
        front: 'pz',
        right: 'px',
        back: 'nz',
        top: 'py',
        bottom: 'ny'
      }
      const tileUrl = (face: string, col: number, row: number, level: number) => {
        const mappedFace = faceMap[face] || face
        return getTileUrl(sceneCode, mappedFace, level, col, row)
      }

      adapterConfig = CubemapTilesAdapter
      panoramaConfig = {
        baseUrl: {
          left: cubemapUrls.left,
          front: cubemapUrls.front,
          right: cubemapUrls.right,
          back: cubemapUrls.back,
          top: cubemapUrls.top,
          bottom: cubemapUrls.bottom,
        },
        flipTopBottom: true,
        levels: [
          { faceSize: 1024, nbTiles: 2 },
          { faceSize: 2048, nbTiles: 4 },
          { faceSize: 2048, nbTiles: 8 },
        ],
        tileUrl,
      }
      console.log('[Pano] 使用 CubemapTilesAdapter 瓦片渲染模式')
      tilePreloader = new TilePreloader(sceneCode)
    } else {
      const panoramaUrl = scene.tile_url || scene.source_url
      if (!panoramaUrl) {
        error.value = '全景图地址不存在'
        return
      }
      panoramaConfig = panoramaUrl
      console.log('[Pano] 使用 Equirectangular 全景图渲染模式')
    }

    const viewerConfig: any = {
      container: containerRef.value,
      panorama: panoramaConfig,
      defaultZoomLvl: 0,
      defaultPitch: scene.initial_pitch || 0,
      defaultYaw: scene.initial_yaw || 0,
      minFov: 30,
      maxFov: 90,
      navbar: false,
      plugins: [
        [MarkersPlugin, {}]
      ]
    }

    if (adapterConfig) {
      viewerConfig.adapter = adapterConfig
    }

    viewer = new Viewer(viewerConfig)

    markersPlugin = viewer.getPlugin(MarkersPlugin) as any

    viewer.addEventListener('ready', () => {
      loading.value = false
      
      if (tilePreloader) {
        tilePreloader.preloadLevel(1).catch(console.error)
      }
    })

    viewer.addEventListener('panorama-error', (err: any) => {
      console.error('[Pano] Panorama load error:', err)

      if (sceneCode) {
        console.log('[Pano] Fallback: 使用 source API 重试')
        initViewerWithFallback(`${baseApiUrl}/api/v1/res/sources/${sceneCode}`)
      } else {
        error.value = '全景图加载失败，请重新上传或切片'
        loading.value = false
      }
    })

    markersPlugin?.addEventListener('select-marker', ({ marker }: any) => {
      handleMarkerClick(marker.config.id)
    })

  } catch (e: any) {
    error.value = e.message || '初始化失败'
  }
}

const initViewerWithFallback = async (panoramaUrl: string) => {
  if (!containerRef.value || !currentScene.value) return

  if (viewer) {
    viewer.destroy()
    viewer = null
  }

  const scene = currentScene.value

  try {
    const viewerConfig: any = {
      container: containerRef.value,
      panorama: panoramaUrl,
      defaultZoomLvl: 0,
      defaultPitch: scene.initial_pitch || 0,
      defaultYaw: scene.initial_yaw || 0,
      minFov: 30,
      maxFov: 90,
      navbar: false,
      plugins: [
        [MarkersPlugin, {}]
      ]
    }

    viewer = new Viewer(viewerConfig)
    markersPlugin = viewer.getPlugin(MarkersPlugin) as any

    viewer.addEventListener('ready', () => {
      loading.value = false
    })

    viewer.addEventListener('panorama-error', () => {
      error.value = '全景图加载失败'
      loading.value = false
    })

    markersPlugin?.addEventListener('select-marker', ({ marker }: any) => {
      handleMarkerClick(marker.config.id)
    })

  } catch (e: any) {
    error.value = e.message || '初始化失败'
    loading.value = false
  }
}

const loadHotspots = async () => {
  if (!currentSceneId.value || !markersPlugin) return

  try {
    const result = await HotspotApi.getHotspotList({ 
      scene_id: currentSceneId.value,
      page: 1,
      page_size: 100
    })
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
    }
  } catch (e) {
    console.error('Failed to load hotspots:', e)
  }
}

const handleMarkerClick = (markerId: string) => {
  const hotspotId = parseInt(markerId.replace('hotspot-', ''))
  const hotspot = currentScene.value?.hotspots?.find((h: any) => h.id === hotspotId)
  
  if (!hotspot) return

  currentHotspot.value = hotspot

  switch (hotspot.type) {
    case 1:
      if ((hotspot as any).target_scene_id) {
        handleSceneSelect((hotspot as any).target_scene_id)
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

const handleSceneSelect = async (sceneId: number) => {
  try {
    loading.value = true
    const result = await SceneApi.getSceneDetail(sceneId)
    if (result.code === 200 && result.data) {
      currentScene.value = result.data
      currentSceneId.value = sceneId
      router.push({ query: { scene: result.data.scene_code } })
      await initViewer()
      loadHotspots()
    }
  } catch (e: any) {
    MessagePlugin.error(e.message || '切换场景失败')
  } finally {
    loading.value = false
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
  if (!viewer || !currentScene.value) return
  ;(viewer as any).setPose?.({
    pitch: currentScene.value.initial_pitch || 0,
    yaw: currentScene.value.initial_yaw || 0,
    zoom: currentScene.value.initial_fov || 50
  })
}

const toggleAutoRotate = () => {
  autoRotate.value = !autoRotate.value
  if (viewer) {
    if (autoRotate.value) {
      ;(viewer as any).startAutoRotate?.({ speed: 0.5 })
    } else {
      ;(viewer as any).stopAutoRotate?.()
    }
  }
}

const zoomIn = () => {
  if (!viewer) return
  const currentZoom = viewer.getZoomLevel()
  viewer.zoom(currentZoom + 10)
}

const zoomOut = () => {
  if (!viewer) return
  const currentZoom = viewer.getZoomLevel()
  viewer.zoom(currentZoom - 10)
}

const copyShareLink = () => {
  const link = window.location.href
  navigator.clipboard.writeText(link).then(() => {
    MessagePlugin.success('链接已复制')
  })
}

const submitQuiz = () => {
  if (selectedQuizOption.value === null || !currentHotspot.value) return
  
  const answer = currentHotspot.value.answer
  const options = quizOptions.value
  
  if (options[selectedQuizOption.value] === answer) {
    MessagePlugin.success('回答正确！')
  } else {
    MessagePlugin.warning('回答错误，正确答案是：' + answer)
  }
  
  quizDialogVisible.value = false
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
  if (viewer) {
    viewer.destroy()
    viewer = null
  }
  if (tilePreloader) {
    tilePreloader.clear()
    tilePreloader = null
  }
})
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
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

<<<<<<< HEAD
.panorama-controls {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 100;
  display: flex;
  gap: 10px;
=======
.pano-topbar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: linear-gradient(180deg, rgba(0,0,0,0.6) 0%, transparent 100%);
  z-index: 100;
}

.pano-topbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.pano-topbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.pano-title {
  color: #fff;
  font-size: 18px;
  font-weight: 500;
}

.back-btn,
.action-btn {
  background: rgba(255, 255, 255, 0.1) !important;
  border: 1px solid rgba(255, 255, 255, 0.2) !important;
  color: #fff !important;
  backdrop-filter: blur(10px);
}

.back-btn:hover,
.action-btn:hover {
  background: rgba(255, 255, 255, 0.2) !important;
}

.pano-controls {
  position: absolute;
  right: 24px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 8px;
  z-index: 100;
}

.pano-controls :deep(.t-button) {
  width: 40px;
  height: 40px;
  padding: 0;
  background: rgba(255, 255, 255, 0.1) !important;
  border: 1px solid rgba(255, 255, 255, 0.2) !important;
  color: #fff !important;
  backdrop-filter: blur(10px);
  font-size: 18px;
}

.pano-controls :deep(.t-button:hover) {
  background: rgba(255, 255, 255, 0.2) !important;
}

.pano-controls :deep(.t-button.t-button--theme-primary) {
  background: #0EA5E9 !important;
  border-color: #0EA5E9 !important;
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

.hotspot-info-content,
.hotspot-quiz-content {
  padding: 8px 0;
}

.hotspot-info-content p {
  color: #4e5969;
  line-height: 1.6;
  margin: 0;
}

.quiz-question {
  font-size: 16px;
  font-weight: 500;
  color: #1d2129;
  margin-bottom: 16px;
}

.quiz-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.quiz-option {
  padding: 12px 16px;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.quiz-option:hover {
  border-color: #0EA5E9;
  background: rgba(14, 165, 233, 0.05);
}

.quiz-option.selected {
  border-color: #0EA5E9;
  background: rgba(14, 165, 233, 0.1);
  color: #0EA5E9;
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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
}
</style>
