<template>
  <div class="canvas-wrapper">
    <div ref="containerRef" class="graph-container"></div>
    <div v-if="layoutComputing.isLayoutComputing.value" class="computing-overlay">
      <div class="computing-spinner"></div>
      <p class="computing-text">正在计算布局...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, provide, reactive, watch, nextTick, inject } from 'vue'
import { Graph as X6Graph, History, Selection, Snapline } from '@antv/x6'
import { ForceLayout } from '@antv/layout'
import { useGraphSceneStore } from '@/stores/graphSceneStore'
import { useGraphEditorStore } from '@/stores/graphEditorStore'
import { SceneNodeStatus } from '@/models/graphEditor/node'
import { createDnd, DND_ACTIONS_KEY } from '@/composables/useGraphDnd'
import { useGraphKeyboard } from '@/composables/useGraphKeyboard'
import { NODE_SIZE_PANO, PORT_RADIUS, COLORS, GRID_CONFIG, ZOOM_LIMITS, HISTORY_STACK_SIZE, SNAPLINE_TOLERANCE } from '@/utils/graphEditor/constants'

const containerRef = ref<HTMLElement | null>(null)
const sceneStore = useGraphSceneStore()
const editorStore = useGraphEditorStore()

let graph: X6Graph | null = null
const canUndo = ref(false)
const canRedo = ref(false)

interface LayoutComputingState {
  isLayoutComputing: { value: boolean }
  setLayoutComputing: (val: boolean) => void
}
const layoutComputing = inject<LayoutComputingState>('layoutComputingState', {
  isLayoutComputing: { value: false },
  setLayoutComputing: () => {}
})

const dndActions = reactive<{ startDrag: ((nodeData: any, e: MouseEvent) => void) | null }>({ startDrag: null })
provide(DND_ACTIONS_KEY, dndActions as any)

const graphActionsObj = reactive({
  getGraph: () => graph,
  undo: () => graph?.undo(),
  redo: () => graph?.redo(),
  zoomIn: () => graph?.zoom(ZOOM_LIMITS.step),
  zoomOut: () => graph?.zoom(-ZOOM_LIMITS.step),
  fitCanvas: () => graph?.zoomToFit({ maxScale: 1 }),
  get canUndo() { return canUndo.value },
  get canRedo() { return canRedo.value },
  uploadBackground(file: File) {
    const reader = new FileReader()
    reader.onload = (e) => {
      graph?.drawBackground({
        color: COLORS.canvas.bg,
        image: e.target?.result as string,
        position: 'center',
        size: 'contain',
        opacity: 0.5,
      })
    }
    reader.readAsDataURL(file)
  },
})
provide('graphActions', graphActionsObj)

const LNG_MIN = 103.605
const LNG_MAX = 103.616
const LAT_MIN = 30.995
const LAT_MAX = 31.011

const CANVAS_PADDING = 100
const CANVAS_WIDTH = 1200
const CANVAS_HEIGHT = 800

function geoToCanvas(longitude: number, latitude: number): { x: number; y: number } {
  const x = CANVAS_PADDING + ((longitude - LNG_MIN) / (LNG_MAX - LNG_MIN)) * (CANVAS_WIDTH - CANVAS_PADDING * 2)
  const y = CANVAS_PADDING + ((LAT_MAX - latitude) / (LAT_MAX - LAT_MIN)) * (CANVAS_HEIGHT - CANVAS_PADDING * 2)
  return { x, y }
}

