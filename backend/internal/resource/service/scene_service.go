package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"webGL-720yun/internal/model"
	"webGL-720yun/internal/resource/dto"
	"webGL-720yun/internal/resource/repository"
	sliceService "webGL-720yun/internal/slice/service"
	"webGL-720yun/pkg/image"
	"webGL-720yun/pkg/logger"
	"webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/redis"
	"webGL-720yun/pkg/utils"
)

type SceneService struct {
	repo           *repository.SceneRepository
	spaceRepo      *repository.SpaceRepository
	minioClient    *minio_client.MinIOClient
	imageProcessor *image.Processor
	sliceQueue     *sliceService.SliceQueue
	redisService   *redis.RedisService
}

func NewSceneService(db *gorm.DB, minioClient *minio_client.MinIOClient, sliceQueue *sliceService.SliceQueue, redisService *redis.RedisService) *SceneService {
	return &SceneService{
		repo:           repository.NewSceneRepository(db),
		spaceRepo:      repository.NewSpaceRepository(db),
		minioClient:    minioClient,
		imageProcessor: image.NewProcessor(),
		sliceQueue:     sliceQueue,
		redisService:   redisService,
	}
}

func (s *SceneService) checkSpacePermission(spaceID uint, userID uint, isAdmin bool) (*model.ResSpace, error) {
	space, err := s.spaceRepo.FindByID(spaceID)
	if err != nil {
		return nil, fmt.Errorf("空间不存在: %w", err)
	}

	if !isAdmin && space.CreatedBy != userID {
		return nil, errors.New("无权限在此空间操作")
	}
	return space, nil
}

func (s *SceneService) validateAndPrepareSceneCode(sceneCode, title string, excludeID uint) (string, error) {
	code, err := s.generateSceneCode(sceneCode, title)
	if err != nil {
		return "", err
	}

	exists, err := s.repo.CheckSceneCodeExists(code, excludeID)
	if err != nil {
		return "", fmt.Errorf("检查场景编码失败: %w", err)
	}
	if exists {
		return "", errors.New("场景编码已存在")
	}

	return code, nil
}

func (s *SceneService) CreateScene(req *dto.CreateSceneRequest, userID uint, isAdmin bool) (*dto.CreateSceneResponse, error) {
	space, err := s.checkSpacePermission(req.SpaceID, userID, isAdmin)
	if err != nil {
		return nil, err
	}

	sceneCode, err := s.validateAndPrepareSceneCode(req.SceneCode, req.Title, 0)
	if err != nil {
		return nil, err
	}

	scene := &model.ResScene{
		SpaceID:      req.SpaceID,
		Title:        req.Title,
		SceneCode:    sceneCode,
		FileID:       req.FileID,
		PanoramaType: req.PanoramaType,
		InitialFOV:   req.InitialFOV,
		InitialPitch: req.InitialPitch,
		InitialYaw:   req.InitialYaw,
		NorthOffset:  req.NorthOffset,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		SortOrder:    req.SortOrder,
		Status:       1,
		SliceStatus:  model.SliceStatusPending,
	}

	if scene.InitialFOV == 0 {
		scene.InitialFOV = 100
	}
	if scene.PanoramaType == "" {
		scene.PanoramaType = "equirectangular"
	}

	// 先创建场景到数据库
	if err := s.repo.Create(scene); err != nil {
		return nil, fmt.Errorf("创建场景失败: %w", err)
	}

	if req.FileID != "" {
		fileInfo, err := s.redisService.GetFileInfo(req.FileID)
		if err != nil {
			return nil, fmt.Errorf("获取文件信息失败: %w", err)
		}
		if fileInfo == nil {
			return nil, errors.New("文件不存在，请先上传")
		}

		sourceURL, _ := fileInfo["source_url"].(string)
		if sourceURL == "" {
			return nil, errors.New("文件上传未完成，请先完成上传")
		}

		thumbURL, _ := fileInfo["thumb_url"].(string)
		var fileSize int64
		if fs, ok := fileInfo["file_size"].(float64); ok {
			fileSize = int64(fs)
		}

		scene.SourceURL = sourceURL
		scene.ThumbnailURL = thumbURL
		scene.SourceFileSize = fileSize

		if w, ok := fileInfo["width"].(float64); ok {
			scene.SourceWidth = int(w)
		}
		if h, ok := fileInfo["height"].(float64); ok {
			scene.SourceHeight = int(h)
		}

		taskID := uuid.New().String()
		scene.TaskID = taskID

		if s.sliceQueue != nil {
			task := &sliceService.SliceTask{
				TaskID:    taskID,
				SceneID:   scene.ID,
				SceneCode: scene.SceneCode,
				FileID:    scene.FileID,
				SpaceName: space.Name,
				SpaceSlug: space.Slug,
				UserID:    userID,
				CreatedAt: time.Now().Unix(),
			}

			if err := s.sliceQueue.PushTask(task); err != nil {
				scene.SliceStatus = model.SliceStatusFailed
				s.repo.Update(scene)
				return nil, fmt.Errorf("推送切片任务失败: %w", err)
			}

			// 保持为 Pending 状态，直到 Worker 开始领任务并更新为 Slicing
			if err := s.repo.Update(scene); err != nil {
				return nil, fmt.Errorf("更新场景状态失败: %w", err)
			}
		}
	}

	return &dto.CreateSceneResponse{
		SceneID:     scene.ID,
		Title:       scene.Title,
		SceneCode:   scene.SceneCode,
		SliceStatus: scene.SliceStatus,
		TaskID:      scene.TaskID,
	}, nil
}

