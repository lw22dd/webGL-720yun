<template>
  <Transition name="minimap-fade">
    <div v-if="visible" class="pano-minimap">
      <div class="minimap-header">
        <span class="minimap-title">空间节点图</span>
        <div class="minimap-zoom-controls">
          <span class="minimap-zoom-label">{{ Math.round(zoomLevel * 100) }}%</span>
          <t-button
            theme="default"
            variant="text"
            shape="square"
            size="small"
            class="minimap-icon-btn"
            title="重置"
            @click="resetView"
          >
            ⟳
          </t-button>
        </div>
        <t-button
          theme="default"
          variant="text"
          shape="square"
          size="small"
          class="minimap-close"
          @click="close"
        >
          ×
        </t-button>
      </div>
      <div class="minimap-body">
        <canvas
          ref="canvasRef"
          class="minimap-canvas"
          @click="handleCanvasClick"
          @wheel="handleWheel"
          @mousedown="handleMouseDown"
          @mousemove="handleMouseMove"
          @mouseup="handleMouseUp"
          @mouseleave="handleMouseUp"
        />
        <div v-if="loading" class="minimap-overlay">
          <t-loading size="small" text="加载中..." />
        </div>
        <div v-else-if="error" class="minimap-overlay minimap-error">
          {{ error }}
        </div>
        <div v-else class="minimap-hint">滚轮缩放 · 拖动平移</div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import type { Viewer } from '@photo-sphere-viewer/core'
import type { SceneDetailResponse } from '@/models/scene.model'
import type { GraphDataResponse, SceneNodeData } from '@/models/scene.model'
import SpaceApi from '@/apis/space.api'

const props = defineProps<{
  visible: boolean
  currentScene: SceneDetailResponse | null
  getViewer: () => Viewer | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'scene-select', sceneId: number): void
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
const graphData = ref<GraphDataResponse | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const currentView = ref<{ yaw: number; fov: number } | null>(null)

const zoomLevel = ref(1)
const panX = ref(0)
const panY = ref(0)
const isPanning = ref(false)
const hoverNodeId = ref<number | null>(null)
let panStartX = 0
let panStartY = 0
let panOriginX = 0
let panOriginY = 0
let dragMoved = false

let positionHandler: ((e: any) => void) | null = null
let boundViewer: Viewer | null = null
let animationLoopId: number | null = null

const CSS_WIDTH = 280
const CSS_HEIGHT = 200
const PADDING = 32
const MIN_ZOOM = 0.5
const MAX_ZOOM = 4
const BASE_NODE_RADIUS = 5
const BASE_CURRENT_RADIUS = 8
const BASE_LABEL_FONT = 10
const BASE_FAN_RATIO = 0.35

function close() {
  emit('update:visible', false)
}

async function loadGraphData(spaceId: number) {
  if (!spaceId || spaceId <= 0) {
    error.value = '当前场景未归属空间，无法显示节点图'
    loading.value = false
    return
  }

  loading.value = true
  error.value = null
  try {
    const result = await SpaceApi.getSpaceGraph(spaceId)
    if (result.code === 200 && result.data) {
      graphData.value = result.data
    } else {
      error.value = result.msg || '加载节点数据失败'
    }
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '加载失败'
  } finally {
    loading.value = false
  }
}

function setupCanvas() {
  const canvas = canvasRef.value
  if (!canvas) return

  const dpr = window.devicePixelRatio || 1
  canvas.width = CSS_WIDTH * dpr
  canvas.height = (CSS_HEIGHT - 40) * dpr // 40px for header
  canvas.style.width = `${CSS_WIDTH}px`
  canvas.style.height = `${CSS_HEIGHT - 40}px`

  const ctx = canvas.getContext('2d')
  if (!ctx) return
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
}

interface Layout {
  scale: number
  centerLng: number
  centerLat: number
  width: number
  height: number
}