function getPortGroups() {
  return {
    top: {
      position: 'top',
      attrs: {
        circle: {
          r: PORT_RADIUS,
          magnet: true,
          stroke: COLORS.port.stroke,
          strokeWidth: 2,
          fill: COLORS.port.fill,
          style: { visibility: 'hidden' },
        },
      },
    },
    bottom: {
      position: 'bottom',
      attrs: {
        circle: {
          r: PORT_RADIUS,
          magnet: true,
          stroke: COLORS.port.stroke,
          strokeWidth: 2,
          fill: COLORS.port.fill,
          style: { visibility: 'hidden' },
        },
      },
    },
    left: {
      position: 'left',
      attrs: {
        circle: {
          r: PORT_RADIUS,
          magnet: true,
          stroke: COLORS.port.stroke,
          strokeWidth: 2,
          fill: COLORS.port.fill,
          style: { visibility: 'hidden' },
        },
      },
    },
    right: {
      position: 'right',
      attrs: {
        circle: {
          r: PORT_RADIUS,
          magnet: true,
          stroke: COLORS.port.stroke,
          strokeWidth: 2,
          fill: COLORS.port.fill,
          style: { visibility: 'hidden' },
        },
      },
    },
  }
}

interface LayoutNode {
  id: number
  size: number
  x: number
  y: number
  title: string
  has_position: boolean
  status: string
}

interface LayoutEdge {
  source: number
  target: number
}

let layoutData: { nodes: LayoutNode[]; edges: LayoutEdge[] } = { nodes: [], edges: [] }

function buildLayoutData() {
  const nodes: LayoutNode[] = []
  const edges: LayoutEdge[] = []

  const container = containerRef.value
  const areaWidth = container ? container.clientWidth : 1200
  const areaHeight = container ? container.clientHeight : 800

  sceneStore.placedNodes.forEach((nodeData) => {
    let x = 0
    let y = 0

    if (nodeData.longitude && nodeData.latitude && nodeData.longitude !== 0 && nodeData.latitude !== 0) {
      const pos = geoToCanvas(nodeData.longitude, nodeData.latitude)
      x = pos.x
      y = pos.y
    } else {
      x = 100 + Math.random() * (areaWidth - 200)
      y = 100 + Math.random() * (areaHeight - 200)
    }

    nodes.push({
      id: nodeData.id,
      size: NODE_SIZE_PANO,
      x,
      y,
      title: nodeData.title,
      has_position: nodeData.has_position,
      status: nodeData.status,
    })
  })

  sceneStore.edges?.forEach((edgeData) => {
    edges.push({
      source: edgeData.source_id,
      target: edgeData.target_id,
    })
  })

  layoutData = { nodes, edges }
}

function getModelFromLayoutData() {
  const model: any = {
    nodes: [],
    edges: [],
  }

  layoutData.nodes.forEach((item) => {
    model.nodes.push({
      id: String(item.id),
      shape: 'circle',
      width: item.size,
      height: item.size,
      x: item.x,
      y: item.y,
      attrs: {
        body: {
          fill: COLORS.node.hasPano,
          stroke: 'transparent',
        },
        label: {
          text: item.title || '',
          fill: '#374151',
          fontSize: 11,
          fontWeight: 500,
          refX: 0,
          refY: item.size / 2 + 14,
          textAnchor: 'middle',
          textVerticalAnchor: 'top',
        },
      },
      data: {
        title: item.title,
        has_position: item.has_position,
        status: item.status,
      },
      ports: {
        groups: getPortGroups(),
        items: [
          { group: 'top', id: 'port-top' },
          { group: 'bottom', id: 'port-bottom' },
          { group: 'left', id: 'port-left' },
          { group: 'right', id: 'port-right' },
        ],
      },
    })
  })

  layoutData.edges.forEach((item) => {
    model.edges.push({
      source: String(item.source),
      target: String(item.target),
      attrs: {
        line: {
          stroke: COLORS.edge.default,
          strokeWidth: 2,
          targetMarker: null,
        },
      },
    })
  })

  return model
}

