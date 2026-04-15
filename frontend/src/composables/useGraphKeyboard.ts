import { onMounted, onUnmounted, inject } from 'vue'
import { useGraphEditorStore } from '@/stores/graphEditorStore'
import { useGraphSceneStore } from '@/stores/graphSceneStore'

export function useGraphKeyboard() {
  const editorStore = useGraphEditorStore()
  const sceneStore = useGraphSceneStore()

  const graphActions = inject<{
    getGraph: () => any
  } | null>('graphActions', null)

  function handleKeydown(e: KeyboardEvent) {
    const tag = (e.target as HTMLElement).tagName
    if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return

    if (e.key === 'Delete' || e.key === 'Backspace') {
      const graph = graphActions?.getGraph()
      if (!graph) return

      if (editorStore.selectedNodeId) {
        graph.removeCell(editorStore.selectedNodeId)
        sceneStore.markNodePending(editorStore.selectedNodeId)
        editorStore.clearSelection()
      } else if (editorStore.selectedEdgeId) {
        graph.removeCell(editorStore.selectedEdgeId)
        sceneStore.removeEdge(editorStore.selectedEdgeId)
        editorStore.clearSelection()
      }
    }
  }

  onMounted(() => {
    window.addEventListener('keydown', handleKeydown)
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeydown)
  })
}
