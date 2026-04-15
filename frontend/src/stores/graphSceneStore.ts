import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { SceneNodeData } from '@/models/graphEditor/node'
import { SceneNodeStatus } from '@/models/graphEditor/node'
import type { SceneEdgeData } from '@/models/graphEditor/edge'
import type { BackgroundMap } from '@/models/graphEditor/scene'
import { fetchSceneData, saveSceneData } from '@/apis/graphEditorMock'
import { generateId } from '@/utils/graphEditor/helpers'

export const useGraphSceneStore = defineStore('graphScene', () => {
  const sceneName = ref('都江堰景区')
  const allNodes = ref<SceneNodeData[]>([])
  const edges = ref<SceneEdgeData[]>([])
  const backgroundMap = ref<BackgroundMap | null>(null)
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

  async function loadScene(sceneId: string) {
    loading.value = true
    try {
      const data = await fetchSceneData(sceneId)
      sceneName.value = data.name
      allNodes.value = data.nodes
      edges.value = data.edges
      backgroundMap.value = data.backgroundMap || null
    } finally {
      loading.value = false
    }
  }

  function markNodePlaced(nodeId: string) {
    const node = allNodes.value.find(n => n.id === nodeId)
    if (node && node.status === SceneNodeStatus.PENDING) {
      node.status = SceneNodeStatus.PLACED
    }
  }

  function markNodePending(nodeId: string) {
    const node = allNodes.value.find(n => n.id === nodeId)
    if (node) {
      node.status = SceneNodeStatus.PENDING
    }
  }

  function addEdge(edge: Omit<SceneEdgeData, 'id'>) {
    edges.value.push({ ...edge, id: generateId('edge') })
  }

  function updateEdge(edgeId: string, data: Partial<SceneEdgeData>) {
    const edge = edges.value.find(e => e.id === edgeId)
    if (edge) {
      Object.assign(edge, data)
    }
  }

  function removeEdge(edgeId: string) {
    const index = edges.value.findIndex(e => e.id === edgeId)
    if (index !== -1) {
      edges.value.splice(index, 1)
    }
  }

  function removeNode(nodeId: string) {
    const index = allNodes.value.findIndex(n => n.id === nodeId)
    if (index !== -1) {
      allNodes.value.splice(index, 1)
    }
    edges.value = edges.value.filter(e => e.sourceId !== nodeId && e.targetId !== nodeId)
  }

  async function saveScene() {
    saving.value = true
    try {
      const data = {
        id: 'scene-001',
        name: sceneName.value,
        nodes: allNodes.value,
        edges: edges.value,
        backgroundMap: backgroundMap.value || undefined,
      }
      return await saveSceneData(data)
    } finally {
      saving.value = false
    }
  }

  function getNodeById(nodeId: string): SceneNodeData | undefined {
    return allNodes.value.find(n => n.id === nodeId)
  }

  return {
    sceneName,
    allNodes,
    edges,
    backgroundMap,
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
