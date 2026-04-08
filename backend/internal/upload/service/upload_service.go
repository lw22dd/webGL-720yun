package service

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
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"webGL-720yun/internal/model"
	"webGL-720yun/internal/upload/dto"
	"webGL-720yun/internal/upload/repository"
	"webGL-720yun/pkg/image"
	"webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/websocket"
)

const (
	ChunkSize            = 5 * 1024 * 1024
	MaxFileSize          = 500 * 1024 * 1024
	MaxConcurrentUploads = 3
)

type UploadService struct {
	uploadRepo     *repository.UploadRepository
	sceneRepo      SceneRepository
	spaceRepo      SpaceRepository
	minioClient    *minio_client.MinIOClient
	imageProcessor *image.Processor
	wsHub          *websocket.Hub
}

type SceneRepository interface {
	FindByMD5(md5 string) (*model.ResScene, error)
	FindBySceneCode(sceneCode string) (*model.ResScene, error)
	Create(scene *model.ResScene) error
	FindByID(id uint) (*model.ResScene, error)
}

type SpaceRepository interface {
	FindByID(id uint) (*model.ResSpace, error)
}

func NewUploadService(
	uploadRepo *repository.UploadRepository,
	sceneRepo SceneRepository,
	spaceRepo SpaceRepository,
	minioClient *minio_client.MinIOClient,
	wsHub *websocket.Hub,
) *UploadService {
	return &UploadService{
		uploadRepo:     uploadRepo,
		sceneRepo:      sceneRepo,
		spaceRepo:      spaceRepo,
		minioClient:    minioClient,
		imageProcessor: image.NewProcessor(),
		wsHub:          wsHub,
	}
}

func (s *UploadService) InitUpload(req *dto.InitUploadRequest, userID uint) (*dto.InitUploadResponse, error) {
	if req.FileSize > MaxFileSize {
		return nil, errors.New("文件大小超过限制（最大500MB）")
	}

	canUpload, err := s.uploadRepo.CanUserUpload(userID)
	if err != nil {
		return nil, fmt.Errorf("检查上传权限失败: %w", err)
	}
	if !canUpload {
		return nil, errors.New("同时上传文件数量超过限制（最多3个）")
	}

	space, err := s.spaceRepo.FindByID(req.SpaceID)
	if err != nil {
		return nil, fmt.Errorf("空间不存在: %w", err)
	}

	md5Key := fmt.Sprintf("%s_%s", req.FileName, req.FileMD5)
	existingScene, err := s.sceneRepo.FindByMD5(md5Key)
	if err == nil && existingScene != nil {
		return &dto.InitUploadResponse{
			SkipUpload: true,
			SceneID:    &existingScene.ID,
		}, nil
	}

	existingScene, err = s.sceneRepo.FindBySceneCode(req.SceneCode)
	if err == nil && existingScene != nil {
		return nil, errors.New("场景编码已存在")
	}

	uploadID := uuid.New().String()
	totalChunks := int((req.FileSize + ChunkSize - 1) / ChunkSize)

	task := &dto.UploadTask{
		UploadID:      uploadID,
		UserID:        userID,
		FileName:      req.FileName,
		FileSize:      req.FileSize,
		FileMD5:       req.FileMD5,
		SpaceID:       req.SpaceID,
		SceneCode:     req.SceneCode,
		Title:         req.Title,
		TotalChunks:   totalChunks,
		ChunkSize:     ChunkSize,
		Status:        repository.TaskStatusPending,
		UploadedBytes: 0,
	}

	if err := s.uploadRepo.CreateTask(task); err != nil {
		return nil, fmt.Errorf("创建上传任务失败: %w", err)
	}

	if _, err := s.uploadRepo.IncrementUserUploadCount(userID); err != nil {
		s.uploadRepo.DeleteTask(uploadID)
		return nil, fmt.Errorf("更新用户上传计数失败: %w", err)
	}

	_ = space

	return &dto.InitUploadResponse{
		UploadID:       uploadID,
		SkipUpload:     false,
		ChunkSize:      ChunkSize,
		TotalChunks:    totalChunks,
		UploadedChunks: []int{},
	}, nil
}

func (s *UploadService) UploadChunk(uploadID string, chunkIndex int, chunkData *multipart.FileHeader, userID uint) (*dto.ChunkUploadResponse, error) {
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
		return &dto.ChunkUploadResponse{
			ChunkIndex:     chunkIndex,
			UploadedChunks: uploadedChunks,
			TotalChunks:    task.TotalChunks,
			UploadID:       uploadID,
		}, nil
	}

	space, err := s.spaceRepo.FindByID(task.SpaceID)
	if err != nil {
		return nil, fmt.Errorf("空间不存在: %w", err)
	}

	chunkObjectName := fmt.Sprintf("spaces/%s/temp/%s/chunk_%d", space.Name, uploadID, chunkIndex)

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

	if s.wsHub != nil {
		percentage := float64(len(uploadedChunks)) / float64(task.TotalChunks) * 100
		progressData := &websocket.ProgressData{
			UploadedChunks: len(uploadedChunks),
			TotalChunks:    task.TotalChunks,
			Percentage:     percentage,
			UploadedBytes:  uploadedBytes,
			TotalBytes:     task.FileSize,
		}
		msg := websocket.NewProgressMessage(uploadID, userID, progressData)
		s.wsHub.SendToUser(userID, msg)
	}

	return &dto.ChunkUploadResponse{
		ChunkIndex:     chunkIndex,
		UploadedChunks: uploadedChunks,
		TotalChunks:    task.TotalChunks,
		UploadID:       uploadID,
	}, nil
}

