package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

const DefaultTileSize = 256

const DefaultJPEGQuality = 85

// LODSpec 定义单个 LOD 层级规格
type LODSpec struct {
	TileCount  int // 每面瓦片数量 (1, 4, 8)
	MaxSize    int // 单张瓦片最大尺寸
	TargetSize int // 该层级整面的目标像素尺寸
}

// DefaultLODSepcs 三级 LOD 规格：Level 0 (480p), Level 1 (1080p), Level 2 (4K)
// 每一级都会重采样到固定面尺寸，源图更大就缩小，源图更小就 Lanczos 放大，
// 从而保证前端拿到的三档观感稳定。
var DefaultLODSepcs = map[int]LODSpec{
	0: {TileCount: 1, MaxSize: 512, TargetSize: 512},  // 480p 底图：1×1，512px
	1: {TileCount: 4, MaxSize: 512, TargetSize: 2048}, // 1080p：4×4，每张 512px
	2: {TileCount: 8, MaxSize: 512, TargetSize: 4096}, // 4K：8×8，每张 512px
}

// Tiler 负责将图片切割成瓦片矩阵
type Tiler struct {
	lodSpecs    map[int]LODSpec
	tileSize    int
	jpegQuality int
}

// TilerOption 配置选项函数
type TilerOption func(*Tiler)

func WithLODSpecs(specs map[int]LODSpec) TilerOption {
	return func(t *Tiler) {
		t.lodSpecs = specs
	}
}

func WithTileSize(size int) TilerOption {
	return func(t *Tiler) {
		t.tileSize = size
	}
}

func WithJPEGQuality(quality int) TilerOption {
	return func(t *Tiler) {
		t.jpegQuality = quality
	}
}

