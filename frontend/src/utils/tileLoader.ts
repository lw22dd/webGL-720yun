const baseApiUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'

export interface CubemapTiles {
  left: string
  front: string
  right: string
  back: string
  top: string
  bottom: string
}

export interface TileRange {
  minX: number
  maxX: number
  minY: number
  maxY: number
}

export interface LODConfig {
  level: number
  tileCount: number
  tileSize: number
}

export const LOD_LEVELS: Record<number, LODConfig> = {
  0: { level: 0, tileCount: 1, tileSize: 1024 },
  1: { level: 1, tileCount: 4, tileSize: 512 },
  2: { level: 2, tileCount: 8, tileSize: 512 }
}

export function getTileUrl(sceneCode: string, face: string, level: number, x: number, y: number): string {
  return `${baseApiUrl}/api/v1/res/tiles/${sceneCode}/${face}/${level}/${x}/${y}`
}

// 获取 Level 0 瓦片 URL（1×1，作为 baseUrl 使用）
export function getBaseTileUrls(sceneCode: string): CubemapTiles {
  const faces = ['nx', 'pz', 'px', 'nz', 'py', 'ny']
  const faceNames = ['left', 'front', 'right', 'back', 'top', 'bottom']
  const result = {} as CubemapTiles
  faceNames.forEach((name, index) => {
    (result as any)[name] = getTileUrl(sceneCode, faces[index], 0, 0, 0)
  })
  return result
}

export function getLevelTileCount(level: number): number {
  return LOD_LEVELS[level]?.tileCount || 2
}

export function getMaxCoord(level: number): number {
  return getLevelTileCount(level) - 1
}

export function isValidTileCoord(level: number, x: number, y: number): boolean {
  const maxCoord = getMaxCoord(level)
  return x >= 0 && x <= maxCoord && y >= 0 && y <= maxCoord
}

export function calculateVisibleTiles(
  _yaw: number,
  _pitch: number,
  _fov: number,
  level: number
): Map<string, TileRange> {
  const tileCount = getLevelTileCount(level)
  const ranges = new Map<string, TileRange>()

  const faces = ['px', 'nx', 'py', 'ny', 'pz', 'nz']
  
  faces.forEach(face => {
    ranges.set(face, {
      minX: 0,
      maxX: tileCount - 1,
      minY: 0,
      maxY: tileCount - 1
    })
  })

  return ranges
}

export class TilePreloader {
  private loadedTiles: Map<string, boolean> = new Map()
  private loadingTiles: Map<string, Promise<void>> = new Map()
  private sceneCode: string
  private enabled: boolean = true

  constructor(sceneCode: string) {
    this.sceneCode = sceneCode
  }

  setEnabled(enabled: boolean): void {
    this.enabled = enabled
  }

  async preloadLevel(level: number): Promise<void> {
    if (!this.enabled) return

    const tileCount = getLevelTileCount(level)
    const faces = ['px', 'nx', 'py', 'ny', 'pz', 'nz']

    const promises: Promise<void>[] = []

    faces.forEach(face => {
      for (let y = 0; y < tileCount; y++) {
        for (let x = 0; x < tileCount; x++) {
          const key = `${face}/${level}/${x}/${y}`
          if (!this.loadedTiles.has(key)) {
            promises.push(this.loadTile(face, level, x, y))
          }
        }
      }
    })

    await Promise.allSettled(promises)
  }

  private async loadTile(face: string, level: number, x: number, y: number): Promise<void> {
    const key = `${face}/${level}/${x}/${y}`
    
    if (this.loadedTiles.has(key)) {
      return
    }

    if (this.loadingTiles.has(key)) {
      return this.loadingTiles.get(key)!
    }

    const url = getTileUrl(this.sceneCode, face, level, x, y)
    const promise = fetch(url, { method: 'GET' })
      .then(res => {
        if (res.ok) {
          this.loadedTiles.set(key, true)
        }
      })
      .catch(() => {
        // 静默处理错误
      })
      .finally(() => {
        this.loadingTiles.delete(key)
      })

    this.loadingTiles.set(key, promise)
    return promise
  }

  isLoaded(face: string, level: number, x: number, y: number): boolean {
    const key = `${face}/${level}/${x}/${y}`
    return this.loadedTiles.has(key)
  }

  clear(): void {
    this.loadedTiles.clear()
    this.loadingTiles.clear()
  }
}