function initGraph() {
  if (!containerRef.value) return

  graph = new X6Graph({
    container: containerRef.value,
    autoResize: true,
    grid: GRID_CONFIG,
    background: { color: COLORS.canvas.bg },
    panning: {
      enabled: true,
      modifiers: ['shift'],
    },
    interacting: {
      nodeMovable: true,
      edgeMovable: true,
      arrowheadMovable: true,
      vertexMovable: true,
      vertexAddable: false,
      vertexDeletable: false,
    },
    mousewheel: {
      enabled: true,
      factor: ZOOM_LIMITS.factor,
      minScale: ZOOM_LIMITS.min,
      maxScale: ZOOM_LIMITS.max,
    },
    connecting: {
      snap: true,
      allowBlank: false,
      allowLoop: false,
      allowMulti: false,
      allowNode: false,
      allowPort: true,
      highlight: true,
      createEdge() {
        return (this as any).createEdge({
          shape: 'edge',
          attrs: {
            line: {
              stroke: COLORS.edge.default,
              strokeWidth: 2,
              targetMarker: {
                name: 'block',
                width: 12,
                height: 8,
              },
            },
          },
          router: {
            name: 'manhattan',
            args: { padding: 20 },
          },
          connector: {
            name: 'rounded',
            args: { radius: 8 },
          },
          data: {
            type: 'walk',
          },
        })
      },
      validateConnection({ sourcePort, targetPort }) {
        if (sourcePort === targetPort) return false
        return true
      },
    },
    highlighting: {
      magnetAvailable: {
        name: 'stroke',
        args: {
          attrs: {
            fill: '#fff',
            stroke: COLORS.port.active,
            strokeWidth: 3,
          },
        },
      },
      magnetAdsorbed: {
        name: 'stroke',
        args: {
          attrs: {
            fill: '#fff',
            stroke: COLORS.node.default,
            strokeWidth: 3,
          },
        },
      },
    },
  })

  graph.use(new History({ enabled: true, stackSize: HISTORY_STACK_SIZE }))
  graph.use(new Selection({
    enabled: true,
    rubberband: true,
    multiple: true,
    showNodeSelectionBox: true,
    modifiers: [],
  }))
  graph.use(new Snapline({
    enabled: true,
    tolerance: SNAPLINE_TOLERANCE,
    sharp: true,
  }))

  const dnd = createDnd(graph)
  dndActions.startDrag = dnd.startDrag
  useGraphKeyboard()

  setupGraphEvents()

  nextTick(() => {
    if (sceneStore.placedNodes.length > 0) {
      renderWithForceLayout()
    }
  })

  watch(() => sceneStore.placedNodes.length, (newLen) => {
    if (newLen > 0 && graph && graph.getNodes().length === 0) {
      renderWithForceLayout()
    }
  })
}

function resolveOverlaps(nodes: LayoutNode[], minDistance: number) {
  const sorted = [...nodes].sort((a, b) => a.x - b.x)
  for (let i = 0; i < sorted.length; i++) {
    for (let j = i + 1; j < sorted.length; j++) {
      const dx = sorted[j].x - sorted[i].x
      const dy = sorted[j].y - sorted[i].y
      const dist = Math.sqrt(dx * dx + dy * dy)
      if (dist < minDistance) {
        const angle = Math.atan2(dy || 1, dx || 1)
        const pushX = (minDistance - dist) * Math.cos(angle)
        const pushY = (minDistance - dist) * Math.sin(angle)
        sorted[j].x += pushX
        sorted[j].y += pushY
      }
    }
  }
  sorted.forEach(n => {
    const orig = nodes.find(o => o.id === n.id)
    if (orig) {
      orig.x = n.x
      orig.y = n.y
    }
  })
}

