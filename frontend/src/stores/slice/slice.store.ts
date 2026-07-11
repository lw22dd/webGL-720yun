import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import wsClient from '@/services/websocket.service'
import { useUserStore } from '@/stores/user.store'

export interface SliceTask {
  taskId: string
  sceneId: number
  sceneCode: string
  status: 'pending' | 'slicing' | 'ready' | 'failed'
  progress: number
  stage: string
  message: string
  tileUrl?: string
  previewUrl?: string
  error?: string
  createdAt: number
}

let handlersRegistered = false

export const useSliceStore = defineStore('slice', () => {
  const sliceTasks = ref<Map<string, SliceTask>>(new Map())
  const wsConnected = ref(false)

  const currentTasks = computed(() => {
    return Array.from(sliceTasks.value.values()).filter(
      task => task.status === 'pending' || task.status === 'slicing'
    )
  })

  const completedTasks = computed(() => {
    return Array.from(sliceTasks.value.values()).filter(
      task => task.status === 'ready' || task.status === 'failed'
    )
  })

  function initWebSocket() {
    const userStore = useUserStore()
    const token = userStore.accessToken

    if (!token) {
      console.error('No access token available')
      return
    }

    if (!handlersRegistered) {
      wsClient.on('slice_progress', handleSliceProgress)
      wsClient.on('slice_complete', handleSliceComplete)
      wsClient.on('slice_error', handleSliceError)
      handlersRegistered = true
    }

    wsClient.connect(token)
      .then(() => {
        wsConnected.value = true
        console.log('Slice WebSocket initialized')
      })
      .catch((err: Error) => {
        console.error('Failed to connect WebSocket:', err)
        wsConnected.value = false
      })
  }

  function handleSliceProgress(data: {
    task_id: string
    scene_id: number
    status: string
    progress: number
    stage: string
    message: string
  }) {
    const task = sliceTasks.value.get(data.task_id)
    if (task) {
      task.status = 'slicing'
      task.progress = data.progress
      task.stage = data.stage
      task.message = data.message
    } else {
      sliceTasks.value.set(data.task_id, {
        taskId: data.task_id,
        sceneId: data.scene_id,
        sceneCode: '',
        status: 'slicing',
        progress: data.progress,
        stage: data.stage,
        message: data.message,
        createdAt: Date.now()
      })
    }
  }

  function handleSliceComplete(data: {
    task_id: string
    scene_id: number
    status: string
    tile_url: string
    preview_url: string
  }) {
    const task = sliceTasks.value.get(data.task_id)
    if (task) {
      task.status = 'ready'
      task.progress = 100
      task.tileUrl = data.tile_url
      task.previewUrl = data.preview_url
      task.message = '切片完成'
    }
  }

  function handleSliceError(data: {
    task_id: string
    scene_id: number
    error: string
  }) {
    const task = sliceTasks.value.get(data.task_id)
    if (task) {
      task.status = 'failed'
      task.error = data.error
      task.message = `切片失败: ${data.error}`
    }
  }

  function addSliceTask(task: SliceTask) {
    sliceTasks.value.set(task.taskId, task)
  }

  function updateSliceTask(taskId: string, updates: Partial<SliceTask>) {
    const task = sliceTasks.value.get(taskId)
    if (task) {
      Object.assign(task, updates)
    }
  }

  function removeSliceTask(taskId: string) {
    sliceTasks.value.delete(taskId)
  }

  function getSliceTask(taskId: string): SliceTask | undefined {
    return sliceTasks.value.get(taskId)
  }

  function getSliceTaskBySceneId(sceneId: number): SliceTask | undefined {
    return Array.from(sliceTasks.value.values()).find(task => task.sceneId === sceneId)
  }

  function clearCompletedSliceTasks() {
    sliceTasks.value.forEach((task, taskId) => {
      if (task.status === 'ready' || task.status === 'failed') {
        sliceTasks.value.delete(taskId)
      }
    })
  }

  function cleanupWebSocket() {
    if (!handlersRegistered) return
    wsClient.off('slice_progress', handleSliceProgress)
    wsClient.off('slice_complete', handleSliceComplete)
    wsClient.off('slice_error', handleSliceError)
    handlersRegistered = false
  }

  function disconnectWebSocket() {
    wsClient.disconnect()
    wsConnected.value = false
  }

  return {
    sliceTasks,
    wsConnected,
    currentTasks,
    completedTasks,
    initWebSocket,
    cleanupWebSocket,
    addSliceTask,
    updateSliceTask,
    removeSliceTask,
    getSliceTask,
    getSliceTaskBySceneId,
    clearCompletedSliceTasks,
    disconnectWebSocket
  }
})