func (s *SceneService) UpdateScene(id uint, req *dto.UpdateSceneRequest, userID uint, isAdmin bool) (*model.ResScene, error) {
	scene, err := s.repo.FindByIDWithSpace(id)
	if err != nil {
		return nil, err
	}

	if !isAdmin && scene.Space.CreatedBy != userID {
		return nil, errors.New("无权限修改此场景")
	}

	if req.Title != "" {
		scene.Title = req.Title
	}
	if req.PanoramaType != "" {
		scene.PanoramaType = req.PanoramaType
	}
	if req.InitialFOV != 0 {
		scene.InitialFOV = req.InitialFOV
	}
	if req.InitialPitch != 0 {
		scene.InitialPitch = req.InitialPitch
	}
	if req.InitialYaw != 0 {
		scene.InitialYaw = req.InitialYaw
	}
	if req.NorthOffset != 0 {
		scene.NorthOffset = req.NorthOffset
	}
	if req.Longitude != 0 {
		scene.Longitude = req.Longitude
	}
	if req.Latitude != 0 {
		scene.Latitude = req.Latitude
	}
	if req.SortOrder != 0 {
		scene.SortOrder = req.SortOrder
	}
	if req.Status != nil {
		scene.Status = *req.Status
	}

	shouldTriggerSlice := false
	if req.FileID != "" && scene.FileID == "" {
		fileInfo, err := s.redisService.GetFileInfo(req.FileID)
		if err != nil {
			logger.Warnf("获取文件信息失败: %v", err)
		} else if fileInfo != nil {
			scene.FileID = req.FileID
			if sourceURL, ok := fileInfo["source_url"].(string); ok {
				scene.SourceURL = sourceURL
			}
			if thumbURL, ok := fileInfo["thumb_url"].(string); ok {
				scene.ThumbnailURL = thumbURL
			}
			if fileSize, ok := fileInfo["file_size"].(float64); ok {
				scene.SourceFileSize = int64(fileSize)
			}
			if width, ok := fileInfo["width"].(float64); ok {
				scene.SourceWidth = int(width)
			}
			if height, ok := fileInfo["height"].(float64); ok {
				scene.SourceHeight = int(height)
			}
			shouldTriggerSlice = true
		}
	}

	if err := s.repo.Update(scene); err != nil {
		return nil, fmt.Errorf("更新场景失败: %w", err)
	}

	// 清理元数据缓存，确保瓦片请求能获取到最新路径
	s.redisService.DeleteCachedSceneMeta(scene.SceneCode)

	if shouldTriggerSlice && s.sliceQueue != nil {
		if scene.SourceURL == "" {
			logger.Warnf("⚠️  场景 [%s] 的 SourceURL 为空，跳过切片任务", scene.SceneCode)
		} else {
			taskID := uuid.New().String()
			scene.TaskID = taskID
			s.repo.Update(scene)

			task := &sliceService.SliceTask{
				TaskID:    taskID,
				SceneID:   scene.ID,
				SceneCode: scene.SceneCode,
				FileID:    scene.FileID,
				SpaceName: scene.Space.Name,
				SpaceSlug: scene.Space.Slug,
				UserID:    userID,
				CreatedAt: time.Now().Unix(),
			}

			if err := s.sliceQueue.PushTask(task); err != nil {
				logger.Warnf("推送切片任务失败: %v", err)
			} else {
				scene.SliceStatus = model.SliceStatusPending
				s.repo.Update(scene)
			}
		}
	}

	return scene, nil
}

