import type { Result } from '@/models/Result'
import type {
  InitUploadRequest,
  InitUploadResponse,
  ChunkUploadResponse,
  MergeUploadResponse,
  UploadStatusResponse,
  UploadOptions,
  UploadTask
} from '@/models/UploadModel'
import Axios from '@/utils/axios'
import { useSceneStore } from '@/stores/sceneStore'

const CHUNK_SIZE = 5 * 1024 * 1024
const MAX_CONCURRENT = 4
const MAX_RETRIES = 3
const RETRY_DELAYS = [1000, 2000, 4000]
const MAX_FILE_SIZE = 500 * 1024 * 1024

class UploadService {
  private chunkSize = CHUNK_SIZE
  private maxConcurrent = MAX_CONCURRENT
  private maxRetries = MAX_RETRIES
  private retryDelays = RETRY_DELAYS

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

  private async uploadChunkWithRetry(
    uploadId: string,
    chunkIndex: number,
    chunk: Blob,
    retries = 0
  ): Promise<Result<ChunkUploadResponse>> {
    try {
      return await this.uploadChunk(uploadId, chunkIndex, chunk)
    } catch (error) {
      if (retries < this.maxRetries) {
        const delay = this.retryDelays[retries] || this.retryDelays[this.retryDelays.length - 1]
        await new Promise(resolve => setTimeout(resolve, delay))
        return this.uploadChunkWithRetry(uploadId, chunkIndex, chunk, retries + 1)
      }
      throw error
    }
  }

  async uploadFile(file: File, options: UploadOptions): Promise<void> {
    const sceneStore = useSceneStore()

    if (file.size > MAX_FILE_SIZE) {
      throw new Error('文件大小超过限制（最大500MB）')
    }

    const fileMd5 = await this.calculateMD5(file)

    const initResponse = await this.initUpload({
      file_name: file.name,
      file_size: file.size,
      file_md5: fileMd5,
      space_id: options.spaceId,
      scene_code: options.sceneCode,
      title: options.title
    })

    if (initResponse.code !== 200 || !initResponse.data) {
      throw new Error(initResponse.msg || '初始化上传失败')
    }

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

    const task: UploadTask = {
      uploadId,
      file,
      fileName: file.name,
      fileSize: file.size,
      fileMd5,
      spaceId: options.spaceId,
      sceneCode: options.sceneCode,
      title: options.title,
      totalChunks,
      uploadedChunks: Array.from(uploadedChunks),
      status: 'uploading',
      percentage: uploadedChunks.size / totalChunks * 100,
      speed: '',
      startTime: Date.now(),
      uploadedBytes: uploadedChunks.size * this.chunkSize
    }

    sceneStore.addUploadTask(task)

    const chunks = this.createChunks(file)
    const pendingChunks = chunks.filter(c => !uploadedChunks.has(c.index))

    const uploadQueue: { index: number; chunk: Blob }[] = [...pendingChunks]
    const activeUploads: Promise<void>[] = []

    const processQueue = async () => {
      while (uploadQueue.length > 0) {
        const item = uploadQueue.shift()
        if (!item) break

        const uploadPromise = this.uploadChunkWithRetry(uploadId, item.index, item.chunk)
          .then(response => {
            if (response.code === 200 && response.data) {
              uploadedChunks.add(item.index)
              task.uploadedChunks = Array.from(uploadedChunks)
              task.percentage = uploadedChunks.size / totalChunks * 100
              task.uploadedBytes = uploadedChunks.size * this.chunkSize

              const elapsed = (Date.now() - task.startTime) / 1000
              const speed = task.uploadedBytes / elapsed
              task.speed = this.formatSpeed(speed)

              sceneStore.updateUploadTask(uploadId, {
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
            console.error(`Chunk ${item.index} upload failed:`, error)
            throw error
          })

        activeUploads.push(uploadPromise)

        if (activeUploads.length >= this.maxConcurrent) {
          await Promise.race(activeUploads)
        }

        const completedUploads = await Promise.allSettled(activeUploads)

        for (let i = 0; i < completedUploads.length; i++) {
          if (completedUploads[i].status === 'rejected') {
            throw new Error(`Chunk upload failed`)
          }
        }
      }

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