function computeLayout(): Layout | null {
  const canvas = canvasRef.value
  if (!canvas) return null

  const nodes = graphData.value?.nodes || []
  const current = props.currentScene

  if (!current) return null

  const allNodes: { longitude: number; latitude: number }[] = [
    ...nodes,
    { longitude: current.longitude, latitude: current.latitude }
  ]

  if (allNodes.length === 0) return null

  const lngs = allNodes.map(n => n.longitude)
  const lats = allNodes.map(n => n.latitude)
  const minLng = Math.min(...lngs)
  const maxLng = Math.max(...lngs)
  const minLat = Math.min(...lats)
  const maxLat = Math.max(...lats)

  const rangeLng = Math.max(maxLng - minLng, 0.0001)
  const rangeLat = Math.max(maxLat - minLat, 0.0001)

  const width = CSS_WIDTH
  const height = CSS_HEIGHT - 40
  const availableW = width - PADDING * 2
  const availableH = height - PADDING * 2

  const baseScale = Math.min(availableW / rangeLng, availableH / rangeLat)
  const scale = baseScale * zoomLevel.value

  return {
    scale,
    centerLng: (minLng + maxLng) / 2,
    centerLat: (minLat + maxLat) / 2,
    width,
    height
  }
}

function toCanvas(lng: number, lat: number, layout: Layout) {
  const dx = (lng - layout.centerLng) * layout.scale
  const dy = -(lat - layout.centerLat) * layout.scale
  return {
    x: layout.width / 2 + dx + panX.value,
    y: layout.height / 2 + dy + panY.value
  }
}

function drawBackground(ctx: CanvasRenderingContext2D, width: number, height: number) {
  ctx.clearRect(0, 0, width, height)
  ctx.fillStyle = 'rgba(15, 23, 42, 0.35)'
  ctx.fillRect(0, 0, width, height)
}

function drawPlaceholder(ctx: CanvasRenderingContext2D, width: number, height: number) {
  ctx.fillStyle = 'rgba(255, 255, 255, 0.5)'
  ctx.font = '13px sans-serif'
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  ctx.fillText('暂无空间节点数据', width / 2, height / 2)
}

