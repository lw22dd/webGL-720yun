import type { SpaceListItem } from '@/models/SpaceModel'
import type { SceneSimple } from '@/models/SpaceModel'

export interface MockScene extends SceneSimple {
  longitude: number
  latitude: number
}

export interface MockSpace extends SpaceListItem {
  scenes: MockScene[]
}

export const mockSpaces: MockSpace[] = [
  {
    id: 1,
    name: '北京故宫',
    cover_url: 'https://picsum.photos/seed/gugong/800/600',
    description: '明清两代的皇家宫殿',
    province: '北京市',
    city: '北京市',
    longitude: 116.397128,
    latitude: 39.916527,
    zoom_level: 15,
    sort_order: 1,
    status: 1,
    scene_count: 2,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    scenes: [
      { id: 101, title: '太和殿', scene_code: 'gugong-taihedian', thumbnail_url: 'https://picsum.photos/seed/th1/200/150', view_count: 1234, sort_order: 1, longitude: 116.397128, latitude: 39.916527 },
      { id: 102, title: '乾清宫', scene_code: 'gugong-qianqinggong', thumbnail_url: 'https://picsum.photos/seed/th2/200/150', view_count: 856, sort_order: 2, longitude: 116.3965, latitude: 39.9175 }
    ]
  },
  {
    id: 2,
    name: '上海外滩',
    cover_url: 'https://picsum.photos/seed/waitan/800/600',
    description: '上海最著名的景观',
    province: '上海市',
    city: '上海市',
    longitude: 121.490317,
    latitude: 31.239967,
    zoom_level: 14,
    sort_order: 2,
    status: 1,
    scene_count: 1,
    created_at: '2024-01-02T00:00:00Z',
    updated_at: '2024-01-02T00:00:00Z',
    scenes: [
      { id: 201, title: '外滩观景台', scene_code: 'waitan-viewpoint', thumbnail_url: 'https://picsum.photos/seed/wt1/200/150', view_count: 2345, sort_order: 1, longitude: 121.490317, latitude: 31.239967 }
    ]
  },
  {
    id: 3,
    name: '杭州西湖',
    cover_url: 'https://picsum.photos/seed/xihu/800/600',
    description: '世界文化遗产',
    province: '浙江省',
    city: '杭州市',
    longitude: 120.15507,
    latitude: 30.246143,
    zoom_level: 14,
    sort_order: 3,
    status: 1,
    scene_count: 1,
    created_at: '2024-01-03T00:00:00Z',
    updated_at: '2024-01-03T00:00:00Z',
    scenes: [
      { id: 301, title: '断桥残雪', scene_code: 'xihu-duanqiao', thumbnail_url: 'https://picsum.photos/seed/xh1/200/150', view_count: 3456, sort_order: 1, longitude: 120.15507, latitude: 30.246143 }
    ]
  },
  {
    id: 4,
    name: '都江堰',
    cover_url: 'https://picsum.photos/seed/dujiangyan/800/600',
    description: '世界文化遗产，古代水利工程',
    province: '四川省',
    city: '成都市',
    longitude: 103.611329,
    latitude: 31.001546,
    zoom_level: 15,
    sort_order: 4,
    status: 1,
    scene_count: 1,
    created_at: '2024-01-04T00:00:00Z',
    updated_at: '2024-01-04T00:00:00Z',
    scenes: [
      { id: 401, title: '鱼嘴分水堤', scene_code: 'dujiangyan-yuzui', thumbnail_url: 'https://picsum.photos/seed/djy1/200/150', view_count: 4567, sort_order: 1, longitude: 103.611329, latitude: 31.001546 }
    ]
  }
]
