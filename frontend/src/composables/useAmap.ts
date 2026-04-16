import { ref, onUnmounted } from 'vue'
import { loadAmapScript } from '@/utils/amap-loader'

export function useAmap() {
  const isReady = ref(false)
  const error = ref<string | null>(null)
  let mapInstance: any = null

  async function initMap(
    container: HTMLElement,
    options: {
      center?: [number, number]
      zoom?: number
      style?: string
      pitch?: number
    } = {}
  ): Promise<any> {
    try {
      await loadAmapScript()

      const defaultOptions = {
        center: options.center || [116.397428, 39.90923],
        zoom: options.zoom || 15,
        style: options.style || 'amap://styles/dark',
        pitch: options.pitch || 0,
        viewMode: '2D',
      }

      mapInstance = new window.AMap.Map(container, defaultOptions)
      isReady.value = true

      return mapInstance
    } catch (e: any) {
      error.value = e.message
      throw e
    }
  }

  function setMapStyle(style: 'satellite' | 'dark' | 'normal') {
    if (!mapInstance) return

    const styleMap: Record<string, string> = {
      satellite: 'amap://satellite',
      dark: 'amap://styles/dark',
      normal: 'amap://styles/normal',
    }

    mapInstance.setMapStyle(styleMap[style] || style)
  }

  function setFeatures(features: string[]) {
    if (!mapInstance) return
    mapInstance.setFeatures(features)
  }

  function setCenter(lng: number, lat: number, zoom?: number) {
    if (!mapInstance) return
    if (zoom !== undefined) {
      mapInstance.setZoomAndCenter(zoom, [lng, lat])
    } else {
      mapInstance.setCenter([lng, lat])
    }
  }

  function destroy() {
    if (mapInstance) {
      mapInstance.destroy()
      mapInstance = null
    }
    isReady.value = false
  }

  onUnmounted(() => {
    destroy()
  })

  return {
    isReady,
    error,
    initMap,
    setMapStyle,
    setFeatures,
    setCenter,
    destroy,
    getMapInstance: () => mapInstance,
  }
}
