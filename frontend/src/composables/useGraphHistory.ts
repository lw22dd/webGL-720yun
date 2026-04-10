import { ref, computed } from 'vue'
import type { GraphDataResponse } from '@/models/SceneModel'

const MAX_HISTORY_SIZE = 50

export function useGraphHistory() {
  const history = ref<GraphDataResponse[]>([])
  const currentIndex = ref(-1)

  const canUndo = computed(() => currentIndex.value > 0)
  const canRedo = computed(() => currentIndex.value < history.value.length - 1)

  const pushState = (state: GraphDataResponse) => {
    if (currentIndex.value < history.value.length - 1) {
      history.value = history.value.slice(0, currentIndex.value + 1)
    }

    const stateCopy = JSON.parse(JSON.stringify(state)) as GraphDataResponse

    if (history.value.length >= MAX_HISTORY_SIZE) {
      history.value.shift()
    } else {
      currentIndex.value++
    }

    history.value.push(stateCopy)
  }

  const undo = (): GraphDataResponse | null => {
    if (!canUndo.value) return null
    currentIndex.value--
    return history.value[currentIndex.value] ? JSON.parse(JSON.stringify(history.value[currentIndex.value])) : null
  }

  const redo = (): GraphDataResponse | null => {
    if (!canRedo.value) return null
    currentIndex.value++
    return history.value[currentIndex.value] ? JSON.parse(JSON.stringify(history.value[currentIndex.value])) : null
  }

  const clear = () => {
    history.value = []
    currentIndex.value = -1
  }

  const getCurrentState = (): GraphDataResponse | null => {
    if (currentIndex.value >= 0 && currentIndex.value < history.value.length) {
      return JSON.parse(JSON.stringify(history.value[currentIndex.value]))
    }
    return null
  }

  return {
    canUndo,
    canRedo,
    pushState,
    undo,
    redo,
    clear,
    getCurrentState
  }
}