func (s *SceneService) DeleteScene(id uint, userID uint, isAdmin bool) error {
	scene, err := s.repo.FindByIDWithSpace(id)
	if err != nil {
		return err
	}

	if !isAdmin && scene.Space.CreatedBy != userID {
		return errors.New("无权限删除此场景")
	}

	err = s.repo.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.DeleteHotspotsBySceneID(id); err != nil {
			return fmt.Errorf("删除热点失败: %w", err)
		}

		if err := s.repo.Delete(id); err != nil {
			return fmt.Errorf("删除场景失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 清理元数据缓存
	s.redisService.DeleteCachedSceneMeta(scene.SceneCode)

	if scene.SourceURL != "" {
		sourceObjectName := s.extractObjectName(scene.SourceURL)
		if sourceObjectName != "" {
			_ = s.minioClient.DeleteObject(sourceObjectName)
		}
	}

	if scene.ThumbnailURL != "" {
		thumbObjectName := s.extractObjectName(scene.ThumbnailURL)
		if thumbObjectName != "" {
			_ = s.minioClient.DeleteObject(thumbObjectName)
		}
	}

	if scene.FileID != "" {
		if err := s.redisService.DeleteFileInfo(scene.FileID); err != nil {
			logger.Warnf("清理Redis文件信息失败: %v", err)
		}
	}

	if scene.SourceFileMD5 != "" {
		if err := s.redisService.DeleteFileMD5(scene.SourceFileMD5); err != nil {
			logger.Warnf("清理Redis MD5缓存失败: %v", err)
		}
	}

	return nil
}

func (s *SceneService) GetSceneList(req *dto.SceneListRequest) (*dto.SceneListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	query := &repository.SceneListQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		SpaceID:  req.SpaceID,
		Title:    req.Title,
		Status:   req.Status,
		Keyword:  req.Keyword,
	}

	scenes, total, err := s.repo.GetList(query)
	if err != nil {
		return nil, fmt.Errorf("获取场景列表失败: %w", err)
	}

	items := make([]*dto.SceneListItem, len(scenes))
	for i, scene := range scenes {
		items[i] = dto.ToSceneListItem(scene)
	}

	return &dto.SceneListResponse{
		PageInfo: utils.PageInfo{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
			Pages:    int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
		},
		Scenes: items,
	}, nil
}

func (s *SceneService) GetSceneDetail(id uint) (*dto.SceneDetailResponse, error) {
	scene, err := s.repo.FindByIDWithHotspots(id)
	if err != nil {
		return nil, err
	}

	return dto.ToSceneDetailResponse(scene), nil
}

func (s *SceneService) BatchImport(spaceID uint, file *multipart.FileHeader, userID uint, isAdmin bool) (*dto.BatchImportResponse, error) {
	_, err := s.checkSpacePermission(spaceID, userID, isAdmin)
	if err != nil {
		return nil, err
	}

	var items []dto.BatchImportItem

	ext := filepath.Ext(file.Filename)
	if ext == ".xlsx" || ext == ".xls" {
		items, err = s.parseExcelFile(file)
	} else if ext == ".json" {
		items, err = s.parseJSONFile(file)
	} else {
		return nil, errors.New("不支持的文件格式，仅支持.xlsx、.xls和.json文件")
	}

	if err != nil {
		return nil, fmt.Errorf("解析文件失败: %w", err)
	}

	response := &dto.BatchImportResponse{
		Results: make([]dto.BatchImportResult, len(items)),
	}

	for i, item := range items {
		// 使用统一的验证和准备逻辑
		sceneCode, err := s.validateAndPrepareSceneCode(item.SceneCode, item.Title, 0)
		if err != nil {
			response.Results[i] = dto.BatchImportResult{
				Index:     i + 1,
				SceneCode: item.SceneCode,
				Title:     item.Title,
				Status:    "failed",
				Message:   err.Error(),
			}
			response.FailedCount++
			response.Errors = append(response.Errors, dto.BatchImportError{
				Index:     i + 1,
				SceneCode: item.SceneCode,
				Error:     err.Error(),
			})
			continue
		}

		result := dto.BatchImportResult{
			Index:     i + 1,
			SceneCode: sceneCode,
			Title:     item.Title,
		}

		scene := &model.ResScene{
			SpaceID:   spaceID,
			Title:     item.Title,
			SceneCode: sceneCode,
			Longitude: item.Longitude,
			Latitude:  item.Latitude,
			Status:    1,
		}

		if err := s.repo.Create(scene); err != nil {
			result.Status = "failed"
			result.Message = "创建失败"
			response.FailedCount++
			response.Errors = append(response.Errors, dto.BatchImportError{
				Index:     i + 1,
				SceneCode: sceneCode,
				Error:     err.Error(),
			})
		} else {
			result.Status = "success"
			result.Message = "创建成功"
			response.SuccessCount++
		}

		response.Results[i] = result
	}

	return response, nil
}

func (s *SceneService) parseExcelFile(file *multipart.FileHeader) ([]dto.BatchImportItem, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, errors.New("Excel文件至少需要包含表头和一行数据")
	}

	var items []dto.BatchImportItem
	for i, row := range rows[1:] {
		if len(row) < 2 {
			continue
		}

		item := dto.BatchImportItem{
			Title:     row[0],
			SceneCode: row[1],
		}

		if len(row) > 2 {
			item.FileName = row[2]
		}
		if len(row) > 3 && row[3] != "" {
			fmt.Sscanf(row[3], "%f", &item.Longitude)
		}
		if len(row) > 4 && row[4] != "" {
			fmt.Sscanf(row[4], "%f", &item.Latitude)
		}

		items = append(items, item)
		_ = i
	}

	return items, nil
}

