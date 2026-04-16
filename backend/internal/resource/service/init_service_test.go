package service

import (
	"os"
	"path/filepath"
	"testing"

	"webGL-720yun/config"
	"webGL-720yun/internal/core/setup"
	"webGL-720yun/internal/model"
	miniocli "webGL-720yun/pkg/minio_client"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		dsn = "root:123456@tcp(127.0.0.1:3306)/720db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping test: cannot connect to MySQL: %v", err)
	}

	return db, func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}

func setupTestMinIO(t *testing.T) *miniocli.MinIOClient {
	t.Helper()

	cfg := &config.MinIOConfig{
		Endpoint:  "localhost:9000",
		AccessKey: "admin",
		SecretKey: "admin123456",
		Bucket:    "pannellum-resources",
		UseSSL:    false,
		Region:    "us-east-1",
	}

	minioClient, err := miniocli.NewMinIOClient(cfg)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to MinIO: %v", err)
	}

	return minioClient
}

func writeFileUTF16(path string, content []byte) error {
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	file, err := windows.CreateFile(p, windows.GENERIC_WRITE, 0, nil, windows.CREATE_ALWAYS, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(file)

	var written uint32
	return windows.WriteFile(file, content, &written, nil)
}

func createTestSeedFilesWithChineseNames(t *testing.T, basePath string, spot setup.ScenicSpotSeed) {
	t.Helper()

	spotDir := filepath.Join(basePath, spot.LocalPath)

	err := os.MkdirAll(spotDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	coverPath := filepath.Join(spotDir, spot.CoverFile)
	if err := writeFileUTF16(coverPath, []byte("fake_cover_data")); err != nil {
		t.Fatalf("Failed to create cover: %v", err)
	}

	for _, scene := range spot.Scenes {
		scenePath := filepath.Join(spotDir, scene.FileName)
		if err := writeFileUTF16(scenePath, []byte("fake_panorama_data")); err != nil {
			t.Fatalf("Failed to create scene: %v", err)
		}
	}
}

func init() {
	config.Init()
}

func TestResourceInitService_RealMySQLAndMinIO(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	minioClient := setupTestMinIO(t)

	testSpot := setup.GetScenicSpotSeeds()[0]
	tempDir := t.TempDir()
	createTestSeedFilesWithChineseNames(t, tempDir, testSpot)

<<<<<<< HEAD
	svc := NewResourceInitService(db, minioClient, tempDir)
=======
	svc := NewResourceInitService(db, minioClient, nil, nil, tempDir)
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503

	err := svc.SeedResourcesIfNeeded()
	assert.NoError(t, err)

	var space model.ResSpace
	err = db.Where("name = ?", testSpot.Name).First(&space).Error
	if assert.NoError(t, err) {
		assert.Equal(t, testSpot.Name, space.Name)
	}

	var sceneCount int64
	db.Model(&model.ResScene{}).Where("space_id = ?", space.ID).Count(&sceneCount)
	assert.Equal(t, int64(len(testSpot.Scenes)), sceneCount, "应该创建所有场景")

	var scene model.ResScene
	err = db.Where("scene_code = ?", testSpot.Scenes[0].SceneCode).First(&scene).Error
	if assert.NoError(t, err) {
		assert.Equal(t, testSpot.Scenes[0].Title, scene.Title)
		assert.NotEmpty(t, scene.SourceURL)
		assert.Equal(t, "equirectangular", scene.PanoramaType)
	}
}
