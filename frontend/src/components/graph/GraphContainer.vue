<template>
  <div class="graph-container" ref="containerRef"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { Graph, Cell, Color } from '@antv/x6'
import type { GraphDataResponse, SceneNodeData, EdgeData } from '@/models/SceneModel'
import { CoordinateTransform } from '@/utils/coordinateTransform'

interface Props {
  graphData: GraphDataResponse | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'nodeClick', node: SceneNodeData): void
  (e: 'nodeMove', nodeId: number, longitude: number, latitude: number): void
  (e: 'edgeClick', edge: EdgeData): void
  (e: 'blankClick'): void
  (e: 'zoomChange', zoom: number): void
}>()

const containerRef = ref<HTMLElement | null>(null)
let graph: Graph | null = null
let coordinateTransform: CoordinateTransform | null = null
const nodeDataMap = new Map<string, SceneNodeData>()
const edgeDataMap = new Map<string, EdgeData>()

const initGraph = () => {
  if (!containerRef.value) return

  const width = containerRef.value.offsetWidth
  const height = containerRef.value.offsetHeight

  graph = new Graph({
    container: containerRef.value,
    width,
    height,
    background: {
      color: '#f5f7fa'
    },
    grid: {
      visible: true,
      type: 'dot',
      size: 20,
      args: {
        color: '#e0e0e0',
        thickness: 1
      }
    },
    panning: {
      enabled: true,
      modifiers: []
    },
    mousewheel: {
      enabled: true,
      modifiers: [],
      minScale: 0.2,
      maxScale: 4
    },
    connecting: {
      anchor: 'center',
      connectionPoint: 'anchor',
      allowBlank: false,
      allowLoop: false,
      allowNode: true,
      allowEdge: false,
      highlight: true,
      snap: true,
      createEdge() {
        return graph?.createEdge({
          shape: 'edge',
          attrs: {
            line: {
              stroke: '#0052d9',
              strokeWidth: 2,
              targetMarker: {
                name: 'block',
                width: 8,
                height: 6
              }
            }
          }
        })
      }
    },
    highlighting: {
      magnetAvailable: {
        name: 'stroke',
        args: {
          attrs: {
            fill: '#fff',
            stroke: '#0052d9',
            strokeWidth: 4
          }
        }
      }
    }
  })

  graph.on('node:click', ({ node }) => {
    const data = nodeDataMap.get(node.id)
    if (data) {
      emit('nodeClick', data)
    }
  })

  graph.on('node:moved', ({ node }) => {
    const data = nodeDataMap.get(node.id)
    if (data && coordinateTransform) {
      const pos = node.getPosition()
      const { lng, lat } = coordinateTransform.screenToLngLat(pos.x + 4, pos.y + 4)
      emit('nodeMove', data.id, lng, lat)
    }
  })

  graph.on('node:mouseenter', ({ node }) => { // 节点鼠标悬停事件
    node.addTools({
      name: 'button',
      args: {
        markup: [
          {
            tagName: 'circle',
            selector: 'button',
            attrs: {
              r: 10, 
              stroke: '#fe854f',
              strokeWidth: 2,
              fill: 'white',
              cursor: 'pointer'
            }
          },
          {
            tagName: 'text',
            textContent: 'Btn',
            selector: 'icon',
            attrs: {
              fill: '#fe854f',
              fontSize: 10, // 文字大小
              textAnchor: 'middle',
              pointerEvents: 'none',
              y: '0.3em'
            }
          }
        ],
        x: 0,
        y: 0,
        offset: { x: 10, y: 10 },
        onClick({ cell }: { cell: Cell }) {
          const fill = Color.randomHex()
          cell.attr({
            body: {
              fill
            },
            label: {
              fill: Color.invert(fill, true)
            }
          })
        }
      }
    })
  })

  graph.on('node:mouseleave', ({ cell }) => {
    cell.removeTools()
  })

  graph.on('edge:click', ({ edge }) => {
    const data = edgeDataMap.get(edge.id)
    if (data) {
      emit('edgeClick', data)
    }
  })

  graph.on('blank:click', () => {
    emit('blankClick')
  })

  graph.on('scale', ({ sx }) => {
    emit('zoomChange', sx)
  })
}