func NewTiler(opts ...TilerOption) *Tiler {
	t := &Tiler{
		lodSpecs:    DefaultLODSepcs,
		tileSize:    DefaultTileSize,
		jpegQuality: DefaultJPEGQuality,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// TileResult 存储单个面的瓦片结果
type TileResult struct {
	FaceName string
	Tiles    map[int][]string // level -> file paths
	Err      error
}

// GenerateTiles 对立方体六个面的图片进行 LOD 瓦片切割
// faceFiles: map[faceName]filePath, 如 {"px": "/tmp/face_px.jpg", ...}
// outputDir: 瓦片输出根目录
// 返回: map[faceName]map[level][]filePath
func (t *Tiler) GenerateTiles(faceFiles map[string]string, outputDir string) (map[string]map[int][]string, error) {
	results := make(chan TileResult, len(faceFiles))

	for faceName, filePath := range faceFiles {
		go func(name, path string) {
			result := TileResult{FaceName: name}
			tiles, err := t.ProcessFace(name, path, outputDir)
			result.Tiles = tiles
			result.Err = err
			results <- result
		}(faceName, filePath)
	}

	tileFiles := make(map[string]map[int][]string)
	for i := 0; i < len(faceFiles); i++ {
		result := <-results
		if result.Err != nil {
			return nil, fmt.Errorf("处理面 %s 失败: %w", result.FaceName, result.Err)
		}
		tileFiles[result.FaceName] = result.Tiles
	}

	return tileFiles, nil
}

// ProcessFace 处理单个面的瓦片切割，三级LOD并行生成
func (t *Tiler) ProcessFace(faceName, filePath, outputDir string) (map[int][]string, error) {
	img, err := imaging.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开面 %s 失败: %w", faceName, err)
	}

	faceDir := filepath.Join(outputDir, faceName)
	if err := os.MkdirAll(faceDir, 0755); err != nil {
		return nil, err
	}

	type levelResult struct {
		level int
		files []string
		err   error
	}

	results := make(chan levelResult, 3)

	for level := 0; level <= 2; level++ {
		level := level
		spec := t.lodSpecs[level]

		go func() {
			targetSize := spec.TargetSize
			if targetSize > spec.MaxSize*spec.TileCount {
				targetSize = spec.MaxSize * spec.TileCount
			}

			// 每一级都重采样到目标尺寸：
			// 大图 Lanczos 缩小（8K → 2048 面），小图 Lanczos 放大（8K 源图的 level2 → 4096 面）。
			// 这样 level 0/1/2 的输出面固定为 512/2048/4096，前端三档 LOD 才能稳定切换。
			levelImg := imaging.Resize(img, targetSize, targetSize, imaging.Lanczos)

			levelDir := filepath.Join(faceDir, fmt.Sprintf("level_%d", level))
			if err := os.MkdirAll(levelDir, 0755); err != nil {
				results <- levelResult{level: level, err: err}
				return
			}

			tileFiles, err := t.sliceImageByLOD(levelImg, levelDir, spec.TileCount)
			results <- levelResult{level: level, files: tileFiles, err: err}
		}()
	}

	tiles := make(map[int][]string)
	for i := 0; i < 3; i++ {
		r := <-results
		if r.err != nil {
			return nil, r.err
		}
		tiles[r.level] = r.files
	}

	return tiles, nil
}

// sliceImageByLOD 按 tileCount×tileCount 网格切割图像
func (t *Tiler) sliceImageByLOD(img image.Image, outputDir string, tileCount int) ([]string, error) {
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	tileWidth := imgWidth / tileCount
	tileHeight := imgHeight / tileCount

	var tileFiles []string

	for y := 0; y < tileCount; y++ {
		for x := 0; x < tileCount; x++ {
			startX := x * tileWidth
			startY := y * tileHeight

			// 使用 imaging.Crop 替代逐像素复制，性能提升显著
			rect := image.Rect(startX, startY, startX+tileWidth, startY+tileHeight)
			tileImg := imaging.Crop(img, rect)

			tileFileName := fmt.Sprintf("tile_%d_%d.jpg", x, y)
			tilePath := filepath.Join(outputDir, tileFileName)

			file, err := os.Create(tilePath)
			if err != nil {
				return nil, err
			}

			if err := jpeg.Encode(file, tileImg, &jpeg.Options{Quality: t.jpegQuality}); err != nil {
				file.Close()
				return nil, err
			}
			file.Close()

			tileFiles = append(tileFiles, tilePath)
		}
	}

	return tileFiles, nil
}

// SliceImage 按固定 TileSize 切割图像（非 LOD 模式）
func (t *Tiler) SliceImage(img image.Image, outputDir string) ([]string, error) {
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	tilesX := (imgWidth + t.tileSize - 1) / t.tileSize
	tilesY := (imgHeight + t.tileSize - 1) / t.tileSize

	var tileFiles []string

	for y := 0; y < tilesY; y++ {
		for x := 0; x < tilesX; x++ {
			startX := x * t.tileSize
			startY := y * t.tileSize

			endX := startX + t.tileSize
			if endX > imgWidth {
				endX = imgWidth
			}

			endY := startY + t.tileSize
			if endY > imgHeight {
				endY = imgHeight
			}

			// 使用 imaging.Crop 替代逐像素复制
			rect := image.Rect(startX, startY, endX, endY)
			tileImg := imaging.Crop(img, rect)

			tileFileName := fmt.Sprintf("tile_%d_%d.jpg", x, y)
			tilePath := filepath.Join(outputDir, tileFileName)

			file, err := os.Create(tilePath)
			if err != nil {
				return nil, err
			}

			if err := jpeg.Encode(file, tileImg, &jpeg.Options{Quality: t.jpegQuality}); err != nil {
				file.Close()
				return nil, err
			}
			file.Close()

			tileFiles = append(tileFiles, tilePath)
		}
	}

	return tileFiles, nil
}

// GetLevelDimension 获取指定层级的像素尺寸
func GetLevelDimension(level int) int {
	return 512 << level
}

// GetTileCount 获取指定层级的瓦片数量
func GetTileCount(level int) int {
	dimension := GetLevelDimension(level)
	return (dimension + DefaultTileSize - 1) / DefaultTileSize
}
