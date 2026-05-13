package minio_client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"webGL-720yun/internal/model"

	"github.com/minio/minio-go/v7"
	"golang.org/x/sync/errgroup"
)

// DownloadSceneSource 下载场景的源文件
func (m *MinIOClient) DownloadSceneSource(ctx context.Context, spaceSlug, fileID, destPath string) error {
	objectName := model.GetSceneSourcePath(spaceSlug, fileID)
	return m.DownloadFile(objectName, destPath)
}

// UploadScenePreview 上传场景预览图
func (m *MinIOClient) UploadScenePreview(ctx context.Context, spaceSlug, sceneCode, localFilePath string) (string, error) {
	objectName := model.GetScenePreviewPath(spaceSlug, sceneCode)
	url, err := m.UploadFile(objectName, localFilePath, "image/jpeg")
	if err != nil {
		return "", fmt.Errorf("上传预览图失败: %w", err)
	}
	return url, nil
}

// UploadSceneTiles 批量上传场景瓦片
// tileFiles: map[faceName]map[level][]filePath
func (m *MinIOClient) UploadSceneTiles(ctx context.Context, spaceSlug, sceneCode string, tileFiles map[string]map[int][]string) error {
	type task struct {
		objectName string
		filePath   string
	}

	var tasks []task
	for faceName, levels := range tileFiles {
		for level, files := range levels {
			for _, file := range files {
				fileName := filepath.Base(file)
				var x, y int
				fmt.Sscanf(fileName, "tile_%d_%d.jpg", &x, &y)

				objectName := model.GetSceneTilePath(spaceSlug, sceneCode, faceName, level, x, y)
				tasks = append(tasks, task{objectName: objectName, filePath: file})
			}
		}
	}

	g, ctx := errgroup.WithContext(ctx)
	const maxConcurrent = 15 // 稍微放宽并发限制
	sem := make(chan struct{}, maxConcurrent)

	for _, t := range tasks {
		t := t // 闭包安全
		g.Go(func() error {
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return ctx.Err()
			}

			return m.uploadFile(ctx, t.objectName, t.filePath)
		})
	}

	return g.Wait()
}

// CheckSceneAssetsExist 检查场景资源是否已存在
func (m *MinIOClient) CheckSceneAssetsExist(ctx context.Context, spaceSlug, sceneCode string) bool {
	previewPath := model.GetScenePreviewPath(spaceSlug, sceneCode)
	exists, err := m.ObjectExists(previewPath)
	if err != nil || !exists {
		return false
	}

	tilePath := model.GetSceneTilePath(spaceSlug, sceneCode, "px", 0, 0, 0)
	exists, err = m.ObjectExists(tilePath)
	return err == nil && exists
}

// uploadFile 内部辅助方法
func (m *MinIOClient) uploadFile(ctx context.Context, objectName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = m.client.PutObject(ctx, m.config.Bucket, objectName, file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	return err
}
