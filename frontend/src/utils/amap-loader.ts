const AMAP_KEY = import.meta.env.VITE_AMAP_KEY || ''
const AMAP_SECURITY_KEY = import.meta.env.VITE_AMAP_SECURITY_KEY || ''

let amapLoaded = false
let amapLoading = false
let amapCallbacks: (() => void)[] = []

declare global {
  interface Window {
    AMap: any
    _AMapSecurityConfig?: {
      securityJsCode: string
    }
  }
}

export function loadAmapScript(): Promise<void> {
  return new Promise((resolve, reject) => {
    if (window.AMap) {
      amapLoaded = true
      resolve()
      return
    }

    if (amapLoading) {
      amapCallbacks.push(() => resolve())
      return
    }

    amapLoading = true

    // 配置安全密钥
    if (AMAP_SECURITY_KEY) {
      window._AMapSecurityConfig = {
        securityJsCode: AMAP_SECURITY_KEY,
      }
    }

    const script = document.createElement('script')
    script.src = `https://webapi.amap.com/maps?v=2.0&key=${AMAP_KEY}`
    script.async = true

    script.onload = () => {
      amapLoaded = true
      amapLoading = false
      amapCallbacks.forEach(cb => cb())
      amapCallbacks = []
      resolve()
    }

    script.onerror = () => {
      amapLoading = false
      reject(new Error('Failed to load AMap SDK'))
    }

    document.head.appendChild(script)
  })
}

export function isAmapLoaded(): boolean {
  return amapLoaded
}

export function getAmapKey(): string {
  return AMAP_KEY
}
