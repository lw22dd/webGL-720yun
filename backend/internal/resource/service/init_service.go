package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"webGL-720yun/internal/core/setup"
	"webGL-720yun/internal/model"
	"webGL-720yun/internal/resource/upload"
	"webGL-720yun/internal/slice"
	"webGL-720yun/pkg/logger"
	miniocli "webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/utils"

	"gorm.io/gorm"
)

type ResourceInitService struct {
	db           *gorm.DB
	minioClient  miniocli.MinIOOperations
	sliceQueue   *slice.SliceQueue
	uploadRepo   *upload.UploadRepository
	seedBasePath string
}

func NewResourceInitService(db *gorm.DB, minioCli miniocli.MinIOOperations, sliceQueue *slice.SliceQueue, uploadRepo *upload.UploadRepository, seedBasePath string) *ResourceInitService {
	return &ResourceInitService{
		db:           db,
		minioClient:  minioCli,
		sliceQueue:   sliceQueue,
		uploadRepo:   uploadRepo,
		seedBasePath: seedBasePath,
	}
}

func (s *ResourceInitService) SeedResourcesIfNeeded() error {
	logger.Info("🚀 开始检查全景资源初始化...")

	scenicSpots := setup.GetScenicSpotSeeds()

	// 先尝试为现有的、Slug为空的记录补全 Slug
	for _, spot := range scenicSpots {
		s.db.Model(&model.ResSpace{}).
			Where("name = ? AND (slug = '' OR slug IS NULL)", spot.Name).
			Update("slug", spot.Slug)
	}

	for _, spot := range scenicSpots {
		if err := s.seedScenicSpotIfNotExists(spot); err != nil {
			logger.Errorf("初始化景区 %s 失败: %v", spot.Name, err)
			continue
		}
	}

	logger.Info("✅ 全景资源初始化完成")
	return nil
}

func (s *ResourceInitService) seedScenicSpotIfNotExists(spot setup.ScenicSpotSeed) error {
	var existingSpot model.ResSpace

	err := s.db.Where("name = ?", spot.Name).First(&existingSpot).Error
	if err == nil {
		logger.Infof("⏭️  景区 [%s] 已存在，跳过创建", spot.Name)
		return s.checkAndSyncScenes(existingSpot.ID, existingSpot.Slug, spot)
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("查询景区失败: %w", err)
	}

	newSpot := model.ResSpace{
		Name:        spot.Name,
		Slug:        spot.Slug,
		Description: spot.Description,
		Province:    spot.Province,
		City:        spot.City,
		Longitude:   spot.Longitude,
		Latitude:    spot.Latitude,
		Status:      1,
	}

	if err := s.db.Create(&newSpot).Error; err != nil {
		return fmt.Errorf("创建景区记录失败: %w", err)
	}
	logger.Infof("✅ 创建景区: %s (ID=%d)", newSpot.Name, newSpot.ID)

	coverLocalPath := filepath.Join(s.seedBasePath, spot.LocalPath, spot.CoverFile)
	if _, err := os.Stat(coverLocalPath); err == nil {
		coverMinIOPath := fmt.Sprintf("spaces/%s/covers/cover.jpg", newSpot.Slug)

		exists, _ := s.minioClient.ObjectExists(coverMinIOPath)
		if !exists {
			coverURL, uploadErr := s.minioClient.UploadFile(
				coverMinIOPath,
				coverLocalPath,
				"image/jpeg",
			)
			if uploadErr != nil {
				logger.Warnf("上传封面失败: %v", uploadErr)
			} else {
				newSpot.CoverURL = coverURL
				s.db.Save(&newSpot)
				logger.Infof("☁️  上传封面: %s → MinIO", spot.CoverFile)
			}
		}
	} else {
		logger.Warnf("封面文件不存在: %s", coverLocalPath)
	}

	for _, sceneSeed := range spot.Scenes {
		if err := s.seedSceneIfNotExists(newSpot.ID, newSpot.Slug, spot.LocalPath, sceneSeed); err != nil {
			logger.Warnf("导入场景 %s 失败: %v", sceneSeed.Title, err)
		}
	}

	return nil
}

func (s *ResourceInitService) checkAndSyncScenes(spaceID uint, spaceName string, spot setup.ScenicSpotSeed) error {
	for _, sceneSeed := range spot.Scenes {
		var existingScene model.ResScene
		err := s.db.Where("scene_code = ? AND space_id = ?", sceneSeed.SceneCode, spaceID).First(&existingScene).Error

		if err == gorm.ErrRecordNotFound {
			if seedErr := s.seedSceneIfNotExists(spaceID, spaceName, spot.LocalPath, sceneSeed); seedErr != nil {
				logger.Warnf("同步场景 %s 失败: %v", sceneSeed.Title, seedErr)
			}
		} else if err != nil {
			logger.Warnf("检查场景 %s 时出错: %v", sceneSeed.Title, err)
		}
	}
	return nil
}

