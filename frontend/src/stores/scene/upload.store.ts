import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UploadTask } from '@/models/upload.model'

function createUploadTasksMap(): Map<string, UploadTask> {
  return new Map()
}

function objectToMap(obj: Record<string, UploadTask> | null | undefined): Map<string, UploadTask> {
  if (!obj) return createUploadTasksMap()
  const map = new Map<string, UploadTask>()
  for (const [key, value] of Object.entries(obj)) {
    map.set(key, value)
  }
  return map
}

export const useUploadStore = defineStore('upload', () => {
  const uploadTasks = ref<Map<string, UploadTask>>(createUploadTasksMap())
  const currentUploadId = ref<string | null>(null)

  const currentUpload = computed(() => {
    if (currentUploadId.value) {
      return uploadTasks.value.get(currentUploadId.value) || null
    }
    return null
  })

  const uploadingTasks = computed(() => {
    return Array.from(uploadTasks.value.values()).filter(
      task => task.status === 'uploading' || task.status === 'pending' || task.status === 'merging' || task.status === 'paused'
    )
  })

  const pausedTasks = computed(() => {
    return Array.from(uploadTasks.value.values()).filter(
      task => task.status === 'paused'
    )
  })

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

  function cancelUpload(uploadId: string) {
    const task = uploadTasks.value.get(uploadId)
    if (task && (task.status === 'uploading' || task.status === 'paused' || task.status === 'pending')) {
      task.status = 'cancelled'
    }
  }

  function clearCompletedTasks() {
    uploadTasks.value.forEach((task, uploadId) => {
      if (task.status === 'completed' || task.status === 'failed' || task.status === 'cancelled') {
        uploadTasks.value.delete(uploadId)
      }
    })
  }

  return {
    uploadTasks,
    currentUploadId,
    currentUpload,
    uploadingTasks,
    pausedTasks,
    addUploadTask,
    updateUploadTask,
    removeUploadTask,
    getUploadTask,
    pauseUpload,
    resumeUpload,
    cancelUpload,
    clearCompletedTasks
  }
}, {
  persist: {
    key: 'upload-tasks',
    storage: localStorage,
    serializer: {
      deserialize: (value) => {
        const parsed = JSON.parse(value)
        return {
          ...parsed,
          uploadTasks: objectToMap(parsed.uploadTasks)
        }
      },
      serialize: (value) => JSON.stringify(value)
    }
  }
})
