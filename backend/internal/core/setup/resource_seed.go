package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"webGL-720yun/internal/model"
	"webGL-720yun/pkg/logger"
	"webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/utils"

	"gorm.io/gorm"
)

// SeedResourcesIfNeeded 初始化全景资源（开发/测试用）
// 返回新创建的场景列表，用于后续处理（如切片）
func SeedResourcesIfNeeded(db *gorm.DB, minioClient *minio_client.MinIOClient, seedBasePath string) ([]model.ResScene, error) {
	logger.Info("🚀 开始初始化全景资源...")

	var createdScenes []model.ResScene
	scenicSpots := GetScenicSpotSeeds()

	// 先尝试为现有的、Slug为空的记录补全 Slug
	for _, spot := range scenicSpots {
		db.Model(&model.ResSpace{}).
			Where("name = ? AND (slug = '' OR slug IS NULL)", spot.Name).
			Update("slug", spot.Slug)
	}

	for _, spot := range scenicSpots {
		scenes, err := seedScenicSpot(db, minioClient, seedBasePath, spot)
		if err != nil {
			logger.Errorf("初始化景区 %s 失败: %v", spot.Name, err)
			continue
		}
		createdScenes = append(createdScenes, scenes...)
	}

	logger.Info("✅ 全景资源初始化完成")
	return createdScenes, nil
}

func seedScenicSpot(db *gorm.DB, minioClient *minio_client.MinIOClient, seedBasePath string, spot ScenicSpotSeed) ([]model.ResScene, error) {
	var existingSpot model.ResSpace
	var createdScenes []model.ResScene

	err := db.Where("name = ?", spot.Name).First(&existingSpot).Error
	if err == nil {
		logger.Infof("⏭️  景区 [%s] 已存在，同步场景", spot.Name)
		return syncScenes(db, minioClient, seedBasePath, existingSpot.ID, existingSpot.Slug, spot)
	}

	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询景区失败: %w", err)
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

	if err := db.Create(&newSpot).Error; err != nil {
		return nil, fmt.Errorf("创建景区记录失败: %w", err)
	}
	logger.Infof("✅ 创建景区: %s (ID=%d)", newSpot.Name, newSpot.ID)

	// 上传封面
	coverLocalPath := filepath.Join(seedBasePath, spot.LocalPath, spot.CoverFile)
	coverURL, err := uploadCover(minioClient, newSpot.Slug, coverLocalPath)
	if err != nil {
		logger.Warnf("上传封面失败: %v", err)
	} else if coverURL != "" {
		newSpot.CoverURL = coverURL
		db.Save(&newSpot)
	}

	for _, sceneSeed := range spot.Scenes {
		scene, err := seedScene(db, minioClient, seedBasePath, newSpot.ID, newSpot.Slug, spot.LocalPath, sceneSeed)
		if err != nil {
			logger.Warnf("导入场景 %s 失败: %v", sceneSeed.Title, err)
		} else if scene.ID > 0 {
			createdScenes = append(createdScenes, *scene)
		}
	}

	return createdScenes, nil
}

func syncScenes(db *gorm.DB, minioClient *minio_client.MinIOClient, seedBasePath string, spaceID uint, spaceName string, spot ScenicSpotSeed) ([]model.ResScene, error) {
	var createdScenes []model.ResScene
	for _, sceneSeed := range spot.Scenes {
		var existingScene model.ResScene
		err := db.Where("scene_code = ? AND space_id = ?", sceneSeed.SceneCode, spaceID).First(&existingScene).Error

		if err == gorm.ErrRecordNotFound {
			scene, seedErr := seedScene(db, minioClient, seedBasePath, spaceID, spaceName, spot.LocalPath, sceneSeed)
			if seedErr != nil {
				logger.Warnf("同步场景 %s 失败: %v", sceneSeed.Title, seedErr)
			} else if scene.ID > 0 {
				createdScenes = append(createdScenes, *scene)
			}
		} else if err == nil {
			// 如果场景已存在但坐标为0，尝试补全坐标
			if existingScene.Longitude == 0 && existingScene.Latitude == 0 {
				if lon, lat, found := GetSceneCoordinate(sceneSeed.Title); found {
					db.Model(&existingScene).Updates(map[string]interface{}{
						"longitude": lon,
						"latitude":  lat,
					})
					logger.Infof("📍 补全场景坐标: %s (%f, %f)", sceneSeed.Title, lon, lat)
				}
			}
		} else {
			logger.Warnf("检查场景 %s 时出错: %v", sceneSeed.Title, err)
		}
	}
	return createdScenes, nil
}

