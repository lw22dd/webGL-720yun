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
        fill: '#0a0e27',
        stroke: '#00f0ff',
        lineWidth: 2,
        shadowColor: '#00f0ff',
        shadowBlur: 20,
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
        stroke: '#00f0ff',
        lineWidth: 1,
        opacity: 0.5,
        lineDash: [5, 5],
        endArrow: {
          path: 'M 0 0 L 6 3 L 6 -3 Z',
          fill: '#00f0ff'
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
        fill: '#0a0e27',
        stroke: '#00f0ff',
        lineWidth: 2,
        shadowColor: '#00f0ff',
        shadowBlur: 15,
        cursor: 'pointer',
        labelText: (d: any) => d.data?.label || '',
        labelFill: '#fff',
        labelFontSize: 12,
        labelFontWeight: 500,
        labelOffsetY: 35
      },
      state: {
        hover: {
          lineWidth: 3,
          shadowBlur: 30,
          fill: 'rgba(0, 240, 255, 0.1)'
        }
      }
    },
    edge: {
      type: 'line',
      style: {
        stroke: '#00f0ff',
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

watch(() => props.scenes, () => {
  destroyGraph()
  nextTick(() => {
    initGraph()
  })
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
  background: rgba(0, 0, 0, 0.3);
  border-radius: 12px;
  border: 1px solid rgba(0, 240, 255, 0.15);
  overflow: hidden;
}

.topology-tooltip {
  position: fixed;
  background: rgba(10, 14, 39, 0.95);
  border: 1px solid rgba(0, 240, 255, 0.3);
  border-radius: 8px;
  padding: 10px 14px;
  pointer-events: none;
  z-index: 1000;
  backdrop-filter: blur(10px);
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
