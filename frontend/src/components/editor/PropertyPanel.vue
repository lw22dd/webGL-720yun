<template>
  <div class="property-panel">
    <template v-if="editorStore.selectedType === 'node'">
      <div class="panel-header">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="8" x2="12" y2="16"></line>
          <line x1="8" y1="12" x2="16" y2="12"></line>
        </svg>
        <span>节点属性</span>
      </div>
      <div class="panel-body" v-if="nodeData">
        <div class="form-group">
          <label>名称</label>
          <input type="text" :value="nodeData.title" @input="updateNodeTitle($event)" />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>经度</label>
            <input type="number" :value="nodeData.longitude" step="0.000001" readonly />
          </div>
          <div class="form-group">
            <label>纬度</label>
            <input type="number" :value="nodeData.latitude" step="0.000001" readonly />
          </div>
        </div>
        <div class="form-group">
          <label>场景编码</label>
          <input type="text" :value="nodeData.scene_code" readonly />
        </div>
        <div class="form-group">
          <label>状态</label>
          <div class="status-badge" :class="nodeData.status">
            {{ statusText }}
          </div>
        </div>
        <div class="form-group">
          <label>位置</label>
          <div class="status-badge" :class="nodeData.has_position ? 'has-position' : 'no-position'">
            {{ nodeData.has_position ? '已定位' : '未定位' }}
          </div>
        </div>
      </div>
    </template>

    <template v-else-if="editorStore.selectedType === 'edge'">
      <div class="panel-header">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="5" y1="12" x2="19" y2="12"></line>
          <polyline points="12 5 19 12 12 19"></polyline>
        </svg>
        <span>连线属性</span>
      </div>
      <div class="panel-body" v-if="edgeData">
        <div class="form-group">
          <label>源节点</label>
          <input type="text" :value="sourceNodeTitle" readonly />
        </div>
        <div class="form-group">
          <label>目标节点</label>
          <input type="text" :value="targetNodeTitle" readonly />
        </div>
        <div class="form-group">
          <label>热点标题</label>
          <input
            type="text"
            :value="edgeData.hotspot_title || ''"
            placeholder="如：前往二王庙"
            @input="updateEdgeHotspotTitle($event)"
          />
        </div>
        <div class="form-group">
          <label>类型</label>
          <div class="direction-toggle">
            <button
              class="toggle-btn"
              :class="{ active: edgeData.type === 'walk' }"
              @click="updateEdgeType('walk')"
            >
              步行
            </button>
            <button
              class="toggle-btn"
              :class="{ active: edgeData.type === 'teleport' }"
              @click="updateEdgeType('teleport')"
            >
              传送
            </button>
          </div>
        </div>
        <div class="form-group">
          <button class="delete-btn" @click="handleDeleteEdge">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
            删除连线
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import { useGraphEditorStore } from '@/stores/graphEditorStore'
import { useGraphSceneStore } from '@/stores/graphSceneStore'
import { SceneNodeStatus } from '@/models/graphEditor/node'
import type { SceneNodeData } from '@/models/graphEditor/node'
import type { SceneEdgeData } from '@/models/graphEditor/edge'

const editorStore = useGraphEditorStore()
const sceneStore = useGraphSceneStore()

const graphActions = inject<{
  getGraph: () => any
} | null>('graphActions', null)

const nodeData = computed<SceneNodeData | null>(() => {
  if (editorStore.selectedType !== 'node') return null
  return sceneStore.getNodeById(Number(editorStore.selectedNodeId)) || null
})

const statusText = computed(() => {
  if (!nodeData.value) return ''
  const map: Record<string, string> = {
    [SceneNodeStatus.PENDING]: '待处理',
    [SceneNodeStatus.PLACED]: '已放置',
    [SceneNodeStatus.CONFIGURED]: '已配置',
  }
  return map[nodeData.value.status] || ''
})

const edgeData = computed<SceneEdgeData | null>(() => {
  if (editorStore.selectedType !== 'edge') return null
  return sceneStore.edges.find(e => e.id === Number(editorStore.selectedEdgeId)) || null
})

const sourceNodeTitle = computed(() => {
  if (!edgeData.value) return ''
  const node = sceneStore.getNodeById(edgeData.value.source_id)
  return node?.title || String(edgeData.value.source_id)
})

