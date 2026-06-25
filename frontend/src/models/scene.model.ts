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
    scene_code?: string
    file_id?: string
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
    file_id?: string
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
    source_url: string
    thumbnail_url: string
    source_width: number
    source_height: number
    view_count: number
    slice_status: string
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
    file_id: string
    source_url: string
    source_width: number
    source_height: number
    source_file_size: number
    source_file_sha256: string
    tile_url: string
    preview_url: string
    cubemap_url: string
    is_converted: boolean
    slice_status: string
    task_id: string
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

export type SpaceInfoForGraph = {
    id: number
    name: string
    longitude: number
    latitude: number
    zoom_level: number
}

export type SceneNodeData = {
    id: number
    title: string
    scene_code: string
    thumbnail_url: string
    longitude: number
    latitude: number
    has_position: boolean
    view_count: number
}

export type EdgeData = {
    id: number
    source_id: number
    target_id: number
    hotspot_id: number
    hotspot_title: string
}

export type GraphDataResponse = {
    space_info: SpaceInfoForGraph
    nodes: SceneNodeData[]
    edges: EdgeData[]
    unplaced: SceneNodeData[]
}

export type UpdatePositionRequest = {
    longitude: number
    latitude: number
}

export type ScenePosition = {
    scene_id: number
    longitude: number
    latitude: number
}

export type BatchUpdatePositionRequest = {
    positions: ScenePosition[]
}

export type SceneNode = {
    id: number
    title: string
    scene_code: string
    longitude: number
    latitude: number
    thumbnail_url?: string
    view_count?: number
    sort_order?: number
    hotspots?: { target_scene_id?: number }[]
}

export type SceneSimple = {
    id: number
    title: string
    scene_code: string
    thumbnail_url?: string
    view_count?: number
    sort_order?: number
    longitude: number
    latitude: number
}

export type SceneItem = {
    id: number
    title: string
    thumbnail_url?: string
    scene_code?: string
}