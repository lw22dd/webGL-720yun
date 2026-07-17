export type HotspotListRequest = {
    page?: number
    page_size?: number
    scene_id?: number
    type?: number
    status?: number
}

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

export type HotspotListResponse = {
    page_info: {
        page: number
        page_size: number
        total: number
    }
    hotspots: HotspotListItem[]
}

export type HotspotForViewer = HotspotDetailResponse & {
  target_scene_id?: number
  question?: string
  options?: string
  answer?: string
  score?: number
}

export type HotspotIconPreset = {
  key: string
  label: string
  svg: string
  appliesTo: number[]
}

export type HotspotEditorForm = {
  id?: number
  scene_id: number
  type: 1 | 2 | 3
  pitch: number
  yaw: number
  title: string
  icon_source: 'preset' | 'custom'
  icon_preset_key?: string
  icon_url?: string
  target_scene_id?: number
  content?: string
  media_type?: 'image' | 'text'
  media_url?: string
  transition_effect?: 'fade' | 'zoom' | 'none'
}