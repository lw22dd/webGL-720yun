package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	BASE_URL      = "http://localhost:7000"
	TEST_DATA_DIR = `d:\lwdd\code\毕设\webGL-720yun\参考\都江堰`
)

type PerfResult struct {
	FileSize     int64
	UploadTime   time.Duration
	CreateTime   time.Duration
	SliceTime    time.Duration
	TotalTime    time.Duration
	Success      bool
	ErrorMessage string
}

func main() {
	fmt.Println("=== Slice Pipeline Performance Test ===")
	fmt.Println()

	token := login()
	if token == "" {
		fmt.Println("Login failed")
		return
	}

	spaceID := getOrCreateSpace(token)
	if spaceID == 0 {
		fmt.Println("Failed to get/create space")
		return
	}

	files, err := filepath.Glob(filepath.Join(TEST_DATA_DIR, "*_sphere.jpg"))
	if err != nil {
		fmt.Printf("Failed to list files: %v\n", err)
		return
	}

	fmt.Printf("Found %d panorama files\n", len(files))
	fmt.Println()

	results := make([]*PerfResult, 0)
	totalStart := time.Now()

	for i, file := range files {
		fmt.Printf("[%d/%d] Testing: %s\n", i+1, len(files), filepath.Base(file))

		result := runPerfTest(token, spaceID, file)
		results = append(results, result)

		if result.Success {
			fmt.Printf("  Upload: %v, Create: %v, Slice: %v, Total: %v\n",
				result.UploadTime.Round(time.Millisecond),
				result.CreateTime.Round(time.Millisecond),
				result.SliceTime.Round(time.Millisecond),
				result.TotalTime.Round(time.Millisecond))
		} else {
			fmt.Printf("  FAILED: %s\n", result.ErrorMessage)
		}
		fmt.Println()
	}

	totalDuration := time.Since(totalStart)

	fmt.Println("=== Performance Summary ===")
	fmt.Printf("Total test duration: %v\n", totalDuration.Round(time.Millisecond))
	fmt.Printf("Files processed: %d\n", len(results))
	fmt.Println()

	successCount := 0
	var totalUploadTime, totalCreateTime, totalSliceTime time.Duration
	var totalFileSize int64
	var uploadTimes, createTimes, sliceTimes []time.Duration
	var fileSizes []int64

	for _, r := range results {
		if r.Success {
			successCount++
			totalUploadTime += r.UploadTime
			totalCreateTime += r.CreateTime
			totalSliceTime += r.SliceTime
			totalFileSize += r.FileSize
			uploadTimes = append(uploadTimes, r.UploadTime)
			createTimes = append(createTimes, r.CreateTime)
			sliceTimes = append(sliceTimes, r.SliceTime)
			fileSizes = append(fileSizes, r.FileSize)
		}
	}

	if successCount > 0 {
		fmt.Printf("Success rate: %d/%d (%.1f%%)\n", successCount, len(results), float64(successCount)/float64(len(results))*100)
		fmt.Printf("Total data processed: %.2f MB\n", float64(totalFileSize)/1024/1024)
		fmt.Println()

		fmt.Println("Upload Performance:")
		fmt.Printf("  Min: %v\n", minDuration(uploadTimes).Round(time.Millisecond))
		fmt.Printf("  Max: %v\n", maxDuration(uploadTimes).Round(time.Millisecond))
		fmt.Printf("  Avg: %v\n", avgDuration(uploadTimes).Round(time.Millisecond))
		fmt.Println()

		fmt.Println("Scene Creation Performance:")
		fmt.Printf("  Min: %v\n", minDuration(createTimes).Round(time.Millisecond))
		fmt.Printf("  Max: %v\n", maxDuration(createTimes).Round(time.Millisecond))
		fmt.Printf("  Avg: %v\n", avgDuration(createTimes).Round(time.Millisecond))
		fmt.Println()

		fmt.Println("Slice Pipeline Performance:")
		fmt.Printf("  Min: %v\n", minDuration(sliceTimes).Round(time.Millisecond))
		fmt.Printf("  Max: %v\n", maxDuration(sliceTimes).Round(time.Millisecond))
		fmt.Printf("  Avg: %v\n", avgDuration(sliceTimes).Round(time.Millisecond))
		fmt.Println()

		fmt.Println("=== Optimization Recommendations ===")

		avgSlice := avgDuration(sliceTimes)
		avgUpload := avgDuration(uploadTimes)

		if avgSlice > 30*time.Second {
			fmt.Println("1. Slice pipeline is slow (>30s average):")
			fmt.Println("   - Consider parallel tile generation for different faces")
			fmt.Println("   - Use GPU acceleration for E2C conversion")
			fmt.Println("   - Implement progressive tile generation")
		}

		if avgUpload > 10*time.Second {
			fmt.Println("2. Upload is slow (>10s average):")
			fmt.Println("   - Consider increasing chunk size for large files")
			fmt.Println("   - Implement parallel chunk uploads")
			fmt.Println("   - Use compression for network transfer")
		}

		totalMB := float64(totalFileSize) / 1024 / 1024
		throughput := totalMB / totalDuration.Seconds()

		fmt.Printf("3. Current throughput: %.2f MB/s\n", throughput)
		fmt.Println("   - For production, aim for >5 MB/s throughput")
	}
}

