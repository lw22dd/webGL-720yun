package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"webGL-720yun/internal/model"
	"webGL-720yun/internal/resource/dto"
	"webGL-720yun/internal/resource/repository"
	"webGL-720yun/pkg/image"
	"webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/utils"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type SceneService struct {
	repo           *repository.SceneRepository
	spaceRepo      *repository.SpaceRepository
	minioClient    *minio_client.MinIOClient
	imageProcessor *image.Processor
}

func NewSceneService(db *gorm.DB, minioClient *minio_client.MinIOClient) *SceneService {
	return &SceneService{
		repo:           repository.NewSceneRepository(db),
		spaceRepo:      repository.NewSpaceRepository(db),
		minioClient:    minioClient,
		imageProcessor: image.NewProcessor(),
	}
}

func (s *SceneService) CreateScene(req *dto.CreateSceneRequest, userID uint, isAdmin bool, panoramaFile *multipart.FileHeader) (*model.ResScene, error) {
	space, err := s.spaceRepo.FindByID(req.SpaceID)
	if err != nil {
		return nil, fmt.Errorf("空间不存在: %w", err)
	}

	if !isAdmin && space.CreatedBy != userID {
		return nil, errors.New("无权限在此空间创建场景")
	}

	exists, err := s.repo.CheckSceneCodeExists(req.SceneCode, 0)
	if err != nil {
		return nil, fmt.Errorf("检查场景编码失败: %w", err)
	}
	if exists {
		return nil, errors.New("场景编码已存在")
	}

	scene := &model.ResScene{
		SpaceID:      req.SpaceID,
		Title:        req.Title,
		SceneCode:    req.SceneCode,
		PanoramaType: req.PanoramaType,
		InitialFOV:   req.InitialFOV,
		InitialPitch: req.InitialPitch,
		InitialYaw:   req.InitialYaw,
		NorthOffset:  req.NorthOffset,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		SortOrder:    req.SortOrder,
		Status:       1,
	}

	if scene.InitialFOV == 0 {
		scene.InitialFOV = 100
	}
	if scene.PanoramaType == "" {
		scene.PanoramaType = "equirectangular"
	}

	if panoramaFile != nil {
		if err := s.processPanoramaFile(scene, panoramaFile, space.Name); err != nil {
			return nil, fmt.Errorf("处理全景图失败: %w", err)
		}
	}

	if err := s.repo.Create(scene); err != nil {
		return nil, fmt.Errorf("创建场景失败: %w", err)
	}

	return scene, nil
}

func (s *SceneService) processPanoramaFile(scene *model.ResScene, file *multipart.FileHeader, spaceName string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("panorama_%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename)))
	dst, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer dst.Close()
	defer os.Remove(tempFile)

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("保存临时文件失败: %w", err)
	}

	_, err = s.imageProcessor.ValidateFormat(tempFile)
	if err != nil {
		return fmt.Errorf("文件格式验证失败: %w", err)
	}

	if err := s.imageProcessor.ValidatePanoramaResolution(tempFile); err != nil {
		return fmt.Errorf("分辨率验证失败: %w", err)
	}

	imageInfo, err := s.imageProcessor.GetImageInfo(tempFile)
	if err != nil {
		return fmt.Errorf("获取图片信息失败: %w", err)
	}

	md5Hash, err := s.imageProcessor.CalculateMD5(tempFile)
	if err != nil {
		return fmt.Errorf("计算MD5失败: %w", err)
	}

	existingScene, err := s.repo.FindByMD5(md5Hash)
	if err != nil {
		return fmt.Errorf("检查MD5失败: %w", err)
	}
	if existingScene != nil {
		scene.SourceURL = existingScene.SourceURL
		scene.ThumbnailURL = existingScene.ThumbnailURL
		scene.SourceWidth = existingScene.SourceWidth
		scene.SourceHeight = existingScene.SourceHeight
		scene.SourceFileSize = existingScene.SourceFileSize
		scene.SourceFileMD5 = md5Hash
		return nil
	}

	sourceObjectName := fmt.Sprintf("spaces/%s/sources/%s_sphere.jpg", spaceName, scene.SceneCode)
	sourceURL, err := s.minioClient.UploadFile(sourceObjectName, tempFile, "image/jpeg")
	if err != nil {
		return fmt.Errorf("上传源文件失败: %w", err)
	}

	thumbFile := filepath.Join(os.TempDir(), fmt.Sprintf("thumb_%d.jpg", time.Now().UnixNano()))
	defer os.Remove(thumbFile)

	if err := s.imageProcessor.GenerateThumbnail(tempFile, thumbFile); err != nil {
		return fmt.Errorf("生成缩略图失败: %w", err)
	}

	thumbObjectName := fmt.Sprintf("spaces/%s/thumbnails/%s_thumb.jpg", spaceName, scene.SceneCode)
	thumbURL, err := s.minioClient.UploadFile(thumbObjectName, thumbFile, "image/jpeg")
	if err != nil {
		return fmt.Errorf("上传缩略图失败: %w", err)
	}

	scene.SourceURL = sourceURL
	scene.ThumbnailURL = thumbURL
	scene.SourceWidth = imageInfo.Width
	scene.SourceHeight = imageInfo.Height
	scene.SourceFileSize = imageInfo.FileSize
	scene.SourceFileMD5 = md5Hash

	return nil
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

	if err := s.repo.Update(scene); err != nil {
		return nil, fmt.Errorf("更新场景失败: %w", err)
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
	space, err := s.spaceRepo.FindByID(spaceID)
	if err != nil {
		return nil, fmt.Errorf("空间不存在: %w", err)
	}

	if !isAdmin && space.CreatedBy != userID {
		return nil, errors.New("无权限在此空间导入场景")
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
		result := dto.BatchImportResult{
			Index:     i + 1,
			SceneCode: item.SceneCode,
			Title:     item.Title,
		}

		exists, _ := s.repo.CheckSceneCodeExists(item.SceneCode, 0)
		if exists {
			result.Status = "failed"
			result.Message = "场景编码已存在"
			response.FailedCount++
			response.Errors = append(response.Errors, dto.BatchImportError{
				Index:     i + 1,
				SceneCode: item.SceneCode,
				Error:     "场景编码已存在",
			})
		} else {
			scene := &model.ResScene{
				SpaceID:   spaceID,
				Title:     item.Title,
				SceneCode: item.SceneCode,
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
					SceneCode: item.SceneCode,
					Error:     err.Error(),
				})
			} else {
				result.Status = "success"
				result.Message = "创建成功"
				response.SuccessCount++
			}
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
	var unplaced []*dto.SceneNodeData
	var edges []*dto.EdgeData

	sceneMap := make(map[uint]bool)
	for _, scene := range scenes {
		sceneMap[scene.ID] = true
		nodeData := dto.ToSceneNodeData(scene)
		if nodeData.HasPosition {
			nodes = append(nodes, nodeData)
		} else {
			unplaced = append(unplaced, nodeData)
		}

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
		Unplaced:  unplaced,
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