async function renderWithForceLayout() {
  if (!graph) return

  layoutComputing.setLayoutComputing(true)
  buildLayoutData()

  const container = containerRef.value!
  const centerX = container.clientWidth / 2
  const centerY = container.clientHeight / 2

  const forceLayout = new ForceLayout({
    center: [centerX, centerY],
    preventOverlap: true,
    nodeSize: 80,
    linkDistance: () => 200,
    nodeStrength: () => -500,
    edgeStrength: () => 0.1,
    maxIteration: 500,
  })

  const layoutNodes = layoutData.nodes.map(n => ({
    id: n.id,
    data: {
      x: n.x,
      y: n.y,
      size: n.size,
    },
  }))

  const layoutEdges = layoutData.edges.map((e, i) => ({
    id: `edge-${i}`,
    source: e.source,
    target: e.target,
    data: {},
  }))

  const layoutGraphData = {
    nodes: layoutNodes,
    edges: layoutEdges,
  }

  try {
    const result = await forceLayout.execute(layoutGraphData as any) as any

    if (result && result.nodes) {
      result.nodes.forEach((node: any) => {
        const layoutNode = layoutData.nodes.find(n => n.id === node.id)
        if (layoutNode && node.data) {
          layoutNode.x = node.data.x
          layoutNode.y = node.data.y
        }
      })
    }

    resolveOverlaps(layoutData.nodes, 100)

    const model = getModelFromLayoutData()
    graph!.fromJSON(model)
    graph!.zoomToFit({ padding: 80 })
    editorStore.setZoom(graph!.zoom())
  } catch (error) {
    console.warn('ForceLayout execution failed, using default positions:', error)
    const model = getModelFromLayoutData()
    graph!.fromJSON(model)
    graph!.zoomToFit({ padding: 80 })
    editorStore.setZoom(graph!.zoom())
  } finally {
    layoutComputing.setLayoutComputing(false)
  }
}

function setupGraphEvents() {
  if (!graph) return

  graph.on('node:mouseenter', ({ node }) => {
    const ports = node.getPorts()
    ports.forEach(port => {
      node.setPortProp(port.id!, 'attrs/circle/style/visibility', 'visible')
    })
  })

  graph.on('node:mouseleave', ({ node }) => {
    const ports = node.getPorts()
    ports.forEach(port => {
      node.setPortProp(port.id!, 'attrs/circle/style/visibility', 'hidden')
    })
  })

  graph.on('node:click', ({ node }) => {
    editorStore.selectNode(node.id)
  })

  graph.on('edge:click', ({ edge }) => {
    editorStore.selectEdge(edge.id)
  })

  graph.on('edge:connected', ({ edge }) => {
    const sourceCell = edge.getSourceCell()
    const targetCell = edge.getTargetCell()
    if (sourceCell && targetCell) {
      const sourceId = Number(sourceCell.id)
      const targetId = Number(targetCell.id)
      const edgeData = edge.getData() || {}
      sceneStore.addEdge({
        source_id: sourceId,
        target_id: targetId,
        hotspot_id: 0,
        hotspot_title: '',
        type: edgeData.type || 'walk',
      })
    }
  })

  graph.on('blank:click', () => {
    editorStore.clearSelection()
  })

  graph.on('node:added', ({ node }) => {
    const data = node.getData()
    if (data && data.status === SceneNodeStatus.PENDING) {
      sceneStore.markNodePlaced(Number(node.id))
      node.setData({ ...data, status: SceneNodeStatus.PLACED })
    }
  })

  graph.on('node:removed', ({ node }) => {
    sceneStore.markNodePending(Number(node.id))
  })

  graph.on('edge:removed', ({ edge }) => {
    sceneStore.removeEdge(Number(edge.id))
  })

  graph.on('scale', () => {
    if (graph) {
      editorStore.setZoom(graph.zoom())
    }
  })

  graph.on('history:change', () => {
    if (graph) {
      canUndo.value = graph.canUndo()
      canRedo.value = graph.canRedo()
    }
  })
}

onMounted(() => {
  initGraph()
})

onUnmounted(() => {
  if (graph) {
    graph.dispose()
    graph = null
  }
})
</script>

<style scoped>
.canvas-wrapper {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.graph-container {
  width: 100%;
  height: 100%;
}

.computing-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.85);
  z-index: 100;
}

.computing-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e5e7eb;
  border-top-color: var(--color-primary, #4f46e5);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

.computing-text {
  margin-top: 12px;
  font-size: 13px;
  color: #6b7280;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