function drawEdges(ctx: CanvasRenderingContext2D, layout: Layout) {
  const edges = graphData.value?.edges || []
  const nodes = graphData.value?.nodes || []
  const currentId = props.currentScene?.id

  if (!edges.length || !nodes.length) return

  // 构建节点查找表
  const nodeMap = new Map<number, SceneNodeData>()
  for (const n of nodes) nodeMap.set(n.id, n)

  // 预处理：识别双向边（A→B 且 B→A）
  // key 使用 "minId-maxId" 规范化，避免重复
  const edgePairs = new Set<string>()
  const bidirectionalSet = new Set<string>()
  for (const e of edges) {
    const key = e.source_id < e.target_id
      ? `${e.source_id}-${e.target_id}`
      : `${e.target_id}-${e.source_id}`
    if (edgePairs.has(key)) {
      bidirectionalSet.add(key)
    } else {
      edgePairs.add(key)
    }
  }

  // 按方向分组：同一方向的边只画一次（去重）
  const drawnSet = new Set<string>()
  // 当前场景的出边集合
  const outgoingSet = new Set<string>()

  for (const edge of edges) {
    const source = nodeMap.get(edge.source_id)
    const target = nodeMap.get(edge.target_id)
    if (!source || !target) continue

    const key = `${edge.source_id}-${edge.target_id}`
    if (drawnSet.has(key)) continue
    drawnSet.add(key)

    const p1 = toCanvas(source.longitude, source.latitude, layout)
    const p2 = toCanvas(target.longitude, target.latitude, layout)

    // 判断边类型
    const pairKey = edge.source_id < edge.target_id
      ? `${edge.source_id}-${edge.target_id}`
      : `${edge.target_id}-${edge.source_id}`
    const isBidirectional = bidirectionalSet.has(pairKey)
    const isOutgoing = edge.source_id === currentId

    if (isOutgoing) {
      outgoingSet.add(key)
    }

    // 颜色策略：
    // - 当前场景出边 → 强调色（琥珀色）
    // - 双向边 → 绿色（可往返）
    // - 普通有向边 → 半透明白
    let strokeColor: string
    let lineWidth: number
    if (isOutgoing) {
      strokeColor = 'rgba(245, 158, 11, 0.85)'
      lineWidth = 1.6
    } else if (isBidirectional) {
      strokeColor = 'rgba(16, 185, 129, 0.7)'
      lineWidth = 1.3
    } else {
      strokeColor = 'rgba(255, 255, 255, 0.3)'
      lineWidth = 1
    }

    ctx.strokeStyle = strokeColor
    ctx.lineWidth = lineWidth

    // 收缩端点：避免箭头插入节点圆心，从节点边缘开始/结束
    const nodeRadius = BASE_NODE_RADIUS * Math.max(0.6, Math.min(1.4, zoomLevel.value))
    const dx = p2.x - p1.x
    const dy = p2.y - p1.y
    const dist = Math.hypot(dx, dy)
    if (dist < nodeRadius * 2) continue // 节点重叠时跳过

    const ux = dx / dist
    const uy = dy / dist
    const startX = p1.x + ux * nodeRadius
    const startY = p1.y + uy * nodeRadius
    const endX = p2.x - ux * nodeRadius
    const endY = p2.y - uy * nodeRadius

    ctx.beginPath()
    ctx.moveTo(startX, startY)
    ctx.lineTo(endX, endY)
    ctx.stroke()

    // 箭头：单向边在终点画箭头；双向边两端都画箭头
    if (isBidirectional) {
      drawArrowHead(ctx, startX, startY, ux, uy, strokeColor)
      drawArrowHead(ctx, endX, endY, -ux, -uy, strokeColor)
    } else {
      drawArrowHead(ctx, endX, endY, -ux, -uy, strokeColor)
    }
  }

  // 对当前场景的出边做额外的脉冲高亮效果
  if (currentId !== undefined && outgoingSet.size > 0) {
    const pulse = (Date.now() % 2000) / 2000
    const alpha = 0.4 + 0.4 * Math.sin(pulse * Math.PI * 2)
    for (const key of outgoingSet) {
      const [srcId, tgtId] = key.split('-').map(Number)
      const source = nodeMap.get(srcId)
      const target = nodeMap.get(tgtId)
      if (!source || !target) continue

      const p1 = toCanvas(source.longitude, source.latitude, layout)
      const p2 = toCanvas(target.longitude, target.latitude, layout)
      const nodeRadius = BASE_NODE_RADIUS * Math.max(0.6, Math.min(1.4, zoomLevel.value))
      const dx = p2.x - p1.x
      const dy = p2.y - p1.y
      const dist = Math.hypot(dx, dy)
      if (dist < nodeRadius * 2) continue

      ctx.strokeStyle = `rgba(245, 158, 11, ${alpha})`
      ctx.lineWidth = 2.5
      ctx.beginPath()
      ctx.moveTo(p1.x, p1.y)
      ctx.lineTo(p2.x, p2.y)
      ctx.stroke()
    }
  }
}

function drawArrowHead(
  ctx: CanvasRenderingContext2D,
  tipX: number,
  tipY: number,
  dirX: number,
  dirY: number,
  color: string
) {
  const arrowSize = 5
  // dir 指向箭头指向的反方向（从 tip 往回画两翼）
  const angle = Math.atan2(dirY, dirX)
  ctx.fillStyle = color
  ctx.beginPath()
  ctx.moveTo(tipX, tipY)
  ctx.lineTo(
    tipX + Math.cos(angle - Math.PI / 6) * arrowSize,
    tipY + Math.sin(angle - Math.PI / 6) * arrowSize
  )
  ctx.lineTo(
    tipX + Math.cos(angle + Math.PI / 6) * arrowSize,
    tipY + Math.sin(angle + Math.PI / 6) * arrowSize
  )
  ctx.closePath()
  ctx.fill()
}

