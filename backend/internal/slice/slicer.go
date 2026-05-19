package slice

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"

	"webGL-720yun/internal/model"
	slicedto "webGL-720yun/internal/slice/dto"
	sliceservice "webGL-720yun/internal/slice/service"
	imgprocessor "webGL-720yun/pkg/image"
	"webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/panorama"
	"webGL-720yun/pkg/progress"
	"webGL-720yun/pkg/websocket"
)

type SliceMetrics struct {
	TaskID       string
	SceneID      uint
	DownloadTime time.Duration
	PreviewTime  time.Duration
	E2CTime      time.Duration
	TilesGenTime time.Duration
	UploadTime   time.Duration
	TotalTime    time.Duration
	ImageSize    int64
	TileCount    int
	Success      bool
	ErrorMessage string
	RetryCount   int
}

type SliceProcessor struct {
	db             *gorm.DB
	minioClient    *minio_client.MinIOClient
	imageProcessor *imgprocessor.Processor
	imageTiler     *imgprocessor.Tiler
	panoramaConv   *panorama.Converter
	wsHub          *websocket.Hub
	metricsChan    chan *SliceMetrics
	queue          *sliceservice.SliceQueue
}

func NewSliceProcessor(
	db *gorm.DB,
	minioClient *minio_client.MinIOClient,
	wsHub *websocket.Hub,
) *SliceProcessor {
	return &SliceProcessor{
		db:             db,
		minioClient:    minioClient,
		imageProcessor: imgprocessor.NewProcessor(),
		imageTiler:     imgprocessor.NewTiler(),
		panoramaConv:   panorama.NewConverter(panorama.DefaultOptions()),
		wsHub:          wsHub,
		metricsChan:    make(chan *SliceMetrics, 100),
	}
}

func (p *SliceProcessor) SetQueue(queue *sliceservice.SliceQueue) {
	p.queue = queue
}

func (p *SliceProcessor) StartMetricsCollector() {
	go func() {
		for metrics := range p.metricsChan {
			p.recordMetrics(metrics)
		}
	}()
}

func (p *SliceProcessor) recordMetrics(metrics *SliceMetrics) {
	log.Printf("[Metrics] TaskID: %s, SceneID: %d, Download: %v, Preview: %v, E2C: %v, Tiles: %v, Upload: %v, Total: %v, Tiles: %d, Retry: %d, Success: %v",
		metrics.TaskID,
		metrics.SceneID,
		metrics.DownloadTime.Round(time.Millisecond),
		metrics.PreviewTime.Round(time.Millisecond),
		metrics.E2CTime.Round(time.Millisecond),
		metrics.TilesGenTime.Round(time.Millisecond),
		metrics.UploadTime.Round(time.Millisecond),
		metrics.TotalTime.Round(time.Millisecond),
		metrics.TileCount,
		metrics.RetryCount,
		metrics.Success,
	)
}

func (p *SliceProcessor) ProcessWithRetry(ctx context.Context, task *sliceservice.SliceTask) error {
	maxRetries := 3
	retryDelay := time.Second * 5

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		err := p.Process(ctx, task)
		if err == nil {
			return nil
		}

		lastErr = err
		log.Printf("[Retry] Slice task failed (attempt %d/%d): task_id=%s, error=%v", i+1, maxRetries, task.TaskID, err)

		if i < maxRetries-1 {
			time.Sleep(retryDelay)
			retryDelay *= 2
		}
	}

	return fmt.Errorf("slice task failed after %d retries: %w", maxRetries, lastErr)
}

