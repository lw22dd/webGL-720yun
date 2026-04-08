/**
 * 空间/景区相关的数据模型
 * 对应后端 ResSpace, ResScene, ResHotspot 模型
 */

/**
 * 空间列表查询参数
 */
export type SpaceListRequest = {
    page?: number
    page_size?: number
    name?: string
    status?: number
    keyword?: string
}

/**
 * 创建空间请求参数
 */
export type CreateSpaceRequest = {
    name: string
    slug: string
    description?: string
    province?: string
    city?: string
    longitude?: number
    latitude?: number
    zoom_level?: number
    sort_order?: number
}

/**
 * 更新空间请求参数
 */
export type UpdateSpaceRequest = {
    name?: string
    description?: string
    province?: string
    city?: string
    longitude?: number
    latitude?: number
    zoom_level?: number
    sort_order?: number
    status?: number
}

/**
 * 空间列表项
 * 对应后端 ResSpace 模型
 */
export type SpaceListItem = {
    id: number
    name: string
    slug: string
    cover_url: string
    description: string
    province: string
    city: string
    longitude: number
    latitude: number
    zoom_level: number
    sort_order: number
    status: number
    scene_count: number
    created_at: string
    updated_at: string
}

/**
 * 场景简单信息（用于空间详情中的场景列表）
 */
export type SceneSimple = {
    id: number
    title: string
    scene_code: string
    thumbnail_url: string
    view_count: number
    sort_order: number
}

/**
 * 热点信息
 * 对应后端 ResHotspot 模型
 */
export type Hotspot = {
    id: number
    scene_id: number
    type: number
    target_scene_id?: number
    pitch: number
    yaw: number
    title: string
    icon_url: string
    style: string
    content?: string
    media_type?: string
    media_url?: string
    question?: string
    options?: string
    answer?: string
    score: number
    transition_effect: string
    sort_order: number
    status: number
}

/**
 * 场景详情
 * 对应后端 ResScene 模型
 */
export type SceneDetail = {
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
    cubemap_url: string
    is_converted: boolean
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
    hotspots?: Hotspot[]
}

/**
 * 空间详情响应
 * 对应后端 ResSpace 模型（包含关联的 Scenes）
 */
export type SpaceDetailResponse = {
    id: number
    name: string
    slug: string
    cover_url: string
    description: string
    province: string
    city: string
    longitude: number
    latitude: number
    zoom_level: number
    sort_order: number
    status: number
    created_at: string
    updated_at: string
    created_by: number
    scenes?: SceneSimple[]
}

/**
 * 空间列表响应
 */
export type SpaceListResponse = {
    page_info: {
        page: number
        page_size: number
        total: number
    }
    spaces: SpaceListItem[]
}

/**
 * 热点类型常量
 * 对应后端 HotspotType 常量
 */
export const HotspotType = {
    SWITCH: 1,
    TEACH: 2,
    QUIZ: 3
} as const

export type HotspotTypeValue = typeof HotspotType[keyof typeof HotspotType]