function drawNodes(ctx: CanvasRenderingContext2D, layout: Layout) {
  const nodes = graphData.value?.nodes || []
  const currentId = props.currentScene?.id

  const nodeRadius = BASE_NODE_RADIUS * Math.max(0.6, Math.min(1.4, zoomLevel.value))
  const currentRadius = BASE_CURRENT_RADIUS * Math.max(0.6, Math.min(1.4, zoomLevel.value))
  const fontSize = BASE_LABEL_FONT

  for (const node of nodes) {
    const pos = toCanvas(node.longitude, node.latitude, layout)
    const isCurrent = node.id === currentId
    const isHover = node.id === hoverNodeId.value
    const r = (isCurrent ? currentRadius : nodeRadius) * (isHover ? 1.4 : 1)

    if (isCurrent) {
      const pulse = (Date.now() % 1500) / 1500
      const pulseRadius = r + pulse * 6
      ctx.beginPath()
      ctx.arc(pos.x, pos.y, pulseRadius, 0, Math.PI * 2)
      ctx.fillStyle = `rgba(14, 165, 233, ${0.4 * (1 - pulse)})`
      ctx.fill()

      ctx.beginPath()
      ctx.arc(pos.x, pos.y, r, 0, Math.PI * 2)
      ctx.fillStyle = '#0EA5E9'
      ctx.fill()
      ctx.strokeStyle = '#fff'
      ctx.lineWidth = 2
      ctx.stroke()
    } else {
      ctx.beginPath()
      ctx.arc(pos.x, pos.y, r, 0, Math.PI * 2)
      ctx.fillStyle = isHover ? 'rgba(14, 165, 233, 0.9)' : 'rgba(255, 255, 255, 0.7)'
      ctx.fill()
      if (isHover) {
        ctx.strokeStyle = '#fff'
        ctx.lineWidth = 1.5
        ctx.stroke()
      }
    }

    ctx.fillStyle = 'rgba(255, 255, 255, 0.85)'
    ctx.font = `${fontSize}px sans-serif`
    ctx.textAlign = 'center'
    ctx.textBaseline = 'top'
    ctx.fillText(node.title, pos.x, pos.y + r + 4)
  }
}

function drawViewportFan(ctx: CanvasRenderingContext2D, layout: Layout) {
  if (!currentView.value || !props.currentScene) return

  const currentId = props.currentScene.id
  const currentNode = graphData.value?.nodes.find(n => n.id === currentId)
  if (!currentNode) return

  const pos = toCanvas(currentNode.longitude, currentNode.latitude, layout)
  const radius = Math.min(layout.width, layout.height) * BASE_FAN_RATIO * Math.max(0.6, Math.min(1.4, zoomLevel.value))

  const northOffsetRad = (props.currentScene.north_offset || 0) * Math.PI / 180
  const trueYaw = currentView.value.yaw + northOffsetRad
  const mapAngle = Math.PI / 2 - trueYaw
  const halfFov = currentView.value.fov / 2

  ctx.beginPath()
  ctx.moveTo(pos.x, pos.y)
  ctx.arc(pos.x, pos.y, radius, mapAngle - halfFov, mapAngle + halfFov)
  ctx.closePath()
  ctx.fillStyle = 'rgba(14, 165, 233, 0.25)'
  ctx.fill()

  ctx.beginPath()
  ctx.moveTo(pos.x, pos.y)
  ctx.lineTo(
    pos.x + Math.cos(mapAngle - halfFov) * radius,
    pos.y + Math.sin(mapAngle - halfFov) * radius
  )
  ctx.moveTo(pos.x, pos.y)
  ctx.lineTo(
    pos.x + Math.cos(mapAngle + halfFov) * radius,
    pos.y + Math.sin(mapAngle + halfFov) * radius
  )
  ctx.strokeStyle = 'rgba(14, 165, 233, 0.6)'
  ctx.lineWidth = 1
  ctx.stroke()
}

function draw() {
  const canvas = canvasRef.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  ensureViewerEvents()

  const layout = computeLayout()
  if (!layout) return

  drawBackground(ctx, layout.width, layout.height)

  if (loading.value || error.value) return

  const nodes = graphData.value?.nodes || []
  if (nodes.length === 0) {
    drawPlaceholder(ctx, layout.width, layout.height)
    return
  }

  drawEdges(ctx, layout)
  drawNodes(ctx, layout)
  drawViewportFan(ctx, layout)
}

function startAnimationLoop() {
  if (animationLoopId) return
  const loop = () => {
    draw()
    animationLoopId = requestAnimationFrame(loop)
  }
  animationLoopId = requestAnimationFrame(loop)
}

function stopAnimationLoop() {
  if (animationLoopId) {
    cancelAnimationFrame(animationLoopId)
    animationLoopId = null
  }
}

