package upload

import (
	"errors"
	"time"

	"webGL-720yun/internal/model"
	"webGL-720yun/pkg/redis"

	"gorm.io/gorm"
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
	MaxConcurrentUploads = 10
)

type UploadRepository struct {
	redis *redis.RedisService
	db    *gorm.DB
}

func NewUploadRepository(redis *redis.RedisService, db *gorm.DB) *UploadRepository {
	return &UploadRepository{redis: redis, db: db}
}

func (r *UploadRepository) CreateTask(task *UploadTask) error {
	taskData := map[string]interface{}{
		"upload_id":      task.UploadID,
		"user_id":        task.UserID,
		"space_id":       task.SpaceID,
		"space_name":     task.SpaceName,
		"space_slug":     task.SpaceSlug,
		"scene_code":     task.SceneCode,
		"file_name":      task.FileName,
		"file_size":      task.FileSize,
		"file_md5":       task.FileMD5,
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
	if v, ok := taskData["space_id"].(float64); ok {
		task.SpaceID = uint(v)
	}
	if v, ok := taskData["space_slug"].(string); ok {
		task.SpaceSlug = v
	}
	if v, ok := taskData["scene_code"].(string); ok {
		task.SceneCode = v
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
		"space_id":       task.SpaceID,
		"space_name":     task.SpaceName,
		"space_slug":     task.SpaceSlug,
		"file_name":      task.FileName,
		"file_size":      task.FileSize,
		"file_md5":       task.FileMD5,
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
		"space_id":       task.SpaceID,
		"space_name":     task.SpaceName,
		"space_slug":     task.SpaceSlug,
		"file_name":      task.FileName,
		"file_size":      task.FileSize,
		"file_md5":       task.FileMD5,
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

func (r *UploadRepository) GetFileIDByMD5(md5 string) (string, error) {
	fileID, err := r.redis.GetFileIDByMD5(md5)
	if err == nil && fileID != "" {
		return fileID, nil
	}

	var scene model.ResScene
	err = r.db.Where("source_file_md5 = ?", md5).First(&scene).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}

	// Hit MySQL, cache back to Redis
	_ = r.SaveFileMD5(md5, scene.FileID)

	return scene.FileID, nil
}

func (r *UploadRepository) SaveFileMD5(md5 string, fileID string) error {
	return r.redis.SaveFileMD5(md5, fileID)
}

func (r *UploadRepository) SaveFileInfo(fileID string, info *FileInfo) error {
	infoData := map[string]interface{}{
		"file_id":    info.FileID,
		"space_id":   info.SpaceID,
		"space_name": info.SpaceName,
		"space_slug": info.SpaceSlug,
		"source_url": info.SourceURL,
		"thumb_url":  info.ThumbURL,
		"file_size":  info.FileSize,
		"width":      info.Width,
		"height":     info.Height,
		"created_at": info.CreatedAt.Unix(),
	}
	return r.redis.SaveFileInfo(fileID, infoData)
}

func (r *UploadRepository) GetFileInfo(fileID string) (*FileInfo, error) {
	infoData, err := r.redis.GetFileInfo(fileID)
	if err == nil && infoData != nil {
		info := &FileInfo{}
		if v, ok := infoData["file_id"].(string); ok {
			info.FileID = v
		}
		if v, ok := infoData["space_id"].(float64); ok {
			info.SpaceID = uint(v)
		}
		if v, ok := infoData["space_name"].(string); ok {
			info.SpaceName = v
		}
		if v, ok := infoData["space_slug"].(string); ok {
			info.SpaceSlug = v
		}
		if v, ok := infoData["source_url"].(string); ok {
			info.SourceURL = v
		}
		if v, ok := infoData["thumb_url"].(string); ok {
			info.ThumbURL = v
		}
		if v, ok := infoData["file_size"].(float64); ok {
			info.FileSize = int64(v)
		}
		if v, ok := infoData["width"].(float64); ok {
			info.Width = int(v)
		}
		if v, ok := infoData["height"].(float64); ok {
			info.Height = int(v)
		}
		if v, ok := infoData["created_at"].(float64); ok {
			info.CreatedAt = time.Unix(int64(v), 0)
		}
		return info, nil
	}

	var scene model.ResScene
	err = r.db.Preload("Space").Where("file_id = ?", fileID).First(&scene).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	info := &FileInfo{
		FileID:    scene.FileID,
		SpaceID:   scene.SpaceID,
		SpaceName: scene.Space.Name,
		SpaceSlug: scene.Space.Slug,
		SourceURL: scene.SourceURL,
		ThumbURL:  scene.ThumbnailURL,
		FileSize:  scene.SourceFileSize,
		Width:     scene.SourceWidth,
		Height:    scene.SourceHeight,
		CreatedAt: scene.CreatedAt,
	}

	_ = r.SaveFileInfo(fileID, info)

	return info, nil
}

func (r *UploadRepository) DeleteFileInfo(fileID string) error {
	return r.redis.DeleteFileInfo(fileID)
}