func (s *ResourceInitService) seedSceneIfNotExists(spaceID uint, spaceName, localPath string, scene setup.SceneSeed) error {
	var existingScene model.ResScene

	err := s.db.Where("scene_code = ?", scene.SceneCode).First(&existingScene).Error
	if err == nil {
		logger.Infof("⏭️  场景 [%s] 已存在，跳过", scene.Title)
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("查询场景失败: %w", err)
	}

	localFilePath := filepath.Join(s.seedBasePath, localPath, scene.FileName)
	// 使用拼音语义化的 SceneCode 构建 MinIO 路径
	minIOObjectPath := fmt.Sprintf("spaces/%s/sources/%s/source.jpg", spaceName, scene.SceneCode)

	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warnf("⚠️  本地文件不存在，跳过: %s", localFilePath)
			return nil
		}
		return fmt.Errorf("获取文件信息失败: %w", err)
	}

	md5Str, md5Err := utils.CalculateFileMD5(localFilePath)
	if md5Err != nil {
		logger.Warnf("⚠️  计算本地文件 MD5 失败，将无法使用秒传: %v", md5Err)
	}

	var panoramaURL string
	var resolvedFileID string

	if md5Str != "" && s.uploadRepo != nil {
		// === 双层校验: Redis -> MySQL ===
		fileID, _ := s.uploadRepo.GetFileIDByMD5(md5Str)
		if fileID != "" {
			fileInfo, _ := s.uploadRepo.GetFileInfo(fileID)
			if fileInfo != nil && fileInfo.SourceURL != "" {
				panoramaURL = fileInfo.SourceURL
				resolvedFileID = fileInfo.FileID
				logger.Infof("♻️  命中 MD5 秒传 [%s]，复用链接", scene.FileName)
			}
		}
	}

	if panoramaURL == "" {
		exists, checkErr := s.minioClient.ObjectExists(minIOObjectPath)
		if checkErr != nil {
			return fmt.Errorf("检查MinIO对象存在性失败: %w", checkErr)
		}

		if !exists {
			url, uploadErr := s.minioClient.UploadFile(
				minIOObjectPath,
				localFilePath,
				"image/jpeg",
			)
			if uploadErr != nil {
				return fmt.Errorf("上传全景图到MinIO失败: %w", uploadErr)
			}
			panoramaURL = url
			logger.Infof("☁️  已上传: %s → MinIO (%.2f MB)", scene.FileName, float64(fileInfo.Size())/1024/1024)
		} else {
			panoramaURL = s.minioClient.GetObjectURL(minIOObjectPath)
			logger.Infof("☁️  MinIO中已存在对象: %s，直接应用", scene.FileName)
		}
		
		// 新上传的文件我们使用 UUID 或 SceneCode 做为 FileID
		resolvedFileID = scene.SceneCode
	}

	newScene := model.ResScene{
		SpaceID:        spaceID,
		Title:          scene.Title,
		SceneCode:      scene.SceneCode,
		PanoramaType:   "equirectangular",
		FileID:         resolvedFileID,
		SourceURL:      panoramaURL,
		SourceFileSize: fileInfo.Size(),
		SourceFileMD5:  md5Str,
		Status:         1,
	}

	if createErr := s.db.Create(&newScene).Error; createErr != nil {
		return fmt.Errorf("创建场景记录失败: %w", createErr)
	}

	if md5Str != "" && s.uploadRepo != nil {
		_ = s.uploadRepo.SaveFileMD5(md5Str, newScene.FileID)
	}

	// 推送切片任务
	if s.sliceQueue != nil {
		task := &slice.SliceTask{
			TaskID:    fmt.Sprintf("seed_%s_%s", spaceName, scene.SceneCode),
			SceneID:   newScene.ID,
			SceneCode: newScene.SceneCode,
			FileID:    newScene.SceneCode,
			SpaceName: spaceName, // During seed, we use slug as name for simplicity or we could fetch the name.
			SpaceSlug: spaceName,
			UserID:    0,
			CreatedAt: time.Now().Unix(),
		}
		if err := s.sliceQueue.PushTask(task); err != nil {
			logger.Warnf("推送种子场景切片任务失败: %v", err)
		} else {
			logger.Infof("🚀 已推送场景切片任务: %s", scene.Title)
		}
	}

	logger.Infof("✅ 创建场景: %s (ID=%d)", newScene.Title, newScene.ID)
	return nil
}