func (s *SceneService) parseJSONFile(file *multipart.FileHeader) ([]dto.BatchImportItem, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	var items []dto.BatchImportItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *SceneService) extractObjectName(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 4 {
		return strings.Join(parts[3:], "/")
	}
	return ""
}

func (s *SceneService) GetSpaceGraphData(spaceID uint) (*dto.GraphDataResponse, error) {
	space, err := s.spaceRepo.FindByID(spaceID)
	if err != nil {
		return nil, fmt.Errorf("空间不存在: %w", err)
	}

	scenes, err := s.repo.FindBySpaceIDWithHotspots(spaceID)
	if err != nil {
		return nil, fmt.Errorf("获取场景列表失败: %w", err)
	}

	var nodes []*dto.SceneNodeData
	var edges []*dto.EdgeData

	sceneMap := make(map[uint]bool)
	for _, scene := range scenes {
		sceneMap[scene.ID] = true
		nodeData := dto.ToSceneNodeData(scene)
		nodes = append(nodes, nodeData)

		for _, hotspot := range scene.Hotspots {
			if hotspot.Type == model.HotspotTypeSwitch && hotspot.TargetSceneID != nil {
				if sceneMap[*hotspot.TargetSceneID] {
					edges = append(edges, dto.ToEdgeData(&hotspot))
				}
			}
		}
	}

	return &dto.GraphDataResponse{
		SpaceInfo: dto.ToSpaceInfoForGraph(space),
		Nodes:     nodes,
		Edges:     edges,
	}, nil
}

func (s *SceneService) UpdateScenePosition(id uint, longitude, latitude float64, userID uint, isAdmin bool) error {
	scene, err := s.repo.FindByIDWithSpace(id)
	if err != nil {
		return err
	}

	if !isAdmin && scene.Space.CreatedBy != userID {
		return errors.New("无权限修改此场景")
	}

	return s.repo.UpdatePosition(id, longitude, latitude)
}

func (s *SceneService) BatchUpdateScenePosition(positions []dto.ScenePosition, userID uint, isAdmin bool) error {
	if len(positions) == 0 {
		return nil
	}

	var updates []repository.ScenePositionUpdate
	for _, pos := range positions {
		scene, err := s.repo.FindByIDWithSpace(pos.SceneID)
		if err != nil {
			continue
		}

		if !isAdmin && scene.Space.CreatedBy != userID {
			continue
		}

		updates = append(updates, repository.ScenePositionUpdate{
			ID:        pos.SceneID,
			Longitude: pos.Longitude,
			Latitude:  pos.Latitude,
		})
	}

	if len(updates) == 0 {
		return errors.New("没有可更新的场景")
	}

	return s.repo.BatchUpdatePosition(updates)
}

