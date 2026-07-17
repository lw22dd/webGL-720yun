import type { HotspotForViewer, HotspotIconPreset } from '@/models/hotspot.model'

// 预设 SVG 图标（内联，无外链）
export const HOTSPOT_ICON_PRESETS: HotspotIconPreset[] = [
  {
    key: 'arrow-blue',
    label: '跳转-蓝',
    appliesTo: [1],
    svg: `<svg viewBox="0 0 40 40" width="40" height="40" xmlns="http://www.w3.org/2000/svg"><circle cx="20" cy="20" r="18" fill="rgba(14,165,233,0.85)" stroke="rgba(255,255,255,0.9)" stroke-width="2"/><path d="M16 24 L24 16 M24 16 L19 16 M24 16 L24 21" stroke="white" stroke-width="2.5" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
  },
  {
    key: 'arrow-orange',
    label: '跳转-橙',
    appliesTo: [1],
    svg: `<svg viewBox="0 0 40 40" width="40" height="40" xmlns="http://www.w3.org/2000/svg"><circle cx="20" cy="20" r="18" fill="rgba(249,115,22,0.85)" stroke="rgba(255,255,255,0.9)" stroke-width="2"/><path d="M16 24 L24 16 M24 16 L19 16 M24 16 L24 21" stroke="white" stroke-width="2.5" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
  },
  {
    key: 'info-green',
    label: '信息-绿',
    appliesTo: [2],
    svg: `<svg viewBox="0 0 40 40" width="40" height="40" xmlns="http://www.w3.org/2000/svg"><circle cx="20" cy="20" r="18" fill="rgba(16,185,129,0.85)" stroke="rgba(255,255,255,0.9)" stroke-width="2"/><circle cx="20" cy="13" r="2" fill="white"/><rect x="18" y="17" width="4" height="11" rx="2" fill="white"/></svg>`,
  },
  {
    key: 'info-white',
    label: '信息-白',
    appliesTo: [2],
    svg: `<svg viewBox="0 0 40 40" width="40" height="40" xmlns="http://www.w3.org/2000/svg"><circle cx="20" cy="20" r="18" fill="rgba(255,255,255,0.85)" stroke="rgba(0,0,0,0.2)" stroke-width="2"/><circle cx="20" cy="13" r="2" fill="#333"/><rect x="18" y="17" width="4" height="11" rx="2" fill="#333"/></svg>`,
  },
  {
    key: 'quiz-orange',
    label: '答题-橙',
    appliesTo: [3],
    svg: `<svg viewBox="0 0 40 40" width="40" height="40" xmlns="http://www.w3.org/2000/svg"><circle cx="20" cy="20" r="18" fill="rgba(245,158,11,0.85)" stroke="rgba(255,255,255,0.9)" stroke-width="2"/><text x="20" y="27" text-anchor="middle" font-size="18" font-weight="bold" fill="white" font-family="sans-serif">?</text></svg>`,
  },
]

// 按 type 的默认图标 key
const DEFAULT_KEY_BY_TYPE: Record<number, string> = {
  1: 'arrow-blue',
  2: 'info-green',
  3: 'quiz-orange',
}

// 引擎渲染入口：自定义 img 优先，其次预设 style，最后 type 默认
export function getHotspotIconHtml(h: HotspotForViewer): string {
  // 1. 自定义上传图片
  if (h.icon_url && /^https?:|^data:/.test(h.icon_url)) {
    return `<img class="hotspot-marker hotspot-marker-img" src="${h.icon_url}" alt="${h.title}" style="width:40px;height:40px;object-fit:contain;cursor:pointer;" />`
  }

  // 2. 预设 style key
  const preset = HOTSPOT_ICON_PRESETS.find((p) => p.key === h.style)
  if (preset) {
    return `<div class="hotspot-marker hotspot-preset" style="width:40px;height:40px;cursor:pointer;">${preset.svg}</div>`
  }

  // 3. type 默认
  const defaultKey = DEFAULT_KEY_BY_TYPE[h.type]
  const defaultPreset = HOTSPOT_ICON_PRESETS.find((p) => p.key === defaultKey)
  if (defaultPreset) {
    return `<div class="hotspot-marker hotspot-preset" style="width:40px;height:40px;cursor:pointer;">${defaultPreset.svg}</div>`
  }

  // 兜底
  return `<div class="hotspot-marker hotspot-info" style="width:36px;height:36px;">ℹ</div>`
}

// 获取适用于指定类型的预设图标
export function getPresetsForType(type: number): HotspotIconPreset[] {
  return HOTSPOT_ICON_PRESETS.filter((p) => p.appliesTo.includes(type))
}
