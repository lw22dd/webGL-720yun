import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UploadTask, ProgressData, CompleteData, ErrorData } from '@/models/UploadModel'
import wsClient from './websocketClient'
import { useUserStore } from './userStore'

export const useSceneStore = defineStore('scene', () => {
  const uploadTasks = ref<Map<string, UploadTask>>(new Map())
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

    if (!token) {
      console.error('No access token available')
      return
    }

    wsClient.on('progress', handleProgress)
    wsClient.on('complete', handleComplete)
    wsClient.on('error', handleError)
    wsClient.on('merge_start', handleMergeStart)
    wsClient.on('merge_progress', handleMergeProgress)
    wsClient.on('disconnected', handleDisconnected)
    wsClient.on('error', handleWsError)

    wsClient.connect(token)
      .then(() => {
        wsConnected.value = true
        console.log('WebSocket initialized')
      })
      .catch(err => {
        console.error('Failed to connect WebSocket:', err)
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

  function handleDisconnected(data: { code: number; reason: string }) {
    wsConnected.value = false
    console.log('WebSocket disconnected:', data)
  }

  function handleWsError(error: any) {
    console.error('WebSocket error:', error)
    wsConnected.value = false
  }

  function addUploadTask(task: UploadTask) {
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
    wsConnected,
    currentUploadId,
    currentUpload,
    uploadingTasks,
    initWebSocket,
    addUploadTask,
    updateUploadTask,
    removeUploadTask,
    getUploadTask,
    pauseUpload,
    resumeUpload,
    clearCompletedTasks,
    disconnectWebSocket
  }
}, {
  persist: {
    key: 'scene-upload-tasks',
    storage: localStorage
  }
})
