package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// SceneInfo 场景信息
var sceneInfoMap = map[string]struct {
	Title string
	Pinyin string
}{
	"255566cb-acc9-45b2-83f2-c1c126a737ab": {"都江堰景区入口", "dujiangyan-jingqu-rukou"},
	"544bfab9-bd61-4c03-bae3-9059d904cd43": {"伏龙观", "fulongguan"},
	"5fd690ad-a814-4dfa-80f5-3d18fd2fd0ee": {"安澜索桥", "anlansuoqiao"},
	"e38409a0-e3e7-4586-ac59-c6f3554d57af": {"二王庙", "erwangmiao"},
	"1f131ddc-056c-4535-b9cf-0c8ba2ddfe71": {"鱼嘴分水堤", "yuzuishuifenti"},
	"2c5a3528-42b1-42c8-b0a0-951b59768064": {"飞沙堰", "feishayan"},
	"691cde37-7d99-48fa-8196-5e65b40ed7ff": {"宝瓶口", "baopingkou"},
	"b364d568-f4af-4092-a0b3-f4218f053de7": {"秦堰楼", "qinyanlou"},
	"bea0a54d-15d4-494f-bdbf-76dfa6f341f8": {"玉垒关", "yuleiguan"},
	"d3e99ace-840f-49e7-ad9c-b941e740e0e7": {"离堆公园", "liduigongyuan"},
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
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

	fmt.Println("=== SceneCode 拼音+UUID 迁移工具 ===")
	fmt.Println()

	// 显示迁移计划
	fmt.Println("【迁移计划】")
	fmt.Println("----------------------------------------")
	fmt.Printf("%-40s -> %s\n", "原 SceneCode", "新 SceneCode (拼音-UUID)")
	fmt.Println("----------------------------------------")

	migrationPlan := make(map[string]string) // old -> new
	for uuid, info := range sceneInfoMap {
		newSceneCode := generateNewSceneCode(info.Pinyin, uuid)
		migrationPlan[uuid] = newSceneCode
		fmt.Printf("%-40s -> %s\n", uuid, newSceneCode)
	}
	fmt.Println("----------------------------------------")
	fmt.Println()

	// 模拟迁移（不实际执行）
	fmt.Println("【模拟迁移】检查文件路径...")
	for oldCode, newCode := range migrationPlan {
		fmt.Printf("\n场景: %s\n", sceneInfoMap[oldCode].Title)
		
		// 检查源文件
		oldSourcePath := fmt.Sprintf("spaces/%s/sources/%s/", spaceName, oldCode)
		newSourcePath := fmt.Sprintf("spaces/%s/sources/%s/", spaceName, newCode)
		
		objects := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
			Prefix:    oldSourcePath,
			Recursive: true,
		})
		
		fileCount := 0
		for obj := range objects {
			if obj.Err != nil {
				continue
			}
			fileCount++
			newKey := strings.Replace(obj.Key, oldSourcePath, newSourcePath, 1)
			fmt.Printf("  移动: %s\n", obj.Key)
			fmt.Printf("    -> %s\n", newKey)
		}
		
		if fileCount == 0 {
			fmt.Printf("  ⚠️  未找到文件\n")
		}
	}

	fmt.Println()
	fmt.Println("【数据库更新SQL】")
	fmt.Println("-- 请根据实际数据库执行以下SQL：")
	fmt.Println()
	for oldCode, newCode := range migrationPlan {
		fmt.Printf("UPDATE res_scene SET scene_code = '%s' WHERE scene_code = '%s';\n", 
			newCode, oldCode)
	}
	fmt.Println()
	fmt.Println("【说明】")
	fmt.Println("1. 此脚本仅展示迁移计划，不会实际执行迁移")
	fmt.Println("2. 实际迁移需要：")
	fmt.Println("   a) 复制 MinIO 文件到新路径")
	fmt.Println("   b) 更新数据库 scene_code 字段")
	fmt.Println("   c) 验证无误后删除旧路径文件")
	fmt.Println("3. 建议先备份数据再执行迁移")
}

// generateNewSceneCode 生成新的 SceneCode：拼音-UUID前8位
func generateNewSceneCode(pinyin, uuid string) string {
	// 清理拼音：小写，替换空格为-
	pinyin = strings.ToLower(pinyin)
	pinyin = strings.ReplaceAll(pinyin, " ", "-")
	pinyin = strings.ReplaceAll(pinyin, "_", "-")
	
	// 移除特殊字符
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	pinyin = reg.ReplaceAllString(pinyin, "")
	
	// 取 UUID 前 8 位
	uuidPrefix := uuid
	if len(uuid) > 8 {
		uuidPrefix = uuid[:8]
	}
	
	// 组合：拼音-uuid前缀
	return fmt.Sprintf("%s-%s", pinyin, uuidPrefix)
}