func (p *SliceProcessor) Process(ctx context.Context, task *sliceservice.SliceTask) error {
	metrics := &SliceMetrics{
		TaskID:  task.TaskID,
		SceneID: task.SceneID,
	}
	totalStart := time.Now()
	defer func() {
		metrics.TotalTime = time.Since(totalStart)
		select {
		case p.metricsChan <- metrics:
		default:
		}
	}()

	progress.PrintStart("切片", task.SceneCode, fmt.Sprintf("(TaskID: %s)", task.TaskID))

	bar := progress.GlobalManager.CreateBar(task.TaskID, 100, "切片处理中", false)
	defer progress.GlobalManager.RemoveBar(task.TaskID)

	// 更新状态为正在处理
	if err := p.db.Model(&model.ResScene{}).Where("id = ?", task.SceneID).Update("slice_status", model.SliceStatusSlicing).Error; err != nil {
		log.Printf("⚠️  更新场景状态为 Slicing 失败: %v", err)
	}

	bar.Set(5)

	// 检查是否已经存在处理好的资源，避免重复切片
	if p.minioClient.CheckSceneAssetsExist(ctx, task.SpaceSlug, task.SceneCode) {
		log.Printf("⏩ 场景 [%s] 的资源已在 MinIO 中存在，跳过切片任务", task.SceneCode)
		previewURL := model.GetScenePreviewPath(task.SpaceSlug, task.SceneCode)
		tileURL := model.GetSceneTileBasePath(task.SpaceSlug, task.SceneCode)
		if err := p.updateSceneStatus(task.SceneID, model.SliceStatusReady, tileURL, previewURL); err != nil {
			log.Printf("⚠️  更新已存在场景状态失败: %v", err)
		}
		p.notifyComplete(task, tileURL, previewURL)
		metrics.Success = true
		return nil
	}

	tempDir := filepath.Join(os.TempDir(), "slice_"+task.TaskID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tempDir)

	var sourceFile string

	if task.LocalSourcePath != "" {
		if _, err := os.Stat(task.LocalSourcePath); err == nil {
			sourceFile = task.LocalSourcePath
			log.Printf("✅ 使用本地源文件: %s", task.LocalSourcePath)
			defer func() {
				os.Remove(task.LocalSourcePath)
				os.RemoveAll(filepath.Dir(task.LocalSourcePath))
			}()
		}
	}

	if sourceFile == "" {
		p.notifyProgress(task, 5, slicedto.StageDownloading, "开始下载源文件...")

		sourceFile = filepath.Join(tempDir, "source.jpg")
		downloadStart := time.Now()
		if err := p.minioClient.DownloadSceneSource(ctx, task.SpaceSlug, task.FileID, sourceFile); err != nil {
			if minio.ToErrorResponse(err).Code == "NoSuchKey" {
				log.Printf("⚠️  检测到幽灵任务: 无法找到源文件 %s/%s，该任务可能属于旧版本，已将其跳过。", task.SpaceSlug, task.FileID)
				metrics.Success = true
				return nil
			}
			metrics.Success = false
			metrics.ErrorMessage = err.Error()
			return fmt.Errorf("下载源文件失败: %w", err)
		}
		metrics.DownloadTime = time.Since(downloadStart)

		p.notifyProgress(task, 10, slicedto.StageDownloading, "源文件下载完成")
		bar.Set(10)
	}

	fileInfo, _ := os.Stat(sourceFile)
	if fileInfo != nil {
		metrics.ImageSize = fileInfo.Size()
	}

	p.notifyProgress(task, 15, slicedto.StagePreview, "开始生成快速预览...")
	bar.Set(15)

	previewStart := time.Now()
	previewURL, err := p.generateAndUploadPreview(ctx, sourceFile, task.SceneCode, task.SpaceSlug)
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("生成快速预览失败: %w", err)
	}
	metrics.PreviewTime = time.Since(previewStart)

	if err := p.updateScenePreview(task.SceneID, previewURL); err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("更新场景预览失败: %w", err)
	}

	p.notifyProgress(task, 25, slicedto.StagePreview, "快速预览已生成，开始后台处理...")
	bar.Set(25)

	p.notifyProgress(task, 30, slicedto.StageE2C, "开始E2C转换...")
	bar.Set(30)

	cubemapDir := filepath.Join(tempDir, "cubemap")
	if err := os.MkdirAll(cubemapDir, 0755); err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("创建cubemap目录失败: %w", err)
	}

	e2cStart := time.Now()
	cubemapFiles, err := p.convertToCubemap(sourceFile, cubemapDir)
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("E2C转换失败: %w", err)
	}
	metrics.E2CTime = time.Since(e2cStart)

	p.notifyProgress(task, 45, slicedto.StageE2C, "E2C转换完成")
	bar.Set(45)

	p.notifyProgress(task, 50, slicedto.StageTiles, "开始生成瓦片...")
	bar.Set(50)

	tilesDir := filepath.Join(tempDir, "tiles")
	if err := os.MkdirAll(tilesDir, 0755); err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("创建tiles目录失败: %w", err)
	}

	tilesStart := time.Now()
	tileFiles, err := p.imageTiler.GenerateTiles(cubemapFiles, tilesDir)
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("生成瓦片失败: %w", err)
	}
	metrics.TilesGenTime = time.Since(tilesStart)

	tileCount := 0
	for _, levels := range tileFiles {
		for _, files := range levels {
			tileCount += len(files)
		}
	}
	metrics.TileCount = tileCount

	p.notifyProgress(task, 70, slicedto.StageTiles, "瓦片生成完成")
	bar.Set(70)

	p.notifyProgress(task, 75, slicedto.StageUploading, "开始上传瓦片...")
	bar.Set(75)

	uploadStart := time.Now()
	if err := p.minioClient.UploadSceneTiles(ctx, task.SpaceSlug, task.SceneCode, tileFiles); err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("上传瓦片失败: %w", err)
	}
	metrics.UploadTime = time.Since(uploadStart)

	tileURL := model.GetSceneTileBasePath(task.SpaceSlug, task.SceneCode)

	p.notifyProgress(task, 95, slicedto.StageUploading, "瓦片上传完成")
	bar.Set(95)

	fmt.Printf("[切片] 场景处理完成: %s\n", task.SceneCode)

	if err := p.updateSceneStatus(task.SceneID, model.SliceStatusReady, tileURL, previewURL); err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("更新场景状态失败: %w", err)
	}

	metrics.Success = true
	p.notifyComplete(task, tileURL, previewURL)
	progress.PrintComplete("切片", task.SceneCode)

	return nil
}

