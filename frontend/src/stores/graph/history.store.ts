import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

interface HistoryItem {
  id: number
  name: string
  type: 'space' | 'scene'
  thumbnail_url?: string
  province?: string
  city?: string
  scene_code?: string
  visitedAt: string
}

export const useHistoryStore = defineStore('history', () => {
  const historyList = ref<HistoryItem[]>([])
  const maxHistoryCount = 100

  const loadFromStorage = () => {
    try {
      const data = localStorage.getItem('browseHistory')
      if (data) {
        historyList.value = JSON.parse(data)
      }
    } catch (e) {
      historyList.value = []
    }
  }

  const saveToStorage = () => {
    localStorage.setItem('browseHistory', JSON.stringify(historyList.value))
  }

  const addHistory = (item: Omit<HistoryItem, 'visitedAt'>) => {
    const newItem: HistoryItem = {
      ...item,
      visitedAt: new Date().toISOString()
    }
    
    const existingIndex = historyList.value.findIndex(
      h => h.id === item.id && h.type === item.type
    )
    
    if (existingIndex > -1) {
      historyList.value.splice(existingIndex, 1)
    }
    
    historyList.value.unshift(newItem)
    
    if (historyList.value.length > maxHistoryCount) {
      historyList.value = historyList.value.slice(0, maxHistoryCount)
    }
    
    saveToStorage()
  }

  const removeHistory = (id: number, type: 'space' | 'scene', visitedAt: string) => {
    const index = historyList.value.findIndex(
      h => h.id === id && h.type === type && h.visitedAt === visitedAt
    )
    if (index > -1) {
      historyList.value.splice(index, 1)
      saveToStorage()
    }
  }

  const clearAll = () => {
    historyList.value = []
    saveToStorage()
  }

  const count = computed(() => historyList.value.length)

  const recentHistory = computed(() => historyList.value.slice(0, 10))

  loadFromStorage()

  return {
    historyList,
    count,
    recentHistory,
    addHistory,
    removeHistory,
    clearAll
  }
})
