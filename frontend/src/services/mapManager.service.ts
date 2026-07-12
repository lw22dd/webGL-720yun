import { Scene as L7Scene } from '@antv/l7'
import { GaodeMap } from '@antv/l7-maps'
import { loadAmapScript } from '@/utils/amap-loader'

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
const IDLE_TIMEOUT = 120_000

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

function _waitForAMapInstance(): Promise<any> {
  return new Promise((resolve) => {
    if (!l7Scene) { resolve(null); return }

    const tryGet = (): any => {
      const mapService = (l7Scene as any).mapService
      if (mapService?.map) return mapService.map
      if (mapService) {
        for (const key of Object.keys(mapService)) {
          const val = (mapService as any)[key]
          if (val && typeof val.setMapStyle === 'function') return val
        }
      }
      return null
    }

    const existing = tryGet()
    if (existing) { resolve(existing); return }

    let timer: ReturnType<typeof setInterval> | null = setInterval(() => {
      const inst = tryGet()
      if (inst) {
        if (timer) { clearInterval(timer); timer = null }
        resolve(inst)
      }
    }, 500)

    l7Scene.on('loaded', () => {
      if (timer) { clearInterval(timer); timer = null }
      resolve(tryGet())
    })

    setTimeout(() => {
      if (timer) { clearInterval(timer); timer = null }
      resolve(tryGet())
    }, 10000)
  })
}

function _cancelIdleTimer() {
  if (idleTimerId) {
    clearTimeout(idleTimerId)
    idleTimerId = null
  }
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
  initPromise = (async () => {
    _ensureHiddenContainer()
    _moveContainerTo(targetEl)

    await loadAmapScript()

    l7Scene = new L7Scene({
      id: persistentContainer!,
      map: new GaodeMap({
        style: 'normal',
        center,
        zoom: zoom || 15,
        pitch: 0,
      }),
      logoVisible: false,
    })

    aMapInstance = await _waitForAMapInstance()

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

    return {
      l7Scene: l7Scene!,
      aMapInstance: aMapInstance!,
      isNewInstance: true,
      isSpaceChanged: false,
      spaceId,
      normalLayer,
      satelliteLayer,
    } as MapHandle
  })()

  try {
    const result = await initPromise
    return result
  } finally {
    initPromise = null
  }
}

export function release() {
  refCount = Math.max(0, refCount - 1)
  if (refCount === 0) {
    _moveContainerToHidden()
    _startIdleTimer()
    console.log('[mapManager] Map released to hidden holder, idle timer started')
  }
}

export function forceDestroy() {
  _cancelIdleTimer()

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
  normalLayer = null
  satelliteLayer = null
  currentSpaceId = null
  refCount = 0
  initPromise = null

  console.log('[mapManager] Map force destroyed')
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
