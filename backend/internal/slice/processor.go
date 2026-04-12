package slice

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strconv"
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

type SliceProcessor struct {
	db             *gorm.DB
	minioClient    *minio_client.MinIOClient
	imageProcessor *imgprocessor.Processor
	panoramaConv   *panorama.Converter
	wsHub          *websocket.Hub
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
	}
}

func (p *SliceProcessor) Process(ctx context.Context, task *SliceTask) error {
	tempDir := filepath.Join(os.TempDir(), "slice_"+task.TaskID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tempDir)

	p.notifyProgress(task, 5, StageDownloading, "开始下载源文件...")

	sourceFile := filepath.Join(tempDir, "source.jpg")
	if err := p.downloadSourceFile(ctx, task.FileID, sourceFile); err != nil {
		return fmt.Errorf("下载源文件失败: %w", err)
	}

	p.notifyProgress(task, 10, StageDownloading, "源文件下载完成")

	p.notifyProgress(task, 15, StageE2C, "开始E2C转换...")

	cubemapDir := filepath.Join(tempDir, "cubemap")
	if err := os.MkdirAll(cubemapDir, 0755); err != nil {
		return fmt.Errorf("创建cubemap目录失败: %w", err)
	}

	cubemapFiles, err := p.convertToCubemap(sourceFile, cubemapDir)
	if err != nil {
		return fmt.Errorf("E2C转换失败: %w", err)
	}

	p.notifyProgress(task, 30, StageE2C, "E2C转换完成")

	p.notifyProgress(task, 35, StageTiles, "开始生成瓦片...")

	tilesDir := filepath.Join(tempDir, "tiles")
	if err := os.MkdirAll(tilesDir, 0755); err != nil {
		return fmt.Errorf("创建tiles目录失败: %w", err)
	}

	tileFiles, err := p.generateTiles(cubemapFiles, tilesDir)
	if err != nil {
		return fmt.Errorf("生成瓦片失败: %w", err)
	}

	p.notifyProgress(task, 60, StageTiles, "瓦片生成完成")

	p.notifyProgress(task, 65, StageUploading, "开始上传瓦片...")

	tileURL, err := p.uploadTiles(ctx, tileFiles, cubemapFiles, task.SceneCode)
	if err != nil {
		return fmt.Errorf("上传瓦片失败: %w", err)
	}

	p.notifyProgress(task, 90, StageUploading, "瓦片上传完成")

	previewFile := filepath.Join(tempDir, "preview.jpg")
	previewURL, err := p.generateAndUploadPreview(ctx, cubemapFiles, previewFile, task.SceneCode)
	if err != nil {
		return fmt.Errorf("生成预览图失败: %w", err)
	}

	if err := p.updateSceneStatus(task.SceneID, model.SliceStatusReady, tileURL, previewURL); err != nil {
		return fmt.Errorf("更新场景状态失败: %w", err)
	}

	p.notifyComplete(task, tileURL, previewURL)

	return nil
}

func (p *SliceProcessor) downloadSourceFile(ctx context.Context, fileID string, destPath string) error {
	client := p.minioClient.GetClient()
	bucket := p.minioClient.GetConfig().Bucket

	objectName := fmt.Sprintf("sources/%s/source.jpg", fileID)

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
	tileFiles := make(map[string]map[int][]string)

	for faceName, faceFile := range cubemapFiles {
		img, err := imaging.Open(faceFile)
		if err != nil {
			return nil, fmt.Errorf("打开面 %s 失败: %w", faceName, err)
		}

		faceDir := filepath.Join(outputDir, faceName)
		if err := os.MkdirAll(faceDir, 0755); err != nil {
			return nil, err
		}

		tileFiles[faceName] = make(map[int][]string)

		level := 0
		currentImg := img

		for currentImg.Bounds().Dx() >= TileSize {
			levelDir := filepath.Join(faceDir, fmt.Sprintf("level_%d", level))
			if err := os.MkdirAll(levelDir, 0755); err != nil {
				return nil, err
			}

			tiles, err := p.sliceImage(currentImg, levelDir, level, faceName)
			if err != nil {
				return nil, err
			}

			tileFiles[faceName][level] = tiles

			newSize := currentImg.Bounds().Dx() / 2
			if newSize < TileSize {
				break
			}

			currentImg = imaging.Resize(currentImg, newSize, newSize, imaging.Lanczos)
			level++
		}
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

func (p *SliceProcessor) uploadTiles(ctx context.Context, tileFiles map[string]map[int][]string, cubemapFiles map[string]string, sceneCode string) (string, error) {
	client := p.minioClient.GetClient()
	bucket := p.minioClient.GetConfig().Bucket

	for faceName, levels := range tileFiles {
		for level, files := range levels {
			for _, file := range files {
				fileName := filepath.Base(file)
				objectName := fmt.Sprintf("tiles/%s/lod/%s/level_%d/%s", sceneCode, faceName, level, fileName)

				if err := p.uploadFile(ctx, client, bucket, objectName, file); err != nil {
					return "", err
				}
			}
		}
	}

	for faceName, file := range cubemapFiles {
		objectName := fmt.Sprintf("tiles/%s/cubemap/%s.jpg", sceneCode, faceName)
		if err := p.uploadFile(ctx, client, bucket, objectName, file); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("tiles/%s/", sceneCode), nil
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

func (p *SliceProcessor) generateAndUploadPreview(ctx context.Context, cubemapFiles map[string]string, previewFile string, sceneCode string) (string, error) {
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
	objectName := fmt.Sprintf("tiles/%s/preview.jpg", sceneCode)

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

func GetTilePath(sceneCode string, face string, level int, row int, col int) string {
	return fmt.Sprintf("tiles/%s/lod/%s/level_%d/tile_%d_%d.jpg", sceneCode, face, level, row, col)
}

func ParseTilePath(path string) (face string, level int, row int, col int, err error) {
	_, err = fmt.Sscanf(path, "tiles/%s/lod/%s/level_%d/tile_%d_%d.jpg",
		new(string), &face, &level, &row, &col)
	return
}

func GetCubemapPath(sceneCode string, face string) string {
	return fmt.Sprintf("tiles/%s/cubemap/%s.jpg", sceneCode, face)
}

func GetPreviewPath(sceneCode string) string {
	return fmt.Sprintf("tiles/%s/preview.jpg", sceneCode)
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
