import { ref, onUnmounted } from 'vue'
import { Scene as L7Scene } from '@antv/l7'
import { GaodeMap } from '@antv/l7-maps'
import { loadAmapScript } from '@/utils/amap-loader'

export type BasemapType = 'vector' | 'satellite'

export interface L7SceneOptions {
  container: HTMLElement
  center: [number, number]
  zoom?: number
  pitch?: number
  style?: BasemapType
}

export function useL7Scene() {
  const isReady = ref(false)
  const error = ref<string | null>(null)

  let l7Scene: L7Scene | null = null
  let aMapInstance: any = null

  async function initScene(options: L7SceneOptions): Promise<L7Scene> {
    try {
      // 先加载高德地图 SDK
      await loadAmapScript()

      const styleMap: Record<BasemapType, string> = {
        vector: 'dark',
        satellite: 'satellite',
      }

      l7Scene = new L7Scene({
        id: options.container as unknown as string | HTMLDivElement,
        map: new GaodeMap({
          style: styleMap[options.style || 'vector'],
          center: options.center,
          zoom: options.zoom || 15,
          pitch: options.pitch || 0,
        }),
        logoVisible: false,
      })

      return new Promise((resolve) => {
        // 尝试获取 AMap 实例的方法
        const tryGetAMapInstance = (): any => {
          // 方法1：通过 L7 mapService
          const mapService = (l7Scene as any).mapService
          if (mapService?.map) {
            return mapService.map
          }
          // 方法2：遍历 mapService 属性查找
          if (mapService) {
            for (const key of Object.keys(mapService)) {
              const val = (mapService as any)[key]
              if (val && typeof val.setMapStyle === 'function') {
                return val
              }
            }
          }
          // 方法3：通过 DOM 查找 AMap 实例
          const container = options.container
          if (container) {
            const amapKey = Object.keys(container).find(k => k.startsWith('__amap__'))
            if (amapKey) {
              return (container as any)[amapKey]
            }
          }
          return null
        }

        // 定时尝试获取 AMap 实例
        const aMapReady = setInterval(() => {
          const instance = tryGetAMapInstance()
          if (instance) {
            clearInterval(aMapReady)
            aMapInstance = instance
            console.log('Got AMap instance via polling:', !!aMapInstance)
            isReady.value = true
            resolve(l7Scene!)
          }
        }, 500)

        // 监听 L7 loaded 事件（备用）
        l7Scene!.on('loaded', () => {
          clearInterval(aMapReady)
          if (!aMapInstance) {
            aMapInstance = tryGetAMapInstance()
          }
          console.log('L7 loaded event fired, aMapInstance:', !!aMapInstance)
          isReady.value = true
          resolve(l7Scene!)
        })

        // 超时保护
        setTimeout(() => {
          if (!isReady.value) {
            clearInterval(aMapReady)
            console.warn('Timeout: forcing map ready')
            isReady.value = true
            resolve(l7Scene!)
          }
        }, 10000)
      })
    } catch (e: any) {
      error.value = e.message
      throw e
    }
  }

  function setBasemap(type: BasemapType) {
    if (!l7Scene) return

    const styleMap: Record<BasemapType, string> = {
      vector: 'dark',
      satellite: 'satellite',
    }

    if (typeof (l7Scene as any).setMapStyle === 'function') {
      ;(l7Scene as any).setMapStyle(styleMap[type])
    } else if (aMapInstance) {
      const amapStyleMap: Record<BasemapType, string> = {
        vector: 'amap://styles/dark',
        satellite: 'amap://satellite',
      }
      aMapInstance.setMapStyle(amapStyleMap[type])
    }
  }

  function setFeatures(features: string[]) {
    if (aMapInstance && typeof aMapInstance.setFeatures === 'function') {
      aMapInstance.setFeatures(features)
    }
  }

  function setCenter(lng: number, lat: number, zoom?: number) {
    if (aMapInstance) {
      if (zoom !== undefined) {
        aMapInstance.setZoomAndCenter(zoom, [lng, lat], false, 500)
      } else {
        aMapInstance.setCenter([lng, lat])
      }
    }
  }

  function destroy() {
    if (l7Scene) {
      try {
        l7Scene.destroy()
      } catch (e) {
        console.error('Error destroying scene:', e)
      }
      l7Scene = null
    }
    if (aMapInstance) {
      try {
        aMapInstance.destroy()
      } catch (e) {
        console.error('Error destroying map:', e)
      }
      aMapInstance = null
    }
    isReady.value = false
  }

  onUnmounted(() => {
    destroy()
  })

  return {
    isReady,
    error,
    initScene,
    setBasemap,
    setFeatures,
    setCenter,
    destroy,
    getScene: () => l7Scene,
    getAMapInstance: () => aMapInstance,
  }
}
