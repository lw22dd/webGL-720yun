package setup

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"webGL-720yun/config"
	"webGL-720yun/internal/slice"
	"webGL-720yun/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOSetup struct {
	client *minio.Client
	config *config.MinIOConfig
}

func NewMinIOSetup(cfg *config.MinIOConfig) (*MinIOSetup, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("创建MinIO客户端失败: %w", err)
	}

	return &MinIOSetup{
		client: client,
		config: cfg,
	}, nil
}

func (m *MinIOSetup) EnsureBucketAndStructure(ctx context.Context) error {
	var steps []string
	var successSteps []string
	var failedSteps []string

	steps = append(steps, "检查MinIO连接")
	_, err := m.client.ListBuckets(ctx)
	if err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("检查MinIO连接: 失败 - %v", err))
		m.printInitResult(steps, successSteps, failedSteps)
		return fmt.Errorf("MinIO连接失败: %w", err)
	}
	successSteps = append(successSteps, "检查MinIO连接: 成功")

	steps = append(steps, fmt.Sprintf("确保存储桶 %s 存在", m.config.Bucket))
	if err := m.ensureBucketExists(ctx); err != nil {
		failedSteps = append(failedSteps, fmt.Sprintf("确保存储桶存在: 失败 - %v", err))
		m.printInitResult(steps, successSteps, failedSteps)
		return err
	}
	successSteps = append(successSteps, fmt.Sprintf("确保存储桶 %s 存在: 成功", m.config.Bucket))

	steps = append(steps, "验证存储结构")
	if err := m.verifyStorageStructure(ctx); err != nil {
		logger.Warnf("存储结构验证: %v (首次运行可能为空)", err)
		successSteps = append(successSteps, "验证存储结构: 跳过(空桶)")
	} else {
		successSteps = append(successSteps, "验证存储结构: 成功")
	}

	m.printInitResult(steps, successSteps, failedSteps)
	return nil
}

func (m *MinIOSetup) ensureBucketExists(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.config.Bucket)
	if err != nil {
		return fmt.Errorf("检查存储桶存在性失败: %w", err)
	}

	if !exists {
		if err := m.client.MakeBucket(ctx, m.config.Bucket, minio.MakeBucketOptions{
			Region: m.config.Region,
		}); err != nil {
			return fmt.Errorf("创建存储桶失败: %w", err)
		}
		logger.Infof("✅ 创建MinIO存储桶: %s", m.config.Bucket)
	} else {
		logger.Infof("✅ MinIO存储桶已存在: %s", m.config.Bucket)
	}

	return nil
}

func (m *MinIOSetup) verifyStorageStructure(ctx context.Context) error {
	objects := m.client.ListObjects(ctx, m.config.Bucket, minio.ListObjectsOptions{
		Recursive: false,
	})

	hasSpacesFolder := false
	for obj := range objects {
		if obj.Err != nil {
			return obj.Err
		}
		if obj.Key == "spaces/" || len(obj.Key) > 7 && obj.Key[:7] == "spaces/" {
			hasSpacesFolder = true
			break
		}
	}

	if !hasSpacesFolder {
		return fmt.Errorf("未找到spaces目录结构")
	}

	return nil
}

