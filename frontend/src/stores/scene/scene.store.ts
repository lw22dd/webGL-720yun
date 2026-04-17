import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UploadTask, ProgressData, CompleteData, ErrorData, SliceProgressData, SliceCompleteData, SliceErrorData } from '@/models/upload.model'
import wsClient from '@/services/websocket.service'
import { useUserStore } from '@/stores/user.store'

export interface SliceTaskInfo {
  taskId: string
  sceneId: number
  sceneCode: string
  status: 'pending' | 'slicing' | 'completed' | 'failed'
  progress: number
  stage: string
  message: string
  tileUrl?: string
  previewUrl?: string
}

function createUploadTasksMap(): Map<string, UploadTask> {
  return new Map()
}

function createSliceTasksMap(): Map<string, SliceTaskInfo> {
  return new Map()
}

function objectToUploadMap(obj: Record<string, UploadTask> | null | undefined): Map<string, UploadTask> {
  if (!obj) return createUploadTasksMap()
  const map = new Map<string, UploadTask>()
  for (const [key, value] of Object.entries(obj)) {
    map.set(key, value)
  }
  return map
}

function objectToSliceMap(obj: Record<string, SliceTaskInfo> | null | undefined): Map<string, SliceTaskInfo> {
  if (!obj) return createSliceTasksMap()
  const map = new Map<string, SliceTaskInfo>()
  for (const [key, value] of Object.entries(obj)) {
    map.set(key, value)
  }
  return map
}

