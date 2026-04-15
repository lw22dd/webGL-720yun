import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { SceneNodeData } from '@/models/graphEditor/node'
import { SceneNodeStatus, type SceneNodeStatusType } from '@/models/graphEditor/node'
import type { SceneEdgeData } from '@/models/graphEditor/edge'
import type { SpaceInfoForGraph } from '@/models/graphEditor/scene'
import SpaceApi from '@/apis/spaceApi'

export const useGraphSceneStore = defineStore('graphScene', () => {
  const spaceInfo = ref<SpaceInfoForGraph | null>(null)
  const allNodes = ref<SceneNodeData[]>([])
  const edges = ref<SceneEdgeData[]>([])
  const loading = ref(false)
  const saving = ref(false)

  const pendingNodes = computed(() =>
    allNodes.value.filter(n => n.status === SceneNodeStatus.PENDING)
  )
  const placedNodes = computed(() =>
    allNodes.value.filter(n => n.status !== SceneNodeStatus.PENDING)
  )
  const placedCount = computed(() => placedNodes.value.length)
  const totalCount = computed(() => allNodes.value.length)

  async function loadScene(spaceId: number) {
    loading.value = true
    try {
      const response = await SpaceApi.getSpaceGraph(spaceId)
      if (response.data) {
        spaceInfo.value = response.data.space_info ?? null
        edges.value = response.data.edges ?? []

        const nodesRaw = response.data.nodes ?? []

        // 根据 hasPosition 字段分离已放置和待处理节点
        const placedNodesRaw = nodesRaw.filter((n: SceneNodeData) => n.has_position)
        const unplacedNodesRaw = nodesRaw.filter((n: SceneNodeData) => !n.has_position)

        const placedNodesWithStatus = placedNodesRaw.map((n: SceneNodeData) => ({
          ...n,
          status: SceneNodeStatus.PLACED as SceneNodeStatusType,
        }))

        const unplacedNodesWithStatus = unplacedNodesRaw.map((n: SceneNodeData) => ({
          ...n,
          status: SceneNodeStatus.PENDING as SceneNodeStatusType,
        }))

        allNodes.value = [...placedNodesWithStatus, ...unplacedNodesWithStatus]

        console.log('[GraphSceneStore] 从后端加载的数据:', {
          spaceInfo: response.data.space_info,
          placedNodes: placedNodesRaw.length,
          unplacedNodes: unplacedNodesRaw.length,
          totalNodes: allNodes.value.length,
          totalEdges: edges.value.length,
        })
      }
    } finally {
      loading.value = false
    }
  }

  function markNodePlaced(nodeId: number) {
    const node = allNodes.value.find(n => n.id === nodeId)
    if (node && node.status === SceneNodeStatus.PENDING) {
      node.status = SceneNodeStatus.PLACED
    }
  }

  function markNodePending(nodeId: number) {
    const node = allNodes.value.find(n => n.id === nodeId)
    if (node) {
      node.status = SceneNodeStatus.PENDING
    }
  }

  function addEdge(edge: Omit<SceneEdgeData, 'id'>) {
    const newId = edges.value.length > 0 
      ? Math.max(...edges.value.map(e => e.id)) + 1 
      : 1
    edges.value.push({ ...edge, id: newId })
  }

  function updateEdge(edgeId: number, data: Partial<SceneEdgeData>) {
    const edge = edges.value.find(e => e.id === edgeId)
    if (edge) {
      Object.assign(edge, data)
    }
  }

  function removeEdge(edgeId: number) {
    const index = edges.value.findIndex(e => e.id === edgeId)
    if (index !== -1) {
      edges.value.splice(index, 1)
    }
  }

  function removeNode(nodeId: number) {
    const index = allNodes.value.findIndex(n => n.id === nodeId)
    if (index !== -1) {
      allNodes.value.splice(index, 1)
    }
    edges.value = edges.value.filter(e => e.source_id !== nodeId && e.target_id !== nodeId)
  }

  async function saveScene() {
    saving.value = true
    try {
      return true
    } finally {
      saving.value = false
    }
  }

  function getNodeById(nodeId: number): SceneNodeData | undefined {
    return allNodes.value.find(n => n.id === nodeId)
  }

  return {
    spaceInfo,
    allNodes,
    edges,
    loading,
    saving,
    pendingNodes,
    placedNodes,
    placedCount,
    totalCount,
    loadScene,
    markNodePlaced,
    markNodePending,
    addEdge,
    updateEdge,
    removeEdge,
    removeNode,
    saveScene,
    getNodeById,
  }
})
