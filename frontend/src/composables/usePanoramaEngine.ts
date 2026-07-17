import { ref, type Ref } from 'vue'
import { Viewer } from '@photo-sphere-viewer/core'
import { CubemapTilesAdapter } from '@photo-sphere-viewer/cubemap-tiles-adapter'
import { MarkersPlugin } from '@photo-sphere-viewer/markers-plugin'
import '@photo-sphere-viewer/core/index.css'
import '@photo-sphere-viewer/markers-plugin/index.css'
import type { SceneDetailResponse } from '@/models/scene.model'
import type { HotspotForViewer } from '@/models/hotspot.model'
import { getHotspotIconHtml } from '@/utils/hotspotIcons'
import { TilePreloader, getBaseTileUrls, getTileUrl } from '@/utils/tileLoader'

const baseApiUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'

export function usePanoramaEngine(container: Ref<HTMLElement | null>) {
  const loading = ref(true)
  const error = ref('')
  const viewerReady = ref(false)

  let viewer: Viewer | null = null
  let markersPlugin: MarkersPlugin | null = null
  let tilePreloader: TilePreloader | null = null
  let onHotspotClickCallback: ((hotspot: HotspotForViewer) => void) | null = null
  let onViewerClickCallback: ((pos: { pitch: number; yaw: number }) => void) | null = null

  let readyHandler: (() => void) | null = null
  let panoramaErrorHandler: (() => void) | null = null
  let selectMarkerHandler: ((e: any) => void) | null = null
  let clickHandler: ((e: any) => void) | null = null

  const initScene = async (scene: SceneDetailResponse, hotspots?: HotspotForViewer[]): Promise<void> => {
    if (!container.value) return

    destroy()

    loading.value = true
    error.value = ''

    const sceneCode = scene.scene_code
    const useCubemap = (scene.slice_status === 'ready' || scene.slice_status === 'completed') && !!sceneCode

    try {
      let panoramaConfig: any
      let adapterConfig: typeof CubemapTilesAdapter | null = null

      if (useCubemap) {
        const baseTileUrls = getBaseTileUrls(sceneCode)
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
          const apiLevel = level + 1
          return getTileUrl(sceneCode, mappedFace, apiLevel, col, row)
        }

        adapterConfig = CubemapTilesAdapter
        panoramaConfig = {
          baseUrl: {
            left: baseTileUrls.left,
            front: baseTileUrls.front,
            right: baseTileUrls.right,
            back: baseTileUrls.back,
            top: baseTileUrls.top,
            bottom: baseTileUrls.bottom,
          },
          flipTopBottom: true,
          levels: [
            { faceSize: 2048, nbTiles: 4 },
            { faceSize: 4096, nbTiles: 8 },
          ],
          tileUrl,
        }
        tilePreloader = new TilePreloader(sceneCode)
      } else {
        const panoramaUrl = scene.tile_url || scene.source_url
        if (!panoramaUrl) {
          error.value = '全景图地址不存在'
          loading.value = false
          return
        }
        panoramaConfig = panoramaUrl
      }

      const viewerConfig: any = {
        container: container.value,
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
      markersPlugin = viewer.getPlugin(MarkersPlugin) as unknown as MarkersPlugin

      readyHandler = () => {
        loading.value = false
        viewerReady.value = true

        if (tilePreloader) {
          tilePreloader.preloadLevel(1).catch(console.error)
        }

        if (hotspots && hotspots.length > 0) {
          loadHotspots(hotspots)
        }
      }
      viewer.addEventListener('ready', readyHandler)

      panoramaErrorHandler = () => {
        if (sceneCode) {
          initViewerWithFallback(`${baseApiUrl}/api/v1/res/sources/${sceneCode}`, scene, hotspots)
        } else {
          error.value = '全景图加载失败，请重新上传或切片'
          loading.value = false
        }
      }
      viewer.addEventListener('panorama-error', panoramaErrorHandler)

      clickHandler = (e: any) => {
        if (!onViewerClickCallback) return
        const { pitch, yaw, rightclick } = e.data
        if (rightclick) return
        onViewerClickCallback({ pitch, yaw })
      }
      viewer.addEventListener('click', clickHandler)

      if (markersPlugin) {
        selectMarkerHandler = (e: any) => {
          const hotspot = e.marker.config.data as HotspotForViewer
          if (hotspot && onHotspotClickCallback) {
            onHotspotClickCallback(hotspot)
          }
        }
        markersPlugin.addEventListener('select-marker', selectMarkerHandler)
      }

    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : '初始化失败'
      loading.value = false
    }
  }

  const initViewerWithFallback = async (
    panoramaUrl: string,
    scene: SceneDetailResponse,
    hotspots?: HotspotForViewer[]
  ): Promise<void> => {
    if (!container.value) return

    destroy()

    try {
      const viewerConfig: any = {
        container: container.value,
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
      markersPlugin = viewer.getPlugin(MarkersPlugin) as unknown as MarkersPlugin

      readyHandler = () => {
        loading.value = false
        viewerReady.value = true

        if (hotspots && hotspots.length > 0) {
          loadHotspots(hotspots)
        }
      }
      viewer.addEventListener('ready', readyHandler)

      panoramaErrorHandler = () => {
        error.value = '全景图加载失败'
        loading.value = false
      }
      viewer.addEventListener('panorama-error', panoramaErrorHandler)

      clickHandler = (e: any) => {
        if (!onViewerClickCallback) return
        const { pitch, yaw, rightclick } = e.data
        if (rightclick) return
        onViewerClickCallback({ pitch, yaw })
      }
      viewer.addEventListener('click', clickHandler)

      if (markersPlugin) {
        selectMarkerHandler = (e: any) => {
          const hotspot = e.marker.config.data as HotspotForViewer
          if (hotspot && onHotspotClickCallback) {
            onHotspotClickCallback(hotspot)
          }
        }
        markersPlugin.addEventListener('select-marker', selectMarkerHandler)
      }

    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : '初始化失败'
      loading.value = false
    }
  }

  const buildMarkerConfig = (hotspot: HotspotForViewer): any => ({
    id: `hotspot-${hotspot.id}`,
    position: { pitch: hotspot.pitch, yaw: hotspot.yaw },
    tooltip: hotspot.title,
    data: hotspot,
    html: getHotspotIconHtml(hotspot),
    size: { width: 40, height: 40 },
    anchor: 'center center',
  })

  const loadHotspots = (hotspots: HotspotForViewer[]): void => {
    if (!markersPlugin) return
    hotspots.forEach((hotspot: HotspotForViewer) => {
      markersPlugin?.addMarker(buildMarkerConfig(hotspot))
    })
  }

  const reloadHotspots = (hotspots: HotspotForViewer[]): void => {
    markersPlugin?.clearMarkers()
    loadHotspots(hotspots)
  }

  const clearHotspots = (): void => {
    markersPlugin?.clearMarkers()
  }

  const zoomIn = (): void => {
    if (!viewer) return
    const currentZoom = viewer.getZoomLevel()
    viewer.zoom(currentZoom + 10)
  }

  const zoomOut = (): void => {
    if (!viewer) return
    const currentZoom = viewer.getZoomLevel()
    viewer.zoom(currentZoom - 10)
  }

  const toggleAutoRotate = (enabled: boolean): void => {
    if (viewer) {
      if (enabled) {
        (viewer as any).startAutoRotate?.({ speed: 0.5 })
      } else {
        (viewer as any).stopAutoRotate?.()
      }
    }
  }

  const resetCamera = (pitch?: number, yaw?: number, fov?: number): void => {
    if (!viewer) return
    viewer.animate({
      pitch: pitch || 0,
      yaw: yaw || 0,
      zoom: fov || 50,
      speed: '8rpm'
    })
  }

  const setOnHotspotClick = (callback: (hotspot: HotspotForViewer) => void): void => {
    onHotspotClickCallback = callback
  }

  const setOnViewerClick = (callback: (pos: { pitch: number; yaw: number }) => void): void => {
    onViewerClickCallback = callback
  }

  const getViewer = (): Viewer | null => viewer

  const preloadNextLevel = (): void => {
    if (tilePreloader) {
      tilePreloader.preloadLevel(1).catch(console.error)
    }
  }

  const destroy = (): void => {
    if (viewer) {
      if (readyHandler) {
        viewer.removeEventListener('ready', readyHandler)
      }
      if (panoramaErrorHandler) {
        viewer.removeEventListener('panorama-error', panoramaErrorHandler)
      }
      if (clickHandler) {
        viewer.removeEventListener('click', clickHandler)
      }
      viewer.destroy()
      viewer = null
    }
    if (markersPlugin) {
      if (selectMarkerHandler) {
        markersPlugin.removeEventListener('select-marker', selectMarkerHandler)
      }
      markersPlugin = null
    }
    if (tilePreloader) {
      tilePreloader.clear()
      tilePreloader = null
    }
    readyHandler = null
    panoramaErrorHandler = null
    selectMarkerHandler = null
    clickHandler = null
    viewerReady.value = false
  }

  return {
    loading,
    error,
    viewerReady,
    initScene,
    destroy,
    loadHotspots,
    reloadHotspots,
    clearHotspots,
    zoomIn,
    zoomOut,
    toggleAutoRotate,
    resetCamera,
    setOnHotspotClick,
    setOnViewerClick,
    getViewer,
    preloadNextLevel,
  }
}