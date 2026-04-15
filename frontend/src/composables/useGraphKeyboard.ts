import { onMounted, onUnmounted } from 'vue'
import { useGraphEditorStore } from '@/stores/graphEditorStore'
import { useGraphSceneStore } from '@/stores/graphSceneStore'

export function useGraphKeyboard() {
  const editorStore = useGraphEditorStore()
  const sceneStore = useGraphSceneStore()

  function handleKeydown(e: KeyboardEvent) {
    const tag = (e.target as HTMLElement).tagName
    if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return

    if (e.key === 'Delete' || e.key === 'Backspace') {
      if (editorStore.selectedNodeId) {
        sceneStore.markNodePending(Number(editorStore.selectedNodeId))
        editorStore.clearSelection()
      } else if (editorStore.selectedEdgeId) {
        sceneStore.removeEdge(Number(editorStore.selectedEdgeId))
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
