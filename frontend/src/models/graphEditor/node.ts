export const SceneNodeStatus = {
  PENDING: 'pending',
  PLACED: 'placed',
  CONFIGURED: 'configured',
} as const

export type SceneNodeStatusType = typeof SceneNodeStatus[keyof typeof SceneNodeStatus]

export interface SceneNodeData {
  id: string
  name: string
  group: string
  status: SceneNodeStatusType
  thumbnail?: string
  panoUrl?: string
  hasPano: boolean
  remark?: string
  longitude?: number
  latitude?: number
}
