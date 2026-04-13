import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UploadTask } from '@/models/UploadModel'

export const useUploadStore = defineStore('upload', () => {
  const uploadTasks = ref<Map<string, UploadTask>>(new Map())
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

  return {
    uploadTasks,
    currentUploadId,
    currentUpload,
    uploadingTasks,
    addUploadTask,
    updateUploadTask,
    removeUploadTask,
    getUploadTask,
    pauseUpload,
    resumeUpload,
    clearCompletedTasks
  }
}, {
  persist: {
    key: 'upload-tasks',
    storage: localStorage
  }
})
