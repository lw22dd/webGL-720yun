const baseApiUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:7000'

/**
 * 瓦片 URL 版本号。
 * 后端瓦片是 immutable 强缓存，只要重新切片（尺寸/切分数变化）就必须 +1，
 * 否则浏览器和 PSV 会继续命中旧图。
 */
const TILE_URL_VERSION = 2

export interface CubemapTiles {
  left: string
  front: string
  right: string
  back: string
  top: string
  bottom: string
}

/** 立方体面名（PSV 侧）→ 后端存储面名 */
const API_FACE_BY_NAME: Record<string, string> = {
  left: 'nx',
  front: 'pz',
  right: 'px',
  back: 'nz',
  top: 'py',
  bottom: 'ny',
}

export interface CubeTileLevel {
  /** 该层级整面的像素尺寸 */
  faceSize: number
  /** 每面切分成的瓦片边数，faceSize / nbTiles = 单张瓦片像素 */
  nbTiles: number
  /** /api/v1/res/tiles/{scene}/{face}/{apiLevel}/{x}/{y} 里的 level 段 */
  apiLevel: number
}

/**
 * 唯一配置源：与后端 pkg/image/LOD_tiler.go 的 DefaultLODSepcs 一一对应。
 * 改这里就等于同时改了 baseUrl、PSV 的 levels 和瓦片 URL 映射。
 *
 * 三档固定面尺寸，后端会负责把源图重采样到对应尺寸：
 *   512  → 约 480p，首屏底图，用户能看出是模糊占位
 *   2048 → 约 1080p，正常转动视角时的基准画质
 *   4096 → 约 4K，放大后才请求，保证局部清晰
 */
const TILE_LEVEL_SPECS: readonly CubeTileLevel[] = [
  { faceSize: 512, nbTiles: 1, apiLevel: 0 },
  { faceSize: 2048, nbTiles: 4, apiLevel: 1 },
  { faceSize: 4096, nbTiles: 8, apiLevel: 2 },
]

/** 低清底图，作为 CubemapTilesAdapter 的 baseUrl */
export const BASE_TILE_LEVEL = TILE_LEVEL_SPECS[0]

/** 至少保留一级 LOD，否则 levels 数组为空会让 adapter 崩溃 */
const MIN_LOD = TILE_LEVEL_SPECS[1]

/**
 * 设计上限：允许使用到的最大的面尺寸。
 * 立方体贴图显存 ≈ 6 × faceSize² × 4B：
 *   512 → 6MB，2048 → 100MB，4096 → 402MB。
 * 答辩机器显存吃紧、或想统一降一档观感时，把这个值改成 2048 即可全局封顶。
 */
export const MAX_FACE_SIZE = 4096

/**
 * 交给 PSV 的 LOD 阶梯。
 *
 * 三档与后端固定对应，和源图大小无关：
 *   2048 面 → 1080p，正常观看
 *   4096 面 → 4K，放大后保持清晰
 *
 * 8K 源图的 level 2 会由后端 Lanczos 重采样到 4096 面，虽然不增加原始细节，
 * 但能在放大时保持锐度，避免 2048 面被拉伸后的模糊感。
 */
export function resolveTileLevels(): CubeTileLevel[] {
  const levels = TILE_LEVEL_SPECS.filter(
    (spec) => spec.apiLevel > 0 && spec.faceSize <= MAX_FACE_SIZE
  )
  return levels.length > 0 ? levels : [MIN_LOD]
}

export function getTileUrl(
  sceneCode: string,
  apiFace: string,
  apiLevel: number,
  x: number,
  y: number
): string {
  return `${baseApiUrl}/api/v1/res/tiles/${sceneCode}/${apiFace}/${apiLevel}/${x}/${y}?v=${TILE_URL_VERSION}`
}

/** 获取低清底图 URL（level 0 的 1×1 瓦片，作为 baseUrl 使用） */
export function getBaseTileUrls(sceneCode: string): CubemapTiles {
  const names = Object.keys(API_FACE_BY_NAME)
  const result = {} as CubemapTiles
  names.forEach((name) => {
    result[name as keyof CubemapTiles] = getTileUrl(
      sceneCode,
      API_FACE_BY_NAME[name],
      BASE_TILE_LEVEL.apiLevel,
      0,
      0
    )
  })
  return result
}

/**
 * 生成 CubemapTilesAdapter 的 tileUrl 回调。
 * PSV 传进来的 level 是 levels 数组下标（调用前会按 faceSize 升序排序），
 * 这里再映射到后端 apiLevel，避免 "level + 1" 这类隐式假设。
 */
export function createTileUrlResolver(
  sceneCode: string,
  levels: readonly CubeTileLevel[]
): (face: string, col: number, row: number, level: number) => string {
  return (face, col, row, level) => {
    const spec = levels[level] ?? levels[levels.length - 1]
    return getTileUrl(sceneCode, API_FACE_BY_NAME[face] || face, spec.apiLevel, col, row)
  }
}

/**
 * 允许把 LOD 判定放大到几倍的物理像素。
 * 3x 手机上按 3 倍物理像素挑层会让首屏直接落到最高档，这里封到 2 倍。
 */
const MAX_LOD_PIXEL_RATIO = 2

/**
 * 生成传给 CubemapTilesAdapter 的 levels。
 *
 * PSV 的 LOD 选择用的是 `viewer.state.size`（CSS 像素），但画布实际按
 * `devicePixelRatio` 渲染。直接传真实 faceSize 时，150% 缩放的屏幕上
 * 2048 面会被判定为“已经覆盖屏幕”，即使放到最大也选不到 4096 那一档，
 * 表现就是放大后没有变化。这里把 faceSize 折算成 CSS 像素下的等效尺寸，
 * 让 PSV 按物理像素判断是否需要更高清的一级。
 *
 * adapter 只用 faceSize 做层级选择（tileSize 字段在 5.14.1 中未参与渲染），
 * 实际瓦片宽高由 nbTiles 决定；请求哪一级仍由 createTileUrlResolver 按
 * 原始 apiLevel 映射，不会改变瓦片本身的分辨率。
 */
export function toViewerLevels(
  levels: readonly CubeTileLevel[]
): Array<{ faceSize: number; nbTiles: number }> {
  const pixelRatio =
    typeof window !== 'undefined' && window.devicePixelRatio > 0
      ? Math.min(window.devicePixelRatio, MAX_LOD_PIXEL_RATIO)
      : 1

  return levels.map(({ faceSize, nbTiles }) => ({
    faceSize: Math.max(1, Math.round(faceSize / pixelRatio)),
    nbTiles,
  }))
}
