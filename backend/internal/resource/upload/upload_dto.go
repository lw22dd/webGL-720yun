package upload

import "time"

type InitUploadRequest struct {
	SpaceID  uint   `json:"space_id" binding:"required"`
	FileName string `json:"filename" binding:"required"`
	FileSize int64  `json:"file_size" binding:"required,min=1"`
	FileMD5  string `json:"file_hash" binding:"required,len=32"`
}

type InitUploadResponse struct {
	UploadID       string `json:"upload_id,omitempty"`
	Instant        bool   `json:"instant"`
	FileID         string `json:"file_id,omitempty"`
	SourceURL      string `json:"source_url,omitempty"`
	ThumbURL       string `json:"thumb_url,omitempty"`
	ChunkSize      int    `json:"chunk_size,omitempty"`
	TotalChunks    int    `json:"total_chunks,omitempty"`
	UploadedChunks []int  `json:"uploaded_chunks,omitempty"`
}

type ChunkUploadRequest struct {
	UploadID   string `form:"upload_id" binding:"required"`
	ChunkIndex int    `form:"chunk_index" binding:"required,min=0"`
	ChunkMD5   string `form:"chunk_hash" binding:"required,len=32"`
}

type ChunkUploadResponse struct {
	ChunkIndex     int    `json:"chunk_index"`
	UploadedChunks []int  `json:"uploaded_chunks"`
	TotalChunks    int    `json:"total_chunks"`
	UploadID       string `json:"upload_id"`
}

type CompleteUploadRequest struct {
	UploadID string `json:"upload_id" binding:"required"`
	FileMD5  string `json:"file_hash" binding:"required,len=32"`
}

type CompleteUploadResponse struct {
	FileID    string `json:"file_id"`
	SourceURL string `json:"source_url"`
	ThumbURL  string `json:"thumb_url"`
	FileSize  int64  `json:"file_size"`
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

type FileInfo struct {
	FileID    string    `json:"file_id"`
	SpaceID   uint      `json:"space_id"`
	SpaceName string    `json:"space_name"`
	SourceURL string    `json:"source_url"`
	ThumbURL  string    `json:"thumb_url"`
	FileSize  int64     `json:"file_size"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	CreatedAt time.Time `json:"created_at"`
}

type UploadTask struct {
	UploadID      string `json:"upload_id"`
	UserID        uint   `json:"user_id"`
	SpaceID       uint   `json:"space_id"`
	SpaceName     string `json:"space_name"`
	FileName      string `json:"file_name"`
	FileSize      int64  `json:"file_size"`
	FileMD5       string `json:"file_md5"`
	TotalChunks   int    `json:"total_chunks"`
	ChunkSize     int    `json:"chunk_size"`
	Status        string `json:"status"`
	UploadedBytes int64  `json:"uploaded_bytes"`
}