func (s *UploadService) MergeChunks(uploadID string, userID uint) (*dto.MergeUploadResponse, error) {
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

	uploadedChunks, err := s.uploadRepo.GetUploadedChunks(uploadID)
	if err != nil {
		return nil, fmt.Errorf("获取已上传分片失败: %w", err)
	}

	if len(uploadedChunks) != task.TotalChunks {
		return nil, fmt.Errorf("分片不完整，已上传 %d/%d", len(uploadedChunks), task.TotalChunks)
	}

	s.uploadRepo.UpdateTaskStatus(uploadID, repository.TaskStatusMerging)

	if s.wsHub != nil {
		msg := websocket.NewMergeStartMessage(uploadID, userID)
		s.wsHub.SendToUser(userID, msg)
	}

	space, err := s.spaceRepo.FindByID(task.SpaceID)
	if err != nil {
		return nil, fmt.Errorf("空间不存在: %w", err)
	}

	tempDir := filepath.Join(os.TempDir(), "upload_"+uploadID)
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

	for i, chunkIndex := range uploadedChunks {
		chunkObjectName := fmt.Sprintf("spaces/%s/temp/%s/chunk_%d", space.Name, uploadID, chunkIndex)

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

		if s.wsHub != nil {
			percentage := float64(i+1) / float64(task.TotalChunks) * 100
			mergeData := &websocket.MergeProgressData{
				Stage:      "merging",
				Percentage: percentage,
				Message:    fmt.Sprintf("正在合并分片 %d/%d", i+1, task.TotalChunks),
			}
			msg := websocket.NewMergeProgressMessage(uploadID, userID, mergeData)
			s.wsHub.SendToUser(userID, msg)
		}
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

	sourceObjectName := fmt.Sprintf("spaces/%s/sources/%s_sphere.jpg", space.Name, task.SceneCode)
	sourceURL, err := s.minioClient.UploadFile(sourceObjectName, mergedFile, "image/jpeg")
	if err != nil {
		return nil, fmt.Errorf("上传源文件失败: %w", err)
	}

	thumbFile := filepath.Join(tempDir, "thumb.jpg")
	if err := s.imageProcessor.GenerateThumbnail(mergedFile, thumbFile); err != nil {
		return nil, fmt.Errorf("生成缩略图失败: %w", err)
	}

	thumbObjectName := fmt.Sprintf("spaces/%s/thumbnails/%s_thumb.jpg", space.Name, task.SceneCode)
	thumbURL, err := s.minioClient.UploadFile(thumbObjectName, thumbFile, "image/jpeg")
	if err != nil {
		return nil, fmt.Errorf("上传缩略图失败: %w", err)
	}

	scene := &model.ResScene{
		SpaceID:        task.SpaceID,
		Title:          task.Title,
		SceneCode:      task.SceneCode,
		PanoramaType:   "equirectangular",
		SourceURL:      sourceURL,
		SourceWidth:    imageInfo.Width,
		SourceHeight:   imageInfo.Height,
		SourceFileSize: imageInfo.FileSize,
		SourceFileMD5:  fmt.Sprintf("%s_%s", task.FileName, task.FileMD5),
		ThumbnailURL:   thumbURL,
		InitialFOV:     100,
		Status:         1,
	}

	if err := s.sceneRepo.Create(scene); err != nil {
		return nil, fmt.Errorf("创建场景记录失败: %w", err)
	}

	for _, chunkIndex := range uploadedChunks {
		chunkObjectName := fmt.Sprintf("spaces/%s/temp/%s/chunk_%d", space.Name, uploadID, chunkIndex)
		client.RemoveObject(ctx, bucket, chunkObjectName, minio.RemoveObjectOptions{})
	}

	s.uploadRepo.UpdateTaskStatus(uploadID, repository.TaskStatusCompleted)
	s.uploadRepo.DeleteTask(uploadID)
	s.uploadRepo.DecrementUserUploadCount(userID)

	if s.wsHub != nil {
		completeData := &websocket.CompleteData{
			SceneID:      scene.ID,
			SourceURL:    sourceURL,
			ThumbnailURL: thumbURL,
		}
		msg := websocket.NewCompleteMessage(uploadID, userID, completeData)
		s.wsHub.SendToUser(userID, msg)
	}

	return &dto.MergeUploadResponse{
		SceneID:      scene.ID,
		SourceURL:    sourceURL,
		ThumbnailURL: thumbURL,
	}, nil
}

func (s *UploadService) GetUploadStatus(uploadID string, userID uint) (*dto.UploadStatusResponse, error) {
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

	return &dto.UploadStatusResponse{
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

	space, err := s.spaceRepo.FindByID(task.SpaceID)
	if err == nil {
		ctx := context.Background()
		client := s.minioClient.GetClient()
		bucket := s.minioClient.GetConfig().Bucket

		for i := 0; i < task.TotalChunks; i++ {
			chunkObjectName := fmt.Sprintf("spaces/%s/temp/%s/chunk_%d", space.Name, uploadID, i)
			client.RemoveObject(ctx, bucket, chunkObjectName, minio.RemoveObjectOptions{})
		}
	}

	s.uploadRepo.DeleteTask(uploadID)
	s.uploadRepo.DecrementUserUploadCount(userID)

	return nil
}

func parseChunkIndex(str string) (int, error) {
	return strconv.Atoi(str)
}
