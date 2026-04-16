package service

import (
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
	"webGL-720yun/pkg/minio_client"
	"webGL-720yun/pkg/utils"

	"gorm.io/gorm"
)

type SpaceService struct {
	repo        *repository.SpaceRepository
	minioClient *minio_client.MinIOClient
}

func NewSpaceService(db *gorm.DB, minioClient *minio_client.MinIOClient) *SpaceService {
	return &SpaceService{
		repo:        repository.NewSpaceRepository(db),
		minioClient: minioClient,
	}
}

func (s *SpaceService) CreateSpace(req *dto.CreateSpaceRequest, createdBy uint, coverFile *multipart.FileHeader) (*model.ResSpace, error) {
	exists, err := s.repo.CheckNameExists(req.Name, 0)
	if err != nil {
		return nil, fmt.Errorf("检查名称失败: %w", err)
	}
	if exists {
		return nil, errors.New("景区名称已存在")
	}

	space := &model.ResSpace{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Province:    req.Province,
		City:        req.City,
		Longitude:   req.Longitude,
		Latitude:    req.Latitude,
		ZoomLevel:   req.ZoomLevel,
		SortOrder:   req.SortOrder,
		Status:      1,
		CreatedBy:   createdBy,
	}

	if coverFile != nil {
		coverURL, err := s.uploadCoverImage(coverFile, space.Slug)
		if err != nil {
			return nil, fmt.Errorf("上传封面图失败: %w", err)
		}
		space.CoverURL = coverURL
	}

	if err := s.repo.Create(space); err != nil {
		return nil, fmt.Errorf("创建空间失败: %w", err)
	}

	return space, nil
}

func (s *SpaceService) UpdateSpace(id uint, req *dto.UpdateSpaceRequest, userID uint, isAdmin bool, coverFile *multipart.FileHeader) (*model.ResSpace, error) {
	space, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !isAdmin && space.CreatedBy != userID {
		return nil, errors.New("无权限修改此空间")
	}

	if req.Name != "" {
		exists, err := s.repo.CheckNameExists(req.Name, id)
		if err != nil {
			return nil, fmt.Errorf("检查名称失败: %w", err)
		}
		if exists {
			return nil, errors.New("景区名称已存在")
		}
		space.Name = req.Name
	}

	if req.Slug != "" {
		space.Slug = req.Slug
	}

	if req.Description != "" {
		space.Description = req.Description
	}
	if req.Province != "" {
		space.Province = req.Province
	}
	if req.City != "" {
		space.City = req.City
	}
	if req.Longitude != 0 {
		space.Longitude = req.Longitude
	}
	if req.Latitude != 0 {
		space.Latitude = req.Latitude
	}
	if req.ZoomLevel != 0 {
		space.ZoomLevel = req.ZoomLevel
	}
	if req.SortOrder != 0 {
		space.SortOrder = req.SortOrder
	}
	if req.Status != nil {
		space.Status = *req.Status
	}

	if coverFile != nil {
		if space.CoverURL != "" {
			oldObjectName := s.extractObjectName(space.CoverURL)
			if oldObjectName != "" {
				_ = s.minioClient.DeleteObject(oldObjectName)
			}
		}

		coverURL, err := s.uploadCoverImage(coverFile, space.Slug)
		if err != nil {
			return nil, fmt.Errorf("上传封面图失败: %w", err)
		}
		space.CoverURL = coverURL
	}

	if err := s.repo.Update(space); err != nil {
		return nil, fmt.Errorf("更新空间失败: %w", err)
	}

	return space, nil
}

func (s *SpaceService) DeleteSpace(id uint, userID uint, isAdmin bool) error {
	space, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if !isAdmin && space.CreatedBy != userID {
		return errors.New("无权限删除此空间")
	}

	err = s.repo.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.DeleteHotspotsBySpaceID(id); err != nil {
			return fmt.Errorf("删除热点失败: %w", err)
		}

		if err := s.repo.DeleteScenesBySpaceID(id); err != nil {
			return fmt.Errorf("删除场景失败: %w", err)
		}

		if err := s.repo.Delete(id); err != nil {
			return fmt.Errorf("删除空间失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 删除 MinIO 中的所有相关文件（源图、瓦片、预览图、封面等）
	// 同时尝试删除基于 slug 和基于 ID 的目录（兼容旧数据）
	slugPrefix := fmt.Sprintf("spaces/%s/", space.Slug)
	idPrefix := fmt.Sprintf("spaces/%d/", space.ID)

	_ = s.minioClient.DeleteObjectsWithPrefix(slugPrefix)
	_ = s.minioClient.DeleteObjectsWithPrefix(idPrefix)

	return nil
}

func (s *SpaceService) DeleteSpaceBatch(ids []uint, userID uint, isAdmin bool) error {
	if len(ids) == 0 {
		return errors.New("ids不能为空")
	}

	for _, id := range ids {
		if err := s.DeleteSpace(id, userID, isAdmin); err != nil {
			return fmt.Errorf("删除空间ID %d 失败: %w", id, err)
		}
	}

	return nil
}

func (s *SpaceService) GetSpaceList(req *dto.SpaceListRequest) (*dto.SpaceListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	query := &repository.SpaceListQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Status:   req.Status,
		Keyword:  req.Keyword,
	}

	spaces, total, err := s.repo.GetList(query)
	if err != nil {
		return nil, fmt.Errorf("获取空间列表失败: %w", err)
	}

	items := make([]*dto.SpaceListItem, len(spaces))
	for i, space := range spaces {
		items[i] = dto.ToSpaceListItem(space)
	}

	return &dto.SpaceListResponse{
		PageInfo: utils.PageInfo{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
			Pages:    int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
		},
		Spaces: items,
	}, nil
}

func (s *SpaceService) GetSpaceDetail(id uint) (*dto.SpaceDetailResponse, error) {
	space, err := s.repo.FindByIDWithScenes(id)
	if err != nil {
		return nil, err
	}

	return dto.ToSpaceDetailResponse(space), nil
}

func (s *SpaceService) GetSpaceByID(id uint) (*model.ResSpace, error) {
	return s.repo.FindByID(id)
}

func (s *SpaceService) uploadCoverImage(file *multipart.FileHeader, slug string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	ext := ".jpg"
	objectName := fmt.Sprintf("spaces/%s/covers/cover%s", slug, ext)

	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("upload_%d%s", time.Now().UnixNano(), ext))
	dst, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer dst.Close()
	defer os.Remove(tempFile)

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("保存临时文件失败: %w", err)
	}

	url, err := s.minioClient.UploadFile(objectName, tempFile, "")
	if err != nil {
		return "", fmt.Errorf("上传到MinIO失败: %w", err)
	}

	return url, nil
}

func (s *SpaceService) extractObjectName(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 4 {
		return strings.Join(parts[3:], "/")
	}
	return ""
}
