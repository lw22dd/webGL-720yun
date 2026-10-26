import { Scene as L7Scene } from '@antv/l7'
import { GaodeMap } from '@antv/l7-maps'
import { loadAmapScript } from '@/utils/amap-loader'
import { resourceMonitor } from '@/utils/resourceMonitor'

export interface MapHandle {
  l7Scene: L7Scene
  aMapInstance: any
  isNewInstance: boolean
  isSpaceChanged: boolean
  spaceId: number
  normalLayer: any
  satelliteLayer: any
}

let persistentContainer: HTMLDivElement | null = null
let hiddenHolder: HTMLDivElement | null = null
let l7Scene: L7Scene | null = null
let aMapInstance: any = null
let currentSpaceId: number | null = null
let normalLayer: any = null
let satelliteLayer: any = null
let refCount = 0
let idleTimerId: ReturnType<typeof setTimeout> | null = null
let initPromise: Promise<MapHandle> | null = null

/**
 * 地图和全景各自持有一份 WebGL 上下文，双方同时存活时显存是叠加的。
 * 释放后只保留很短的一段宽限期（够用户误点返回即可），避免离开地图后
 * 还挂着 2 分钟的瓦片缓存和 GPU 资源。
 */
const IDLE_TIMEOUT = 10_000

/** forceDestroy 时自增，用来让在途的异步初始化作废（否则会复活已销毁的地图） */
let instanceToken = 0

function _ensureHiddenContainer() {
  if (!hiddenHolder) {
    hiddenHolder = document.createElement('div')
    hiddenHolder.style.cssText =
      'visibility:hidden;position:fixed;top:0;left:0;width:100vw;height:100vh;z-index:-1;pointer-events:none'
    hiddenHolder.setAttribute('data-map-hidden-holder', '')
    document.body.appendChild(hiddenHolder)
  }
  if (!persistentContainer) {
    persistentContainer = document.createElement('div')
    persistentContainer.style.cssText =
      'width:100%;height:100%;position:absolute;top:0;left:0'
    persistentContainer.setAttribute('data-map-persistent-container', '')
  }
}

function _moveContainerTo(targetEl: HTMLElement) {
  if (persistentContainer && persistentContainer.parentElement !== targetEl) {
    targetEl.appendChild(persistentContainer)
  }
}

function _moveContainerToHidden() {
  _ensureHiddenContainer()
  if (persistentContainer && hiddenHolder && persistentContainer.parentElement !== hiddenHolder) {
    hiddenHolder.appendChild(persistentContainer)
  }
}

function _waitForAMapInstance(scene: L7Scene): Promise<any> {
  return new Promise((resolve) => {
    const tryGet = (): any => {
      const mapService = (scene as any).mapService
      if (!mapService) return null
      if (mapService.map) return mapService.map
      for (const key of Object.keys(mapService)) {
        const val = (mapService as any)[key]
        if (val && typeof val.setMapStyle === 'function') return val
      }
      return null
    }

    const existing = tryGet()
    if (existing) { resolve(existing); return }

    let settled = false
    let timer: ReturnType<typeof setInterval> | null = null
    let deadline: ReturnType<typeof setTimeout> | null = null

    const finish = (value: any) => {
      if (settled) return
      settled = true
      if (timer) clearInterval(timer)
      if (deadline) clearTimeout(deadline)
      resolve(value)
    }

    timer = setInterval(() => {
      const inst = tryGet()
      if (inst) finish(inst)
    }, 500)

    deadline = setTimeout(() => finish(tryGet()), 10000)

    scene.on('loaded', () => finish(tryGet()))
  })
}

function _cancelIdleTimer() {
  if (idleTimerId) {
    clearTimeout(idleTimerId)
    idleTimerId = null
  }
}

/**
 * AMap / L7 的 destroy() 只拆 DOM，不会丢弃 WebGL 上下文（和 PSV 一样）。
 * 它们没暴露 renderer，只能在销毁前从容器里把 canvas 收集起来，销毁后逐个
 * loseContext。否则上下文要等 GC 才回收，长会话演示时很容易堆到浏览器的
 * 上下文上限（约 16 个），触发 “oldest context will be lost”。
 */
function _collectCanvases(root: HTMLElement | null): HTMLCanvasElement[] {
  if (!root) return []
  return Array.from(root.querySelectorAll('canvas'))
}

function _releaseCanvasContexts(canvases: HTMLCanvasElement[]): void {
  canvases.forEach((canvas) => {
    try {
      // 同一个 canvas 只能有一种上下文，两次探测里必有一次命中
      const gl = (canvas.getContext('webgl2') ??
        canvas.getContext('webgl')) as WebGLRenderingContext | null
      gl?.getExtension('WEBGL_lose_context')?.loseContext()
    } catch {
      // 上下文已丢失或已被回收，无需处理
    }
  })
}

function _startIdleTimer() {
  _cancelIdleTimer()
  idleTimerId = setTimeout(() => {
    console.log('[mapManager] Idle timeout, destroying cached map')
    forceDestroy()
  }, IDLE_TIMEOUT)
}

