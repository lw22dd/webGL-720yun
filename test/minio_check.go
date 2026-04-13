package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	endpoint := "localhost:9000"
	accessKey := "admin"
	secretKey := "admin123456"
	bucketName := "resource-scene"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("连接MinIO失败: %v", err)
	}

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("检查Bucket失败: %v", err)
	}
	if !exists {
		fmt.Printf("Bucket '%s' 不存在\n", bucketName)
		return
	}
	fmt.Printf("✅ 连接到 MinIO 成功！Bucket: %s\n\n", bucketName)

	fmt.Println("=== 存储桶列表 ===")
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		log.Fatalf("列出存储桶失败: %v", err)
	}
	for _, b := range buckets {
		fmt.Printf("📁 %s\n", b.Name)
	}
	fmt.Println()

	fmt.Println("=== 存储结构分析 ===")
	objects := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	spaces := make(map[string]*SpaceInfo)
	for obj := range objects {
		if obj.Err != nil {
			log.Printf("获取对象失败: %v", obj.Err)
			continue
		}

		if !strings.HasPrefix(obj.Key, "spaces/") {
			fmt.Printf("⚠️  非标准路径: %s\n", obj.Key)
			continue
		}

		parts := strings.Split(obj.Key, "/")
		if len(parts) < 3 {
			continue
		}

		spaceName := parts[1]
		if spaces[spaceName] == nil {
			spaces[spaceName] = &SpaceInfo{
				Name:     spaceName,
				Sources:  make(map[string][]string),
				Previews: make(map[string][]string),
				Tiles:    make(map[string][]string),
			}
		}

		space := spaces[spaceName]
		remainingPath := strings.Join(parts[2:], "/")

		if strings.HasPrefix(remainingPath, "sources/") {
			sceneCode := extractSceneCode(remainingPath, "sources/")
			if sceneCode != "" {
				space.Sources[sceneCode] = append(space.Sources[sceneCode], remainingPath)
			}
		} else if strings.HasPrefix(remainingPath, "previews/") {
			sceneCode := extractSceneCode(remainingPath, "previews/")
			if sceneCode != "" {
				space.Previews[sceneCode] = append(space.Previews[sceneCode], remainingPath)
			}
		} else if strings.HasPrefix(remainingPath, "tiles/") {
			sceneCode := extractSceneCode(remainingPath, "tiles/")
			if sceneCode != "" {
				space.Tiles[sceneCode] = append(space.Tiles[sceneCode], remainingPath)
			}
		} else if strings.HasPrefix(remainingPath, "covers/") {
			space.Cover = remainingPath
		}
	}

	fmt.Printf("共发现 %d 个空间:\n\n", len(spaces))
	for _, space := range spaces {
		printSpaceInfo(space)
	}

	fmt.Println("\n=== 存储结构验证 ===")
	validateStorageStructure(spaces)
}

type SpaceInfo struct {
	Name     string
	Cover    string
	Sources  map[string][]string
	Previews map[string][]string
	Tiles    map[string][]string
}

func extractSceneCode(path, prefix string) string {
	afterPrefix := strings.TrimPrefix(path, prefix)
	parts := strings.Split(afterPrefix, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func printSpaceInfo(space *SpaceInfo) {
	fmt.Printf("📁 空间: %s\n", space.Name)
	if space.Cover != "" {
		fmt.Printf("   🖼️  封面: %s\n", space.Cover)
	}

	fmt.Printf("   📂 源文件 (%d 个场景):\n", len(space.Sources))
	for sceneCode, files := range space.Sources {
		fmt.Printf("      └─ %s: %d 文件\n", sceneCode, len(files))
	}

	fmt.Printf("   📂 预览图 (%d 个场景):\n", len(space.Previews))
	for sceneCode, files := range space.Previews {
		fmt.Printf("      └─ %s: %v\n", sceneCode, files)
	}

	fmt.Printf("   📂 切片 (%d 个场景):\n", len(space.Tiles))
	for sceneCode, files := range space.Tiles {
		fmt.Printf("      └─ %s: %d 切片文件\n", sceneCode, len(files))
	}
	fmt.Println()
}

func validateStorageStructure(spaces map[string]*SpaceInfo) {
	allValid := true

	for _, space := range spaces {
		for sceneCode := range space.Sources {
			if _, hasPreview := space.Previews[sceneCode]; !hasPreview {
				fmt.Printf("⚠️  [%s/%s] 缺少预览图\n", space.Name, sceneCode)
				allValid = false
			}
			if _, hasTiles := space.Tiles[sceneCode]; !hasTiles {
				fmt.Printf("⚠️  [%s/%s] 缺少切片\n", space.Name, sceneCode)
				allValid = false
			}
		}

		for sceneCode := range space.Previews {
			if _, hasSource := space.Sources[sceneCode]; !hasSource {
				fmt.Printf("⚠️  [%s/%s] 预览图无对应源文件\n", space.Name, sceneCode)
				allValid = false
			}
		}

		for sceneCode := range space.Tiles {
			if _, hasSource := space.Sources[sceneCode]; !hasSource {
				fmt.Printf("⚠️  [%s/%s] 切片无对应源文件\n", space.Name, sceneCode)
				allValid = false
			}
		}
	}

	if allValid {
		fmt.Println("✅ 存储结构验证通过")
	} else {
		fmt.Println("❌ 存储结构存在问题")
	}
}
