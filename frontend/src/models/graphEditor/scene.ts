import type { SceneNodeData } from './node'
import type { SceneEdgeData } from './edge'

export interface SpaceInfoForGraph {
  id: number
  name: string
  longitude: number
  latitude: number
  zoom_level: number
}

export interface GraphDataResponse {
  space_info: SpaceInfoForGraph
  nodes: SceneNodeData[]
  edges: SceneEdgeData[]
}
