<template>
  <div class="scene-topology">
    <canvas ref="canvasRef" class="topology-canvas"></canvas>
    <div class="topology-tooltip" v-if="hoveredNode" :style="tooltipStyle">
      <div class="tooltip-title">{{ hoveredNode.title }}</div>
      <div class="tooltip-views">{{ hoveredNode.view_count || 0 }} 次浏览</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import type { SceneNode } from '@/models/scene.model'

interface Props {
  scenes: SceneNode[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'sceneClick', scene: SceneNode): void
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
const hoveredNode = ref<SceneNode | null>(null)
const tooltipStyle = ref({ left: '0px', top: '0px' })

const NODE_RADIUS = 24
const NODE_COLORS = ['#0EA5E9', '#2563EB', '#10B981', '#F59E0B', '#8B5CF6', '#EC4899', '#DC2626', '#059669']

interface NodePosition {
  scene: SceneNode
  x: number
  y: number
  color: string
}

let nodePositions: NodePosition[] = []
let animationFrame: number | null = null

function buildAdjacencyList(scenes: SceneNode[]): Map<number, number[]> {
  const adj = new Map<number, number[]>()
  scenes.forEach(scene => {
    scene.hotspots?.forEach(hotspot => {
      if (hotspot.target_scene_id) {
        const list = adj.get(scene.id) || []
        list.push(hotspot.target_scene_id)
        adj.set(scene.id, list)
      }
    })
  })
  return adj
}

function calculatePositions(width: number, height: number) {
  if (!props.scenes || props.scenes.length === 0) {
    nodePositions = []
    return
  }

  const centerX = width / 2
  const centerY = height / 2
  const radius = Math.min(width, height) / 3

  nodePositions = props.scenes.map((scene, index) => {
    const angle = (2 * Math.PI * index) / props.scenes.length - Math.PI / 2
    return {
      scene,
      x: centerX + radius * Math.cos(angle),
      y: centerY + radius * Math.sin(angle),
      color: NODE_COLORS[index % NODE_COLORS.length],
    }
  })
}

function drawTopology() {
  const canvas = canvasRef.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const width = canvas.width
  const height = canvas.height

  ctx.clearRect(0, 0, width, height)

  const adj = buildAdjacencyList(props.scenes)

  ctx.strokeStyle = 'rgba(14, 165, 233, 0.3)'
  ctx.lineWidth = 1
  ctx.setLineDash([5, 5])

  adj.forEach((targets, sourceId) => {
    const sourcePos = nodePositions.find(n => n.scene.id === sourceId)
    if (!sourcePos) return

    targets.forEach(targetId => {
      const targetPos = nodePositions.find(n => n.scene.id === targetId)
      if (!targetPos) return

      ctx.beginPath()
      ctx.moveTo(sourcePos.x, sourcePos.y)
      ctx.lineTo(targetPos.x, targetPos.y)
      ctx.stroke()

      const angle = Math.atan2(targetPos.y - sourcePos.y, targetPos.x - sourcePos.x)
      const arrowX = targetPos.x - NODE_RADIUS * Math.cos(angle)
      const arrowY = targetPos.y - NODE_RADIUS * Math.sin(angle)
      const arrowSize = 6

      ctx.beginPath()
      ctx.moveTo(arrowX, arrowY)
      ctx.lineTo(
        arrowX - arrowSize * Math.cos(angle - Math.PI / 6),
        arrowY - arrowSize * Math.sin(angle - Math.PI / 6)
      )
      ctx.lineTo(
        arrowX - arrowSize * Math.cos(angle + Math.PI / 6),
        arrowY - arrowSize * Math.sin(angle + Math.PI / 6)
      )
      ctx.closePath()
      ctx.fillStyle = 'rgba(14, 165, 233, 0.5)'
      ctx.fill()
    })
  })

  ctx.setLineDash([])

  nodePositions.forEach((pos, index) => {
    const isHovered = hoveredNode.value?.id === pos.scene.id

    ctx.beginPath()
    ctx.arc(pos.x, pos.y, NODE_RADIUS + (isHovered ? 4 : 0), 0, 2 * Math.PI)
    ctx.fillStyle = isHovered ? '#1a4a6e' : '#1E3A5F'
    ctx.fill()
    ctx.strokeStyle = pos.color
    ctx.lineWidth = isHovered ? 3 : 2
    ctx.stroke()

    ctx.fillStyle = '#fff'
    ctx.font = '11px "Noto Sans SC", sans-serif'
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillText(String(index + 1), pos.x, pos.y)

    ctx.fillStyle = 'rgba(255, 255, 255, 0.8)'
    ctx.font = '12px "Noto Sans SC", sans-serif'
    ctx.textAlign = 'center'
    ctx.textBaseline = 'top'
    const title = pos.scene.title.length > 8 ? pos.scene.title.slice(0, 8) + '...' : pos.scene.title
    ctx.fillText(title, pos.x, pos.y + NODE_RADIUS + 8)
  })
}

function handleMouseMove(e: MouseEvent) {
  const canvas = canvasRef.value
  if (!canvas) return

  const rect = canvas.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top

  let found: NodePosition | null = null
  for (const pos of nodePositions) {
    const dx = x - pos.x
    const dy = y - pos.y
    if (dx * dx + dy * dy <= NODE_RADIUS * NODE_RADIUS) {
      found = pos
      break
    }
  }

  if (found) {
    hoveredNode.value = found.scene
    tooltipStyle.value = {
      left: `${e.clientX + 10}px`,
      top: `${e.clientY + 10}px`,
    }
    canvas.style.cursor = 'pointer'
  } else {
    hoveredNode.value = null
    canvas.style.cursor = 'default'
  }

  drawTopology()
}

function handleClick(e: MouseEvent) {
  const canvas = canvasRef.value
  if (!canvas) return

  const rect = canvas.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top

  for (const pos of nodePositions) {
    const dx = x - pos.x
    const dy = y - pos.y
    if (dx * dx + dy * dy <= NODE_RADIUS * NODE_RADIUS) {
      emit('sceneClick', pos.scene)
      break
    }
  }
}

function resizeCanvas() {
  const canvas = canvasRef.value
  if (!canvas) return

  const container = canvas.parentElement
  if (!container) return

  canvas.width = container.clientWidth
  canvas.height = container.clientHeight

  calculatePositions(canvas.width, canvas.height)
  drawTopology()
}

watch(() => props.scenes, () => {
  nextTick(() => {
    resizeCanvas()
  })
}, { deep: true })

onMounted(() => {
  const canvas = canvasRef.value
  if (!canvas) return

  canvas.addEventListener('mousemove', handleMouseMove)
  canvas.addEventListener('click', handleClick)
  window.addEventListener('resize', resizeCanvas)

  nextTick(() => {
    resizeCanvas()
  })
})

onUnmounted(() => {
  const canvas = canvasRef.value
  if (canvas) {
    canvas.removeEventListener('mousemove', handleMouseMove)
    canvas.removeEventListener('click', handleClick)
  }
  window.removeEventListener('resize', resizeCanvas)
  if (animationFrame) {
    cancelAnimationFrame(animationFrame)
  }
})
</script>

<style scoped>
.scene-topology {
  position: relative;
  width: 100%;
}

.topology-canvas {
  width: 100%;
  height: 300px;
  background: rgba(30, 58, 95, 0.3);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.topology-tooltip {
  position: fixed;
  background: rgba(26, 35, 50, 0.95);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 10px 14px;
  pointer-events: none;
  z-index: 1000;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
}

.tooltip-title {
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.tooltip-views {
  color: rgba(255, 255, 255, 0.5);
  font-size: 12px;
}
</style>
