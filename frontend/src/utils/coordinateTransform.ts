export interface CoordinateBounds {
  minLng: number
  maxLng: number
  minLat: number
  maxLat: number
}

export class CoordinateTransform {
  private referenceLng: number = 0
  private referenceLat: number = 0
  private scale: number = 10000
  private offsetX: number = 0
  private offsetY: number = 0

  setReference(lng: number, lat: number): void {
    this.referenceLng = lng
    this.referenceLat = lat
  }

  setScale(scale: number): void {
    this.scale = scale
  }

  setOffset(x: number, y: number): void {
    this.offsetX = x
    this.offsetY = y
  }

  lngLatToScreen(lng: number, lat: number): { x: number; y: number } {
    const x = (lng - this.referenceLng) * this.scale * Math.cos((this.referenceLat * Math.PI) / 180)
    const y = -(lat - this.referenceLat) * this.scale
    return {
      x: x + this.offsetX,
      y: y + this.offsetY
    }
  }

  screenToLngLat(x: number, y: number): { lng: number; lat: number } {
    const adjustedX = x - this.offsetX
    const adjustedY = y - this.offsetY
    return {
      lng: this.referenceLng + adjustedX / (this.scale * Math.cos((this.referenceLat * Math.PI) / 180)),
      lat: this.referenceLat - adjustedY / this.scale
    }
  }

  static calculateBounds(nodes: Array<{ longitude: number; latitude: number }>): CoordinateBounds {
    if (nodes.length === 0) {
      return { minLng: 0, maxLng: 0, minLat: 0, maxLat: 0 }
    }

    let minLng = Infinity
    let maxLng = -Infinity
    let minLat = Infinity
    let maxLat = -Infinity

    for (const node of nodes) {
      if (node.longitude !== 0 || node.latitude !== 0) {
        minLng = Math.min(minLng, node.longitude)
        maxLng = Math.max(maxLng, node.longitude)
        minLat = Math.min(minLat, node.latitude)
        maxLat = Math.max(maxLat, node.latitude)
      }
    }

    if (minLng === Infinity) {
      return { minLng: 0, maxLng: 0, minLat: 0, maxLat: 0 }
    }

    return { minLng, maxLng, minLat, maxLat }
  }

  static calculateCenter(bounds: CoordinateBounds): { lng: number; lat: number } {
    return {
      lng: (bounds.minLng + bounds.maxLng) / 2,
      lat: (bounds.minLat + bounds.maxLat) / 2
    }
  }

  static calculateOptimalScale(bounds: CoordinateBounds, canvasWidth: number, canvasHeight: number, padding: number = 100): number {
    const lngRange = bounds.maxLng - bounds.minLng
    const latRange = bounds.maxLat - bounds.minLat

    if (lngRange === 0 && latRange === 0) {
      return 10000
    }

    const availableWidth = canvasWidth - padding * 2
    const availableHeight = canvasHeight - padding * 2

    const scaleByLng = availableWidth / (lngRange || 1)
    const scaleByLat = availableHeight / (latRange || 1)

    const avgLat = (bounds.minLat + bounds.maxLat) / 2
    const latCorrection = Math.cos((avgLat * Math.PI) / 180)

    const scale = Math.min(scaleByLng / latCorrection, scaleByLat) * 0.8

    return Math.max(1000, Math.min(scale, 100000))
  }
}

export function formatCoordinate(lng: number, lat: number): string {
  return `${lng.toFixed(6)}, ${lat.toFixed(6)}`
}

export function parseCoordinate(str: string): { lng: number; lat: number } | null {
  const parts = str.split(',').map(s => s.trim())
  if (parts.length !== 2) return null

  const lng = parseFloat(parts[0])
  const lat = parseFloat(parts[1])

  if (isNaN(lng) || isNaN(lat)) return null
  if (lng < -180 || lng > 180 || lat < -90 || lat > 90) return null

  return { lng, lat }
}
