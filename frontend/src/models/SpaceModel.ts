export type SpaceListRequest = {
    page?: number
    page_size?: number
    name?: string
    status?: number
    keyword?: string
}

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

export type SceneSimple = {
    id: number
    title: string
    scene_code: string
    thumbnail_url: string
    view_count: number
    sort_order: number
}

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

export type SpaceListResponse = {
    page_info: {
        page: number
        page_size: number
        total: number
    }
    spaces: SpaceListItem[]
}