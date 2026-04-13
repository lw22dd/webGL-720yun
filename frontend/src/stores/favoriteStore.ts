import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

interface FavoriteItem {
  id: number
  name: string
  type: 'space' | 'scene'
  cover_url?: string
  province?: string
  city?: string
  description?: string
  scene_count?: number
  scene_code?: string
  addedAt: string
}

export const useFavoriteStore = defineStore('favorite', () => {
  const favorites = ref<FavoriteItem[]>([])

  const loadFromStorage = () => {
    try {
      const data = localStorage.getItem('favorites')
      if (data) {
        favorites.value = JSON.parse(data)
      }
    } catch (e) {
      favorites.value = []
    }
  }

  const saveToStorage = () => {
    localStorage.setItem('favorites', JSON.stringify(favorites.value))
  }

  const addFavorite = (item: Omit<FavoriteItem, 'addedAt'>) => {
    const exists = favorites.value.some(f => f.id === item.id && f.type === item.type)
    if (!exists) {
      favorites.value.push({
        ...item,
        addedAt: new Date().toISOString()
      })
      saveToStorage()
    }
  }

  const removeFavorite = (id: number, type: 'space' | 'scene') => {
    const index = favorites.value.findIndex(f => f.id === id && f.type === type)
    if (index > -1) {
      favorites.value.splice(index, 1)
      saveToStorage()
    }
  }

  const isFavorite = (id: number, type: 'space' | 'scene') => {
    return favorites.value.some(f => f.id === id && f.type === type)
  }

  const toggleFavorite = (item: Omit<FavoriteItem, 'addedAt'>) => {
    if (isFavorite(item.id, item.type)) {
      removeFavorite(item.id, item.type)
      return false
    } else {
      addFavorite(item)
      return true
    }
  }

  const clearAll = () => {
    favorites.value = []
    saveToStorage()
  }

  const count = computed(() => favorites.value.length)

  loadFromStorage()

  return {
    favorites,
    count,
    addFavorite,
    removeFavorite,
    isFavorite,
    toggleFavorite,
    clearAll
  }
})