func runPerfTest(token string, spaceID uint, filePath string) *PerfResult {
	result := &PerfResult{}
	fileInfo, _ := os.Stat(filePath)
	result.FileSize = fileInfo.Size()

	fileMD5, _ := getFileMD5(filePath)

	totalStart := time.Now()

	uploadStart := time.Now()
	fileID, err := uploadFile(token, spaceID, filePath, fileMD5)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Upload failed: %v", err)
		return result
	}
	result.UploadTime = time.Since(uploadStart)

	createStart := time.Now()
	sceneID, err := createScene(token, spaceID, fileID, filePath)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Create scene failed: %v", err)
		return result
	}
	result.CreateTime = time.Since(createStart)

	sliceStart := time.Now()
	status := waitForSliceComplete(token, sceneID, 5*time.Minute)
	if status == "ready" {
		result.SliceTime = time.Since(sliceStart)
		result.TotalTime = time.Since(totalStart)
		result.Success = true
		return result
	} else if status == "failed" {
		result.ErrorMessage = "Slice task failed"
		return result
	}

	result.ErrorMessage = "Slice timeout"
	return result
}

func waitForSliceComplete(token string, sceneID uint, timeout time.Duration) string {
	checkInterval := 2 * time.Second
	start := time.Now()

	for time.Since(start) < timeout {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/resource/scenes/%d", BASE_URL, sceneID), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			time.Sleep(checkInterval)
			continue
		}
		defer resp.Body.Close()

		var detailResp struct {
			Data struct {
				SliceStatus string `json:"slice_status"`
			} `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&detailResp)

		if detailResp.Data.SliceStatus == "ready" {
			return "ready"
		} else if detailResp.Data.SliceStatus == "failed" {
			return "failed"
		}
		time.Sleep(checkInterval)
	}

	return "timeout"
}

func login() string {
	body := map[string]string{"username": "admin", "password": "admin123"}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(BASE_URL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var loginResp struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&loginResp)
	return loginResp.Data.AccessToken
}

func getOrCreateSpace(token string) uint {
	req, _ := http.NewRequest("GET", BASE_URL+"/api/v1/resource/spaces?page=1&page_size=100", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	var spacesResp struct {
		Data struct {
			Spaces []struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			} `json:"spaces"`
		} `json:"data"`
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &spacesResp)

	for _, s := range spacesResp.Data.Spaces {
		if s.Name == "dujiangyan" || s.Name == "Dujiangyan" {
			return s.ID
		}
	}

	return 0
}

func getFileMD5(filePath string) (string, error) {
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

func uploadFile(token string, spaceID uint, filePath, fileMD5 string) (string, error) {
	fileInfo, _ := os.Stat(filePath)

	initBody := map[string]interface{}{
		"space_id":  spaceID,
		"filename":  filepath.Base(filePath),
		"file_size": fileInfo.Size(),
		"file_hash": fileMD5,
	}
	jsonBody, _ := json.Marshal(initBody)

	req, _ := http.NewRequest("POST", BASE_URL+"/api/v1/upload/init", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var initResp struct {
		Data struct {
			UploadID    string `json:"upload_id"`
			Instant     bool   `json:"instant"`
			FileID      string `json:"file_id"`
			ChunkSize   int    `json:"chunk_size"`
			TotalChunks int    `json:"total_chunks"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&initResp)

	if initResp.Data.Instant {
		return initResp.Data.FileID, nil
	}

	file, _ := os.Open(filePath)
	defer file.Close()

	chunkSize := initResp.Data.ChunkSize
	totalChunks := initResp.Data.TotalChunks
	uploadID := initResp.Data.UploadID
	buffer := make([]byte, chunkSize)

	for i := 0; i < totalChunks; i++ {
		n, _ := file.Read(buffer)
		chunkData := buffer[:n]
		chunkMD5 := md5.Sum(chunkData)
		chunkMD5Str := hex.EncodeToString(chunkMD5[:])

		if err := uploadChunk(token, uploadID, i, chunkMD5Str, chunkData); err != nil {
			return "", err
		}
	}

	completeBody := map[string]string{
		"upload_id": uploadID,
		"file_hash": fileMD5,
	}
	jsonBody, _ = json.Marshal(completeBody)

	req, _ = http.NewRequest("POST", BASE_URL+"/api/v1/upload/complete", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var completeResp struct {
		Data struct {
			FileID string `json:"file_id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&completeResp)

	return completeResp.Data.FileID, nil
}

func uploadChunk(token, uploadID string, chunkIndex int, chunkMD5 string, chunkData []byte) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("upload_id", uploadID)
	writer.WriteField("chunk_index", fmt.Sprintf("%d", chunkIndex))
	writer.WriteField("chunk_hash", chunkMD5)

	part, _ := writer.CreateFormFile("chunk_data", "chunk.bin")
	part.Write(chunkData)
	writer.Close()

	req, _ := http.NewRequest("POST", BASE_URL+"/api/v1/upload/chunk", body)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func createScene(token string, spaceID uint, fileID, filePath string) (uint, error) {
	sceneCode := fmt.Sprintf("perf_%d", time.Now().UnixNano())
	title := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	body := map[string]interface{}{
		"space_id":       spaceID,
		"title":          title,
		"scene_code":     sceneCode,
		"file_id":        fileID,
		"panorama_type":  "equirectangular",
		"initial_fov":    90,
		"initial_pitch":  0,
		"initial_yaw":    0,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", BASE_URL+"/api/v1/resource/scenes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var sceneResp struct {
		Data struct {
			SceneID uint `json:"scene_id"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&sceneResp)

	return sceneResp.Data.SceneID, nil
}

func minDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	min := durations[0]
	for _, d := range durations {
		if d < min {
			min = d
		}
	}
	return min
}

func maxDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	max := durations[0]
	for _, d := range durations {
		if d > max {
			max = d
		}
	}
	return max
}

func avgDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}
