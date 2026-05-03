import AMapLoader from '@amap/amap-jsapi-loader'

const AMAP_KEY = import.meta.env.VITE_AMAP_KEY || '38e34e74351ae3554aee1a94fdee9c62'
const AMAP_SECURITY_KEY = import.meta.env.VITE_AMAP_SECURITY_KEY || '7c140de5a9222a71cc6894bb50300d13'

let amapInstance: any = null
let amapLoaderPromise: Promise<any> | null = null

export const initAMap = async (container: HTMLElement, options?: any) => {
  if (amapInstance) {
    return amapInstance
  }

  if (!amapLoaderPromise) {
    // 配置安全密钥 - JS API 2.0 必须在加载前配置
    if (AMAP_SECURITY_KEY) {
      window._AMapSecurityConfig = {
        securityJsCode: AMAP_SECURITY_KEY,
      }
    }

    amapLoaderPromise = AMapLoader.load({
      key: AMAP_KEY,
      version: '2.0',
      plugins: ['AMap.Geocoder', 'AMap.AutoComplete', 'AMap.PlaceSearch'],
      AMapUI: {
        version: '1.1',
        plugins: []
      }
    })
  }

  const AMap = await amapLoaderPromise

  amapInstance = new AMap.Map(container, {
    zoom: 5,
    center: [105.397428, 38.04623],
    mapStyle: 'amap://styles/darkblue',
    viewMode: '3D',
    pitch: 30,
    ...options
  })

  return amapInstance
}

export const createMarker = (map: any, position: [number, number], options?: any) => {
  const AMap = window.AMap
  if (!AMap) return null
  const marker = new AMap.Marker({
    position: position,
    anchor: 'bottom-center',
    ...options
  })
  marker.setMap(map)
  return marker
}

export const geocode = async (address: string): Promise<{ lng: number; lat: number } | null> => {
  const AMap = await amapLoaderPromise
  if (!AMap) return null
  return new Promise((resolve) => {
    const geocoder = new AMap.Geocoder()
    geocoder.getLocation(address, (status: string, result: any) => {
      if (status === 'complete' && result.geocodes.length) {
        const location = result.geocodes[0].location
        resolve({ lng: location.lng, lat: location.lat })
      } else {
        resolve(null)
      }
    })
  })
}

export const autoComplete = async (keyword: string): Promise<any[]> => {
  const AMap = await amapLoaderPromise
  if (!AMap) return []
  return new Promise((resolve) => {
    const auto = new AMap.AutoComplete({
      city: '全国'
    })
    auto.search(keyword, (status: string, result: any) => {
      if (status === 'complete' && result.tips) {
        resolve(result.tips.filter((tip: any) => tip.location))
      } else {
        resolve([])
      }
    })
  })
}

export const placeSearch = async (keyword: string): Promise<any[]> => {
  const AMap = await amapLoaderPromise
  if (!AMap) return []
  return new Promise((resolve) => {
    const ps = new AMap.PlaceSearch({
      pageSize: 10,
      pageIndex: 1,
      city: '全国'
    })
    ps.search(keyword, (status: string, result: any) => {
      if (status === 'complete' && result.poiList && result.poiList.pois) {
        resolve(result.poiList.pois)
      } else {
        resolve([])
      }
    })
  })
}

export const setMarkerAnimation = (marker: any, type: 'breathe' | 'glow' | 'none') => {
  if (type === 'breathe') {
    marker.setAnimation('AMAP_ANIMATION_BOUNCE')
  } else if (type === 'none') {
    marker.setAnimation('AMAP_ANIMATION_NONE')
  }
}

export const createCustomMarker = (_map: any, position: [number, number], content: HTMLElement) => {
  const AMap = window.AMap
  if (!AMap) return null
  return new AMap.Marker({
    position: position,
    content: content,
    anchor: 'bottom-center',
    offset: new AMap.Pixel(0, 0)
  })
}

export const lngLatToXY = (map: any, lng: number, lat: number) => {
  const pixel = map.lngLatToContainer([lng, lat])
  return { x: pixel.getX(), y: pixel.getY() }
}

export const xyToLngLat = (map: any, x: number, y: number) => {
  const lnglat = map.containerToLngLat([x, y])
  return { lng: lnglat.getLng(), lat: lnglat.getLat() }
}

export const setMapCenter = (map: any, lng: number, lat: number, zoom?: number) => {
  map.setCenter([lng, lat])
  if (zoom) {
    map.setZoom(zoom)
  }
}

export const getMapInstance = () => amapInstance

export const destroyMap = () => {
  if (amapInstance) {
    amapInstance.destroy()
    amapInstance = null
    amapLoaderPromise = null
  }
}

declare global {
  interface Window {
    AMap: any
    _AMapSecurityConfig?: {
      securityJsCode: string
    }
  }
}