func (p *SliceProcessor) generateAndUploadPreview(ctx context.Context, sourceFile, sceneCode, spaceSlug string) (string, error) {
	tempPreviewFile := filepath.Join(os.TempDir(), fmt.Sprintf("quick_preview_%s.jpg", sceneCode))
	defer os.Remove(tempPreviewFile)

	if err := p.imageProcessor.GeneratePreview(sourceFile, tempPreviewFile, 1024, 512); err != nil {
		return "", fmt.Errorf("生成预览文件失败: %w", err)
	}

	return p.minioClient.UploadScenePreview(ctx, spaceSlug, sceneCode, tempPreviewFile)
}

func (p *SliceProcessor) updateScenePreview(sceneID uint, previewURL string) error {
	return p.db.Model(&model.ResScene{}).Where("id = ?", sceneID).Updates(map[string]interface{}{
		"preview_url":  previewURL,
		"slice_status": model.SliceStatusSlicing,
		"updated_at":   time.Now(),
	}).Error
}

func (p *SliceProcessor) convertToCubemap(sourceFile string, outputDir string) (map[string]string, error) {
	faceNames := panorama.FaceNameStrings()

	opts := &panorama.Options{
		Size:      0,
		Layout:    panorama.LayoutNone,
		FaceNames: faceNames,
	}

	converter := panorama.NewConverter(opts)

	outputBase := filepath.Join(outputDir, "face")
	if err := converter.Convert(sourceFile, outputBase); err != nil {
		return nil, err
	}

	cubemapFiles := make(map[string]string)
	for _, name := range faceNames {
		cubemapFiles[name] = fmt.Sprintf("%s_%s.jpg", outputBase, name)
	}

	return cubemapFiles, nil
}

func (p *SliceProcessor) updateSceneStatus(sceneID uint, status string, tileURL string, previewURL string) error {
	return p.db.Model(&model.ResScene{}).Where("id = ?", sceneID).Updates(map[string]interface{}{
		"slice_status": status,
		"tile_url":     tileURL,
		"preview_url":  previewURL,
		"is_converted": true,
		"updated_at":   time.Now(),
	}).Error
}

func (p *SliceProcessor) notifyProgress(task *sliceservice.SliceTask, progress int, stage string, message string) {
	if p.wsHub == nil {
		return
	}

	var data *slicedto.SliceProgressData
	if p.queue != nil {
		queueStats, err := p.queue.GetQueueStats(task.UserID)
		if err == nil {
			data = slicedto.NewSliceProgressDataWithQueue(
				task.TaskID, task.SceneID, task.SceneCode,
				progress, stage, message,
				queueStats.QueueAheadCount,
				queueStats.EstimatedWaitSec,
				queueStats.ActiveUsers,
			)
		}
	}

	if data == nil {
		data = slicedto.NewSliceProgressData(task.TaskID, task.SceneID, task.SceneCode, progress, stage, message)
	}

	p.wsHub.Broadcast(&websocket.Message{
		Type:   websocket.MessageTypeSliceProgress,
		UserID: task.UserID,
		Data:   data,
	})
}

func (p *SliceProcessor) notifyComplete(task *sliceservice.SliceTask, tileURL string, previewURL string) {
	if p.wsHub == nil {
		log.Printf("[切片] wsHub is nil, cannot send complete notification")
		return
	}

	data := slicedto.NewSliceCompleteData(task.TaskID, task.SceneID, task.SceneCode, tileURL, previewURL)
	log.Printf("[切片] Sending complete notification: task_id=%s, scene_code=%s, user_id=%d", task.TaskID, task.SceneCode, task.UserID)

	p.wsHub.Broadcast(&websocket.Message{
		Type:   websocket.MessageTypeSliceComplete,
		UserID: task.UserID,
		Data:   data,
	})
}

func (p *SliceProcessor) notifyError(task *sliceservice.SliceTask, errorMsg string) {
	if p.wsHub == nil {
		return
	}

	data := slicedto.NewSliceErrorData(task.TaskID, task.SceneID, task.SceneCode, errorMsg)

	p.wsHub.Broadcast(&websocket.Message{
		Type:   websocket.MessageTypeSliceError,
		UserID: task.UserID,
		Data:   data,
	})
}