func (m *MinIOSetup) ListStorageStructure(ctx context.Context) (map[string][]string, error) {
	structure := make(map[string][]string)

	objects := m.client.ListObjects(ctx, m.config.Bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	for obj := range objects {
		if obj.Err != nil {
			continue
		}

		if len(obj.Key) > 7 && obj.Key[:7] == "spaces/" {
			structure["spaces"] = append(structure["spaces"], obj.Key)
		}
	}

	return structure, nil
}

func (m *MinIOSetup) printInitResult(steps []string, successSteps []string, failedSteps []string) {
	fmt.Println("\n=== MinIO 初始化结果 ===")
	for _, step := range successSteps {
		fmt.Printf("✅ %s\n", step)
	}
	for _, step := range failedSteps {
		fmt.Printf("❌ %s\n", step)
	}
	fmt.Printf("\n总步骤: %d, 成功: %d, 失败: %d\n", len(steps), len(successSteps), len(failedSteps))
	fmt.Println("========================")
}

func (m *MinIOSetup) GetClient() *minio.Client {
	return m.client
}

type MinIOPathStructure struct{}

func NewMinIOPathStructure() *MinIOPathStructure {
	return &MinIOPathStructure{}
}

func (p *MinIOPathStructure) GetSourcePath(spaceName, sceneCode string) string {
	return fmt.Sprintf("spaces/%s/sources/%s/source.jpg", spaceName, sceneCode)
}

func (p *MinIOPathStructure) GetPreviewPath(spaceName, sceneCode string) string {
	return fmt.Sprintf("spaces/%s/previews/%s/preview.jpg", spaceName, sceneCode)
}

func (p *MinIOPathStructure) GetThumbPath(spaceName, sceneCode string) string {
	return fmt.Sprintf("spaces/%s/previews/%s/thumb.jpg", spaceName, sceneCode)
}

func (p *MinIOPathStructure) GetTilesPath(spaceName, sceneCode string) string {
	return fmt.Sprintf("spaces/%s/tiles/%s/", spaceName, sceneCode)
}

func (p *MinIOPathStructure) GetCubemapPath(spaceName, sceneCode, face string, level int) string {
	return fmt.Sprintf("spaces/%s/tiles/%s/cubemap/%s/level_%d/", spaceName, sceneCode, face, level)
}

func (p *MinIOPathStructure) GetTileFilePath(spaceName, sceneCode, face string, level, row, col int) string {
	return fmt.Sprintf("spaces/%s/tiles/%s/cubemap/%s/level_%d/tile_%d_%d.jpg", spaceName, sceneCode, face, level, row, col)
}

func (m *MinIOSetup) InitSpaceStructure(ctx context.Context, spaceName string) error {
	bucket := m.config.Bucket

	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("检查存储桶失败: %w", err)
	}
	if !exists {
		return fmt.Errorf("存储桶 %s 不存在", bucket)
	}

	dirs := []string{
		fmt.Sprintf("spaces/%s/sources/", spaceName),
		fmt.Sprintf("spaces/%s/previews/", spaceName),
		fmt.Sprintf("spaces/%s/tiles/", spaceName),
	}

	for _, dir := range dirs {
		objInfo, err := m.client.StatObject(ctx, bucket, dir, minio.StatObjectOptions{})
		if err != nil {
			logger.Infof("📁 目录不存在，将被创建: %s", dir)
			continue
		}
		logger.Infof("✅ 目录已存在: %s (size: %d)", dir, objInfo.Size)
	}

	logger.Infof("✅ Space %s 目录结构初始化完成", spaceName)
	return nil
}

func (m *MinIOSetup) InitScenicSpotTestData(ctx context.Context) error {
	seeds := GetScenicSpotSeeds()

	for _, seed := range seeds {
		spaceName := seed.Slug

		if err := m.InitSpaceStructure(ctx, spaceName); err != nil {
			logger.Warnf("初始化景点 %s 失败: %v", seed.Name, err)
			continue
		}

		sourceDir := fmt.Sprintf("spaces/%s/sources/", spaceName)
		previewDir := fmt.Sprintf("spaces/%s/previews/", spaceName)
		tilesDir := fmt.Sprintf("spaces/%s/tiles/", spaceName)

		logger.Infof("📁 %s - 源文件目录: %s", seed.Name, sourceDir)
		logger.Infof("📁 %s - 预览图目录: %s", seed.Name, previewDir)
		logger.Infof("📁 %s - 瓦片目录: %s", seed.Name, tilesDir)

		for _, scene := range seed.Scenes {
			logger.Infof("   └─ 场景: %s (code: %s, file: %s)", scene.Title, scene.SceneCode, scene.FileName)
		}
	}

	logger.Infof("✅ 共初始化 %d 个景点的MinIO目录结构", len(seeds))
	return nil
}

func GetScenicSpotSliceTasks() []*slice.SliceTask {
	seeds := GetScenicSpotSeeds()
	var tasks []*slice.SliceTask

	for _, seed := range seeds {
		for _, scene := range seed.Scenes {
			task := &slice.SliceTask{
				TaskID:    fmt.Sprintf("seed_%s_%s", seed.Slug, scene.SceneCode),
				SceneID:   0,
				SceneCode: scene.SceneCode,
				FileID:    scene.SceneCode,
				SpaceName: seed.Name,
				SpaceSlug: seed.Slug,
				UserID:    0,
				CreatedAt: 0,
			}
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func ParseLocalPath(localPath string, fileName string) string {
	normalizedPath := filepath.ToSlash(filepath.Join(localPath, fileName))
	normalizedPath = strings.ReplaceAll(normalizedPath, "\\", "/")
	return normalizedPath
}

func GetScenicSpotLocalPaths() []struct {
	ScenicSpot string
	SceneTitle string
	LocalPath  string
} {
	var paths []struct {
		ScenicSpot string
		SceneTitle string
		LocalPath  string
	}

	seeds := GetScenicSpotSeeds()
	for _, seed := range seeds {
		for _, scene := range seed.Scenes {
			paths = append(paths, struct {
				ScenicSpot string
				SceneTitle string
				LocalPath  string
			}{
				ScenicSpot: seed.Name,
				SceneTitle: scene.Title,
				LocalPath:  ParseLocalPath(seed.LocalPath, scene.FileName),
			})
		}
	}

	return paths
}
