export type InitUploadRequest = {
  space_id: number
  filename: string
  file_size: number
  file_hash: string
}

export type InitUploadResponse = {
  upload_id?: string
  instant: boolean
  file_id?: string
  source_url?: string
  thumb_url?: string
  chunk_size?: number
  total_chunks?: number
  uploaded_chunks?: number[]
}

export type ChunkUploadResponse = {
  chunk_index: number
  uploaded_chunks: number[]
  total_chunks: number
  upload_id: string
}

export type CompleteUploadRequest = {
  upload_id: string
  file_hash: string
}

export type CompleteUploadResponse = {
  file_id: string
  source_url: string
  thumb_url: string
  file_size: number
}

export type UploadStatusResponse = {
  upload_id: string
  status: 'pending' | 'uploading' | 'merging' | 'completed' | 'failed' | 'paused' | 'cancelled'
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
  totalChunks: number
  uploadedChunks: number[]
  status: 'pending' | 'uploading' | 'merging' | 'completed' | 'failed' | 'paused' | 'cancelled'
  percentage: number
  speed: string
  startTime: number
  uploadedBytes: number
}

export type WebSocketMessage = {
  type: 'progress' | 'complete' | 'error' | 'merge_start' | 'merge_progress' | 'slice_progress' | 'slice_complete' | 'slice_error'
  upload_id?: string
  user_id: number
  data: ProgressData | CompleteData | ErrorData | MergeProgressData | SliceProgressData | SliceCompleteData | SliceErrorData
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

export type SliceProgressData = {
  task_id: string
  scene_id: number
  status: string
  progress: number
  stage: string
  message: string
  timestamp: number
}

export type SliceCompleteData = {
  task_id: string
  scene_id: number
  status: string
  tile_url: string
  preview_url: string
  timestamp: number
}

export type SliceErrorData = {
  task_id: string
  scene_id: number
  error: string
  timestamp: number
}

export type UploadOptions = {
  space_id?: number
  scene_code?: string
  title?: string
  onInit?: (uploadId: string) => void
  onProgress?: (progress: ProgressData) => void
  onComplete?: (data: { file_id: string; source_url: string; thumb_url: string }) => void
  onError?: (error: ErrorData) => void
}

export type FileInfo = {
  file_id: string
  source_url: string
  thumb_url: string
  file_size: number
  width: number
  height: number
  created_at: string
}
