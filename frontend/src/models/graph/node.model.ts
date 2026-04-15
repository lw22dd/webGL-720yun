export const SceneNodeStatus = {
  PENDING: 'pending',
  PLACED: 'placed',
  CONFIGURED: 'configured',
} as const

export type SceneNodeStatusType = typeof SceneNodeStatus[keyof typeof SceneNodeStatus]

export interface SceneNodeData {
  id: number
  title: string
  scene_code: string
  thumbnail_url: string
  longitude: number
  latitude: number
  has_position: boolean
  view_count: number
  status: SceneNodeStatusType
}
