package upload

import (
	"time"

	"webGL-720yun/pkg/redis"
)

const (
	TaskStatusPending   = "pending"
	TaskStatusUploading = "uploading"
	TaskStatusMerging   = "merging"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"
)

const (
	TaskExpiration       = 48 * time.Hour
	MaxConcurrentUploads = 3
)

type UploadRepository struct {
	redis *redis.RedisService
}

func NewUploadRepository(redis *redis.RedisService) *UploadRepository {
	return &UploadRepository{redis: redis}
}

func (r *UploadRepository) CreateTask(task *UploadTask) error {
	taskData := map[string]interface{}{
		"upload_id":      task.UploadID,
		"user_id":        task.UserID,
		"file_name":      task.FileName,
		"file_size":      task.FileSize,
		"file_md5":       task.FileMD5,
		"space_id":       task.SpaceID,
		"scene_code":     task.SceneCode,
		"title":          task.Title,
		"total_chunks":   task.TotalChunks,
		"chunk_size":     task.ChunkSize,
		"status":         task.Status,
		"uploaded_bytes": task.UploadedBytes,
	}
	return r.redis.CreateUploadTask(task.UploadID, taskData, TaskExpiration)
}

func (r *UploadRepository) GetTask(uploadID string) (*UploadTask, error) {
	taskData, err := r.redis.GetUploadTask(uploadID)
	if err != nil {
		return nil, err
	}
	if taskData == nil {
		return nil, nil
	}

	task := &UploadTask{}
	if v, ok := taskData["upload_id"].(string); ok {
		task.UploadID = v
	}
	if v, ok := taskData["user_id"].(float64); ok {
		task.UserID = uint(v)
	}
	if v, ok := taskData["file_name"].(string); ok {
		task.FileName = v
	}
	if v, ok := taskData["file_size"].(float64); ok {
		task.FileSize = int64(v)
	}
	if v, ok := taskData["file_md5"].(string); ok {
		task.FileMD5 = v
	}
	if v, ok := taskData["space_id"].(float64); ok {
		task.SpaceID = uint(v)
	}
	if v, ok := taskData["scene_code"].(string); ok {
		task.SceneCode = v
	}
	if v, ok := taskData["title"].(string); ok {
		task.Title = v
	}
	if v, ok := taskData["total_chunks"].(float64); ok {
		task.TotalChunks = int(v)
	}
	if v, ok := taskData["chunk_size"].(float64); ok {
		task.ChunkSize = int(v)
	}
	if v, ok := taskData["status"].(string); ok {
		task.Status = v
	}
	if v, ok := taskData["uploaded_bytes"].(float64); ok {
		task.UploadedBytes = int64(v)
	}

	return task, nil
}

func (r *UploadRepository) UpdateTaskStatus(uploadID string, status string) error {
	task, err := r.GetTask(uploadID)
	if err != nil {
		return err
	}
	if task == nil {
		return nil
	}

	taskData := map[string]interface{}{
		"upload_id":      task.UploadID,
		"user_id":        task.UserID,
		"file_name":      task.FileName,
		"file_size":      task.FileSize,
		"file_md5":       task.FileMD5,
		"space_id":       task.SpaceID,
		"scene_code":     task.SceneCode,
		"title":          task.Title,
		"total_chunks":   task.TotalChunks,
		"chunk_size":     task.ChunkSize,
		"status":         status,
		"uploaded_bytes": task.UploadedBytes,
	}
	return r.redis.UpdateUploadTask(uploadID, taskData)
}

func (r *UploadRepository) UpdateTaskUploadedBytes(uploadID string, bytes int64) error {
	task, err := r.GetTask(uploadID)
	if err != nil {
		return err
	}
	if task == nil {
		return nil
	}

	taskData := map[string]interface{}{
		"upload_id":      task.UploadID,
		"user_id":        task.UserID,
		"file_name":      task.FileName,
		"file_size":      task.FileSize,
		"file_md5":       task.FileMD5,
		"space_id":       task.SpaceID,
		"scene_code":     task.SceneCode,
		"title":          task.Title,
		"total_chunks":   task.TotalChunks,
		"chunk_size":     task.ChunkSize,
		"status":         task.Status,
		"uploaded_bytes": bytes,
	}
	return r.redis.UpdateUploadTask(uploadID, taskData)
}

func (r *UploadRepository) DeleteTask(uploadID string) error {
	return r.redis.DeleteUploadTask(uploadID)
}

func (r *UploadRepository) AddUploadedChunk(uploadID string, chunkIndex int) error {
	return r.redis.AddUploadedChunk(uploadID, chunkIndex)
}

func (r *UploadRepository) GetUploadedChunks(uploadID string) ([]int, error) {
	return r.redis.GetUploadedChunks(uploadID)
}

func (r *UploadRepository) IsChunkUploaded(uploadID string, chunkIndex int) (bool, error) {
	return r.redis.IsChunkUploaded(uploadID, chunkIndex)
}

func (r *UploadRepository) GetUploadedChunkCount(uploadID string) (int, error) {
	return r.redis.GetUploadedChunkCount(uploadID)
}

func (r *UploadRepository) CanUserUpload(userID uint) (bool, error) {
	count, err := r.redis.GetUserUploadCount(userID)
	if err != nil {
		return false, err
	}
	return count < MaxConcurrentUploads, nil
}

func (r *UploadRepository) IncrementUserUploadCount(userID uint) (int, error) {
	return r.redis.IncrementUserUploadCount(userID)
}

func (r *UploadRepository) DecrementUserUploadCount(userID uint) (int, error) {
	return r.redis.DecrementUserUploadCount(userID)
}

func (r *UploadRepository) GetUserUploadCount(userID uint) (int, error) {
	return r.redis.GetUserUploadCount(userID)
}
