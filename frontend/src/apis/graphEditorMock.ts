import type { SceneData } from '@/models/graphEditor/scene'
import { SceneNodeStatus } from '@/models/graphEditor/node'
import { EdgeDirection } from '@/models/graphEditor/edge'

export const mockSceneData: SceneData = {
  id: 'scene-001',
  name: '都江堰景区',
  nodes: [
    { id: '3576333', name: '玉垒关观赏平台', group: '景点', status: SceneNodeStatus.CONFIGURED, hasPano: true, longitude: 103.614500, latitude: 31.000500 },
    { id: '3576350', name: '玉垒关', group: '古迹', status: SceneNodeStatus.CONFIGURED, hasPano: true },
    { id: '3576369', name: '玉垒殿', group: '古迹', status: SceneNodeStatus.CONFIGURED, hasPano: true },
    { id: '3576391', name: '禹王宫', group: '古迹', status: SceneNodeStatus.CONFIGURED, hasPano: true, longitude: 103.611862, latitude: 31.003596 },
    { id: '3576420', name: '鱼嘴分水堤1', group: '水利设施', status: SceneNodeStatus.CONFIGURED, hasPano: true, longitude: 103.605000, latitude: 31.009000 },
    { id: '3576400', name: '鱼嘴分水堤2', group: '水利设施', status: SceneNodeStatus.CONFIGURED, hasPano: true, longitude: 103.607000, latitude: 31.003000 },
    { id: '3576562', name: '安澜索桥', group: '桥梁', status: SceneNodeStatus.CONFIGURED, hasPano: true, longitude: 103.609000, latitude: 31.006000 },
    { id: '3576650', name: '二王庙入口', group: '古迹', status: SceneNodeStatus.CONFIGURED, hasPano: true },
    { id: '3576658', name: '二王庙内景', group: '古迹', status: SceneNodeStatus.CONFIGURED, hasPano: true },
    { id: '3576681', name: '二王庙道观', group: '古迹', status: SceneNodeStatus.CONFIGURED, hasPano: true },
    { id: '3576708', name: '二王庙', group: '古迹', status: SceneNodeStatus.CONFIGURED, hasPano: true, longitude: 103.610500, latitude: 31.006000 },
    { id: '3576638', name: '飞沙堰', group: '水利设施', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.608000, latitude: 31.001000 },
    { id: '3576597', name: '伏龙观内景', group: '古迹', status: SceneNodeStatus.PLACED, hasPano: true },
    { id: '3576605', name: '伏龙观观景台', group: '古迹', status: SceneNodeStatus.PLACED, hasPano: true },
    { id: '3576621', name: '伏龙观', group: '古迹', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.612200, latitude: 30.999000 },
    { id: '3576546', name: '秦堰楼', group: '古迹', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.607000, latitude: 31.010000 },
    { id: '3576780', name: '南桥-夜景', group: '桥梁', status: SceneNodeStatus.PLACED, hasPano: true },
    { id: '3576761', name: '都江堰景区大门', group: '景点', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.615379, latitude: 30.995792 },
    { id: '3576741', name: '都江堰后门-白景', group: '景点', status: SceneNodeStatus.PLACED, hasPano: true, longitude: 103.607628, latitude: 31.008030 },
    { id: '3576750', name: '东苑', group: '景点', status: SceneNodeStatus.PLACED, hasPano: false },
    { id: '3576575', name: '马王殿', group: '古迹', status: SceneNodeStatus.PLACED, hasPano: false },
    { id: '3576433', name: '宣威门', group: '古迹', status: SceneNodeStatus.PENDING, hasPano: false },
    { id: '3576465', name: '天下爱情第一桥', group: '桥梁', status: SceneNodeStatus.PENDING, hasPano: false },
    { id: '3576455', name: '西关', group: '古迹', status: SceneNodeStatus.PENDING, hasPano: true, longitude: 103.612581, latitude: 30.999904 },
    { id: '3576475', name: '太极殿', group: '古迹', status: SceneNodeStatus.PENDING, hasPano: true },
    { id: '3576487', name: '石林', group: '自然景观', status: SceneNodeStatus.PENDING, hasPano: false },
    { id: '3576515', name: '善水阁', group: '建筑', status: SceneNodeStatus.PENDING, hasPano: true, longitude: 103.611500, latitude: 31.004500 },
    { id: '3576503', name: '商铺街', group: '商业', status: SceneNodeStatus.PENDING, hasPano: false, longitude: 103.612500, latitude: 31.003000 },
    { id: '3576494', name: '神木艺术馆', group: '展馆', status: SceneNodeStatus.PENDING, hasPano: true, longitude: 103.613500, latitude: 31.001500 },
    { id: '3576567', name: '门犀亭', group: '古迹', status: SceneNodeStatus.PENDING, hasPano: false },
    { id: '3576534', name: '三官殿', group: '古迹', status: SceneNodeStatus.PENDING, hasPano: true },
    { id: '3576583', name: '魁星点斗', group: '古迹', status: SceneNodeStatus.PENDING, hasPano: false },
    { id: '3576591', name: '观景台', group: '景点', status: SceneNodeStatus.PENDING, hasPano: true },
    { id: '3576767', name: '财神殿', group: '古迹', status: SceneNodeStatus.PENDING, hasPano: false },
    { id: '3576775', name: '安澜桥', group: '桥梁', status: SceneNodeStatus.PENDING, hasPano: false },
    { id: '3576730', name: '斗犀亭', group: '古迹', status: SceneNodeStatus.PENDING, hasPano: false },
    { id: '3576798', name: '都江堰后门-夜景', group: '景点', status: SceneNodeStatus.PENDING, hasPano: false },
  ],
  edges: [
    { id: 'edge-001', sourceId: '3576562', targetId: '3576708', direction: EdgeDirection.BIDIRECTIONAL, label: '步行 5 分钟', type: 'walk' },
    { id: 'edge-002', sourceId: '3576400', targetId: '3576638', direction: EdgeDirection.BIDIRECTIONAL, label: '步行 3 分钟', type: 'walk' },
    { id: 'edge-003', sourceId: '3576638', targetId: '3576621', direction: EdgeDirection.UNIDIRECTIONAL, label: '步行 2 分钟', type: 'walk' },
    { id: 'edge-004', sourceId: '3576597', targetId: '3576621', direction: EdgeDirection.BIDIRECTIONAL, label: '步行 1 分钟', type: 'walk' },
    { id: 'edge-005', sourceId: '3576420', targetId: '3576400', direction: EdgeDirection.BIDIRECTIONAL, label: '步行 4 分钟', type: 'walk' },
  ],
}

export async function fetchSceneData(sceneId: string): Promise<SceneData> {
  console.log(`[Mock] 获取场景数据: ${sceneId}`)
  return new Promise((resolve) => {
    setTimeout(() => resolve({ ...mockSceneData }), 300)
  })
}

export async function saveSceneData(data: SceneData): Promise<boolean> {
  console.log('[Mock] 保存场景数据:', data)
  return new Promise((resolve) => {
    setTimeout(() => resolve(true), 500)
  })
}