function syncCurrentView() {
  const viewer = props.getViewer()
  if (!viewer) return

  const position = viewer.getPosition()
  const zoomLevel = viewer.getZoomLevel()
  const fovDeg = 90 - (90 - 30) * (zoomLevel / 100)

  currentView.value = {
    yaw: position.yaw,
    fov: fovDeg * Math.PI / 180
  }
}

function ensureViewerEvents() {
  const viewer = props.getViewer()
  if (!viewer) {
    if (boundViewer) unbindViewerEvents()
    return
  }

  if (boundViewer === viewer && positionHandler) return

  if (boundViewer && positionHandler) {
    boundViewer.removeEventListener('position-updated', positionHandler)
  }

  boundViewer = viewer

  positionHandler = (_e: any) => {
    const v = props.getViewer()
    if (!v) return
    const position = v.getPosition()
    const zoomLevel = v.getZoomLevel()
    const fovDeg = 90 - (90 - 30) * (zoomLevel / 100)
    currentView.value = {
      yaw: position.yaw,
      fov: fovDeg * Math.PI / 180
    }
  }

  viewer.addEventListener('position-updated', positionHandler)
}

function unbindViewerEvents() {
  if (boundViewer && positionHandler) {
    boundViewer.removeEventListener('position-updated', positionHandler)
  }
  boundViewer = null
  positionHandler = null
}

function handleCanvasClick(event: MouseEvent) {
  if (dragMoved) {
    dragMoved = false
    return
  }
  const canvas = canvasRef.value
  if (!canvas || !graphData.value?.nodes.length) return

  const rect = canvas.getBoundingClientRect()
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top

  const layout = computeLayout()
  if (!layout) return

  const hitRadius = 14 * Math.max(0.7, Math.min(1.5, zoomLevel.value))
  const nodes = graphData.value.nodes
  for (const node of nodes) {
    const pos = toCanvas(node.longitude, node.latitude, layout)
    const dist = Math.hypot(pos.x - x, pos.y - y)
    if (dist <= hitRadius) {
      emit('scene-select', node.id)
      return
    }
  }
}

function clampZoom(value: number) {
  return Math.max(MIN_ZOOM, Math.min(MAX_ZOOM, value))
}

function resetView() {
  zoomLevel.value = 1
  panX.value = 0
  panY.value = 0
}

function handleWheel(event: WheelEvent) {
  event.preventDefault()
  const canvas = canvasRef.value
  if (!canvas) return

  const rect = canvas.getBoundingClientRect()
  const mouseX = event.clientX - rect.left
  const mouseY = event.clientY - rect.top

  const oldZoom = zoomLevel.value
  const factor = event.deltaY < 0 ? 1.1 : 1 / 1.1
  const newZoom = clampZoom(oldZoom * factor)
  if (newZoom === oldZoom) return

  const centerX = rect.width / 2
  const centerY = rect.height / 2
  const worldX = (mouseX - centerX - panX.value) / oldZoom
  const worldY = (mouseY - centerY - panY.value) / oldZoom

  zoomLevel.value = newZoom
  panX.value = mouseX - centerX - worldX * newZoom
  panY.value = mouseY - centerY - worldY * newZoom
}

function handleMouseDown(event: MouseEvent) {
  if (event.button !== 0) return
  isPanning.value = true
  dragMoved = false
  panStartX = event.clientX
  panStartY = event.clientY
  panOriginX = panX.value
  panOriginY = panY.value
}

function handleMouseMove(event: MouseEvent) {
  if (isPanning.value) {
    const dx = event.clientX - panStartX
    const dy = event.clientY - panStartY
    if (Math.hypot(dx, dy) > 3) dragMoved = true
    panX.value = panOriginX + dx
    panY.value = panOriginY + dy
    return
  }

  updateHover(event)
}

