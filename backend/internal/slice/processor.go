package slice

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"

	"webGL-720yun/internal/model"
	imgprocessor "webGL-720yun/pkg/image"
	"webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/panorama"
	"webGL-720yun/pkg/websocket"
)

const (
	TileSize = 256
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
	panoramaConv   *panorama.Converter
	wsHub          *websocket.Hub
	metricsChan    chan *SliceMetrics
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
		panoramaConv:   panorama.NewConverter(panorama.DefaultOptions()),
		wsHub:          wsHub,
		metricsChan:    make(chan *SliceMetrics, 100),
	}
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

func (p *SliceProcessor) ProcessWithRetry(ctx context.Context, task *SliceTask) error {
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

func (p *SliceProcessor) Process(ctx context.Context, task *SliceTask) error {
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

	// 检查是否已经存在处理好的资源，避免重复切片
	if p.checkAssetsExist(ctx, task.SpaceSlug, task.SceneCode) {
		log.Printf("⏩ 场景 [%s] 的资源已在 MinIO 中存在，跳过切片任务", task.SceneCode)
		previewURL := GetPreviewPath(task.SpaceSlug, task.SceneCode)
		tileURL := fmt.Sprintf("spaces/%s/tiles/%s/", task.SpaceSlug, task.SceneCode)
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

	p.notifyProgress(task, 5, StageDownloading, "开始下载源文件...")

	downloadStart := time.Now()
	sourceFile := filepath.Join(tempDir, "source.jpg")
	if err := p.downloadSourceFile(ctx, task.FileID, task.SpaceSlug, sourceFile); err != nil {
		// 检查是否是源文件不存在（NoSuchKey），这通常意味着是一个旧的残留任务（Ghost Task）
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			log.Printf("⚠️  检测到幽灵任务: 无法找到源文件 %s/%s，该任务可能属于旧版本，已将其跳过。", task.SpaceSlug, task.FileID)
			metrics.Success = true // 标记为处理完成（忽略）
			return nil
		}
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("下载源文件失败: %w", err)
	}
	metrics.DownloadTime = time.Since(downloadStart)

	fileInfo, _ := os.Stat(sourceFile)
	if fileInfo != nil {
		metrics.ImageSize = fileInfo.Size()
	}

	p.notifyProgress(task, 10, StageDownloading, "源文件下载完成")

	p.notifyProgress(task, 15, StagePreview, "开始生成快速预览...")

	previewStart := time.Now()
	previewURL, err := p.generateQuickPreview(ctx, sourceFile, task.SceneCode, task.SpaceSlug)
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

	p.notifyProgress(task, 25, StagePreview, "快速预览已生成，开始后台处理...")

	p.notifyProgress(task, 30, StageE2C, "开始E2C转换...")

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

	p.notifyProgress(task, 45, StageE2C, "E2C转换完成")

	p.notifyProgress(task, 50, StageTiles, "开始生成瓦片...")

	tilesDir := filepath.Join(tempDir, "tiles")
	if err := os.MkdirAll(tilesDir, 0755); err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("创建tiles目录失败: %w", err)
	}

	tilesStart := time.Now()
	tileFiles, err := p.generateTiles(cubemapFiles, tilesDir)
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

	p.notifyProgress(task, 70, StageTiles, "瓦片生成完成")

	p.notifyProgress(task, 75, StageUploading, "开始上传瓦片...")

	uploadStart := time.Now()
	tileURL, err := p.uploadTiles(ctx, tileFiles, cubemapFiles, task.SceneCode, task.SpaceSlug)
	if err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("上传瓦片失败: %w", err)
	}
	metrics.UploadTime = time.Since(uploadStart)

	p.notifyProgress(task, 95, StageUploading, "瓦片上传完成")

	if err := p.updateSceneStatus(task.SceneID, model.SliceStatusReady, tileURL, previewURL); err != nil {
		metrics.Success = false
		metrics.ErrorMessage = err.Error()
		return fmt.Errorf("更新场景状态失败: %w", err)
	}

	metrics.Success = true
	p.notifyComplete(task, tileURL, previewURL)

	return nil
}

