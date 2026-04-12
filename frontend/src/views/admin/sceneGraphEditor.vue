<template>
  <div class="scene-graph-editor">
    <div class="editor-header">
      <div class="header-left">
        <t-button variant="text" @click="goBack">
          <template #icon><t-icon name="chevron-left" /></template>
          返回
        </t-button>
        <div class="space-name">{{ graphData?.space_info?.name || '加载中...' }}</div>
      </div>
      <div class="header-center">
        <GraphToolbar
          :can-undo="canUndo"
          :can-redo="canRedo"
          :zoom="zoom"
          @undo="handleUndo"
          @redo="handleRedo"
          @zoom-in="handleZoomIn"
          @zoom-out="handleZoomOut"
          @zoom-reset="handleZoomReset"
          @fit-content="handleFitContent"
          @save="handleSave"
        />
      </div>
      <div class="header-right">
        <t-tag theme="primary" variant="light">
          已放置: {{ placedCount }} / {{ totalCount }}
        </t-tag>
      </div>
    </div>

    <div class="editor-main">
      <div class="graph-area" ref="graphAreaRef">
        <GraphContainer
          ref="graphContainerRef"
          :graph-data="graphData"
          @node-click="handleNodeClick"
          @node-move="handleNodeMove"
          @edge-click="handleEdgeClick"
          @blank-click="handleBlankClick"
          @zoom-change="handleZoomChange"
        />
      </div>

      <div class="unplaced-area" v-if="unplacedNodes.length > 0">
        <div class="unplaced-header">
          <t-icon name="map-location" />
          <span>待处理节点 ({{ unplacedNodes.length }})</span>
        </div>
        <div class="unplaced-list">
          <div
            v-for="node in unplacedNodes"
            :key="node.id"
            class="unplaced-item"
            draggable="true"
            @dragstart="handleDragStart($event, node)"
          >
            <div class="unplaced-thumb">
              <img v-if="node.thumbnail_url" :src="node.thumbnail_url" :alt="node.title" />
              <t-icon v-else name="image" size="24" />
            </div>
            <div class="unplaced-info">
              <div class="unplaced-title">{{ node.title }}</div>
              <div class="unplaced-code">{{ node.scene_code }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <PropertyPanel
      v-if="selectedNode || selectedEdge"
      :node="selectedNode"
      :edge="selectedEdge"
      @update-node="handleUpdateNode"
      @update-edge="handleUpdateEdge"
      @close="handleClosePanel"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import GraphContainer from '@/components/graph/GraphContainer.vue'
import GraphToolbar from '@/components/graph/GraphToolbar.vue'
import PropertyPanel from '@/components/graph/PropertyPanel.vue'
import type { GraphDataResponse, SceneNodeData, EdgeData, ScenePosition } from '@/models/SceneModel'
import { useGraphHistory } from '@/composables/useGraphHistory'
import { mockSpaces } from '@/utils/mockData'

const route = useRoute()
const router = useRouter()

const spaceId = computed(() => Number(route.params.id))

const graphData = ref<GraphDataResponse | null>(null)
const loading = ref(false)
const zoom = ref(1)
const selectedNode = ref<SceneNodeData | null>(null)
const selectedEdge = ref<EdgeData | null>(null)

const graphContainerRef = ref<InstanceType<typeof GraphContainer> | null>(null)
const graphAreaRef = ref<HTMLElement | null>(null)

const { canUndo, canRedo, pushState, undo, redo, clear: clearHistory } = useGraphHistory()

const placedCount = computed(() => graphData.value?.nodes?.length || 0)
const totalCount = computed(() => {
  const data = graphData.value
  if (!data) return 0
  return (data.nodes?.length || 0) + (data.unplaced?.length || 0)
})
const unplacedNodes = computed(() => graphData.value?.unplaced || [])

const loadGraphData = () => {
  const space = mockSpaces.find(s => s.id === spaceId.value)
  if (!space) {
    MessagePlugin.warning('空间不存在')
    return
  }

  const nodes: SceneNodeData[] = []
  const unplaced: SceneNodeData[] = []
  const edges: EdgeData[] = []

  space.scenes.forEach((scene, index) => {
    const nodeData: SceneNodeData = {
      id: scene.id,
      title: scene.title,
      scene_code: scene.scene_code,
      thumbnail_url: scene.thumbnail_url,
      longitude: scene.longitude,
      latitude: scene.latitude,
      has_position: scene.longitude !== 0 && scene.latitude !== 0,
      view_count: scene.view_count
    }

    if (nodeData.has_position) {
      nodes.push(nodeData)
    } else {
      unplaced.push(nodeData)
    }

  })

  // 不创建边，只显示散点

  graphData.value = {
    space_info: {
      id: space.id,
      name: space.name,
      longitude: space.longitude,
      latitude: space.latitude,
      zoom_level: space.zoom_level
    },
    nodes,
    edges,
    unplaced
  }
  clearHistory()
  pushState(graphData.value)
}

const goBack = () => {
  router.push('/admin/spaces')
}

const handleNodeClick = (node: SceneNodeData) => {
  selectedNode.value = node
  selectedEdge.value = null
}

const handleEdgeClick = (edge: EdgeData) => {
  selectedEdge.value = edge
  selectedNode.value = null
}

const handleBlankClick = () => {
  selectedNode.value = null
  selectedEdge.value = null
}

const handleNodeMove = async (nodeId: number, longitude: number, latitude: number) => {
  if (graphData.value) {
    const nodeIndex = graphData.value.nodes.findIndex(n => n.id === nodeId)
    if (nodeIndex !== -1) {
      graphData.value.nodes[nodeIndex].longitude = longitude
      graphData.value.nodes[nodeIndex].latitude = latitude
    }
    const unplacedIndex = graphData.value.unplaced.findIndex(n => n.id === nodeId)
    if (unplacedIndex !== -1) {
      const movedNode = graphData.value.unplaced.splice(unplacedIndex, 1)[0]
      movedNode.longitude = longitude
      movedNode.latitude = latitude
      movedNode.has_position = true
      graphData.value.nodes.push(movedNode)
    }
    pushState(graphData.value)
  }
  MessagePlugin.success('坐标已更新（测试模式）')
}

const handleZoomChange = (newZoom: number) => {
  zoom.value = newZoom
}

const handleUndo = () => {
  const state = undo()
  if (state) {
    graphData.value = state
  }
}

const handleRedo = () => {
  const state = redo()
  if (state) {
    graphData.value = state
  }
}

const handleZoomIn = () => {
  graphContainerRef.value?.zoomIn()
}

const handleZoomOut = () => {
  graphContainerRef.value?.zoomOut()
}

const handleZoomReset = () => {
  graphContainerRef.value?.zoomTo(1)
}

const handleFitContent = () => {
  graphContainerRef.value?.fitContent()
}

const handleSave = async () => {
  if (!graphData.value) return

  const positions: ScenePosition[] = graphData.value.nodes
    .filter(n => n.has_position)
    .map(n => ({
      scene_id: n.id,
      longitude: n.longitude,
      latitude: n.latitude
    }))

  if (positions.length === 0) {
    MessagePlugin.warning('没有需要保存的坐标')
    return
  }

  console.log('保存坐标（测试模式）:', positions)
  MessagePlugin.success('保存成功（测试模式）')
}

const handleDragStart = (event: DragEvent, node: SceneNodeData) => {
  event.dataTransfer?.setData('application/json', JSON.stringify(node))
}

const handleUpdateNode = (updates: Partial<SceneNodeData>) => {
  if (!selectedNode.value || !graphData.value) return
  const index = graphData.value.nodes.findIndex(n => n.id === selectedNode.value!.id)
  if (index !== -1) {
    Object.assign(graphData.value.nodes[index], updates)
    pushState(graphData.value)
  }
}

const handleUpdateEdge = (updates: Partial<EdgeData>) => {
  if (!selectedEdge.value || !graphData.value) return
  const index = graphData.value.edges.findIndex(e => e.id === selectedEdge.value!.id)
  if (index !== -1) {
    Object.assign(graphData.value.edges[index], updates)
    pushState(graphData.value)
  }
}

const handleClosePanel = () => {
  selectedNode.value = null
  selectedEdge.value = null
}

onMounted(() => {
  loadGraphData()
})
</script>

<style scoped>
.scene-graph-editor {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--td-bg-color-page);
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: var(--td-bg-color-container);
  border-bottom: 1px solid var(--td-component-border);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.space-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--td-text-color-primary);
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.editor-main {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.graph-area {
  max-width: 80%;
  flex: 1;
  position: relative;
  overflow: hidden;
}

.unplaced-area {
  width: 420px;
  min-width: 320px;
  background-color: var(--td-bg-color-container);
  border-left: 1px solid var(--td-component-border);
  display: flex;
  flex-direction: column;
  margin-right: 20px;
}

.unplaced-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--td-component-border);
  font-weight: 500;
  color: var(--td-text-color-primary);
}

.unplaced-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.unplaced-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px;
  background-color: var(--td-bg-color-container-hover);
  border: 1px dashed var(--td-component-border);
  border-radius: 6px;
  cursor: grab;
  transition: all 0.2s;
}

.unplaced-item:hover {
  border-color: var(--td-brand-color);
  background-color: var(--td-bg-color-specialcomponent);
}

.unplaced-item:active {
  cursor: grabbing;
}

.unplaced-thumb {
  width: 48px;
  height: 32px;
  border-radius: 4px;
  overflow: hidden;
  background-color: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.unplaced-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.unplaced-thumb :deep(.t-icon) {
  color: #86909c;
}

.unplaced-info {
  flex: 1;
  min-width: 0;
}

.unplaced-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--td-text-color-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.unplaced-code {
  font-size: 11px;
  color: var(--td-text-color-placeholder);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
