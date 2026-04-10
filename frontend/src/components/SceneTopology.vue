<template>
  <div class="scene-topology">
    <div ref="topologyContainer" class="topology-container"></div>
    <div class="topology-tooltip" v-if="hoveredNode" :style="tooltipStyle">
      <div class="tooltip-title">{{ hoveredNode.title }}</div>
      <div class="tooltip-views">{{ hoveredNode.view_count }} 次浏览</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { Graph, type NodeConfig } from '@antv/g6'
import type { MockScene } from '@/utils/mockData'

interface Props {
  scenes: MockScene[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'sceneClick', scene: MockScene): void
}>()

const topologyContainer = ref<HTMLElement | null>(null)
const hoveredNode = ref<MockScene | null>(null)
const tooltipStyle = ref({ left: '0px', top: '0px' })

let graph: Graph | null = null

const initGraph = () => {
  if (!topologyContainer.value || props.scenes.length === 0) return

  const width = topologyContainer.value.offsetWidth
  const height = 300

  const centerLng = props.scenes[0].longitude
  const centerLat = props.scenes[0].latitude

  const nodes: NodeConfig[] = props.scenes.map((scene, index) => {
    const offsetLng = (scene.longitude - centerLng) * 10000
    const offsetLat = (scene.latitude - centerLat) * 10000

    return {
      id: `node-${scene.id}`,
      data: {
        label: scene.title,
        scene
      },
      style: {
        x: width / 2 + offsetLng,
        y: height / 2 - offsetLat,
        size: 60,
        fill: '#ffffff',
        stroke: '#0052d9',
        lineWidth: 2,
        shadowColor: 'rgba(0, 82, 217, 0.3)',
        shadowBlur: 15,
        cursor: 'pointer'
      }
    }
  })

  const edges: any[] = []
  for (let i = 0; i < nodes.length - 1; i++) {
    edges.push({
      id: `edge-${i}`,
      source: nodes[i].id,
      target: nodes[i + 1].id,
      style: {
        stroke: '#4080ff',
        lineWidth: 1,
        opacity: 0.5,
        lineDash: [5, 5],
        endArrow: {
          path: 'M 0 0 L 6 3 L 6 -3 Z',
          fill: '#4080ff'
        }
      }
    })
  }

  graph = new Graph({
    container: topologyContainer.value,
    width,
    height,
    modes: {
      default: ['drag-canvas', 'zoom-canvas']
    },
    data: {
      nodes,
      edges
    },
    node: {
      type: 'circle',
      style: {
        size: 60,
        fill: '#ffffff',
        stroke: '#0052d9',
        lineWidth: 2,
        shadowColor: 'rgba(0, 82, 217, 0.3)',
        shadowBlur: 15,
        cursor: 'pointer',
        labelText: (d: any) => d.data?.label || '',
        labelFill: '#1d2129',
        labelFontSize: 12,
        labelFontWeight: 500,
        labelOffsetY: 35
      },
      state: {
        hover: {
          lineWidth: 3,
          shadowBlur: 25,
          fill: '#e6f0ff'
        }
      }
    },
    edge: {
      type: 'line',
      style: {
        stroke: '#4080ff',
        lineWidth: 1,
        opacity: 0.5,
        lineDash: [5, 5]
      }
    },
    behaviors: ['drag-canvas', 'zoom-canvas']
  })

  graph.on('node:mouseenter', (e: any) => {
    const nodeId = e.itemId
    const nodeData = props.scenes.find(s => `node-${s.id}` === nodeId)
    if (nodeData) {
      hoveredNode.value = nodeData
      graph?.setElementState(nodeId, 'hover')
    }
  })

  graph.on('node:mouseleave', (e: any) => {
    const nodeId = e.itemId
    hoveredNode.value = null
    graph?.setElementState(nodeId, [])
  })

  graph.on('node:click', (e: any) => {
    const nodeId = e.itemId
    const scene = props.scenes.find(s => `node-${s.id}` === nodeId)
    if (scene) {
      emit('sceneClick', scene)
    }
  })

  graph.on('mousemove', (e: any) => {
    tooltipStyle.value = {
      left: `${e.client.x + 10}px`,
      top: `${e.client.y + 10}px`
    }
  })

  graph.render()
}

const destroyGraph = () => {
  if (graph) {
    graph.destroy()
    graph = null
  }
}

watch(() => props.scenes, (newScenes) => {
  if (newScenes && newScenes.length > 0) {
    setTimeout(() => {
      destroyGraph()
      nextTick(() => {
        initGraph()
      })
    }, 350)
  }
}, { deep: true })

onMounted(() => {
  nextTick(() => {
    initGraph()
  })
})

onUnmounted(() => {
  destroyGraph()
})
</script>

<style scoped>
.scene-topology {
  position: relative;
  width: 100%;
}

.topology-container {
  width: 100%;
  height: 300px;
  background: #f5f7fa;
  border-radius: 12px;
  border: 1px solid #e5e6eb;
  overflow: hidden;
}

.topology-tooltip {
  position: fixed;
  background: #ffffff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  padding: 10px 14px;
  pointer-events: none;
  z-index: 1000;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.tooltip-title {
  color: #1d2129;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.tooltip-views {
  color: #86909c;
  font-size: 12px;
}
</style>
