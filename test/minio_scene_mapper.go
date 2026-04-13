package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// SceneInfo 场景信息映射
var sceneNameMap = map[string]string{
	// 你可以从数据库导出这些映射，或在这里手动维护
	"255566cb-acc9-45b2-83f2-c1c126a737ab": "都江堰景区入口",
	"544bfab9-bd61-4c03-bae3-9059d904cd43": "伏龙观",
	"5fd690ad-a814-4dfa-80f5-3d18fd2fd0ee": "安澜索桥",
	"e38409a0-e3e7-4586-ac59-c6f3554d57af": "二王庙",
	"1f131ddc-056c-4535-b9cf-0c8ba2ddfe71": "鱼嘴分水堤",
	"2c5a3528-42b1-42c8-b0a0-951b59768064": "飞沙堰",
	"691cde37-7d99-48fa-8196-5e65b40ed7ff": "宝瓶口",
	"b364d568-f4af-4092-a0b3-f4218f053de7": "秦堰楼",
	"bea0a54d-15d4-494f-bdbf-76dfa6f341f8": "玉垒关",
	"d3e99ace-840f-49e7-ad9c-b941e740e0e7": "离堆公园",
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	endpoint := "localhost:9000"
	accessKey := "admin"
	secretKey := "admin123456"
	bucketName := "resource-scene"
	spaceName := "dujiangyan"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("连接MinIO失败: %v", err)
	}

	fmt.Println("=== MinIO 场景名称映射表 ===")
	fmt.Printf("存储桶: %s\n", bucketName)
	fmt.Printf("空间: %s\n\n", spaceName)

	// 列出 sources 目录下的所有场景
	prefix := fmt.Sprintf("spaces/%s/sources/", spaceName)
	objects := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: false,
	})

	fmt.Println("场景列表:")
	fmt.Println("----------------------------------------")
	fmt.Printf("%-40s | %s\n", "文件夹名 (SceneCode)", "场景名称 (Title)")
	fmt.Println("----------------------------------------")

	for obj := range objects {
		if obj.Err != nil {
			continue
		}

		// 提取场景代码（文件夹名）
		sceneCode := extractSceneCode(obj.Key, prefix)
		if sceneCode == "" {
			continue
		}

		// 获取场景名称
		sceneName := sceneNameMap[sceneCode]
		if sceneName == "" {
			sceneName = "【未命名场景】"
		}

		fmt.Printf("%-40s | %s\n", sceneCode, sceneName)
	}

	fmt.Println("----------------------------------------")
	fmt.Println("\n提示: 你可以从数据库 res_scene 表导出 scene_code -> title 的映射")
	fmt.Println("SQL: SELECT scene_code, title FROM res_scene WHERE space_id = ?")
}

func extractSceneCode(key, prefix string) string {
	if len(key) <= len(prefix) {
		return ""
	}
	// 去掉前缀和末尾的 /
	code := key[len(prefix):]
	if len(code) > 0 && code[len(code)-1] == '/' {
		code = code[:len(code)-1]
	}
	return code
}