func (p *SliceProcessor) generateQuickPreview(ctx context.Context, sourceFile string, sceneCode string, spaceSlug string) (string, error) {
	img, err := imaging.Open(sourceFile)
	if err != nil {
		return "", fmt.Errorf("打开源文件失败: %w", err)
	}

	preview := imaging.Resize(img, 1024, 512, imaging.Lanczos)

	tempPreviewFile := filepath.Join(os.TempDir(), fmt.Sprintf("quick_preview_%s.jpg", sceneCode))
	defer os.Remove(tempPreviewFile)

	if err := imaging.Save(preview, tempPreviewFile, imaging.JPEGQuality(85)); err != nil {
		return "", fmt.Errorf("保存预览文件失败: %w", err)
	}

	client := p.minioClient.GetClient()
	bucket := p.minioClient.GetConfig().Bucket
	objectName := fmt.Sprintf("spaces/%s/previews/%s/preview.jpg", spaceSlug, sceneCode)

	file, err := os.Open(tempPreviewFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	_, err = client.PutObject(ctx, bucket, objectName, file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (p *SliceProcessor) updateScenePreview(sceneID uint, previewURL string) error {
	return p.db.Model(&model.ResScene{}).Where("id = ?", sceneID).Updates(map[string]interface{}{
		"preview_url":  previewURL,
		"slice_status": model.SliceStatusSlicing,
		"updated_at":   time.Now(),
	}).Error
}

func (p *SliceProcessor) downloadSourceFile(ctx context.Context, fileID string, spaceSlug string, destPath string) error {
	client := p.minioClient.GetClient()
	bucket := p.minioClient.GetConfig().Bucket

	objectName := fmt.Sprintf("spaces/%s/sources/%s/source.jpg", spaceSlug, fileID)

	obj, err := client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer obj.Close()

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.ReadFrom(obj)
	return err
}

func (p *SliceProcessor) convertToCubemap(sourceFile string, outputDir string) (map[string]string, error) {
	faceNames := []string{"px", "nx", "py", "ny", "pz", "nz"}

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

func (p *SliceProcessor) generateTiles(cubemapFiles map[string]string, outputDir string) (map[string]map[int][]string, error) {
	type faceTilesResult struct {
		faceName string
		tiles    map[int][]string
		err      error
	}

	resultChan := make(chan faceTilesResult, len(cubemapFiles))
	var wg sync.WaitGroup

	for faceName, faceFile := range cubemapFiles {
		wg.Add(1)
		go func(name, file string) {
			defer wg.Done()

			img, err := imaging.Open(file)
			if err != nil {
				resultChan <- faceTilesResult{
					faceName: name,
					err:      fmt.Errorf("打开面 %s 失败: %w", name, err),
				}
				return
			}

			faceDir := filepath.Join(outputDir, name)
			if err := os.MkdirAll(faceDir, 0755); err != nil {
				resultChan <- faceTilesResult{
					faceName: name,
					err:      err,
				}
				return
			}

			tiles := make(map[int][]string)

			level := 0
			currentImg := img

			for currentImg.Bounds().Dx() >= TileSize {
				levelDir := filepath.Join(faceDir, fmt.Sprintf("level_%d", level))
				if err := os.MkdirAll(levelDir, 0755); err != nil {
					resultChan <- faceTilesResult{
						faceName: name,
						err:      err,
					}
					return
				}

				tileFiles, err := p.sliceImage(currentImg, levelDir, level, name)
				if err != nil {
					resultChan <- faceTilesResult{
						faceName: name,
						err:      err,
					}
					return
				}

				tiles[level] = tileFiles

				newSize := currentImg.Bounds().Dx() / 2
				if newSize < TileSize {
					break
				}

				currentImg = imaging.Resize(currentImg, newSize, newSize, imaging.Lanczos)
				level++
			}

			resultChan <- faceTilesResult{
				faceName: name,
				tiles:    tiles,
			}
		}(faceName, faceFile)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	tileFiles := make(map[string]map[int][]string)

	for result := range resultChan {
		if result.err != nil {
			return nil, result.err
		}
		tileFiles[result.faceName] = result.tiles
	}

	return tileFiles, nil
}

func (p *SliceProcessor) sliceImage(img image.Image, outputDir string, level int, faceName string) ([]string, error) {
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	tilesX := (imgWidth + TileSize - 1) / TileSize
	tilesY := (imgHeight + TileSize - 1) / TileSize

	var tileFiles []string

	for y := 0; y < tilesY; y++ {
		for x := 0; x < tilesX; x++ {
			startX := x * TileSize
			startY := y * TileSize

			endX := startX + TileSize
			if endX > imgWidth {
				endX = imgWidth
			}

			endY := startY + TileSize
			if endY > imgHeight {
				endY = imgHeight
			}

			tileWidth := endX - startX
			tileHeight := endY - startY

			tileImg := imaging.New(TileSize, TileSize, image.Black)

			for ty := 0; ty < tileHeight; ty++ {
				for tx := 0; tx < tileWidth; tx++ {
					tileImg.Set(tx, ty, img.At(startX+tx, startY+ty))
				}
			}

			tileFileName := fmt.Sprintf("tile_%d_%d.jpg", y, x)
			tilePath := filepath.Join(outputDir, tileFileName)

			file, err := os.Create(tilePath)
			if err != nil {
				return nil, err
			}

			if err := jpeg.Encode(file, tileImg, &jpeg.Options{Quality: 85}); err != nil {
				file.Close()
				return nil, err
			}
			file.Close()

			tileFiles = append(tileFiles, tilePath)
		}
	}

	return tileFiles, nil
}

func (p *SliceProcessor) uploadTiles(ctx context.Context, tileFiles map[string]map[int][]string, cubemapFiles map[string]string, sceneCode string, spaceSlug string) (string, error) {
	client := p.minioClient.GetClient()
	bucket := p.minioClient.GetConfig().Bucket

	type uploadTask struct {
		objectName string
		filePath   string
	}

	var uploadTasks []uploadTask

	for faceName, levels := range tileFiles {
		for level, files := range levels {
			for _, file := range files {
				fileName := filepath.Base(file)
				objectName := fmt.Sprintf("spaces/%s/tiles/%s/cubemap/%s/level_%d/%s", spaceSlug, sceneCode, faceName, level, fileName)
				uploadTasks = append(uploadTasks, uploadTask{
					objectName: objectName,
					filePath:   file,
				})
			}
		}
	}

	for faceName, file := range cubemapFiles {
		objectName := fmt.Sprintf("spaces/%s/tiles/%s/cubemap/%s.jpg", spaceSlug, sceneCode, faceName)
		uploadTasks = append(uploadTasks, uploadTask{
			objectName: objectName,
			filePath:   file,
		})
	}

	maxConcurrent := 10
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	errChan := make(chan error, len(uploadTasks))

	for _, task := range uploadTasks {
		wg.Add(1)
		go func(t uploadTask) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if err := p.uploadFile(ctx, client, bucket, t.objectName, t.filePath); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}(task)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("spaces/%s/tiles/%s/", spaceSlug, sceneCode), nil
}

func (p *SliceProcessor) uploadFile(ctx context.Context, client *minio.Client, bucket string, objectName string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = client.PutObject(ctx, bucket, objectName, file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	return err
}

func (p *SliceProcessor) generateAndUploadPreview(ctx context.Context, cubemapFiles map[string]string, previewFile string, sceneCode string, spaceSlug string) (string, error) {
	pxFile, ok := cubemapFiles["px"]
	if !ok {
		return "", fmt.Errorf("找不到px面文件")
	}

	img, err := imaging.Open(pxFile)
	if err != nil {
		return "", err
	}

	previewImg := imaging.Resize(img, 512, 512, imaging.Lanczos)

	if err := imaging.Save(previewImg, previewFile); err != nil {
		return "", err
	}

	client := p.minioClient.GetClient()
	bucket := p.minioClient.GetConfig().Bucket
	objectName := fmt.Sprintf("spaces/%s/previews/%s/preview.jpg", spaceSlug, sceneCode)

	if err := p.uploadFile(ctx, client, bucket, objectName, previewFile); err != nil {
		return "", err
	}

	return objectName, nil
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

func (p *SliceProcessor) checkAssetsExist(ctx context.Context, spaceSlug, sceneCode string) bool {
	// 检查预览图是否存在
	previewPath := GetPreviewPath(spaceSlug, sceneCode)
	exists, err := p.minioClient.ObjectExists(previewPath)
	if err != nil || !exists {
		return false
	}

	// 检查瓦片目录是否包含基本文件（检查 px 面的 level_0）
	tilePath := fmt.Sprintf("spaces/%s/tiles/%s/cubemap/px/level_0/tile_0_0.jpg", spaceSlug, sceneCode)
	exists, err = p.minioClient.ObjectExists(tilePath)
	return err == nil && exists
}

func (p *SliceProcessor) notifyProgress(task *SliceTask, progress int, stage string, message string) {
	if p.wsHub == nil {
		return
	}

	data := NewSliceProgressData(task.TaskID, task.SceneID, progress, stage, message)

	p.wsHub.Broadcast(&websocket.Message{
		Type:   websocket.MessageTypeSliceProgress,
		UserID: task.UserID,
		Data:   data,
	})
}

func (p *SliceProcessor) notifyComplete(task *SliceTask, tileURL string, previewURL string) {
	if p.wsHub == nil {
		return
	}

	data := NewSliceCompleteData(task.TaskID, task.SceneID, tileURL, previewURL)

	p.wsHub.Broadcast(&websocket.Message{
		Type:   websocket.MessageTypeSliceComplete,
		UserID: task.UserID,
		Data:   data,
	})
}

func (p *SliceProcessor) notifyError(task *SliceTask, errorMsg string) {
	if p.wsHub == nil {
		return
	}

	data := NewSliceErrorData(task.TaskID, task.SceneID, errorMsg)

	p.wsHub.Broadcast(&websocket.Message{
		Type:   websocket.MessageTypeSliceError,
		UserID: task.UserID,
		Data:   data,
	})
}

func GetTilePath(spaceSlug, sceneCode string, face string, level int, row int, col int) string {
	return fmt.Sprintf("spaces/%s/tiles/%s/cubemap/%s/level_%d/tile_%d_%d.jpg", spaceSlug, sceneCode, face, level, row, col)
}

func ParseTilePath(path string) (spaceSlug, sceneCode, face string, level, row, col int, err error) {
	_, err = fmt.Sscanf(path, "spaces/%s/tiles/%s/cubemap/%s/level_%d/tile_%d_%d.jpg",
		&spaceSlug, &sceneCode, &face, &level, &row, &col)
	return
}

func GetCubemapPath(spaceSlug, sceneCode string, face string) string {
	return fmt.Sprintf("spaces/%s/tiles/%s/cubemap/%s.jpg", spaceSlug, sceneCode, face)
}

func GetPreviewPath(spaceSlug, sceneCode string) string {
	return fmt.Sprintf("spaces/%s/previews/%s/preview.jpg", spaceSlug, sceneCode)
}

func GetLevelDimension(level int) int {
	return 512 << level
}

func GetTileCount(level int) int {
	dimension := GetLevelDimension(level)
	return (dimension + TileSize - 1) / TileSize
}

func TileXYToLonLat(face string, tileX, tileY, level int) (lon, lat float64) {
	tileCount := GetTileCount(level)

	normX := float64(tileX) / float64(tileCount)
	normY := float64(tileY) / float64(tileCount)

	lon = (normX - 0.5) * 360
	lat = (0.5 - normY) * 180

	return lon, lat
}

func LonLatToTileXY(face string, lon, lat float64, level int) (tileX, tileY int) {
	tileCount := GetTileCount(level)

	normX := (lon / 360) + 0.5
	normY := 0.5 - (lat / 180)

	tileX = int(normX * float64(tileCount))
	tileY = int(normY * float64(tileCount))

	tileX = tileX % tileCount
	if tileX < 0 {
		tileX += tileCount
	}

	tileY = tileY % tileCount
	if tileY < 0 {
		tileY += tileCount
	}

	return tileX, tileY
}

func init() {
	_ = strconv.Itoa
}
