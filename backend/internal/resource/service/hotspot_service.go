package service

import (
	"errors"
	"fmt"

	"webGL-720yun/internal/model"
	"webGL-720yun/internal/resource/dto"
	"webGL-720yun/internal/resource/repository"
	"webGL-720yun/pkg/utils"

	"gorm.io/gorm"
)

type HotspotService struct {
	repo      *repository.HotspotRepository
	sceneRepo *repository.SceneRepository
}

func NewHotspotService(db *gorm.DB) *HotspotService {
	return &HotspotService{
		repo:      repository.NewHotspotRepository(db),
		sceneRepo: repository.NewSceneRepository(db),
	}
}

func (s *HotspotService) CreateHotspot(req *dto.CreateHotspotRequest, userID uint, isAdmin bool) (*model.ResHotspot, error) {
	scene, err := s.sceneRepo.FindByIDWithSpace(req.SceneID)
	if err != nil {
		return nil, fmt.Errorf("场景不存在: %w", err)
	}

	if !isAdmin && scene.Space.CreatedBy != userID {
		return nil, errors.New("无权限在此场景创建热点")
	}

	if err := s.validateHotspotType(req); err != nil {
		return nil, err
	}

	hotspot := &model.ResHotspot{
		SceneID:          req.SceneID,
		Type:             req.Type,
		TargetSceneID:    req.TargetSceneID,
		Pitch:            req.Pitch,
		Yaw:              req.Yaw,
		Title:            req.Title,
		IconURL:          req.IconURL,
		Style:            req.Style,
		Content:          req.Content,
		MediaType:        req.MediaType,
		MediaURL:         req.MediaURL,
		Question:         req.Question,
		Options:          req.Options,
		Answer:           req.Answer,
		Score:            req.Score,
		TransitionEffect: req.TransitionEffect,
		SortOrder:        req.SortOrder,
		Status:           1,
	}

	if hotspot.Style == "" {
		hotspot.Style = "default"
	}
	if hotspot.TransitionEffect == "" {
		hotspot.TransitionEffect = "fade"
	}
	if hotspot.Score == 0 {
		hotspot.Score = 10
	}

	if err := s.repo.Create(hotspot); err != nil {
		return nil, fmt.Errorf("创建热点失败: %w", err)
	}

	return hotspot, nil
}

func (s *HotspotService) validateHotspotType(req *dto.CreateHotspotRequest) error {
	switch req.Type {
	case model.HotspotTypeSwitch:
		if req.TargetSceneID == nil {
			return errors.New("场景切换热点必须指定目标场景ID")
		}
	case model.HotspotTypeTeach:
		if req.Content == "" {
			return errors.New("教学热点必须包含内容")
		}
	case model.HotspotTypeQuiz:
		if req.Question == "" || req.Options == "" || req.Answer == "" {
			return errors.New("答题热点必须包含问题、选项和答案")
		}
	default:
		return errors.New("无效的热点类型")
	}
	return nil
}

func (s *HotspotService) UpdateHotspot(id uint, req *dto.UpdateHotspotRequest, userID uint, isAdmin bool) (*model.ResHotspot, error) {
	hotspot, err := s.repo.FindByIDWithScene(id)
	if err != nil {
		return nil, err
	}

	if !isAdmin && hotspot.Scene.Space.CreatedBy != userID {
		return nil, errors.New("无权限修改此热点")
	}

	if req.TargetSceneID != nil {
		hotspot.TargetSceneID = req.TargetSceneID
	}
	if req.Pitch != 0 {
		hotspot.Pitch = req.Pitch
	}
	if req.Yaw != 0 {
		hotspot.Yaw = req.Yaw
	}
	if req.Title != "" {
		hotspot.Title = req.Title
	}
	if req.IconURL != "" {
		hotspot.IconURL = req.IconURL
	}
	if req.Style != "" {
		hotspot.Style = req.Style
	}
	if req.Content != "" {
		hotspot.Content = req.Content
	}
	if req.MediaType != "" {
		hotspot.MediaType = req.MediaType
	}
	if req.MediaURL != "" {
		hotspot.MediaURL = req.MediaURL
	}
	if req.Question != "" {
		hotspot.Question = req.Question
	}
	if req.Options != "" {
		hotspot.Options = req.Options
	}
	if req.Answer != "" {
		hotspot.Answer = req.Answer
	}
	if req.Score != 0 {
		hotspot.Score = req.Score
	}
	if req.TransitionEffect != "" {
		hotspot.TransitionEffect = req.TransitionEffect
	}
	if req.SortOrder != 0 {
		hotspot.SortOrder = req.SortOrder
	}
	if req.Status != nil {
		hotspot.Status = *req.Status
	}

	if err := s.repo.Update(hotspot); err != nil {
		return nil, fmt.Errorf("更新热点失败: %w", err)
	}

	return hotspot, nil
}

func (s *HotspotService) DeleteHotspot(id uint, userID uint, isAdmin bool) error {
	hotspot, err := s.repo.FindByIDWithScene(id)
	if err != nil {
		return err
	}

	if !isAdmin && hotspot.Scene.Space.CreatedBy != userID {
		return errors.New("无权限删除此热点")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("删除热点失败: %w", err)
	}

	return nil
}

func (s *HotspotService) GetHotspotList(req *dto.HotspotListRequest) (*dto.HotspotListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	query := &repository.HotspotListQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		SceneID:  req.SceneID,
		Type:     req.Type,
		Status:   req.Status,
	}

	hotspots, total, err := s.repo.GetList(query)
	if err != nil {
		return nil, fmt.Errorf("获取热点列表失败: %w", err)
	}

	items := make([]*dto.HotspotListItem, len(hotspots))
	for i, hotspot := range hotspots {
		items[i] = dto.ToHotspotListItem(hotspot)
	}

	return &dto.HotspotListResponse{
		PageInfo: utils.PageInfo{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
			Pages:    int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
		},
		Hotspots: items,
	}, nil
}

func (s *HotspotService) GetHotspotDetail(id uint) (*dto.HotspotDetailResponse, error) {
	hotspot, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return dto.ToHotspotDetailResponse(hotspot), nil
}
