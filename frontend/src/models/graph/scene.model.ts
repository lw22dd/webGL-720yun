import type { SceneNodeData } from './node.model'
import type { SceneEdgeData } from './edge.model'

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
