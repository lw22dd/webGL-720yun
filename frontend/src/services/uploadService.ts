import type { Result } from '@/models/result.model'
import type {
  InitUploadRequest,
  InitUploadResponse,
  ChunkUploadResponse,
  CompleteUploadResponse,
  UploadStatusResponse,
  UploadOptions,
  UploadTask
} from '@/models/upload.model'
import UploadApi from '@/services/api/upload.api'
import { useUploadStore } from '@/stores/scene/upload.store'

const CHUNK_SIZE = 5 * 1024 * 1024
const MAX_CONCURRENT = 4
const MAX_RETRIES = 3
const RETRY_DELAYS = [1000, 2000, 4000]
const MAX_FILE_SIZE = 500 * 1024 * 1024
const PAUSE_CHECK_INTERVAL = 100

class UploadService {
  private chunkSize = CHUNK_SIZE
  private maxConcurrent = MAX_CONCURRENT
  private maxRetries = MAX_RETRIES
  private retryDelays = RETRY_DELAYS
  private abortController: AbortController | null = null
  private pausedUploads: Map<string, boolean> = new Map()
  private uploadPromises: Map<string, Promise<void>> = new Map()

  async calculateMD5(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = async () => {
        try {
          const buffer = reader.result as ArrayBuffer
          const hashBuffer = await crypto.subtle.digest('SHA-256', buffer)
          const hashArray = Array.from(new Uint8Array(hashBuffer))
          const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
          resolve(hashHex.substring(0, 32))
        } catch (error) {
          reject(error)
        }
      }
      reader.onerror = () => reject(reader.error)
      reader.readAsArrayBuffer(file)
    })
  }

  async initUpload(params: InitUploadRequest): Promise<Result<InitUploadResponse>> {
    return await UploadApi.initUpload(params)
  }

  async uploadChunk(
    uploadId: string,
    chunkIndex: number,
    chunk: Blob,
    chunkMd5: string,
    signal?: AbortSignal
  ): Promise<Result<ChunkUploadResponse>> {
    return await UploadApi.uploadChunk(uploadId, chunkIndex, chunk, chunkMd5, signal)
  }

  async completeUpload(uploadId: string, fileHash: string): Promise<Result<CompleteUploadResponse>> {
    return await UploadApi.completeUpload(uploadId, fileHash)
  }

  async getUploadStatus(uploadId: string): Promise<Result<UploadStatusResponse>> {
    return await UploadApi.getUploadStatus(uploadId)
  }

  async cancelUpload(uploadId: string): Promise<Result<{ message: string }>> {
    this.abortController?.abort()
    this.pausedUploads.delete(uploadId)
    this.uploadPromises.delete(uploadId)
    return await UploadApi.cancelUpload(uploadId)
  }

  async getFileInfo(fileId: string): Promise<Result<{
    file_id: string
    source_url: string
    thumb_url: string
    file_size: number
    width: number
    height: number
  }>> {
    return await UploadApi.getFileInfo(fileId)
  }

  pauseUpload(uploadId: string): void {
    this.pausedUploads.set(uploadId, true)
    const uploadStore = useUploadStore()
    uploadStore.pauseUpload(uploadId)
  }

  resumeUpload(uploadId: string): void {
    this.pausedUploads.set(uploadId, false)
    const uploadStore = useUploadStore()
    uploadStore.resumeUpload(uploadId)
  }

  isPaused(uploadId: string): boolean {
    return this.pausedUploads.get(uploadId) === true
  }

  private createChunks(file: File): { index: number; chunk: Blob }[] {
    const chunks: { index: number; chunk: Blob }[] = []
    let start = 0
    let index = 0

    while (start < file.size) {
      const end = Math.min(start + this.chunkSize, file.size)
      chunks.push({
        index,
        chunk: file.slice(start, end)
      })
      start = end
      index++
    }

    return chunks
  }

  private async waitWhilePaused(uploadId: string): Promise<void> {
    while (this.pausedUploads.get(uploadId) === true) {
      await new Promise(resolve => setTimeout(resolve, PAUSE_CHECK_INTERVAL))
    }
  }

  private async uploadChunkWithRetry(
    uploadId: string,
    chunkIndex: number,
    chunk: Blob,
    chunkMd5: string,
    retries = 0,
    signal?: AbortSignal
  ): Promise<Result<ChunkUploadResponse>> {
    try {
      return await this.uploadChunk(uploadId, chunkIndex, chunk, chunkMd5, signal)
    } catch (error: any) {
      if (error.message === 'UPLOAD_CANCELLED') {
        throw error
      }
      if (retries < this.maxRetries) {
        const delay = this.retryDelays[retries] || this.retryDelays[this.retryDelays.length - 1]
        await new Promise(resolve => setTimeout(resolve, delay))
        return this.uploadChunkWithRetry(uploadId, chunkIndex, chunk, chunkMd5, retries + 1, signal)
      }
      throw error
    }
  }

  async uploadFile(file: File, options: UploadOptions): Promise<{ file_id: string; source_url: string; thumb_url: string; upload_id?: string }> {
    const uploadStore = useUploadStore()

    if (file.size > MAX_FILE_SIZE) {
      throw new Error('文件大小超过限制（最大500MB）')
    }

    this.abortController = new AbortController()
    const { signal } = this.abortController

    const fileMd5 = await this.calculateMD5(file)

    const initResponse = await this.initUpload({
      space_id: (options as any).space_id,
      scene_code: (options as any).scene_code,
      filename: file.name,
      file_size: file.size,
      file_hash: fileMd5
    })

    if (signal.aborted) {
      throw new Error('UPLOAD_CANCELLED')
    }

    if (initResponse.code !== 200 || !initResponse.data) {
      throw new Error(initResponse.msg || '初始化上传失败')
    }

    if (initResponse.data.instant) {
      return {
        file_id: initResponse.data.file_id!,
        source_url: initResponse.data.source_url!,
        thumb_url: initResponse.data.thumb_url!,
        upload_id: initResponse.data.upload_id
      }
    }

    const uploadId = initResponse.data.upload_id!
    
    options.onInit?.(uploadId)
    
    const totalChunks = initResponse.data.total_chunks!
    const uploadedChunks = new Set<number>(initResponse.data.uploaded_chunks || [])
    this.pausedUploads.set(uploadId, false)

    const task: UploadTask = {
      uploadId,
      file,
      fileName: file.name,
      fileSize: file.size,
      fileMd5,
      totalChunks,
      uploadedChunks: Array.from(uploadedChunks),
      status: 'uploading',
      percentage: uploadedChunks.size / totalChunks * 100,
      speed: '',
      startTime: Date.now(),
      uploadedBytes: uploadedChunks.size * this.chunkSize
    }

    uploadStore.addUploadTask(task)

    const chunks = this.createChunks(file)
    const pendingChunks = chunks.filter(c => !uploadedChunks.has(c.index))

    const uploadQueue: { index: number; chunk: Blob }[] = [...pendingChunks]
    const activeUploads: Promise<void>[] = []
    let isCancelled = false

    const processQueue = async (): Promise<void> => {
      while (uploadQueue.length > 0 && !isCancelled) {
        if (signal.aborted) {
          isCancelled = true
          break
        }

        await this.waitWhilePaused(uploadId)
        if (this.pausedUploads.get(uploadId) === undefined) {
          isCancelled = true
          break
        }

        const item = uploadQueue.shift()
        if (!item) break

        if (signal.aborted) {
          isCancelled = true
          break
        }

        const chunkMd5 = await this.calculateChunkMD5(item.chunk)

        const uploadPromise = this.uploadChunkWithRetry(uploadId, item.index, item.chunk, chunkMd5, 0, signal)
          .then(response => {
            if (response.code === 200 && response.data) {
              uploadedChunks.add(item.index)
              task.uploadedChunks = Array.from(uploadedChunks)
              task.percentage = uploadedChunks.size / totalChunks * 100
              task.uploadedBytes = uploadedChunks.size * this.chunkSize

              const elapsed = (Date.now() - task.startTime) / 1000
              const speed = task.uploadedBytes / elapsed
              task.speed = this.formatSpeed(speed)

              uploadStore.updateUploadTask(uploadId, {
                uploadedChunks: task.uploadedChunks,
                percentage: task.percentage,
                uploadedBytes: task.uploadedBytes,
                speed: task.speed
              })

              options.onProgress?.({
                uploaded_chunks: uploadedChunks.size,
                total_chunks: totalChunks,
                percentage: task.percentage,
                speed: task.speed,
                uploaded_bytes: task.uploadedBytes,
                total_bytes: file.size
              })
            }
          })
          .catch(error => {
            if (error.message === 'UPLOAD_CANCELLED') {
              isCancelled = true
              return
            }
            console.error(`Chunk ${item.index} upload failed:`, error)
            throw error
          })

        this.uploadPromises.set(`${uploadId}_${item.index}`, uploadPromise as any)
        activeUploads.push(uploadPromise)

        if (activeUploads.length >= this.maxConcurrent) {
          await Promise.race(activeUploads)
        }

        if (this.pausedUploads.get(uploadId) === true) {
          await Promise.all(activeUploads)
          continue
        }

        const completedUploads = await Promise.allSettled(activeUploads)
        activeUploads.length = 0

        for (const result of completedUploads) {
          if (result.status === 'rejected' && result.reason?.message !== 'UPLOAD_CANCELLED') {
            throw new Error(`Chunk upload failed`)
          }
        }
      }

      if (!isCancelled) {
        await Promise.all(activeUploads)
      }
    }

    try {
      await processQueue()

      if (signal.aborted || isCancelled) {
        uploadStore.updateUploadTask(uploadId, { status: 'cancelled' })
        throw new Error('UPLOAD_CANCELLED')
      }

      if (uploadedChunks.size !== totalChunks) {
        uploadStore.updateUploadTask(uploadId, { status: 'failed' })
        throw new Error('上传不完整')
      }

      uploadStore.updateUploadTask(uploadId, { status: 'merging' })

      const completeResponse = await this.completeUpload(uploadId, fileMd5)

      if (completeResponse.code !== 200 || !completeResponse.data) {
        uploadStore.updateUploadTask(uploadId, { status: 'failed' })
        throw new Error(completeResponse.msg || '完成上传失败')
      }

      uploadStore.updateUploadTask(uploadId, {
        status: 'completed',
        percentage: 100
      })

      this.pausedUploads.delete(uploadId)
      this.uploadPromises.delete(uploadId)
      this.abortController = null

      return {
        file_id: completeResponse.data.file_id,
        source_url: completeResponse.data.source_url,
        thumb_url: completeResponse.data.thumb_url,
        upload_id: uploadId
      }
    } catch (error: any) {
      if (error.message === 'UPLOAD_CANCELLED') {
        try {
          await this.cancelUpload(uploadId)
        } catch (e) {
          console.error('Cancel upload failed:', e)
        }
      }
      this.pausedUploads.delete(uploadId)
      this.uploadPromises.delete(uploadId)
      this.abortController = null
      throw error
    }
  }

  private async calculateChunkMD5(chunk: Blob): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = async () => {
        try {
          const buffer = reader.result as ArrayBuffer
          const hashBuffer = await crypto.subtle.digest('SHA-256', buffer)
          const hashArray = Array.from(new Uint8Array(hashBuffer))
          const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
          resolve(hashHex.substring(0, 32))
        } catch (error) {
          reject(error)
        }
      }
      reader.onerror = () => reject(reader.error)
      reader.readAsArrayBuffer(chunk)
    })
  }

  private formatSpeed(bytesPerSecond: number): string {
    if (bytesPerSecond < 1024) {
      return `${bytesPerSecond.toFixed(0)} B/s`
    } else if (bytesPerSecond < 1024 * 1024) {
      return `${(bytesPerSecond / 1024).toFixed(2)} KB/s`
    } else {
      return `${(bytesPerSecond / (1024 * 1024)).toFixed(2)} MB/s`
    }
  }

  formatFileSize(bytes: number): string {
    if (bytes < 1024) {
      return `${bytes} B`
    } else if (bytes < 1024 * 1024) {
      return `${(bytes / 1024).toFixed(2)} KB`
    } else if (bytes < 1024 * 1024 * 1024) {
      return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
    } else {
      return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`
    }
  }
}

export default new UploadService()