export async function acquire(
  targetEl: HTMLElement,
  spaceId: number,
  center: [number, number],
  zoom?: number,
): Promise<MapHandle> {
  _cancelIdleTimer()

  // Concurrent init: return same promise
  if (initPromise) {
    return initPromise.then((handle) => {
      const changed = handle.spaceId !== spaceId
      if (changed) {
        if (l7Scene) l7Scene.removeAllMarkers()
        if (aMapInstance) {
          aMapInstance.setCenter(center)
          if (zoom) aMapInstance.setZoom(zoom)
        }
        currentSpaceId = spaceId
      }
      _moveContainerTo(targetEl)
      resize()
      refCount++
      return {
        l7Scene: l7Scene!,
        aMapInstance: aMapInstance!,
        isNewInstance: false,
        isSpaceChanged: changed,
        spaceId,
        normalLayer,
        satelliteLayer,
      }
    })
  }

  // Map already cached
  if (l7Scene && aMapInstance) {
    const isSpaceChanged = currentSpaceId !== spaceId
    if (isSpaceChanged) {
      l7Scene.removeAllMarkers()
      aMapInstance.setCenter(center)
      if (zoom) aMapInstance.setZoom(zoom)
      currentSpaceId = spaceId
    }
    _moveContainerTo(targetEl)
    resize()
    refCount++
    return {
      l7Scene,
      aMapInstance,
      isNewInstance: false,
      isSpaceChanged,
      spaceId,
      normalLayer,
      satelliteLayer,
    }
  }

  // First-ever creation
  const token = ++instanceToken
  const promise = (async () => {
    _ensureHiddenContainer()
    _moveContainerTo(targetEl)

    await loadAmapScript()
    if (token !== instanceToken) {
      throw new Error('地图初始化已取消')
    }

    const scene = new L7Scene({
      id: persistentContainer!,
      map: new GaodeMap({
        style: 'normal',
        center,
        zoom: zoom || 15,
        pitch: 0,
      }),
      logoVisible: false,
    })
    l7Scene = scene

    const mapInstance = await _waitForAMapInstance(scene)
    if (token !== instanceToken) {
      // 等待期间被 forceDestroy 了，别把引用写回模块状态
      throw new Error('地图初始化已取消')
    }
    aMapInstance = mapInstance

    if (aMapInstance) {
      const layers = aMapInstance.getLayers()
      if (layers && layers.length > 0) {
        normalLayer = layers[0]
      }
      if (typeof aMapInstance.setFeatures === 'function') {
        aMapInstance.setFeatures(['bg', 'road', 'building'])
      }
    }

    currentSpaceId = spaceId
    refCount++

    console.log('[mapManager] Map created, spaceId:', spaceId)
    resourceMonitor.report('map:create', { spaceId })

    return {
      l7Scene: scene,
      aMapInstance: aMapInstance!,
      isNewInstance: true,
      isSpaceChanged: false,
      spaceId,
      normalLayer,
      satelliteLayer,
    } as MapHandle
  })()

  initPromise = promise
  try {
    return await promise
  } finally {
    // 只清理自己这一轮的 promise，避免把新一轮初始化干掉
    if (initPromise === promise) {
      initPromise = null
    }
  }
}

export function release() {
  refCount = Math.max(0, refCount - 1)
  // 已经 forceDestroy 过的场景不要再起一个空转的 idle 定时器
  if (refCount === 0 && l7Scene) {
    _moveContainerToHidden()
    _startIdleTimer()
    console.log('[mapManager] Map released to hidden holder, idle timer started')
  }
}

export function forceDestroy() {
  _cancelIdleTimer()
  instanceToken += 1

  const hadInstance = !!l7Scene || !!aMapInstance
  const canvases = _collectCanvases(persistentContainer)

  if (aMapInstance) {
    try { aMapInstance.destroy() } catch (e) { console.warn('[mapManager] Error destroying AMap:', e) }
    aMapInstance = null
  }
  if (l7Scene) {
    try { l7Scene.destroy() } catch (e) { console.warn('[mapManager] Error destroying L7:', e) }
    l7Scene = null
  }
  if (persistentContainer) {
    persistentContainer.remove()
    persistentContainer = null
  }
  _releaseCanvasContexts(canvases)
  normalLayer = null
  satelliteLayer = null
  currentSpaceId = null
  refCount = 0
  initPromise = null

  if (hadInstance) {
    console.log('[mapManager] Map force destroyed')
    // 推迟一拍：webglcontextlost 是异步派发的，立刻采样会多算已释放的上下文
    setTimeout(() => resourceMonitor.report('map:destroy'), 0)
  }
}

export function setNormalLayer(layer: any) {
  normalLayer = layer
}

export function setSatelliteLayer(layer: any) {
  satelliteLayer = layer
}

export function resize() {
  if (aMapInstance) {
    requestAnimationFrame(() => {
      setTimeout(() => {
        aMapInstance?.resize?.()
      }, 50)
    })
  }
}

export function getScene(): L7Scene | null {
  return l7Scene
}

export function getAMapInstance(): any {
  return aMapInstance
}

export function getCurrentSpaceId(): number | null {
  return currentSpaceId
}

export function isCached(): boolean {
  return l7Scene !== null && aMapInstance !== null
}
