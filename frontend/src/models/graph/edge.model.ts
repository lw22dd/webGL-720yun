export interface SceneEdgeData {
  id: number
  source_id: number
  target_id: number
  hotspot_id: number
  hotspot_title: string
  type: 'walk' | 'teleport'
}
