/**
 * 热点相关的数据模型
 * 对应后端 ResHotspot 模型
 */

/**
 * 热点类型常量
 * 1-场景切换热点，2-教学热点，3-答题热点
 */
export const HotspotType = {
    SWITCH: 1,
    TEACH: 2,
    QUIZ: 3
} as const

export type HotspotTypeValue = typeof HotspotType[keyof typeof HotspotType]

/**
 * 媒体类型常量
 */
export const MediaType = {
    IMAGE: 'image',
    VIDEO: 'video',
    TEXT: 'text'
} as const

export type MediaTypeValue = typeof MediaType[keyof typeof MediaType]

/**
 * 转场效果常量
 */
export const TransitionEffect = {
    FADE: 'fade',
    ZOOM: 'zoom',
    NONE: 'none'
} as const

export type TransitionEffectValue = typeof TransitionEffect[keyof typeof TransitionEffect]

/**
 * 热点列表查询参数
 */
export type HotspotListRequest = {
    page?: number
    page_size?: number
    scene_id?: number
    type?: number
    status?: number
}

/**
 * 创建热点请求参数
 * type: 1-场景切换热点，2-教学热点，3-答题热点
 */
export type CreateHotspotRequest = {
    scene_id: number
    type: 1 | 2 | 3
    target_scene_id?: number
    pitch: number
    yaw: number
    title: string
    icon_url?: string
    style?: string
    content?: string
    media_type?: 'image' | 'video' | 'text'
    media_url?: string
    question?: string
    options?: string
    answer?: string
    score?: number
    transition_effect?: 'fade' | 'zoom' | 'none'
    sort_order?: number
}

/**
 * 更新热点请求参数
 */
export type UpdateHotspotRequest = {
    target_scene_id?: number
    pitch?: number
    yaw?: number
    title?: string
    icon_url?: string
    style?: string
    content?: string
    media_type?: 'image' | 'video' | 'text'
    media_url?: string
    question?: string
    options?: string
    answer?: string
    score?: number
    transition_effect?: 'fade' | 'zoom' | 'none'
    sort_order?: number
    status?: number
}

/**
 * 热点列表项
 */
export type HotspotListItem = {
    id: number
    scene_id: number
    type: number
    target_scene_id?: number
    pitch: number
    yaw: number
    title: string
    icon_url: string
    style: string
    transition_effect: string
    sort_order: number
    status: number
    created_at: string
}

/**
 * 热点详情响应
 */
export type HotspotDetailResponse = {
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
    created_at: string
}

/**
 * 热点列表响应
 */
export type HotspotListResponse = {
    page_info: {
        page: number
        page_size: number
        total: number
    }
    hotspots: HotspotListItem[]
}
