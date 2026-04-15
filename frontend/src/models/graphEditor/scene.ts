import type { SceneNodeData } from './node'
import type { SceneEdgeData } from './edge'

export interface BackgroundMap {
  url: string
  width: number
  height: number
  x: number
  y: number
  opacity: number
}

export interface SceneData {
  id: string
  name: string
  nodes: SceneNodeData[]
  edges: SceneEdgeData[]
  backgroundMap?: BackgroundMap
}
