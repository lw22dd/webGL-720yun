package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	endpoint := "localhost:9000"
	accessKey := "admin"
	secretKey := "admin123456"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("连接MinIO失败: %v", err)
	}

	bucketsToDelete := []string{"panoramas", "thumbnails", "tiles"}
	oldBucket := "pannellum-resources"
	newBucket := "resource-scene"

	fmt.Println("=== Step 1: Delete unused buckets ===")
	for _, bucket := range bucketsToDelete {
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			log.Printf("检查Bucket %s 失败: %v", bucket, err)
			continue
		}
		if !exists {
			fmt.Printf("Bucket %s 不存在，跳过\n", bucket)
			continue
		}

		objectsCh := client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true})
		var objects []minio.ObjectInfo
		for obj := range objectsCh {
			if obj.Err != nil {
				log.Printf("列出 %s 对象失败: %v", bucket, obj.Err)
				continue
			}
			objects = append(objects, obj)
		}

		if len(objects) > 0 {
			fmt.Printf("Bucket %s 有 %d 个对象，正在删除...\n", bucket, len(objects))
			for _, obj := range objects {
				if err := client.RemoveObject(ctx, bucket, obj.Key, minio.RemoveObjectOptions{}); err != nil {
					log.Printf("删除对象 %s/%s 失败: %v", bucket, obj.Key, err)
				}
			}
		}

		if err := client.RemoveBucket(ctx, bucket); err != nil {
			log.Printf("删除Bucket %s 失败: %v", bucket, err)
		} else {
			fmt.Printf("✅ 已删除 Bucket: %s\n", bucket)
		}
	}

	fmt.Println("\n=== Step 2: Rename pannellum-resources to resource-scene ===")

	oldExists, err := client.BucketExists(ctx, oldBucket)
	if err != nil {
		log.Fatalf("检查旧Bucket失败: %v", err)
	}
	if !oldExists {
		fmt.Printf("旧Bucket %s 不存在\n", oldBucket)
		return
	}

	newExists, err := client.BucketExists(ctx, newBucket)
	if err != nil {
		log.Fatalf("检查新Bucket失败: %v", err)
	}
	if newExists {
		fmt.Printf("新Bucket %s 已存在\n", newBucket)
	} else {
		if err := client.MakeBucket(ctx, newBucket, minio.MakeBucketOptions{}); err != nil {
			log.Fatalf("创建新Bucket失败: %v", err)
		}
		fmt.Printf("✅ 创建新Bucket: %s\n", newBucket)
	}

	fmt.Println("正在复制对象...")
	objectsCh := client.ListObjects(ctx, oldBucket, minio.ListObjectsOptions{Recursive: true})
	copiedCount := 0
	for obj := range objectsCh {
		if obj.Err != nil {
			log.Printf("列出对象失败: %v", obj.Err)
			continue
		}

		_, err := client.CopyObject(ctx, minio.CopyDestOptions{
			Bucket: newBucket,
			Object: obj.Key,
		}, minio.CopySrcOptions{
			Bucket: oldBucket,
			Object: obj.Key,
		})
		if err != nil {
			log.Printf("复制对象 %s 失败: %v", obj.Key, err)
			continue
		}
		copiedCount++
		if copiedCount%50 == 0 {
			fmt.Printf("已复制 %d 个对象...\n", copiedCount)
		}
	}
	fmt.Printf("✅ 共复制 %d 个对象\n", copiedCount)

	fmt.Println("正在删除旧Bucket中的对象...")
	objectsCh = client.ListObjects(ctx, oldBucket, minio.ListObjectsOptions{Recursive: true})
	for obj := range objectsCh {
		if obj.Err != nil {
			continue
		}
		client.RemoveObject(ctx, oldBucket, obj.Key, minio.RemoveObjectOptions{})
	}

	if err := client.RemoveBucket(ctx, oldBucket); err != nil {
		log.Printf("删除旧Bucket失败: %v", err)
	} else {
		fmt.Printf("✅ 已删除旧Bucket: %s\n", oldBucket)
	}

	fmt.Println("\n=== 完成 ===")
	fmt.Println("当前Bucket列表:")
	buckets, _ := client.ListBuckets(ctx)
	for _, b := range buckets {
		fmt.Printf("📁 %s\n", b.Name)
	}
}
