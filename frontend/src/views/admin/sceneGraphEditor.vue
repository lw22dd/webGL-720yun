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
        <GraphToolbar :can-undo="canUndo" :can-redo="canRedo" :zoom="zoom" @undo="handleUndo" @redo="handleRedo"
          @zoom-in="handleZoomIn" @zoom-out="handleZoomOut" @zoom-reset="handleZoomReset"
          @fit-content="handleFitContent" @save="handleSave" />
      </div>
      <div class="header-right">
        <t-tag theme="primary" variant="light">
          已放置: {{ placedCount }} / {{ totalCount }}
        </t-tag>
      </div>
    </div>

    <div class="editor-main">
      <div class="graph-area" ref="graphAreaRef">
        <GraphContainer ref="graphContainerRef" :graph-data="graphData" @node-click="handleNodeClick"
          @node-move="handleNodeMove" @edge-click="handleEdgeClick" @blank-click="handleBlankClick"
          @zoom-change="handleZoomChange" />
      </div>

      <div class="unplaced-area" v-if="unplacedNodes.length > 0">
        <div class="unplaced-header">
          <t-icon name="map-location" />
          <span>待处理节点 ({{ unplacedNodes.length }})</span>
        </div>
        <div class="unplaced-list">
          <div v-for="node in unplacedNodes" :key="node.id" class="unplaced-item" draggable="true"
            @dragstart="handleDragStart($event, node)">
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

    <PropertyPanel v-if="selectedNode || selectedEdge" :node="selectedNode" :edge="selectedEdge"
      @update-node="handleUpdateNode" @update-edge="handleUpdateEdge" @close="handleClosePanel" />
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
import SceneApi from '@/apis/sceneApi'

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

const loadGraphData = async () => {
  loading.value = true
  try {
    const res = await SceneApi.getSpaceGraph(spaceId.value)
    if (res.code === 200 && res.data) {
      console.log(res.data)
      graphData.value = res.data
      clearHistory()
      pushState(res.data)
    } else {
      MessagePlugin.error(res.msg || '获取图数据失败')
    }
  } catch (error) {
    console.error('加载图数据失败:', error)
    MessagePlugin.error('加载图数据失败')
  } finally {
    loading.value = false
  }
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

  loading.value = true
  try {
    const res = await SceneApi.batchUpdateScenePosition({ positions })
    if (res.code === 200) {
      MessagePlugin.success('保存成功')
    } else {
      MessagePlugin.error(res.msg || '保存失败')
    }
  } catch (error) {
    console.error('保存坐标失败:', error)
    MessagePlugin.error('保存坐标失败')
  } finally {
    loading.value = false
  }
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
  overflow: hidden;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: var(--td-bg-color-container);
  border-bottom: 1px solid var(--td-component-border);
  flex-shrink: 0;
  min-width: 0;
  max-width: 100%;
}

.editor-header .header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.space-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--td-text-color-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 300px;
}

.editor-header .header-center {
  flex: 0 0 auto;
  display: flex;
  justify-content: center;
  min-width: 0;
  max-width: 600px;
  overflow: hidden;
}

.editor-header .header-right {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
  min-width: 0;
}

.editor-main {
  flex: 1;
  display: flex;
  overflow: hidden;
  min-width: 0;
  max-width: 100%;
}

.graph-area {
  flex: 1;
  position: relative;
  overflow: hidden;
  min-width: 0;
  max-width: calc(100% - 400px);
}

.unplaced-area {
  width: 400px;
  min-width: 300px;
  max-width: 400px;
  background-color: var(--td-bg-color-container);
  border-left: 1px solid var(--td-component-border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  margin-right: 0;
  box-sizing: border-box;
}

.unplaced-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--td-component-border);
  font-weight: 500;
  color: var(--td-text-color-primary);
  flex-shrink: 0;
}

.unplaced-list {
  flex: 1;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
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
  min-width: 0;
  max-width: 100%;
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
  overflow: hidden;
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

/* 响应式适配 */
@media screen and (max-width: 1400px) {
  .space-name {
    max-width: 200px;
  }

  .graph-area {
    max-width: calc(100% - 360px);
  }

  .unplaced-area {
    width: 360px;
    max-width: 360px;
  }
}

@media screen and (max-width: 1200px) {
  .editor-header .header-center {
    max-width: 400px;
  }

  .graph-area {
    max-width: calc(100% - 320px);
  }

  .unplaced-area {
    width: 320px;
    min-width: 300px;
    max-width: 320px;
  }
}

@media screen and (max-width: 992px) {
  .editor-main {
    flex-direction: column;
  }

  .graph-area {
    max-width: 100%;
    height: 60%;
  }

  .unplaced-area {
    width: 100%;
    max-width: 100%;
    height: 40%;
    border-left: none;
    border-top: 1px solid var(--td-component-border);
    margin-right: 0;
  }

  .unplaced-list {
    flex-direction: row;
    flex-wrap: wrap;
  }

  .unplaced-item {
    width: calc(50% - 4px);
  }
}

@media screen and (max-width: 768px) {
  .editor-header {
    padding: 8px 12px;
    flex-wrap: wrap;
    gap: 8px;
    height: auto;
  }

  .editor-header .header-left,
  .editor-header .header-center,
  .editor-header .header-right {
    flex: none;
    width: 100%;
    justify-content: center;
  }

  .space-name {
    font-size: 14px;
    max-width: none;
  }

  .unplaced-item {
    width: 100%;
  }
}
</style>