function updateHover(event: MouseEvent) {
  const canvas = canvasRef.value
  if (!canvas || !graphData.value?.nodes.length) {
    if (hoverNodeId.value !== null) hoverNodeId.value = null
    return
  }

  const rect = canvas.getBoundingClientRect()
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top

  const layout = computeLayout()
  if (!layout) return

  const hitRadius = 14 * Math.max(0.7, Math.min(1.5, zoomLevel.value))
  let foundId: number | null = null
  for (const node of graphData.value.nodes) {
    const pos = toCanvas(node.longitude, node.latitude, layout)
    const dist = Math.hypot(pos.x - x, pos.y - y)
    if (dist <= hitRadius) {
      foundId = node.id
      break
    }
  }

  if (hoverNodeId.value !== foundId) {
    hoverNodeId.value = foundId
    canvas.style.cursor = foundId !== null ? 'pointer' : 'grab'
  }
}

function handleMouseUp() {
  isPanning.value = false
}

onMounted(() => {
  setupCanvas()
  if (props.visible && props.currentScene) {
    startAnimationLoop()
    loadGraphData(props.currentScene.space_id).then(() => {
      syncCurrentView()
      ensureViewerEvents()
      draw()
    })
  }
})

onUnmounted(() => {
  unbindViewerEvents()
  stopAnimationLoop()
})

watch(() => props.visible, (visible) => {
  if (visible && props.currentScene) {
    nextTick(() => {
      setupCanvas()
      startAnimationLoop()
      if (!graphData.value || graphData.value.space_info.id !== props.currentScene!.space_id) {
        loadGraphData(props.currentScene!.space_id).then(() => {
          syncCurrentView()
          ensureViewerEvents()
          draw()
        })
      } else {
        syncCurrentView()
        ensureViewerEvents()
        draw()
      }
    })
  } else if (!visible) {
    unbindViewerEvents()
    stopAnimationLoop()
  }
})

watch(() => props.currentScene, (scene, oldScene) => {
  if (!scene || !props.visible) return

  startAnimationLoop()
  if (!oldScene || scene.space_id !== oldScene.space_id) {
    loadGraphData(scene.space_id).then(() => {
      syncCurrentView()
      ensureViewerEvents()
      draw()
    })
  } else {
    syncCurrentView()
    draw()
  }
})

watch(graphData, () => {
  draw()
})
</script>

<style scoped>
.pano-minimap {
  position: absolute;
  top: 72px;
  right: 24px;
  width: 280px;
  height: 200px;
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.18);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
  z-index: 90;
  overflow: hidden;
}

.minimap-header {
  height: 40px;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 8px 0 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.15);
  flex-shrink: 0;
  background: rgba(255, 255, 255, 0.04);
}

.minimap-title {
  font-size: 13px;
  font-weight: 500;
  color: #fff;
  white-space: nowrap;
  letter-spacing: 0.2px;
}

.minimap-zoom-controls {
  display: flex;
  align-items: center;
  gap: 2px;
  margin-left: auto;
}

.minimap-zoom-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.75);
  min-width: 38px;
  text-align: center;
  font-variant-numeric: tabular-nums;
  user-select: none;
}

.minimap-icon-btn,
.minimap-close {
  background: transparent !important;
  border: none !important;
  color: rgba(255, 255, 255, 0.8) !important;
  font-size: 16px !important;
  line-height: 1;
  min-width: 22px !important;
  height: 22px !important;
  padding: 0 !important;
  border-radius: 4px !important;
  transition: color 0.15s ease, background 0.15s ease;
}

.minimap-icon-btn:hover,
.minimap-close:hover {
  color: #fff !important;
  background: rgba(255, 255, 255, 0.12) !important;
}

.minimap-close {
  font-size: 18px !important;
}

.minimap-body {
  position: relative;
  flex: 1;
  min-height: 0;
}

.minimap-canvas {
  display: block;
  width: 100%;
  height: 100%;
  cursor: grab;
  touch-action: none;
}

.minimap-canvas:active {
  cursor: grabbing;
}

.minimap-hint {
  position: absolute;
  bottom: 4px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.45);
  pointer-events: none;
  white-space: nowrap;
  user-select: none;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

.minimap-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.25);
  backdrop-filter: blur(2px);
  -webkit-backdrop-filter: blur(2px);
  pointer-events: none;
}

.minimap-error {
  color: rgba(255, 120, 120, 0.95);
  font-size: 12px;
  padding: 0 16px;
  text-align: center;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

.minimap-fade-enter-active,
.minimap-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.minimap-fade-enter-from,
.minimap-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
