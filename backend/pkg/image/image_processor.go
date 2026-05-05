package image

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

var (
	ErrInvalidFormat     = errors.New("无效的图片格式")
	ErrInvalidResolution = errors.New("图片分辨率不符合要求，宽高比应为 2:1")
	ErrFileTooLarge      = errors.New("文件大小超过限制")
)

type ImageInfo struct {
	Width       int
	Height      int
	FileSize    int64
	Format      string
	MD5         string
	AspectRatio float64
}

type Processor struct {
	maxFileSize    int64
	thumbnailWidth int
	thumbnailHeight int
}

func NewProcessor() *Processor {
	return &Processor{
		maxFileSize:     500 * 1024 * 1024,
		thumbnailWidth:  256,
		thumbnailHeight: 128,
	}
}

func (p *Processor) SetMaxFileSize(size int64) {
	p.maxFileSize = size
}

func (p *Processor) ValidateFormat(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 8)
	n, err := file.Read(buffer)
	if err != nil || n < 8 {
		return "", ErrInvalidFormat
	}

	if buffer[0] == 0xFF && buffer[1] == 0xD8 && buffer[2] == 0xFF {
		return "jpeg", nil
	}

	if buffer[0] == 0x89 && buffer[1] == 0x50 && buffer[2] == 0x4E && buffer[3] == 0x47 {
		return "png", nil
	}

	return "", ErrInvalidFormat
}

func (p *Processor) ValidateFormatByReader(reader io.Reader) (string, error) {
	buffer := make([]byte, 8)
	n, err := reader.Read(buffer)
	if err != nil || n < 8 {
		return "", ErrInvalidFormat
	}

	if buffer[0] == 0xFF && buffer[1] == 0xD8 && buffer[2] == 0xFF {
		return "jpeg", nil
	}

	if buffer[0] == 0x89 && buffer[1] == 0x50 && buffer[2] == 0x4E && buffer[3] == 0x47 {
		return "png", nil
	}

	return "", ErrInvalidFormat
}

func (p *Processor) GetImageInfo(filePath string) (*ImageInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	if fileInfo.Size() > p.maxFileSize {
		return nil, ErrFileTooLarge
	}

	img, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("解析图片失败: %w", err)
	}

	md5Hash, err := p.CalculateMD5(filePath)
	if err != nil {
		return nil, fmt.Errorf("计算MD5失败: %w", err)
	}

	aspectRatio := float64(img.Width) / float64(img.Height)

	return &ImageInfo{
		Width:       img.Width,
		Height:      img.Height,
		FileSize:    fileInfo.Size(),
		Format:      format,
		MD5:         md5Hash,
		AspectRatio: aspectRatio,
	}, nil
}

func (p *Processor) ValidatePanoramaResolution(filePath string) error {
	info, err := p.GetImageInfo(filePath)
	if err != nil {
		return err
	}

	aspectRatio := float64(info.Width) / float64(info.Height)
	if aspectRatio < 1.99 || aspectRatio > 2.01 {
		return fmt.Errorf("%w (当前比例: %.2f:1)", ErrInvalidResolution, aspectRatio)
	}

	return nil
}

func (p *Processor) GenerateThumbnail(srcPath, dstPath string) error {
	img, err := imaging.Open(srcPath)
	if err != nil {
		return fmt.Errorf("打开源图片失败: %w", err)
	}

	thumbnail := imaging.Resize(img, p.thumbnailWidth, p.thumbnailHeight, imaging.Lanczos)

	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(dstPath))
	switch ext {
	case ".jpg", ".jpeg":
		err = imaging.Save(thumbnail, dstPath, imaging.JPEGQuality(85))
	case ".png":
		err = imaging.Save(thumbnail, dstPath)
	default:
		err = imaging.Save(thumbnail, dstPath, imaging.JPEGQuality(85))
	}

	if err != nil {
		return fmt.Errorf("保存缩略图失败: %w", err)
	}

	return nil
}

func (p *Processor) GenerateThumbnailFromReader(reader io.Reader, dstPath string) error {
	img, _, err := image.Decode(reader)
	if err != nil {
		return fmt.Errorf("解码图片失败: %w", err)
	}

	thumbnail := imaging.Resize(img, p.thumbnailWidth, p.thumbnailHeight, imaging.Lanczos)

	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dstFile.Close()

	ext := strings.ToLower(filepath.Ext(dstPath))
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(dstFile, thumbnail, &jpeg.Options{Quality: 85})
	case ".png":
		err = png.Encode(dstFile, thumbnail)
	default:
		err = jpeg.Encode(dstFile, thumbnail, &jpeg.Options{Quality: 85})
	}

	if err != nil {
		return fmt.Errorf("保存缩略图失败: %w", err)
	}

	return nil
}

func (p *Processor) CalculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("计算MD5失败: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (p *Processor) CalculateMD5FromReader(reader io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", fmt.Errorf("计算MD5失败: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (p *Processor) ValidateAndProcess(filePath string) (*ImageInfo, error) {
	format, err := p.ValidateFormat(filePath)
	if err != nil {
		return nil, err
	}

	if err := p.ValidatePanoramaResolution(filePath); err != nil {
		return nil, err
	}

	info, err := p.GetImageInfo(filePath)
	if err != nil {
		return nil, err
	}

	info.Format = format

	return info, nil
}

func (p *Processor) ResizeImage(srcPath, dstPath string, width, height int) error {
	img, err := imaging.Open(srcPath)
	if err != nil {
		return fmt.Errorf("打开源图片失败: %w", err)
	}

	resized := imaging.Resize(img, width, height, imaging.Lanczos)

	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(dstPath))
	switch ext {
	case ".jpg", ".jpeg":
		err = imaging.Save(resized, dstPath, imaging.JPEGQuality(85))
	case ".png":
		err = imaging.Save(resized, dstPath)
	default:
		err = imaging.Save(resized, dstPath, imaging.JPEGQuality(85))
	}

	if err != nil {
		return fmt.Errorf("保存图片失败: %w", err)
	}

	return nil
}

func (p *Processor) GeneratePreview(srcPath, dstPath string, width, height int) error {
	img, err := imaging.Open(srcPath)
	if err != nil {
		return fmt.Errorf("打开源图片失败: %w", err)
	}

	preview := imaging.Resize(img, width, height, imaging.Lanczos)

	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(dstPath))
	switch ext {
	case ".jpg", ".jpeg":
		err = imaging.Save(preview, dstPath, imaging.JPEGQuality(85))
	case ".png":
		err = imaging.Save(preview, dstPath)
	default:
		err = imaging.Save(preview, dstPath, imaging.JPEGQuality(85))
	}

	if err != nil {
		return fmt.Errorf("保存预览图失败: %w", err)
	}

	return nil
}

func (p *Processor) CropImage(srcPath, dstPath string, width, height int) error {
	img, err := imaging.Open(srcPath)
	if err != nil {
		return fmt.Errorf("打开源图片失败: %w", err)
	}

	cropped := imaging.CropCenter(img, width, height)

	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(dstPath))
	switch ext {
	case ".jpg", ".jpeg":
		err = imaging.Save(cropped, dstPath, imaging.JPEGQuality(85))
	case ".png":
		err = imaging.Save(cropped, dstPath)
	default:
		err = imaging.Save(cropped, dstPath, imaging.JPEGQuality(85))
	}

	if err != nil {
		return fmt.Errorf("保存图片失败: %w", err)
	}

	return nil
}

func GetFileExtension(format string) string {
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		return ".jpg"
	case "png":
		return ".png"
	default:
		return ".jpg"
	}
}
