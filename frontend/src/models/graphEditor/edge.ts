export const EdgeDirection = {
  UNIDIRECTIONAL: 'unidirectional',
  BIDIRECTIONAL: 'bidirectional',
} as const

export type EdgeDirectionType = typeof EdgeDirection[keyof typeof EdgeDirection]

export interface SceneEdgeData {
  id: string
  sourceId: string
  targetId: string
  direction: EdgeDirectionType
  label?: string
  type: 'walk' | 'teleport'
}