const renderGraph = () => {
  if (!graph || !props.graphData) return

  graph.clearCells()
  nodeDataMap.clear()
  edgeDataMap.clear()

  const { nodes, edges } = props.graphData

  if (nodes.length === 0) return

  const bounds = CoordinateTransform.calculateBounds(nodes)
  const center = CoordinateTransform.calculateCenter(bounds)

  const width = containerRef.value?.offsetWidth || 800
  const height = containerRef.value?.offsetHeight || 600

  coordinateTransform = new CoordinateTransform()
  coordinateTransform.setReference(center.lng, center.lat)
  coordinateTransform.setOffset(width / 2, height / 2)

  const optimalScale = CoordinateTransform.calculateOptimalScale(bounds, width, height, 100)
  coordinateTransform.setScale(optimalScale*5)

  nodes.forEach(node => {
    const { x, y } = coordinateTransform!.lngLatToScreen(node.longitude, node.latitude)

    const graphNode = graph!.addNode({ // 节点
      id: `node-${node.id}`,
      shape: 'circle',
      x: x - 15, // 节点中心坐标
      y: y - 15, // 节点中心坐标
      width: 30,
      height: 30,
      attrs: {
        body: {
          fill: '#0052d9',
          stroke: '#0052d9',
          strokeWidth: 1
        },
        label: {
          text: node.title,
          fill: '#1d2129',
          fontSize: 16,
          fontWeight: 500,
          textAnchor: 'middle',
          textVerticalAnchor: 'bottom',
          refX: '50%', // 文字水平居中
          refY: '0%', // 文字垂直居中
        }
      }
    })

    nodeDataMap.set(graphNode.id, node)
  })

  edges.forEach(edge => {
    const sourceNode = graph!.getCellById(`node-${edge.source_id}`)
    const targetNode = graph!.getCellById(`node-${edge.target_id}`)

    if (sourceNode && targetNode) {
      const graphEdge = graph!.addEdge({
        id: `edge-${edge.id}`,
        source: { cell: sourceNode.id },
        target: { cell: targetNode.id },
        attrs: {
          line: {
            stroke: '#4080ff',
            strokeWidth: 2,
            strokeDasharray: '10 5',//虚线
            targetMarker: {
              name: 'block',
              width: 8,
              height: 6
            }
          }
        },
        labels: edge.hotspot_title ? [
          {
            attrs: {
              text: {
                text: edge.hotspot_title,
                fill: '#666',
                fontSize: 10
              },
              rect: {
                fill: '#fff',
                stroke: '#ddd',
                strokeWidth: 1,
                rx: 4,
                ry: 4
              }
            },
            position: {
              distance: 0.5,
              offset: { x: 0, y: -10 }
            }
          }
        ] : []
      })

      edgeDataMap.set(graphEdge.id, edge)
    }
  })

  graph!.centerContent()
}

const zoomIn = () => {
  graph?.zoom(0.1)
}

const zoomOut = () => {
  graph?.zoom(-0.1)
}

const zoomTo = (scale: number) => {
  graph?.zoomTo(scale)
}

const fitContent = () => {
  graph?.zoomToFit({ padding: 50 })
}

const handleResize = () => {
  if (graph && containerRef.value) {
    graph.resize(containerRef.value.offsetWidth, containerRef.value.offsetHeight)
  }
}

defineExpose({
  zoomIn,
  zoomOut,
  zoomTo,
  fitContent
})

watch(() => props.graphData, () => {
  nextTick(() => {
    renderGraph()
  })
}, { deep: true })

onMounted(() => {
  initGraph()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (graph) {
    graph.dispose()
    graph = null
  }
})
</script>

<style scoped>
.graph-container {
  
  width: 100%;
  height: 100%;
  position: relative;
}
</style>