func seedScene(db *gorm.DB, minioClient *minio_client.MinIOClient, seedBasePath string, spaceID uint, spaceName, localPath string, scene SceneSeed) (*model.ResScene, error) {
	var existingScene model.ResScene

	err := db.Where("scene_code = ?", scene.SceneCode).First(&existingScene).Error
	if err == nil {
		logger.Infof("⏭️  场景 [%s] 已存在，跳过", scene.Title)
		return &existingScene, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询场景失败: %w", err)
	}

	localFilePath := filepath.Join(seedBasePath, localPath, scene.FileName)

	// 检查本地文件是否存在
	if _, err := os.Stat(localFilePath); err != nil {
		if os.IsNotExist(err) {
			logger.Warnf("⚠️  本地文件不存在，跳过: %s", localFilePath)
			return nil, nil
		}
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 直接上传源文件到 MinIO
	minIOObjectPath := model.GetSceneSourcePath(spaceName, scene.SceneCode)
	panoramaURL, uploadErr := minioClient.UploadFile(
		minIOObjectPath,
		localFilePath,
		"image/jpeg",
	)
	if uploadErr != nil {
		return nil, fmt.Errorf("上传全景图到MinIO失败: %w", uploadErr)
	}
	logger.Infof("☁️  已上传: %s → MinIO", filepath.Base(localFilePath))

	newScene := model.ResScene{
		SpaceID:      spaceID,
		Title:        scene.Title,
		SceneCode:    scene.SceneCode,
		PanoramaType: "equirectangular",
		FileID:       scene.SceneCode,
		SourceURL:    panoramaURL,
		Status:       1,
		InitialFOV:   100,
	}

	// 设置坐标
	if lon, lat, found := GetSceneCoordinate(scene.Title); found {
		newScene.Longitude = lon
		newScene.Latitude = lat
	}

	if createErr := db.Create(&newScene).Error; createErr != nil {
		return nil, fmt.Errorf("创建场景记录失败: %w", createErr)
	}

	// 保存 MD5 映射
	md5Str, _ := utils.CalculateFileMD5(localFilePath)
	if md5Str != "" {
		logger.Infof("📦 文件 MD5: %s", md5Str)
	}

	logger.Infof("✅ 创建场景: %s (ID=%d)", newScene.Title, newScene.ID)
	return &newScene, nil
}

// uploadCover 上传景区封面到 MinIO
func uploadCover(minioClient *minio_client.MinIOClient, spaceSlug, coverLocalPath string) (string, error) {
	if _, err := os.Stat(coverLocalPath); err != nil {
		if os.IsNotExist(err) {
			logger.Warnf("封面文件不存在: %s", coverLocalPath)
			return "", nil
		}
		return "", fmt.Errorf("获取封面文件信息失败: %w", err)
	}

	coverMinIOPath := model.GetSpaceCoverPath(spaceSlug)

	coverURL, uploadErr := minioClient.UploadFile(
		coverMinIOPath,
		coverLocalPath,
		"image/jpeg",
	)
	if uploadErr != nil {
		return "", fmt.Errorf("上传封面失败: %w", uploadErr)
	}

	logger.Infof("☁️  上传封面: %s → MinIO", filepath.Base(coverLocalPath))
	return coverURL, nil
}
