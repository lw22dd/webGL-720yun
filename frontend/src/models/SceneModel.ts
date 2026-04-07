export type SceneListRequest = {
    page?: number
    page_size?: number
    space_id?: number
    title?: string
    status?: number
    keyword?: string
}

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

export type HotspotSimple = {
    id: number
    type: number
    title: string
    pitch: number
    yaw: number
    icon_url: string
}

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

export type SceneListResponse = {
    page_info: {
        page: number
        page_size: number
        total: number
    }
    scenes: SceneListItem[]
}

export type BatchImportResult = {
    index: number
    scene_code: string
    title: string
    status: string
    message: string
}

export type BatchImportError = {
    index: number
    scene_code: string
    error: string
}

export type BatchImportResponse = {
    success_count: number
    failed_count: number
    results: BatchImportResult[]
    errors?: BatchImportError[]
}