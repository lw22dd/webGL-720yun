import type { SpaceListItem } from '@/models/SpaceModel'

export interface SceneSimple {
  id: number
  title: string
  scene_code: string
  thumbnail_url: string
  view_count: number
  sort_order: number
  longitude: number
  latitude: number
}

export interface MockScene extends SceneSimple {
}

export interface MockSpace extends SpaceListItem {
  scenes: MockScene[]
}

export const mockSpaces: MockSpace[] = [
  {
    id: 1,
    name: '都江堰景区',
    cover_url: 'https://picsum.photos/seed/dujiangyan/800/600',
    description: '世界文化遗产，古代水利工程杰作，始建于公元前256年',
    province: '四川省',
    city: '成都市',
    longitude: 103.611379,
    latitude: 31.001752,
    zoom_level: 15,
    sort_order: 1,
    status: 1,
    scene_count: 37,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    scenes: [
      // 景区大门 → 离堆（小索桥）段
      { id: 101, title: '景区大门', scene_code: 'djy-gate', thumbnail_url: 'https://picsum.photos/seed/gate/200/150', view_count: 400, sort_order: 1, longitude: 103.615379, latitude: 30.995792 },
      { id: 102, title: '卧铁', scene_code: 'djy-wotie', thumbnail_url: 'https://picsum.photos/seed/wotie/200/150', view_count: 150, sort_order: 2, longitude: 103.614500, latitude: 30.996500 },
      { id: 103, title: '清溪园', scene_code: 'djy-qingxi', thumbnail_url: 'https://picsum.photos/seed/qingxi/200/150', view_count: 120, sort_order: 3, longitude: 103.613800, latitude: 30.997200 },
      { id: 104, title: '天府源茶馆', scene_code: 'djy-teahouse', thumbnail_url: 'https://picsum.photos/seed/teahouse/200/150', view_count: 80, sort_order: 4, longitude: 103.613200, latitude: 30.997800 },
      { id: 105, title: '堰功道', scene_code: 'djy-yangong', thumbnail_url: 'https://picsum.photos/seed/yangong/200/150', view_count: 100, sort_order: 5, longitude: 103.612800, latitude: 30.998200 },
      { id: 106, title: '张松银杏', scene_code: 'djy-yinxing', thumbnail_url: 'https://picsum.photos/seed/yinxing/200/150', view_count: 90, sort_order: 6, longitude: 103.612500, latitude: 30.998600 },
      { id: 107, title: '伏龙观', scene_code: 'djy-fulong', thumbnail_url: 'https://picsum.photos/seed/fulong/200/150', view_count: 200, sort_order: 7, longitude: 103.612200, latitude: 30.999000 },
      { id: 108, title: '离堆', scene_code: 'djy-lidui', thumbnail_url: 'https://picsum.photos/seed/lidui/200/150', view_count: 250, sort_order: 8, longitude: 103.611800, latitude: 30.999500 },
      
      // 飞沙堰 → 金刚堤 → 安澜索桥 → 鱼嘴
      { id: 109, title: '飞沙堰', scene_code: 'djy-feisha', thumbnail_url: 'https://picsum.photos/seed/feisha/200/150', view_count: 180, sort_order: 9, longitude: 103.608000, latitude: 31.001000 },
      { id: 110, title: '金刚堤', scene_code: 'djy-jingang', thumbnail_url: 'https://picsum.photos/seed/jingang/200/150', view_count: 160, sort_order: 10, longitude: 103.607000, latitude: 31.003000 },
      { id: 111, title: '安澜索桥', scene_code: 'djy-anlan', thumbnail_url: 'https://picsum.photos/seed/anlan/200/150', view_count: 300, sort_order: 11, longitude: 103.609000, latitude: 31.006000 },
      { id: 112, title: '鱼嘴', scene_code: 'djy-yuzui', thumbnail_url: 'https://picsum.photos/seed/yuzui/200/150', view_count: 220, sort_order: 12, longitude: 103.605000, latitude: 31.009000 },
      
      // 安澜索桥 → 秦堰楼 → 二王庙 → 灵动森林 → 敬修之牌坊 → 松茂古道
      { id: 113, title: '秦堰楼', scene_code: 'djy-qinyan', thumbnail_url: 'https://picsum.photos/seed/qinyan/200/150', view_count: 170, sort_order: 13, longitude: 103.607000, latitude: 31.010000 },
      { id: 114, title: '二王庙', scene_code: 'djy-erwang', thumbnail_url: 'https://picsum.photos/seed/erwang/200/150', view_count: 190, sort_order: 14, longitude: 103.610500, latitude: 31.006000 },
      { id: 115, title: '灵动森林', scene_code: 'djy-senlin', thumbnail_url: 'https://picsum.photos/seed/senlin/200/150', view_count: 110, sort_order: 15, longitude: 103.611500, latitude: 31.004500 },
      { id: 116, title: '敬修之牌坊', scene_code: 'djy-paifang', thumbnail_url: 'https://picsum.photos/seed/paifang/200/150', view_count: 130, sort_order: 16, longitude: 103.612500, latitude: 31.003000 },
      { id: 117, title: '松茂古道', scene_code: 'djy-songmao', thumbnail_url: 'https://picsum.photos/seed/songmao/200/150', view_count: 140, sort_order: 17, longitude: 103.613500, latitude: 31.001500 },
      
      // 松茂古道 ↔ 玉垒阁 → 城隍庙 → 十殿 → 玉垒山广场
      { id: 118, title: '玉垒阁', scene_code: 'djy-yulei', thumbnail_url: 'https://picsum.photos/seed/yulei/200/150', view_count: 160, sort_order: 18, longitude: 103.614500, latitude: 31.000500 },
      { id: 119, title: '城隍庙', scene_code: 'djy-chenghuang', thumbnail_url: 'https://picsum.photos/seed/chenghuang/200/150', view_count: 120, sort_order: 19, longitude: 103.615000, latitude: 31.000000 },
      { id: 120, title: '十殿', scene_code: 'djy-shidian', thumbnail_url: 'https://picsum.photos/seed/shidian/200/150', view_count: 100, sort_order: 20, longitude: 103.615500, latitude: 30.999500 },
      { id: 121, title: '玉垒山广场', scene_code: 'djy-square', thumbnail_url: 'https://picsum.photos/seed/square/200/150', view_count: 180, sort_order: 21, longitude: 103.616000, latitude: 30.999000 },
      
      // 待处理区域节点（无坐标，不连接边）
      { id: 122, title: '玉垒关', scene_code: '3576350', thumbnail_url: 'https://picsum.photos/seed/3576350/200/150', view_count: 95, sort_order: 22, longitude: 0, latitude: 0 },
      { id: 123, title: '玉垒殿', scene_code: '3576369', thumbnail_url: 'https://picsum.photos/seed/3576369/200/150', view_count: 50, sort_order: 23, longitude: 0, latitude: 0 },
      { id: 124, title: '宣威门', scene_code: '3576433', thumbnail_url: 'https://picsum.photos/seed/3576433/200/150', view_count: 30, sort_order: 24, longitude: 0, latitude: 0 },
      { id: 125, title: '太极殿', scene_code: '3576475', thumbnail_url: 'https://picsum.photos/seed/3576475/200/150', view_count: 40, sort_order: 25, longitude: 0, latitude: 0 },
      { id: 126, title: '石林', scene_code: '3576487', thumbnail_url: 'https://picsum.photos/seed/3576487/200/150', view_count: 60, sort_order: 26, longitude: 0, latitude: 0 },
      { id: 127, title: '善水阁', scene_code: '3576515', thumbnail_url: 'https://picsum.photos/seed/3576515/200/150', view_count: 35, sort_order: 27, longitude: 0, latitude: 0 },
      { id: 128, title: '商铺街', scene_code: '3576503', thumbnail_url: 'https://picsum.photos/seed/3576503/200/150', view_count: 45, sort_order: 28, longitude: 0, latitude: 0 },
      { id: 129, title: '神木艺术馆', scene_code: '3576494', thumbnail_url: 'https://picsum.photos/seed/3576494/200/150', view_count: 55, sort_order: 29, longitude: 0, latitude: 0 },
      { id: 130, title: '门犀亭', scene_code: '3576567', thumbnail_url: 'https://picsum.photos/seed/3576567/200/150', view_count: 25, sort_order: 30, longitude: 0, latitude: 0 },
      { id: 131, title: '三官殿', scene_code: '3576534', thumbnail_url: 'https://picsum.photos/seed/3576534/200/150', view_count: 30, sort_order: 31, longitude: 0, latitude: 0 },
      { id: 132, title: '魁星点斗', scene_code: '3576583', thumbnail_url: 'https://picsum.photos/seed/3576583/200/150', view_count: 20, sort_order: 32, longitude: 0, latitude: 0 },
      { id: 133, title: '东苑', scene_code: '3576750', thumbnail_url: 'https://picsum.photos/seed/3576750/200/150', view_count: 40, sort_order: 33, longitude: 0, latitude: 0 },
      { id: 134, title: '马王殿', scene_code: '3576575', thumbnail_url: 'https://picsum.photos/seed/3576575/200/150', view_count: 35, sort_order: 34, longitude: 0, latitude: 0 },
      { id: 135, title: '斗犀亭', scene_code: '3576730', thumbnail_url: 'https://picsum.photos/seed/3576730/200/150', view_count: 25, sort_order: 35, longitude: 0, latitude: 0 },
      { id: 136, title: '财神殿', scene_code: '3576767', thumbnail_url: 'https://picsum.photos/seed/3576767/200/150', view_count: 45, sort_order: 36, longitude: 0, latitude: 0 },
      { id: 137, title: '南桥-夜景', scene_code: '3576780', thumbnail_url: 'https://picsum.photos/seed/3576780/200/150', view_count: 320, sort_order: 37, longitude: 0, latitude: 0 },
    ]
  },
  {
    id: 2,
    name: '北京故宫',
    cover_url: 'https://picsum.photos/seed/gugong/800/600',
    description: '明清两代的皇家宫殿',
    province: '北京市',
    city: '北京市',
    longitude: 116.397128,
    latitude: 39.916527,
    zoom_level: 15,
    sort_order: 2,
    status: 1,
    scene_count: 2,
    created_at: '2024-01-02T00:00:00Z',
    updated_at: '2024-01-02T00:00:00Z',
    scenes: [
      { id: 201, title: '太和殿', scene_code: 'gugong-taihedian', thumbnail_url: 'https://picsum.photos/seed/th1/200/150', view_count: 1234, sort_order: 1, longitude: 116.397128, latitude: 39.916527 },
      { id: 202, title: '乾清宫', scene_code: 'gugong-qianqinggong', thumbnail_url: 'https://picsum.photos/seed/th2/200/150', view_count: 856, sort_order: 2, longitude: 116.3965, latitude: 39.9175 }
    ]
  },
  {
    id: 3,
    name: '上海外滩',
    cover_url: 'https://picsum.photos/seed/waitan/800/600',
    description: '上海最著名的景观',
    province: '上海市',
    city: '上海市',
    longitude: 121.490317,
    latitude: 31.239967,
    zoom_level: 14,
    sort_order: 3,
    status: 1,
    scene_count: 1,
    created_at: '2024-01-03T00:00:00Z',
    updated_at: '2024-01-03T00:00:00Z',
    scenes: [
      { id: 301, title: '外滩观景台', scene_code: 'waitan-viewpoint', thumbnail_url: 'https://picsum.photos/seed/wt1/200/150', view_count: 2345, sort_order: 1, longitude: 121.490317, latitude: 31.239967 }
    ]
  }
]
