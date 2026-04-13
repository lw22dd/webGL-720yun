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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	fmt.Println("=== 列出所有 Bucket ===")

	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		log.Fatalf("列出Bucket失败: %v", err)
	}

	if len(buckets) == 0 {
		fmt.Println("没有找到任何 Bucket")
		return
	}

	fmt.Println("Bucket 列表:")
	fmt.Println("-------------------")
	for _, bucket := range buckets {
		fmt.Printf("📁 %s\n", bucket.Name)
	}
}