const targetNodeTitle = computed(() => {
  if (!edgeData.value) return ''
  const node = sceneStore.getNodeById(edgeData.value.target_id)
  return node?.title || String(edgeData.value.target_id)
})

function updateNodeTitle(e: Event) {
  const value = (e.target as HTMLInputElement).value
  if (!nodeData.value) return
  const node = sceneStore.allNodes.find(n => n.id === nodeData.value!.id)
  if (node) {
    node.title = value
    const graph = graphActions?.getGraph()
    if (graph) {
      const x6Node = graph.getCellById(String(nodeData.value.id))
      if (x6Node) {
        x6Node.setData({ ...x6Node.getData(), title: value })
        x6Node.setAttrByPath('label/text', value)
      }
    }
  }
}

function updateEdgeHotspotTitle(e: Event) {
  const value = (e.target as HTMLInputElement).value
  if (!edgeData.value) return
  sceneStore.updateEdge(edgeData.value.id, { hotspot_title: value || '' })
  const graph = graphActions?.getGraph()
  if (graph) {
    const x6Edge = graph.getCellById(String(edgeData.value.id))
    if (x6Edge) {
      x6Edge.setData({ ...x6Edge.getData(), hotspot_title: value || '' })
      if (value) {
        x6Edge.setLabels([{ attrs: { label: { text: value } } }])
      } else {
        x6Edge.setLabels([])
      }
    }
  }
}

function updateEdgeType(type: 'walk' | 'teleport') {
  if (!edgeData.value) return
  sceneStore.updateEdge(edgeData.value.id, { type })
  const graph = graphActions?.getGraph()
  if (graph) {
    const x6Edge = graph.getCellById(String(edgeData.value.id))
    if (x6Edge) {
      x6Edge.setData({ ...x6Edge.getData(), type })
    }
  }
}

function handleDeleteEdge() {
  if (!edgeData.value) return
  const graph = graphActions?.getGraph()
  if (graph) {
    graph.removeCell(String(edgeData.value.id))
  }
  editorStore.clearSelection()
}
</script>

<style scoped>
.property-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.panel-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px;
  font-size: var(--font-size-base);
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}

.panel-header svg {
  color: #6b7280;
}

.panel-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-group label {
  font-size: var(--font-size-xs);
  color: #6b7280;
  font-weight: 500;
}

.form-group input {
  padding: 7px 10px;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  font-size: var(--font-size-sm);
  color: #1f2937;
  background: #fff;
  transition: border-color var(--transition-fast);
}

.form-group input:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.form-group input[readonly] {
  background: var(--bg-panel);
  color: #6b7280;
}

.form-row {
  display: flex;
  gap: 10px;
}

.form-row .form-group {
  flex: 1;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 12px;
  font-size: var(--font-size-xs);
  font-weight: 500;
  width: fit-content;
}

.status-badge.pending {
  background: #fef3c7;
  color: #92400e;
}

.status-badge.placed {
  background: #dbeafe;
  color: #1e40af;
}

.status-badge.configured {
  background: #d1fae5;
  color: #065f46;
}

.status-badge.has-position {
  background: #d1fae5;
  color: #065f46;
}

.status-badge.no-position {
  background: #fee2e2;
  color: #991b1b;
}

.direction-toggle {
  display: flex;
  gap: 6px;
}

.toggle-btn {
  flex: 1;
  padding: 6px 10px;
  border-radius: var(--border-radius-sm);
  font-size: var(--font-size-xs);
  color: #6b7280;
  background: var(--bg-panel);
  border: 1px solid var(--border-color);
  transition: all var(--transition-fast);
}

.toggle-btn:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.toggle-btn.active {
  background: var(--color-primary-light);
  border-color: var(--color-primary);
  color: var(--color-primary);
  font-weight: 500;
}

.delete-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px;
  border-radius: var(--border-radius-sm);
  font-size: var(--font-size-sm);
  color: var(--color-danger);
  background: #fef2f2;
  border: 1px solid #fecaca;
  transition: all var(--transition-fast);
  margin-top: 8px;
}

.delete-btn:hover {
  background: #fee2e2;
  border-color: #f87171;
}
</style>
