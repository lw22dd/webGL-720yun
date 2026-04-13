package upload

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"

	"webGL-720yun/internal/resource/repository"
	"webGL-720yun/pkg/image"
	"webGL-720yun/pkg/minio_client"
)

const (
	ChunkSize   = 5 * 1024 * 1024
	MaxFileSize = 500 * 1024 * 1024
)

type UploadService struct {
	uploadRepo     *UploadRepository
	spaceRepo      *repository.SpaceRepository
	minioClient    *minio_client.MinIOClient
	imageProcessor *image.Processor
}

func NewUploadService(
	uploadRepo *UploadRepository,
	db *gorm.DB,
	minioClient *minio_client.MinIOClient,
) *UploadService {
	return &UploadService{
		uploadRepo:     uploadRepo,
		spaceRepo:      repository.NewSpaceRepository(db),
		minioClient:    minioClient,
		imageProcessor: image.NewProcessor(),
	}
}

func (s *UploadService) InitUpload(req *InitUploadRequest, userID uint) (*InitUploadResponse, error) {
	if req.FileSize > MaxFileSize {
		return nil, errors.New("文件大小超过限制（最大500MB）")
	}

	space, err := s.spaceRepo.FindByID(req.SpaceID)
	if err != nil {
		return nil, fmt.Errorf("空间不存在: %w", err)
	}

	canUpload, err := s.uploadRepo.CanUserUpload(userID)
	if err != nil {
		return nil, fmt.Errorf("检查上传权限失败: %w", err)
	}
	if !canUpload {
		return nil, errors.New("同时上传文件数量超过限制（最多3个）")
	}

	fileID, err := s.uploadRepo.GetFileIDByMD5(req.FileMD5)
	if err != nil {
		return nil, fmt.Errorf("检查文件MD5失败: %w", err)
	}
	if fileID != "" {
		fileInfo, err := s.uploadRepo.GetFileInfo(fileID)
		if err != nil {
			return nil, fmt.Errorf("获取文件信息失败: %w", err)
		}
		if fileInfo != nil && fileInfo.SpaceID == req.SpaceID {
			return &InitUploadResponse{
				Instant:   true,
				FileID:    fileID,
				SourceURL: fileInfo.SourceURL,
				ThumbURL:  fileInfo.ThumbURL,
			}, nil
		}
	}

	uploadID := uuid.New().String()
	totalChunks := int((req.FileSize + ChunkSize - 1) / ChunkSize)

	task := &UploadTask{
		UploadID:      uploadID,
		UserID:        userID,
		SpaceID:       req.SpaceID,
		SpaceName:     space.Name,
		FileName:      req.FileName,
		FileSize:      req.FileSize,
		FileMD5:       req.FileMD5,
		TotalChunks:   totalChunks,
		ChunkSize:     ChunkSize,
		Status:        TaskStatusPending,
		UploadedBytes: 0,
	}

	if err := s.uploadRepo.CreateTask(task); err != nil {
		return nil, fmt.Errorf("创建上传任务失败: %w", err)
	}

	if _, err := s.uploadRepo.IncrementUserUploadCount(userID); err != nil {
		s.uploadRepo.DeleteTask(uploadID)
		return nil, fmt.Errorf("更新用户上传计数失败: %w", err)
	}

	return &InitUploadResponse{
		UploadID:       uploadID,
		Instant:        false,
		ChunkSize:      ChunkSize,
		TotalChunks:    totalChunks,
		UploadedChunks: []int{},
	}, nil
}

