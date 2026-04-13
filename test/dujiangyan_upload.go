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

type LoginResponse struct {
	Code int `json:"code"`
	Data struct {
		AccessToken string `json:"access_token"`
	} `json:"data"`
}

type SpaceResponse struct {
	Code int `json:"code"`
	Data struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"data"`
}

type SpacesListResponse struct {
	Code int `json:"code"`
	Data struct {
		Spaces []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"spaces"`
	} `json:"data"`
}

type InitUploadResponse struct {
	Code int `json:"code"`
	Data struct {
		UploadID    string `json:"upload_id"`
		Instant     bool   `json:"instant"`
		FileID      string `json:"file_id"`
		ChunkSize   int    `json:"chunk_size"`
		TotalChunks int    `json:"total_chunks"`
	} `json:"data"`
}

type CompleteUploadResponse struct {
	Code int `json:"code"`
	Data struct {
		FileID    string `json:"file_id"`
		SourceURL string `json:"source_url"`
	} `json:"data"`
}

type SceneResponse struct {
	Code int `json:"code"`
	Data struct {
		SceneID     uint   `json:"scene_id"`
		TaskID      string `json:"task_id"`
		SliceStatus string `json:"slice_status"`
	} `json:"data"`
}

type SceneDetailResponse struct {
	Code int `json:"code"`
	Data struct {
		SliceStatus string `json:"slice_status"`
		Title       string `json:"title"`
	} `json:"data"`
}

func main() {
	fmt.Println("=== Dujiangyan Panorama Upload Test ===")

	token := login()
	if token == "" {
		fmt.Println("Login failed")
		return
	}
	fmt.Println("Login successful")

	spaceID := getOrCreateSpace(token)
	if spaceID == 0 {
		fmt.Println("Failed to get/create space")
		return
	}
	fmt.Printf("Using Space ID: %d\n", spaceID)

	files, err := filepath.Glob(filepath.Join(TEST_DATA_DIR, "*_sphere.jpg"))
	if err != nil {
		fmt.Printf("Failed to list files: %v\n", err)
		return
	}

	if len(files) > 3 {
		files = files[:3]
	}
	fmt.Printf("Found %d panorama files to test\n", len(files))

	uploadedScenes := make([]struct {
		File    string
		SceneID uint
	}, 0)

	for _, file := range files {
		fmt.Printf("\n=== Processing: %s ===\n", filepath.Base(file))

		fileID, err := uploadFile(token, spaceID, file)
		if err != nil {
			fmt.Printf("Upload failed: %v\n", err)
			continue
		}
		fmt.Printf("Upload completed, File ID: %s\n", fileID)

		sceneID, err := createScene(token, spaceID, fileID, file)
		if err != nil {
			fmt.Printf("Create scene failed: %v\n", err)
			continue
		}
		fmt.Printf("Scene created, ID: %d\n", sceneID)

		uploadedScenes = append(uploadedScenes, struct {
			File    string
			SceneID uint
		}{File: filepath.Base(file), SceneID: sceneID})
	}

	fmt.Println("\n=== Waiting for slice tasks ===")
	time.Sleep(5 * time.Second)

	fmt.Println("\n=== Checking slice status ===")
	for _, scene := range uploadedScenes {
		status := getSceneStatus(token, scene.SceneID)
		fmt.Printf("%s: slice_status=%s\n", scene.File, status)
	}

	fmt.Println("\n=== Test Summary ===")
	fmt.Printf("Space ID: %d\n", spaceID)
	fmt.Printf("Uploaded scenes: %d\n", len(uploadedScenes))
}

func login() string {
	body := map[string]string{"username": "admin", "password": "admin123"}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(BASE_URL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Login request failed: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResp)
	return loginResp.Data.AccessToken
}

func getOrCreateSpace(token string) uint {
	req, _ := http.NewRequest("GET", BASE_URL+"/api/v1/resource/spaces?page=1&page_size=100", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Get spaces request failed: %v\n", err)
		return 0
	}
	defer resp.Body.Close()

	var spacesResp SpacesListResponse
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Spaces response: %s\n", string(body))
	json.Unmarshal(body, &spacesResp)

	for _, s := range spacesResp.Data.Spaces {
		if s.Name == "dujiangyan" || s.Name == "Dujiangyan" {
			fmt.Printf("Found existing space: %s (ID=%d)\n", s.Name, s.ID)
			return s.ID
		}
	}

	body2 := map[string]interface{}{
		"name":        "Dujiangyan",
		"description": "Dujiangyan Test Space",
		"province":    "Sichuan",
		"city":        "Chengdu",
		"longitude":   103.611379,
		"latitude":    31.001752,
		"zoom_level":  15,
		"sort_order":  1,
	}
	jsonBody, _ := json.Marshal(body2)

	req, _ = http.NewRequest("POST", BASE_URL+"/api/v1/resource/spaces", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Create space request failed: %v\n", err)
		return 0
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Create space response: %s\n", string(respBody))

	var spaceResp SpaceResponse
	json.Unmarshal(respBody, &spaceResp)
	return spaceResp.Data.ID
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

func uploadFile(token string, spaceID uint, filePath string) (string, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	fileMD5, err := getFileMD5(filePath)
	if err != nil {
		return "", err
	}

	fmt.Printf("File size: %.2f MB, MD5: %s\n", float64(fileInfo.Size())/1024/1024, fileMD5)

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

	var initResp InitUploadResponse
	json.NewDecoder(resp.Body).Decode(&initResp)

	if initResp.Data.Instant {
		fmt.Println("Instant upload (file already exists)")
		return initResp.Data.FileID, nil
	}

	uploadID := initResp.Data.UploadID
	chunkSize := initResp.Data.ChunkSize
	totalChunks := initResp.Data.TotalChunks

	fmt.Printf("Upload initialized: ID=%s, Chunks=%d\n", uploadID, totalChunks)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, chunkSize)
	for i := 0; i < totalChunks; i++ {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return "", err
		}

		chunkData := buffer[:n]
		chunkMD5 := md5.Sum(chunkData)
		chunkMD5Str := hex.EncodeToString(chunkMD5[:])

		progress := float64(i+1) / float64(totalChunks) * 100
		fmt.Printf("  Chunk %d/%d (%.1f%%)\n", i+1, totalChunks, progress)

		if err := uploadChunk(token, uploadID, i, chunkMD5Str, chunkData); err != nil {
			return "", fmt.Errorf("chunk %d upload failed: %v", i, err)
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

	var completeResp CompleteUploadResponse
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
	baseName := filepath.Base(filePath)
	sceneCode := fmt.Sprintf("test_%d", time.Now().UnixNano())
	title := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	body := map[string]interface{}{
		"space_id":      spaceID,
		"title":         title,
		"scene_code":    sceneCode,
		"file_id":       fileID,
		"panorama_type": "equirectangular",
		"initial_fov":   90,
		"initial_pitch": 0,
		"initial_yaw":   0,
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

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Create scene response: %s\n", string(respBody))

	var sceneResp SceneResponse
	json.Unmarshal(respBody, &sceneResp)

	return sceneResp.Data.SceneID, nil
}

func getSceneStatus(token string, sceneID uint) string {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/resource/scenes/%d", BASE_URL, sceneID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "error"
	}
	defer resp.Body.Close()

	var detailResp SceneDetailResponse
	json.NewDecoder(resp.Body).Decode(&detailResp)
	return detailResp.Data.SliceStatus
}
