import { SceneNodeStatus } from '@/models/graph/node'
import { EdgeDirection } from '@/models/graph/edge'
import type { GraphSceneData } from '@/models/graph/scene'

/** Mock 场景数据 */
export const mockSceneData: GraphSceneData = {
  id: 'scene-001',
  name: '都江堰景区',
  nodes: [
    // 已放置节点
    { id: '3576333', name: '玉垒关观赏平台', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.6145, latitude: 31.0005 },
    { id: '3576391', name: '禹王宫', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.611862, latitude: 31.003596 },
    { id: '3576420', name: '鱼嘴分水堤1', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.605, latitude: 31.009 },
    { id: '3576400', name: '鱼嘴分水堤2', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.607, latitude: 31.003 },
    { id: '3576562', name: '安澜索桥', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.609, latitude: 31.006 },
    { id: '3576708', name: '二王庙', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.6105, latitude: 31.006 },
    { id: '3576638', name: '飞沙堰', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.608, latitude: 31.001 },
    { id: '3576621', name: '伏龙观', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.6122, latitude: 30.999 },
    { id: '3576546', name: '秦堰楼', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.607, latitude: 31.01 },
    { id: '3576761', name: '都江堰大门', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.615379, latitude: 30.995792 },

    // 待处理节点
    { id: '3576433', name: '宣威门', status: SceneNodeStatus.PENDING, hasPano: false, longitude: 0, latitude: 0 },
    { id: '3576465', name: '爱情桥', status: SceneNodeStatus.PENDING, hasPano: false, longitude: 0, latitude: 0 },
    { id: '3576455', name: '西关', status: SceneNodeStatus.PENDING, hasPano: true, longitude: 0, latitude: 0 },
    { id: '3576475', name: '太极殿', status: SceneNodeStatus.PENDING, hasPano: true, longitude: 0, latitude: 0 },
    { id: '3576487', name: '石林', status: SceneNodeStatus.PENDING, hasPano: false, longitude: 0, latitude: 0 },
    { id: '3576515', name: '善水阁', status: SceneNodeStatus.PENDING, hasPano: true, longitude: 0, latitude: 0 },
  ],
  edges: [
    { id: 'edge-001', sourceId: '3576562', targetId: '3576708', direction: EdgeDirection.BIDIRECTIONAL, label: '步行 5 分钟', type: 'walk' },
    { id: 'edge-002', sourceId: '3576400', targetId: '3576638', direction: EdgeDirection.BIDIRECTIONAL, label: '步行 3 分钟', type: 'walk' },
  ],
}

export async function fetchGraphSceneData(id: string): Promise<GraphSceneData> {
  return new Promise((resolve) => {
    setTimeout(() => resolve({ ...mockSceneData }), 500)
  })
}
