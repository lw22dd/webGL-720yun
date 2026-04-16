package minio_client

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"webGL-720yun/config"
	"webGL-720yun/pkg/logger"
)

type MinIOClient struct {
	client *minio.Client
	config *config.MinIOConfig
}

func NewMinIOClient(cfg *config.MinIOConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("创建MinIO客户端失败: %w", err)
	}

	m := &MinIOClient{
		client: client,
		config: cfg,
	}

	if err := m.ensureBucketExists(context.Background()); err != nil {
		return nil, fmt.Errorf("确保Bucket存在失败: %w", err)
	}

	return m, nil
}

func (m *MinIOClient) ensureBucketExists(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.config.Bucket)
	if err != nil {
		return fmt.Errorf("检查Bucket存在性失败: %w", err)
	}
	if !exists {
		if err := m.client.MakeBucket(ctx, m.config.Bucket, minio.MakeBucketOptions{Region: m.config.Region}); err != nil {
			return fmt.Errorf("创建Bucket失败: %w", err)
		}
		logger.Infof("✅ 创建MinIO Bucket: %s", m.config.Bucket)
	}
	return nil
}

func (m *MinIOClient) UploadFile(objectName, localFilePath, contentType string) (string, error) {
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(localFilePath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	file, err := os.Open(localFilePath)
	if err != nil {
		return "", fmt.Errorf("打开本地文件失败 [%s]: %w", localFilePath, err)
	}
	defer file.Close()

 fileInfo, _ := file.Stat()
 fileSize := fileInfo.Size()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	info, err := m.client.PutObject(ctx, m.config.Bucket, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		return "", fmt.Errorf("上传文件到MinIO失败 [%s]: %w", objectName, err)
	}

	url := m.GetObjectURL(objectName)
	logger.Infof("☁️  上传成功: %s (%d bytes)", objectName, info.Size)
	return url, nil
}

func (m *MinIOClient) DownloadFile(objectName, localFilePath string) error {
	dir := filepath.Dir(localFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	obj, err := m.client.GetObject(ctx, m.config.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("从MinIO获取对象失败 [%s]: %w", objectName, err)
	}
	defer obj.Close()

	localFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败 [%s]: %w", localFilePath, err)
	}
	defer localFile.Close()

	if _, err := io.Copy(localFile, obj); err != nil {
		return fmt.Errorf("写入本地文件失败: %w", err)
	}

	logger.Infof("📥 下载成功: %s → %s", objectName, localFilePath)
	return nil
}

func (m *MinIOClient) ObjectExists(objectName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.client.StatObject(ctx, m.config.Bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("检查对象存在性失败 [%s]: %w", objectName, err)
	}
	return true, nil
}

func (m *MinIOClient) DeleteObject(objectName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.client.RemoveObject(ctx, m.config.Bucket, objectName, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("删除对象失败 [%s]: %w", objectName, err)
	}

	logger.Infof("🗑️  删除成功: %s", objectName)
	return nil
}

func (m *MinIOClient) DeleteObjectsWithPrefix(prefix string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	objectsCh := m.client.ListObjects(ctx, m.config.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for obj := range objectsCh {
		if obj.Err != nil {
			return fmt.Errorf("列出对象失败 [%s]: %w", prefix, obj.Err)
		}
		if err := m.client.RemoveObject(ctx, m.config.Bucket, obj.Key, minio.RemoveObjectOptions{}); err != nil {
			return fmt.Errorf("删除对象失败 [%s]: %w", obj.Key, err)
		}
	}

	logger.Infof("🗑️  递归删除成功: %s", prefix)
	return nil
}

func (m *MinIOClient) GetObjectURL(objectName string) string {
	scheme := "http"
	if m.config.UseSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, m.config.Endpoint, m.config.Bucket, objectName)
}

func (m *MinIOClient) GetPresignedURL(objectName string, expiry time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, m.config.Bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败 [%s]: %w", objectName, err)
	}
	return presignedURL.String(), nil
}

func (m *MinIOClient) ListObjects(prefix string, recursive bool) ([]minio.ObjectInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var objects []minio.ObjectInfo
	objectCh := m.client.ListObjects(ctx, m.config.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	})

	for obj := range objectCh {
		if obj.Err != nil {
			return objects, fmt.Errorf("列出对象失败: %w", obj.Err)
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

// GetObjectStream 流式获取 MinIO 对象，返回 *minio.Object（实现 io.ReadCloser）
// 用于瓦片/预览图等资源的零拷贝流式转发，不会将整个文件读入内存
func (m *MinIOClient) GetObjectStream(ctx context.Context, objectName string) (*minio.Object, error) {
	obj, err := m.client.GetObject(ctx, m.config.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取对象流失败 [%s]: %w", objectName, err)
	}
	return obj, nil
}

// StitchTilesToStream 从多个瓦片拼接成完整图片并返回流
// tiles: 瓦片路径列表，按从左到右、从上到下的顺序排列
// cols, rows: 瓦片的列数和行数
func (m *MinIOClient) StitchTilesToStream(ctx context.Context, tilePaths []string, cols, rows int) (io.ReadCloser, error) {
	if len(tilePaths) != cols*rows {
		return nil, fmt.Errorf("瓦片数量 %d 与布局 %dx%d 不匹配", len(tilePaths), cols, rows)
	}

	// 读取所有瓦片
	var tiles []image.Image
	var tileWidth, tileHeight int

	for i, path := range tilePaths {
		obj, err := m.client.GetObject(ctx, m.config.Bucket, path, minio.GetObjectOptions{})
		if err != nil {
			return nil, fmt.Errorf("获取瓦片失败 [%s]: %w", path, err)
		}

		img, err := imaging.Decode(obj)
		obj.Close()
		if err != nil {
			return nil, fmt.Errorf("解码瓦片失败 [%s]: %w", path, err)
		}

		if i == 0 {
			tileWidth = img.Bounds().Dx()
			tileHeight = img.Bounds().Dy()
		}

		tiles = append(tiles, img)
	}

	// 创建拼接后的图片
	resultWidth := tileWidth * cols
	resultHeight := tileHeight * rows
	resultImg := imaging.New(resultWidth, resultHeight, image.Black)

	// 拼接瓦片
	for i, tile := range tiles {
		row := i / cols
		col := i % cols
		x := col * tileWidth
		y := row * tileHeight
		resultImg = imaging.Paste(resultImg, tile, image.Pt(x, y))
	}

	// 编码为 JPEG 并写入 buffer
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, resultImg, &jpeg.Options{Quality: 85}); err != nil {
		return nil, fmt.Errorf("编码拼接图片失败: %w", err)
	}

	// 返回一个 io.ReadCloser
	return io.NopCloser(&buf), nil
}

func (m *MinIOClient) GetClient() *minio.Client {
	return m.client
}

func (m *MinIOClient) GetConfig() *config.MinIOConfig {
	return m.config
}

func (m *MinIOClient) UploadDirectory(prefix, localDir string) error {
	return filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		// 计算相对路径
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}
		
		// 构建对象名称
		objectName := filepath.Join(prefix, relPath)
		
		// 上传文件
		_, err = m.UploadFile(objectName, path, "")
		if err != nil {
			return err
		}
		
		return nil
	})
}