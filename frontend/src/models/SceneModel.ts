/**
 * 场景相关的数据模型
 * 对应后端 ResScene 模型
 */

/**
 * 场景列表查询参数
 */
export type SceneListRequest = {
    page?: number
    page_size?: number
    space_id?: number
    title?: string
    status?: number
    keyword?: string
}

/**
 * 创建场景请求参数
 */
export type CreateSceneRequest = {
    space_id: number
    title: string
    scene_code: string
    panorama_type?: 'equirectangular' | 'cubemap'
    initial_fov?: number
    initial_pitch?: number
    initial_yaw?: number
    north_offset?: number
    longitude?: number
    latitude?: number
    sort_order?: number
}

/**
 * 更新场景请求参数
 */
export type UpdateSceneRequest = {
    title?: string
    panorama_type?: 'equirectangular' | 'cubemap'
    initial_fov?: number
    initial_pitch?: number
    initial_yaw?: number
    north_offset?: number
    longitude?: number
    latitude?: number
    sort_order?: number
    status?: number
}

/**
 * 热点简单信息（用于场景详情中的热点列表）
 */
export type HotspotSimple = {
    id: number
    type: number
    title: string
    pitch: number
    yaw: number
    icon_url: string
}

/**
 * 场景列表项
 */
export type SceneListItem = {
    id: number
    space_id: number
    title: string
    scene_code: string
    panorama_type: string
    thumbnail_url: string
    source_width: number
    source_height: number
    view_count: number
    sort_order: number
    status: number
    created_at: string
    updated_at: string
}

/**
 * 场景详情响应
 * 对应后端 ResScene 模型（包含关联的 Hotspots）
 */
export type SceneDetailResponse = {
    id: number
    space_id: number
    title: string
    scene_code: string
    panorama_type: string
    source_url: string
    source_width: number
    source_height: number
    source_file_size: number
    source_file_md5: string
    thumbnail_url: string
    cover_image_url: string
    initial_fov: number
    initial_pitch: number
    initial_yaw: number
    north_offset: number
    longitude: number
    latitude: number
    view_count: number
    sort_order: number
    status: number
    created_at: string
    updated_at: string
    hotspots?: HotspotSimple[]
}

/**
 * 场景列表响应
 */
export type SceneListResponse = {
    page_info: {
        page: number
        page_size: number
        total: number
    }
    scenes: SceneListItem[]
}

/**
 * 批量导入结果项
 */
export type BatchImportResult = {
    index: number
    scene_code: string
    title: string
    status: string
    message: string
}

/**
 * 批量导入错误项
 */
export type BatchImportError = {
    index: number
    scene_code: string
    error: string
}

/**
 * 批量导入响应
 */
export type BatchImportResponse = {
    success_count: number
    failed_count: number
    results: BatchImportResult[]
    errors?: BatchImportError[]
}

/**
 * 全景图类型常量
 */
export const PanoramaType = {
    EQUIRECTANGULAR: 'equirectangular',
    CUBEMAP: 'cubemap'
} as const

export type PanoramaTypeValue = typeof PanoramaType[keyof typeof PanoramaType]
