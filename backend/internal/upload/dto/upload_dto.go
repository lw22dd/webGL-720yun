package dto

type InitUploadRequest struct {
	FileName  string `json:"file_name" binding:"required"`
	FileSize  int64  `json:"file_size" binding:"required,min=1"`
	FileMD5   string `json:"file_md5" binding:"required,len=32"`
	SpaceID   uint   `json:"space_id" binding:"required"`
	SceneCode string `json:"scene_code" binding:"required"`
	Title     string `json:"title"`
}

type InitUploadResponse struct {
	UploadID        string `json:"upload_id"`
	SkipUpload      bool   `json:"skip_upload"`
	SceneID         *uint  `json:"scene_id,omitempty"`
	ChunkSize       int    `json:"chunk_size"`
	TotalChunks     int    `json:"total_chunks"`
	UploadedChunks  []int  `json:"uploaded_chunks"`
	UploadTaskExist bool   `json:"upload_task_exist"`
}

type ChunkUploadResponse struct {
	ChunkIndex     int    `json:"chunk_index"`
	UploadedChunks []int  `json:"uploaded_chunks"`
	TotalChunks    int    `json:"total_chunks"`
	UploadID       string `json:"upload_id"`
}

type MergeUploadRequest struct {
	UploadID string `json:"upload_id" binding:"required"`
}

type MergeUploadResponse struct {
	SceneID      uint   `json:"scene_id"`
	SourceURL    string `json:"source_url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type UploadStatusResponse struct {
	UploadID       string `json:"upload_id"`
	Status         string `json:"status"`
	UploadedChunks []int  `json:"uploaded_chunks"`
	TotalChunks    int    `json:"total_chunks"`
	Percentage     int    `json:"percentage"`
	FileName       string `json:"file_name"`
	FileSize       int64  `json:"file_size"`
}

type UploadTask struct {
	UploadID      string `json:"upload_id"`
	UserID        uint   `json:"user_id"`
	FileName      string `json:"file_name"`
	FileSize      int64  `json:"file_size"`
	FileMD5       string `json:"file_md5"`
	SpaceID       uint   `json:"space_id"`
	SceneCode     string `json:"scene_code"`
	Title         string `json:"title"`
	TotalChunks   int    `json:"total_chunks"`
	ChunkSize     int    `json:"chunk_size"`
	Status        string `json:"status"`
	UploadedBytes int64  `json:"uploaded_bytes"`
}
