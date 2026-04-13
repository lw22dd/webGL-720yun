package setup

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"webGL-720yun/config"
	"webGL-720yun/pkg/logger"
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
	fmt.Println("========================\n")
}

func (m *MinIOSetup) GetClient() *minio.Client {
	return m.client
}
