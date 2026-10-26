package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"webGL-720yun/config"
	"webGL-720yun/internal/model"
	"webGL-720yun/internal/slice"
	sliceservice "webGL-720yun/internal/slice/service"
	"webGL-720yun/pkg/database"
	"webGL-720yun/pkg/logger"
	"webGL-720yun/pkg/minio_client"
)

// 重新生成指定场景（不传参数则全部场景）的三级 LOD 瓦片。
//
// 用法（在 backend 目录下）：
//
//	go run ./cmd/reslice
//	go run ./cmd/reslice dongrukouguangchang-3574b6c9 doujiangyanjingqudamen-45739910
//
// 只删除 preview 用于绕过 SliceProcessor 的“资源已存在就跳过”检查，
// 旧瓦片在生成过程中仍然可用，新瓦片按相同对象名覆盖，避免出现 404。
func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}
	logger.Setup(&config.Conf.Logger)

	db, err := database.NewMySQL(&config.Conf.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	minioClient, err := minio_client.NewMinIOClient(&config.MinIOConfig{
		Endpoint:  config.Conf.MinIO.Endpoint,
		AccessKey: config.Conf.MinIO.AccessKey,
		SecretKey: config.Conf.MinIO.SecretKey,
		Bucket:    config.Conf.MinIO.Bucket,
		UseSSL:    config.Conf.MinIO.UseSSL,
		Region:    config.Conf.MinIO.Region,
	})
	if err != nil {
		log.Fatalf("初始化 MinIO 失败: %v", err)
	}

	query := db.DB.Preload("Space").Order("id")
	if len(os.Args) > 1 {
		query = query.Where("scene_code IN ?", os.Args[1:])
	}

	var scenes []model.ResScene
	if err := query.Find(&scenes).Error; err != nil {
		log.Fatalf("查询场景失败: %v", err)
	}
	if len(scenes) == 0 {
		log.Fatal("没有找到需要重新切片的场景")
	}

	processor := slice.NewSliceProcessor(db.DB, minioClient, nil)
	ctx := context.Background()

	failed := 0
	skipped := 0
	for i, scene := range scenes {
		log.Printf("[%d/%d] 重新切片 %s", i+1, len(scenes), scene.SceneCode)

		sourceKey := scene.FileID
		if sourceKey == "" {
			sourceKey = scene.SceneCode
		}
		sourcePath := model.GetSceneSourcePath(scene.Space.Slug, sourceKey)
		exists, err := minioClient.ObjectExists(sourcePath)
		if err != nil || !exists {
			log.Printf("  跳过：源文件不存在 %s（err=%v）", sourcePath, err)
			skipped++
			continue
		}

		previewPath := model.GetScenePreviewPath(scene.Space.Slug, scene.SceneCode)
		if err := minioClient.DeleteObject(previewPath); err != nil {
			log.Printf("  删除旧预览失败（继续）: %v", err)
		}

		task := &sliceservice.SliceTask{
			TaskID:    fmt.Sprintf("reslice-%s-%d", scene.SceneCode, time.Now().UnixNano()),
			SceneID:   scene.ID,
			SceneCode: scene.SceneCode,
			FileID:    sourceKey,
			SpaceName: scene.Space.Name,
			SpaceSlug: scene.Space.Slug,
			UserID:    0,
		}

		if err := processor.ProcessWithRetry(ctx, task); err != nil {
			failed++
			log.Printf("  ❌ 失败: %v", err)
			continue
		}
		log.Printf("  ✅ 完成")
	}

	if failed > 0 {
		log.Fatalf("完成 %d 个，失败 %d 个，跳过 %d 个", len(scenes)-failed-skipped, failed, skipped)
	}
	log.Printf("全部完成：成功 %d 个，跳过 %d 个", len(scenes)-skipped, skipped)
}