func (s *SceneService) generateSceneCode(sceneCode, title string) (string, error) {
	// 如果提供了 SceneCode，清理并验证
	if sceneCode != "" {
		// 清理格式
		sceneCode = utils.SanitizeSceneCode(sceneCode)

		// 验证长度
		if len(sceneCode) < 2 {
			return "", errors.New("场景编码太短，至少需要 2 个字符")
		}
		if len(sceneCode) > 100 {
			return "", errors.New("场景编码太长，最多 100 个字符")
		}

		return sceneCode, nil
	}

	// 未提供 SceneCode，自动生成
	if title == "" {
		return "", errors.New("标题不能为空，无法生成场景编码")
	}

	// 生成拼音-UUID 格式
	return utils.GenerateSceneCode(title), nil
}

// ==================== 资源流式读取 ====================

// validFaces cubemap 六个面的合法名称
var validFaces = map[string]bool{
	"px": true, "nx": true,
	"py": true, "ny": true,
	"pz": true, "nz": true,
}

// IsValidFace 校验 cubemap 面名称是否合法
func IsValidFace(face string) bool {
	return validFaces[face]
}

// ResourceStreamResult 资源流式读取的结果
type ResourceStreamResult struct {
	Stream      io.ReadCloser
	Size        int64
	ETag        string
	ContentType string
}

// getSpaceSlug 通过 sceneCode 获取 SpaceSlug，优先从缓存读取
func (s *SceneService) getSpaceSlug(sceneCode string) (string, error) {
	// 1. 尝试从 Redis 缓存读取
	meta, err := s.redisService.GetCachedSceneMeta(sceneCode)
	if err == nil && meta != nil {
		if slug, ok := meta["space_slug"].(string); ok {
			return slug, nil
		}
	}

	// 2. 缓存未命中，从 MySQL 读取
	scene, err := s.repo.FindBySceneCodeWithSpace(sceneCode)
	if err != nil {
		return "", fmt.Errorf("scene not found: %w", err)
	}

	spaceSlug := scene.Space.Slug

	// 3. 异步写入缓存 (过期时间 1 小时)
	go func() {
		s.redisService.CacheSceneMeta(sceneCode, map[string]string{
			"space_slug": spaceSlug,
		}, 1*time.Hour)
	}()

	return spaceSlug, nil
}

// GetTileStream 流式获取瓦片图片
func (s *SceneService) GetTileStream(ctx context.Context, sceneCode, face string, level, x, y int) (*ResourceStreamResult, error) {
	spaceName, err := s.getSpaceSlug(sceneCode)
	if err != nil {
		return nil, err
	}

	// 使用统一的路径生成函数，保持 x, y 顺序一致
	objectPath := model.GetSceneTilePath(spaceName, sceneCode, face, level, x, y)
	return s.getObjectStream(ctx, objectPath, "image/jpeg")
}

// GetPreviewStream 流式获取预览图
func (s *SceneService) GetPreviewStream(ctx context.Context, sceneCode string) (*ResourceStreamResult, error) {
	spaceName, err := s.getSpaceSlug(sceneCode)
	if err != nil {
		return nil, err
	}

	objectPath := model.GetScenePreviewPath(spaceName, sceneCode)
	return s.getObjectStream(ctx, objectPath, "image/jpeg")
}

// GetSourceStream 流式获取场景原始全景图
func (s *SceneService) GetSourceStream(ctx context.Context, sceneCode string) (*ResourceStreamResult, error) {
	spaceName, err := s.getSpaceSlug(sceneCode)
	if err != nil {
		return nil, err
	}

	objectPath := model.GetSceneSourcePath(spaceName, sceneCode)
	return s.getObjectStream(ctx, objectPath, "image/jpeg")
}

// GetCoverStream 流式获取空间封面
func (s *SceneService) GetCoverStream(ctx context.Context, spaceName string) (*ResourceStreamResult, error) {
	objectPath := model.GetSpaceCoverPath(spaceName)
	return s.getObjectStream(ctx, objectPath, "image/jpeg")
}

// getObjectStream 通用的 MinIO 对象流式获取
func (s *SceneService) getObjectStream(ctx context.Context, objectPath, contentType string) (*ResourceStreamResult, error) {
	obj, err := s.minioClient.GetObjectStream(ctx, objectPath)
	if err != nil {
		return nil, fmt.Errorf("get object failed: %w", err)
	}

	// Stat 获取对象元信息（大小、ETag 等）
	info, err := obj.Stat()
	if err != nil {
		obj.Close()
		return nil, fmt.Errorf("object not found: %w", err)
	}

	return &ResourceStreamResult{
		Stream:      obj,
		Size:        info.Size,
		ETag:        info.ETag,
		ContentType: contentType,
	}, nil
}
