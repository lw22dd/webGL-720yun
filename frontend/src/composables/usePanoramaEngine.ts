import { ref, type Ref } from 'vue'
import { Cache, Viewer } from '@photo-sphere-viewer/core'
import { CubemapTilesAdapter } from '@photo-sphere-viewer/cubemap-tiles-adapter'
import { MarkersPlugin } from '@photo-sphere-viewer/markers-plugin'
import '@photo-sphere-viewer/core/index.css'
import '@photo-sphere-viewer/markers-plugin/index.css'
import type { SceneDetailResponse } from '@/models/scene.model'
import type { HotspotForViewer } from '@/models/hotspot.model'
import { getHotspotIconHtml } from '@/utils/hotspotIcons'
import {
  createTileUrlResolver,
  getBaseTileUrls,
  resolveTileLevels,
  toViewerLevels,
  type CubeTileLevel,
} from '@/utils/tileLoader'
import { resourceMonitor } from '@/utils/resourceMonitor'

const baseApiUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'

/**
 * PSV 的全局图片缓存默认是「最多 10 个全景、10 分钟不访问才清理」。
 * 它缓存的是已解码的 HTMLImageElement，内存 ≈ 宽×高×4B，而且完全不计入
 * performance.memory —— 一个 8K 场景把所有瓦片看一遍就是 100~200MB，
 * 20 个场景叠起来就是 GB 级。这才是长会话真正的内存大头，比 GPU 纹理大得多。
 *
 * 压到 2 个全景 / 60 秒：当前场景 + 上一个（供切换淡入）够用；
 * 重新访问时后端返回 immutable 强缓存，不会真的再走 MinIO。
 */
Cache.maxItems = 2
Cache.ttl = 60

/** 当前 PSV 图片缓存里还剩多少个全景、多少张图，写进诊断日志 */
const psvCacheStats = (): Record<string, number> => {
  const keys = Object.keys(Cache.items)
  let images = 0
  keys.forEach((key) => {
    images += Object.keys(Cache.items[key].files).length
  })
  return { cachePanoramas: keys.length, cacheImages: images }
}

/**
 * 设计目标：最大放大时屏幕上的 1 个像素对应源图 1 个像素（1:1）。
 *
 * 注意 PSV 的 minFov/maxFov 是**垂直**视野（源码 state.vFov = zoomLevelToFov），
 * 水平视野要按宽高比换算：hFov = 2·atan(tan(vFov/2) × aspect)。
 * 立方体一面覆盖 90°、有 faceSize 个像素，所以屏幕可用的水平源像素 = faceSize × hFov / 90。
 * 1920×1080（aspect 1.778）下，1:1 需要的面尺寸：
 *   minFov 30° → hFov 50.9° → 需 3392 面
 *   minFov 45° → hFov 72.7° → 需 2377 面
 *   minFov 57° → hFov 88.0° → 需 1964 面
 *
 * 取 45°：8K 源图（1963 面）下是 1.21× 插值，仍可接受；16K 源图（4096 面）是 1.72× 超采样。
 * 想要「绝不插值」就把 minFov 提告到 57°（代价是能放大的幅度变小）。
 */
const MIN_FOV = 45
const MAX_FOV = 90

