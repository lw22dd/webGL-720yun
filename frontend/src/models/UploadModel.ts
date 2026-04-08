export type InitUploadRequest = {
  file_name: string
  file_size: number
  file_md5: string
  space_id: number
  scene_code: string
  title?: string
}

export type InitUploadResponse = {
  upload_id: string
  skip_upload: boolean
  scene_id?: number
  chunk_size: number
  total_chunks: number
  uploaded_chunks: number[]
  upload_task_exist: boolean
}

export type ChunkUploadResponse = {
  chunk_index: number
  uploaded_chunks: number[]
  total_chunks: number
  upload_id: string
}

export type MergeUploadRequest = {
  upload_id: string
}

export type MergeUploadResponse = {
  scene_id: number
  source_url: string
  thumbnail_url: string
}

export type UploadStatusResponse = {
  upload_id: string
  status: 'pending' | 'uploading' | 'merging' | 'completed' | 'failed'
  uploaded_chunks: number[]
  total_chunks: number
  percentage: number
  file_name: string
  file_size: number
}

export type UploadTask = {
  uploadId: string
  file: File
  fileName: string
  fileSize: number
  fileMd5: string
  spaceId: number
  sceneCode: string
  title: string
  totalChunks: number
  uploadedChunks: number[]
  status: 'pending' | 'uploading' | 'merging' | 'completed' | 'failed' | 'paused'
  percentage: number
  speed: string
  startTime: number
  uploadedBytes: number
}

export type WebSocketMessage = {
  type: 'progress' | 'complete' | 'error' | 'merge_start' | 'merge_progress'
  upload_id: string
  user_id: number
  data: ProgressData | CompleteData | ErrorData | MergeProgressData
  timestamp: number
}

export type ProgressData = {
  uploaded_chunks: number
  total_chunks: number
  percentage: number
  speed: string
  uploaded_bytes: number
  total_bytes: number
}

export type CompleteData = {
  scene_id: number
  source_url: string
  thumbnail_url: string
}

export type ErrorData = {
  code: number
  message: string
}

export type MergeProgressData = {
  stage: string
  percentage: number
  message: string
}

export type UploadOptions = {
  spaceId: number
  sceneCode: string
  title: string
  onProgress?: (progress: ProgressData) => void
  onComplete?: (data: CompleteData) => void
  onError?: (error: ErrorData) => void
}
