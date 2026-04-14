package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// CalculateFileMD5 计算文件的 MD5 值
func CalculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