export function usePanoramaEngine(container: Ref<HTMLElement | null>) {
  const loading = ref(true)
  const error = ref('')
  const viewerReady = ref(false)

  let viewer: Viewer | null = null
  let markersPlugin: MarkersPlugin | null = null
  let tileLevels: CubeTileLevel[] = []
  let onHotspotClickCallback: ((hotspot: HotspotForViewer) => void) | null = null
  let onViewerClickCallback: ((pos: { pitch: number; yaw: number }) => void) | null = null

  let readyHandler: (() => void) | null = null
  let panoramaErrorHandler: (() => void) | null = null
  let selectMarkerHandler: ((e: any) => void) | null = null
  let clickHandler: ((e: any) => void) | null = null

  /** 每次 destroy 自增，让在途的异步回调（延迟 fallback、预热）失效 */
  let generation = 0
  let fallbackScheduled = false

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
    hotspots.forEach((hotspot) => {
      markersPlugin?.addMarker(buildMarkerConfig(hotspot))
    })
  }

  /** cubemap 与等距柱状图两条路径共用的 viewer 装配 */
  const mountViewer = (
    panorama: unknown,
    scene: SceneDetailResponse,
    hotspots: HotspotForViewer[] | undefined,
    options: { adapter?: typeof CubemapTilesAdapter; onPanoramaError: () => void }
  ): void => {
    const viewerConfig: any = {
      container: container.value,
      panorama,
      defaultZoomLvl: 0,
      defaultPitch: scene.initial_pitch || 0,
      defaultYaw: scene.initial_yaw || 0,
      minFov: MIN_FOV,
      maxFov: MAX_FOV,
      navbar: false,
      plugins: [[MarkersPlugin, {}]],
    }
    if (options.adapter) {
      viewerConfig.adapter = options.adapter
    }

    viewer = new Viewer(viewerConfig)
    markersPlugin = viewer.getPlugin(MarkersPlugin) as unknown as MarkersPlugin

    const gen = generation

    readyHandler = () => {
      if (gen !== generation) return
      loading.value = false
      viewerReady.value = true

      // 可见瓦片交给 adapter 按视锥自己加载即可；
      // 之前在这里整级预热会把当前档全部 96 张瓦片都拉一遍，
      // 既和 adapter 的首屏请求重复，也白白增加首屏流量。
      if (hotspots && hotspots.length > 0) {
        loadHotspots(hotspots)
      }
    }
    viewer.addEventListener('ready', readyHandler)

    panoramaErrorHandler = () => {
      if (gen !== generation) return
      options.onPanoramaError()
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
  }

  const initScene = async (
    scene: SceneDetailResponse,
    hotspots?: HotspotForViewer[]
  ): Promise<void> => {
    if (!container.value) return

    destroy()

    loading.value = true
    error.value = ''

    const sceneCode = scene.scene_code
    const useCubemap =
      (scene.slice_status === 'ready' || scene.slice_status === 'completed') && !!sceneCode

    try {
      if (!useCubemap) {
        // tile_url 存的是 MinIO 前缀而不是图片地址，且 MinIO 桶不开匿名读（实测 403），
        // 所以这里统一走后端的流式接口。
        const panoramaUrl = sceneCode
          ? `${baseApiUrl}/api/v1/res/sources/${sceneCode}`
          : scene.source_url
        if (!panoramaUrl) {
          error.value = '全景图地址不存在'
          loading.value = false
          return
        }

        mountViewer(panoramaUrl, scene, hotspots, {
          onPanoramaError: () => {
            error.value = '全景图加载失败'
            loading.value = false
          },
        })
        resourceMonitor.report('pano:create', { mode: 'equirect' })
        return
      }

      // 三档固定为 480p / 1080p / 4K，后端负责把源图重采样到对应面尺寸；
      // 首屏只加载 level 0 底图，可见瓦片由 adapter 按视锥拉 level 1，放大后再拉 level 2。
      tileLevels = resolveTileLevels()
      const expectedGeneration = generation

      const panoramaConfig: any = {
        baseUrl: getBaseTileUrls(sceneCode),
        flipTopBottom: true,
        levels: toViewerLevels(tileLevels),
        tileUrl: createTileUrlResolver(sceneCode, tileLevels),
      }

      mountViewer(panoramaConfig, scene, hotspots, {
        adapter: CubemapTilesAdapter,
        onPanoramaError: () => {
          if (fallbackScheduled) return
          fallbackScheduled = true
          // 不能在 panorama-error 的事件分发里同步 destroy + 重建：
          // 旧 renderer 还没释放，同一个容器上会瞬时存在两个 WebGL 上下文
          setTimeout(() => {
            fallbackScheduled = false
            if (expectedGeneration !== generation || !container.value) return
            initViewerWithFallback(
              `${baseApiUrl}/api/v1/res/sources/${sceneCode}`,
              scene,
              hotspots
            )
          }, 0)
        },
      })

      resourceMonitor.report('pano:create', {
        mode: 'cubemap',
        levels: tileLevels.map((level) => level.faceSize),
        sourceWidth: scene.source_width,
      })
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : '初始化失败'
      loading.value = false
    }
  }

  /** 切片资源不可用时，退回原始等距柱状图 */
  const initViewerWithFallback = async (
    panoramaUrl: string,
    scene: SceneDetailResponse,
    hotspots?: HotspotForViewer[]
  ): Promise<void> => {
    if (!container.value) return

    destroy()

    try {
      mountViewer(panoramaUrl, scene, hotspots, {
        onPanoramaError: () => {
          error.value = '全景图加载失败'
          loading.value = false
        },
      })
      resourceMonitor.report('pano:create', { mode: 'fallback' })
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : '初始化失败'
      loading.value = false
    }
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

  /**
   * PSV 的 renderer.destroy() 只 dispose three 侧的资源，不会真正丢弃 WebGL 上下文，
   * 上下文要等 canvas 被 GC 才回收。长会话里连续切 20+ 场景会撞上浏览器约 16 个
   * 上下文的上限，旧上下文被强制丢弃，表现出来就是黑屏/掉帧。
   * 这里显式 loseContext，让显存立刻归还。
   */
  const releaseWebGLContext = (target: Viewer): void => {
    const threeRenderer = (target as any).renderer?.renderer
    try {
      threeRenderer?.forceContextLoss?.()
    } catch {
      // 自定义 renderer 或上下文已经丢失时忽略
    }
  }

  const destroy = (): void => {
    generation += 1
    fallbackScheduled = false

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
      const target = viewer
      try {
        target.destroy()
      } catch (e) {
        console.warn('[panorama] destroy viewer failed:', e)
      } finally {
        // 即使 destroy 抛错也要把引用和 GPU 上下文清掉，否则下一个场景会叠加一份
        viewer = null
        releaseWebGLContext(target)
        // 只保留最近 1~2 个全景的解码图片，避免长会话里缓存无限堆积
        Cache.purge()
        // 推迟一拍再采样：webglcontextlost 是异步派发的，立刻采样会多算一个上下文
        setTimeout(() => resourceMonitor.report('pano:destroy', psvCacheStats()), 0)
      }
    }
    if (markersPlugin) {
      if (selectMarkerHandler) {
        markersPlugin.removeEventListener('select-marker', selectMarkerHandler)
      }
      markersPlugin = null
    }
    tileLevels = []
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
  }
}