func (s *UploadService) UploadChunk(uploadID string, chunkIndex int, chunkData *multipart.FileHeader, userID uint) (*ChunkUploadResponse, error) {
	task, err := s.uploadRepo.GetTask(uploadID)
	if err != nil {
		return nil, fmt.Errorf("获取上传任务失败: %w", err)
	}
	if task == nil {
		return nil, errors.New("上传任务不存在")
	}
	if task.UserID != userID {
		return nil, errors.New("无权限操作此任务")
	}

	isUploaded, err := s.uploadRepo.IsChunkUploaded(uploadID, chunkIndex)
	if err != nil {
		return nil, fmt.Errorf("检查分片状态失败: %w", err)
	}
	if isUploaded {
		uploadedChunks, _ := s.uploadRepo.GetUploadedChunks(uploadID)
		return &ChunkUploadResponse{
			ChunkIndex:     chunkIndex,
			UploadedChunks: uploadedChunks,
			TotalChunks:    task.TotalChunks,
			UploadID:       uploadID,
		}, nil
	}

	chunkObjectName := fmt.Sprintf("temp/%s/chunk_%d", uploadID, chunkIndex)

	src, err := chunkData.Open()
	if err != nil {
		return nil, fmt.Errorf("打开分片文件失败: %w", err)
	}
	defer src.Close()

	chunkBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取分片数据失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := s.minioClient.GetClient()
	bucket := s.minioClient.GetConfig().Bucket

	_, err = client.PutObject(ctx, bucket, chunkObjectName, bytes.NewReader(chunkBytes), int64(len(chunkBytes)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return nil, fmt.Errorf("上传分片到MinIO失败: %w", err)
	}

	if err := s.uploadRepo.AddUploadedChunk(uploadID, chunkIndex); err != nil {
		return nil, fmt.Errorf("记录分片失败: %w", err)
	}

	uploadedBytes := task.UploadedBytes + int64(len(chunkBytes))
	s.uploadRepo.UpdateTaskUploadedBytes(uploadID, uploadedBytes)

	uploadedChunks, _ := s.uploadRepo.GetUploadedChunks(uploadID)

	return &ChunkUploadResponse{
		ChunkIndex:     chunkIndex,
		UploadedChunks: uploadedChunks,
		TotalChunks:    task.TotalChunks,
		UploadID:       uploadID,
	}, nil
}

func (s *UploadService) CompleteUpload(req *CompleteUploadRequest, userID uint) (*CompleteUploadResponse, error) {
	task, err := s.uploadRepo.GetTask(req.UploadID)
	if err != nil {
		return nil, fmt.Errorf("获取上传任务失败: %w", err)
	}
	if task == nil {
		return nil, errors.New("上传任务不存在")
	}
	if task.UserID != userID {
		return nil, errors.New("无权限操作此任务")
	}

	uploadedChunks, err := s.uploadRepo.GetUploadedChunks(req.UploadID)
	if err != nil {
		return nil, fmt.Errorf("获取已上传分片失败: %w", err)
	}

	if len(uploadedChunks) != task.TotalChunks {
		return nil, fmt.Errorf("分片不完整，已上传 %d/%d", len(uploadedChunks), task.TotalChunks)
	}

	s.uploadRepo.UpdateTaskStatus(req.UploadID, TaskStatusMerging)

	tempDir := filepath.Join(os.TempDir(), "upload_"+req.UploadID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tempDir)

	sort.Ints(uploadedChunks)

	mergedFile := filepath.Join(tempDir, "merged.jpg")
	mergedWriter, err := os.Create(mergedFile)
	if err != nil {
		return nil, fmt.Errorf("创建合并文件失败: %w", err)
	}
	defer mergedWriter.Close()

	ctx := context.Background()
	client := s.minioClient.GetClient()
	bucket := s.minioClient.GetConfig().Bucket

	for _, chunkIndex := range uploadedChunks {
		chunkObjectName := fmt.Sprintf("temp/%s/chunk_%d", req.UploadID, chunkIndex)

		obj, err := client.GetObject(ctx, bucket, chunkObjectName, minio.GetObjectOptions{})
		if err != nil {
			mergedWriter.Close()
			return nil, fmt.Errorf("获取分片 %d 失败: %w", chunkIndex, err)
		}

		if _, err := io.Copy(mergedWriter, obj); err != nil {
			obj.Close()
			return nil, fmt.Errorf("合并分片 %d 失败: %w", chunkIndex, err)
		}
		obj.Close()
	}
	mergedWriter.Close()

	if _, err := s.imageProcessor.ValidateFormat(mergedFile); err != nil {
		return nil, fmt.Errorf("文件格式验证失败: %w", err)
	}

	if err := s.imageProcessor.ValidatePanoramaResolution(mergedFile); err != nil {
		return nil, fmt.Errorf("分辨率验证失败: %w", err)
	}

	imageInfo, err := s.imageProcessor.GetImageInfo(mergedFile)
	if err != nil {
		return nil, fmt.Errorf("获取图片信息失败: %w", err)
	}

	fileID := uuid.New().String()
	sourceObjectName := fmt.Sprintf("spaces/%s/sources/%s/source.jpg", task.SpaceName, fileID)
	sourceURL, err := s.minioClient.UploadFile(sourceObjectName, mergedFile, "image/jpeg")
	if err != nil {
		return nil, fmt.Errorf("上传源文件失败: %w", err)
	}

	thumbFile := filepath.Join(tempDir, "thumb.jpg")
	if err := s.imageProcessor.GenerateThumbnail(mergedFile, thumbFile); err != nil {
		return nil, fmt.Errorf("生成缩略图失败: %w", err)
	}

	thumbObjectName := fmt.Sprintf("spaces/%s/previews/%s/thumb.jpg", task.SpaceName, fileID)
	thumbURL, err := s.minioClient.UploadFile(thumbObjectName, thumbFile, "image/jpeg")
	if err != nil {
		return nil, fmt.Errorf("上传缩略图失败: %w", err)
	}

	for _, chunkIndex := range uploadedChunks {
		chunkObjectName := fmt.Sprintf("temp/%s/chunk_%d", req.UploadID, chunkIndex)
		client.RemoveObject(ctx, bucket, chunkObjectName, minio.RemoveObjectOptions{})
	}

	if err := s.uploadRepo.SaveFileMD5(task.FileMD5, fileID); err != nil {
		return nil, fmt.Errorf("保存文件MD5映射失败: %w", err)
	}

	if err := s.uploadRepo.SaveFileInfo(fileID, &FileInfo{
		FileID:    fileID,
		SpaceID:   task.SpaceID,
		SpaceName: task.SpaceName,
		SourceURL: sourceURL,
		ThumbURL:  thumbURL,
		FileSize:  imageInfo.FileSize,
		Width:     imageInfo.Width,
		Height:    imageInfo.Height,
		CreatedAt: time.Now(),
	}); err != nil {
		return nil, fmt.Errorf("保存文件信息失败: %w", err)
	}

	s.uploadRepo.UpdateTaskStatus(req.UploadID, TaskStatusCompleted)
	s.uploadRepo.DeleteTask(req.UploadID)
	s.uploadRepo.DecrementUserUploadCount(userID)

	return &CompleteUploadResponse{
		FileID:    fileID,
		SourceURL: sourceURL,
		ThumbURL:  thumbURL,
		FileSize:  imageInfo.FileSize,
	}, nil
}

func (s *UploadService) GetUploadStatus(uploadID string, userID uint) (*UploadStatusResponse, error) {
	task, err := s.uploadRepo.GetTask(uploadID)
	if err != nil {
		return nil, fmt.Errorf("获取上传任务失败: %w", err)
	}
	if task == nil {
		return nil, errors.New("上传任务不存在")
	}
	if task.UserID != userID {
		return nil, errors.New("无权限查看此任务")
	}

	uploadedChunks, _ := s.uploadRepo.GetUploadedChunks(uploadID)
	percentage := 0
	if task.TotalChunks > 0 {
		percentage = len(uploadedChunks) * 100 / task.TotalChunks
	}

	return &UploadStatusResponse{
		UploadID:       uploadID,
		Status:         task.Status,
		UploadedChunks: uploadedChunks,
		TotalChunks:    task.TotalChunks,
		Percentage:     percentage,
		FileName:       task.FileName,
		FileSize:       task.FileSize,
	}, nil
}

func (s *UploadService) CancelUpload(uploadID string, userID uint) error {
	task, err := s.uploadRepo.GetTask(uploadID)
	if err != nil {
		return fmt.Errorf("获取上传任务失败: %w", err)
	}
	if task == nil {
		return errors.New("上传任务不存在")
	}
	if task.UserID != userID {
		return errors.New("无权限操作此任务")
	}

	ctx := context.Background()
	client := s.minioClient.GetClient()
	bucket := s.minioClient.GetConfig().Bucket

	for i := 0; i < task.TotalChunks; i++ {
		chunkObjectName := fmt.Sprintf("temp/%s/chunk_%d", uploadID, i)
		client.RemoveObject(ctx, bucket, chunkObjectName, minio.RemoveObjectOptions{})
	}

	s.uploadRepo.DeleteTask(uploadID)
	s.uploadRepo.DecrementUserUploadCount(userID)

	return nil
}

func (s *UploadService) GetFileInfo(fileID string) (*FileInfo, error) {
	return s.uploadRepo.GetFileInfo(fileID)
}