export const useSceneStore = defineStore('scene', () => {
  const uploadTasks = ref<Map<string, UploadTask>>(createUploadTasksMap())
  const sliceTasks = ref<Map<string, SliceTaskInfo>>(createSliceTasksMap())
  const wsConnected = ref(false)
  const currentUploadId = ref<string | null>(null)

  const currentUpload = computed(() => {
    if (currentUploadId.value) {
      return uploadTasks.value.get(currentUploadId.value) || null
    }
    return null
  })

  const uploadingTasks = computed(() => {
    return Array.from(uploadTasks.value.values()).filter(
      task => task.status === 'uploading' || task.status === 'pending' || task.status === 'merging'
    )
  })

  function initWebSocket() {
    const userStore = useUserStore()
    const token = userStore.accessToken

    console.log('[SceneStore] initWebSocket called, token exists:', !!token)

    if (!token) {
      console.error('[SceneStore] No access token available')
      return
    }

    wsClient.on('progress', handleProgress)
    wsClient.on('complete', handleComplete)
    wsClient.on('error', handleError)
    wsClient.on('merge_start', handleMergeStart)
    wsClient.on('merge_progress', handleMergeProgress)
    wsClient.on('slice_progress', handleSliceProgress)
    wsClient.on('slice_complete', handleSliceComplete)
    wsClient.on('slice_error', handleSliceError)
    wsClient.on('disconnected', handleDisconnected)
    wsClient.on('error', handleWsError)

    console.log('[SceneStore] Connecting WebSocket...')
    wsClient.connect(token)
      .then(() => {
        wsConnected.value = true
        console.log('[SceneStore] WebSocket connected successfully')
      })
      .catch(err => {
        console.error('[SceneStore] Failed to connect WebSocket:', err)
        wsConnected.value = false
      })
  }

  function handleProgress(data: ProgressData & { upload_id: string }) {
    const task = uploadTasks.value.get(data.upload_id)
    if (task) {
      task.uploadedChunks = Array.from({ length: data.uploaded_chunks }, (_, i) => i)
      task.percentage = data.percentage
      task.speed = data.speed || ''
      task.uploadedBytes = data.uploaded_bytes
      task.status = 'uploading'
    }
  }

  function handleComplete(data: CompleteData & { upload_id: string }) {
    const task = uploadTasks.value.get(data.upload_id)
    if (task) {
      task.status = 'completed'
      task.percentage = 100
    }
  }

  function handleError(data: ErrorData & { upload_id: string }) {
    const task = uploadTasks.value.get(data.upload_id)
    if (task) {
      task.status = 'failed'
    }
    console.error('Upload error:', data.message)
  }

  function handleMergeStart(data: { upload_id: string }) {
    const task = uploadTasks.value.get(data.upload_id)
    if (task) {
      task.status = 'merging'
    }
  }

  function handleMergeProgress(data: { upload_id: string; percentage: number; message: string }) {
    const task = uploadTasks.value.get(data.upload_id)
    if (task) {
      task.percentage = data.percentage
    }
  }

  function handleSliceProgress(data: SliceProgressData) {
    console.log('[Store] Slice progress:', data)
    const task = sliceTasks.value.get(data.task_id)
    if (task) {
      const updatedTask = {
        ...task,
        progress: data.progress,
        stage: data.stage,
        message: data.message,
        status: 'slicing' as const
      }
      sliceTasks.value = new Map(sliceTasks.value.set(data.task_id, updatedTask))
    } else {
      console.log('[Store] Task not found for progress:', data.task_id)
    }
  }

  function handleSliceComplete(data: SliceCompleteData) {
    console.log('[Store] Slice complete:', data)
    const task = sliceTasks.value.get(data.task_id)
    if (task) {
      const updatedTask = {
        ...task,
        status: 'completed' as const,
        progress: 100,
        tileUrl: data.tile_url,
        previewUrl: data.preview_url,
        message: '切片完成'
      }
      sliceTasks.value = new Map(sliceTasks.value.set(data.task_id, updatedTask))
    } else {
      console.log('[Store] Task not found for complete:', data.task_id)
    }
  }

  function handleSliceError(data: SliceErrorData) {
    console.log('[Store] Slice error:', data)
    const task = sliceTasks.value.get(data.task_id)
    if (task) {
      const updatedTask = {
        ...task,
        status: 'failed' as const,
        message: data.error
      }
      sliceTasks.value = new Map(sliceTasks.value.set(data.task_id, updatedTask))
    } else {
      console.log('[Store] Task not found for error:', data.task_id)
    }
  }

  function handleDisconnected(data: { code: number; reason: string }) {
    wsConnected.value = false
    console.log('WebSocket disconnected:', data)
  }

  function handleWsError(error: any) {
    console.error('WebSocket error:', error)
    wsConnected.value = false
  }

  function addUploadTask(task: UploadTask) {
    if (!(uploadTasks.value instanceof Map)) {
      uploadTasks.value = createUploadTasksMap()
    }
    uploadTasks.value.set(task.uploadId, task)
    currentUploadId.value = task.uploadId
  }

  function updateUploadTask(uploadId: string, updates: Partial<UploadTask>) {
    const task = uploadTasks.value.get(uploadId)
    if (task) {
      Object.assign(task, updates)
    }
  }

  function removeUploadTask(uploadId: string) {
    uploadTasks.value.delete(uploadId)
    if (currentUploadId.value === uploadId) {
      currentUploadId.value = null
    }
  }

  function getUploadTask(uploadId: string): UploadTask | undefined {
    return uploadTasks.value.get(uploadId)
  }

  function addSliceTask(taskInfo: SliceTaskInfo) {
    sliceTasks.value = new Map(sliceTasks.value.set(taskInfo.taskId, taskInfo))
  }

  function getSliceTask(taskId: string): SliceTaskInfo | undefined {
    return sliceTasks.value.get(taskId)
  }

  function updateSliceTask(taskId: string, updates: Partial<SliceTaskInfo>) {
    const task = sliceTasks.value.get(taskId)
    if (task) {
      const updatedTask = { ...task, ...updates }
      sliceTasks.value = new Map(sliceTasks.value.set(taskId, updatedTask))
    }
  }

  function removeSliceTask(taskId: string) {
    sliceTasks.value.delete(taskId)
  }

  function pauseUpload(uploadId: string) {
    const task = uploadTasks.value.get(uploadId)
    if (task && task.status === 'uploading') {
      task.status = 'paused'
    }
  }

  function resumeUpload(uploadId: string) {
    const task = uploadTasks.value.get(uploadId)
    if (task && task.status === 'paused') {
      task.status = 'uploading'
    }
  }

  function clearCompletedTasks() {
    uploadTasks.value.forEach((task, uploadId) => {
      if (task.status === 'completed' || task.status === 'failed') {
        uploadTasks.value.delete(uploadId)
      }
    })
  }

  function disconnectWebSocket() {
    wsClient.disconnect()
    wsConnected.value = false
  }

  return {
    uploadTasks,
    sliceTasks,
    wsConnected,
    currentUploadId,
    currentUpload,
    uploadingTasks,
    initWebSocket,
    addUploadTask,
    updateUploadTask,
    removeUploadTask,
    getUploadTask,
    addSliceTask,
    getSliceTask,
    updateSliceTask,
    removeSliceTask,
    pauseUpload,
    resumeUpload,
    clearCompletedTasks,
    disconnectWebSocket
  }
}, {
  persist: {
    key: 'scene-upload-tasks',
    storage: localStorage,
    serializer: {
      deserialize: (value) => {
        const parsed = JSON.parse(value)
        return {
          ...parsed,
          uploadTasks: objectToUploadMap(parsed.uploadTasks),
          sliceTasks: objectToSliceMap(parsed.sliceTasks)
        }
      },
      serialize: (value) => JSON.stringify(value)
    }
  }
})
