import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useGraphEditorStore = defineStore('graphEditor', () => {
  const selectedNodeId = ref<string | null>(null)
  const selectedEdgeId = ref<string | null>(null)
  const zoomLevel = ref(1)
  const sidebarCollapsed = ref(false)

  const selectedType = computed<'node' | 'edge' | null>(() => {
    if (selectedNodeId.value) return 'node'
    if (selectedEdgeId.value) return 'edge'
    return null
  })

  const hasSelection = computed(() => selectedNodeId.value !== null || selectedEdgeId.value !== null)

  function selectNode(id: string | null) {
    selectedNodeId.value = id
    selectedEdgeId.value = null
  }

  function selectEdge(id: string | null) {
    selectedEdgeId.value = id
    selectedNodeId.value = null
  }

  function clearSelection() {
    selectedNodeId.value = null
    selectedEdgeId.value = null
  }

  function setZoom(level: number) {
    zoomLevel.value = Math.round(level * 100) / 100
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  return {
    selectedNodeId,
    selectedEdgeId,
    selectedType,
    hasSelection,
    zoomLevel,
    sidebarCollapsed,
    selectNode,
    selectEdge,
    clearSelection,
    setZoom,
    toggleSidebar,
  }
})
