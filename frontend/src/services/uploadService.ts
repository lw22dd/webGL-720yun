<<<<<<< HEAD
import type { Result } from '@/models/Result'
=======
import type { Result } from '@/models/result.model'
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
import type {
  InitUploadRequest,
  InitUploadResponse,
  ChunkUploadResponse,
<<<<<<< HEAD
  MergeUploadResponse,
  UploadStatusResponse,
  UploadOptions,
  UploadTask
} from '@/models/UploadModel'
import Axios from '@/utils/axios'
import { useSceneStore } from '@/stores/sceneStore'
=======
  CompleteUploadResponse,
  UploadStatusResponse,
  UploadOptions,
  UploadTask
} from '@/models/upload.model'
import UploadApi from '@/services/api/upload.api'
import { useUploadStore } from '@/stores/scene/upload.store'
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503

const CHUNK_SIZE = 5 * 1024 * 1024
const MAX_CONCURRENT = 4
const MAX_RETRIES = 3
const RETRY_DELAYS = [1000, 2000, 4000]
const MAX_FILE_SIZE = 500 * 1024 * 1024
<<<<<<< HEAD
=======
const PAUSE_CHECK_INTERVAL = 100
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503

class UploadService {
  private chunkSize = CHUNK_SIZE
  private maxConcurrent = MAX_CONCURRENT
  private maxRetries = MAX_RETRIES
  private retryDelays = RETRY_DELAYS
<<<<<<< HEAD
=======
  private abortController: AbortController | null = null
  private pausedUploads: Map<string, boolean> = new Map()
  private uploadPromises: Map<string, Promise<void>> = new Map()
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503

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
<<<<<<< HEAD
    return await Axios.post('/upload/init', params)
  }

  async uploadChunk(uploadId: string, chunkIndex: number, chunk: Blob): Promise<Result<ChunkUploadResponse>> {
    const formData = new FormData()
    formData.append('upload_id', uploadId)
    formData.append('chunk_index', chunkIndex.toString())
    formData.append('chunk_data', chunk)

    return await Axios.post('/upload/chunk', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      timeout: 60000
    })
  }

  async mergeChunks(uploadId: string): Promise<Result<MergeUploadResponse>> {
    return await Axios.post('/upload/merge', { upload_id: uploadId })
  }

  async getUploadStatus(uploadId: string): Promise<Result<UploadStatusResponse>> {
    return await Axios.get(`/upload/status/${uploadId}`)
  }

  async cancelUpload(uploadId: string): Promise<Result<{ message: string }>> {
    return await Axios.delete(`/upload/${uploadId}`)
=======
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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
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

<<<<<<< HEAD
=======
  private async waitWhilePaused(uploadId: string): Promise<void> {
    while (this.pausedUploads.get(uploadId) === true) {
      await new Promise(resolve => setTimeout(resolve, PAUSE_CHECK_INTERVAL))
    }
  }

>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
  private async uploadChunkWithRetry(
    uploadId: string,
    chunkIndex: number,
    chunk: Blob,
<<<<<<< HEAD
    retries = 0
  ): Promise<Result<ChunkUploadResponse>> {
    try {
      return await this.uploadChunk(uploadId, chunkIndex, chunk)
    } catch (error) {
      if (retries < this.maxRetries) {
        const delay = this.retryDelays[retries] || this.retryDelays[this.retryDelays.length - 1]
        await new Promise(resolve => setTimeout(resolve, delay))
        return this.uploadChunkWithRetry(uploadId, chunkIndex, chunk, retries + 1)
=======
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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
      }
      throw error
    }
  }

<<<<<<< HEAD
  async uploadFile(file: File, options: UploadOptions): Promise<void> {
    const sceneStore = useSceneStore()
=======
  async uploadFile(file: File, options: UploadOptions): Promise<{ file_id: string; source_url: string; thumb_url: string; upload_id?: string }> {
    const uploadStore = useUploadStore()
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503

    if (file.size > MAX_FILE_SIZE) {
      throw new Error('文件大小超过限制（最大500MB）')
    }

<<<<<<< HEAD
    const fileMd5 = await this.calculateMD5(file)

    const initResponse = await this.initUpload({
      file_name: file.name,
      file_size: file.size,
      file_md5: fileMd5,
      space_id: options.spaceId,
      scene_code: options.sceneCode,
      title: options.title
    })

=======
    this.abortController = new AbortController()
    const { signal } = this.abortController

    const fileMd5 = await this.calculateMD5(file)

    const initResponse = await this.initUpload({
      space_id: (options as any).space_id,
      filename: file.name,
      file_size: file.size,
      file_hash: fileMd5
    })

    if (signal.aborted) {
      throw new Error('UPLOAD_CANCELLED')
    }

>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
    if (initResponse.code !== 200 || !initResponse.data) {
      throw new Error(initResponse.msg || '初始化上传失败')
    }

<<<<<<< HEAD
    if (initResponse.data.skip_upload) {
      options.onComplete?.({
        scene_id: initResponse.data.scene_id!,
        source_url: '',
        thumbnail_url: ''
      })
      return
    }

    const uploadId = initResponse.data.upload_id
    const totalChunks = initResponse.data.total_chunks
    const uploadedChunks = new Set(initResponse.data.uploaded_chunks || [])
=======
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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503

    const task: UploadTask = {
      uploadId,
      file,
      fileName: file.name,
      fileSize: file.size,
      fileMd5,
<<<<<<< HEAD
      spaceId: options.spaceId,
      sceneCode: options.sceneCode,
      title: options.title,
=======
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
      totalChunks,
      uploadedChunks: Array.from(uploadedChunks),
      status: 'uploading',
      percentage: uploadedChunks.size / totalChunks * 100,
      speed: '',
      startTime: Date.now(),
      uploadedBytes: uploadedChunks.size * this.chunkSize
    }

<<<<<<< HEAD
    sceneStore.addUploadTask(task)
=======
    uploadStore.addUploadTask(task)
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503

    const chunks = this.createChunks(file)
    const pendingChunks = chunks.filter(c => !uploadedChunks.has(c.index))

    const uploadQueue: { index: number; chunk: Blob }[] = [...pendingChunks]
    const activeUploads: Promise<void>[] = []
<<<<<<< HEAD

    const processQueue = async () => {
      while (uploadQueue.length > 0) {
        const item = uploadQueue.shift()
        if (!item) break

        const uploadPromise = this.uploadChunkWithRetry(uploadId, item.index, item.chunk)
=======
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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
          .then(response => {
            if (response.code === 200 && response.data) {
              uploadedChunks.add(item.index)
              task.uploadedChunks = Array.from(uploadedChunks)
              task.percentage = uploadedChunks.size / totalChunks * 100
              task.uploadedBytes = uploadedChunks.size * this.chunkSize

              const elapsed = (Date.now() - task.startTime) / 1000
              const speed = task.uploadedBytes / elapsed
              task.speed = this.formatSpeed(speed)

<<<<<<< HEAD
              sceneStore.updateUploadTask(uploadId, {
=======
              uploadStore.updateUploadTask(uploadId, {
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
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
<<<<<<< HEAD
=======
            if (error.message === 'UPLOAD_CANCELLED') {
              isCancelled = true
              return
            }
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
            console.error(`Chunk ${item.index} upload failed:`, error)
            throw error
          })

<<<<<<< HEAD
=======
        this.uploadPromises.set(`${uploadId}_${item.index}`, uploadPromise as any)
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
        activeUploads.push(uploadPromise)

        if (activeUploads.length >= this.maxConcurrent) {
          await Promise.race(activeUploads)
        }

<<<<<<< HEAD
        const completedUploads = await Promise.allSettled(activeUploads)

        for (let i = 0; i < completedUploads.length; i++) {
          if (completedUploads[i].status === 'rejected') {
=======
        if (this.pausedUploads.get(uploadId) === true) {
          await Promise.all(activeUploads)
          continue
        }

        const completedUploads = await Promise.allSettled(activeUploads)
        activeUploads.length = 0

        for (const result of completedUploads) {
          if (result.status === 'rejected' && result.reason?.message !== 'UPLOAD_CANCELLED') {
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
            throw new Error(`Chunk upload failed`)
          }
        }
      }

<<<<<<< HEAD
      await Promise.all(activeUploads)
    }

    await processQueue()

    if (uploadedChunks.size !== totalChunks) {
      throw new Error('上传不完整')
    }

    sceneStore.updateUploadTask(uploadId, { status: 'merging' })

    const mergeResponse = await this.mergeChunks(uploadId)

    if (mergeResponse.code !== 200 || !mergeResponse.data) {
      sceneStore.updateUploadTask(uploadId, { status: 'failed' })
      throw new Error(mergeResponse.msg || '合并文件失败')
    }

    sceneStore.updateUploadTask(uploadId, {
      status: 'completed',
      percentage: 100
    })

    options.onComplete?.(mergeResponse.data)
=======
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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
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
