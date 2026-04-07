package repository

import (
	"errors"

	"webGL-720yun/internal/model"

	"gorm.io/gorm"
)

type HotspotRepository struct {
	db *gorm.DB
}

func NewHotspotRepository(db *gorm.DB) *HotspotRepository {
	return &HotspotRepository{db: db}
}

func (r *HotspotRepository) Create(hotspot *model.ResHotspot) error {
	return r.db.Create(hotspot).Error
}

func (r *HotspotRepository) Update(hotspot *model.ResHotspot) error {
	return r.db.Save(hotspot).Error
}

func (r *HotspotRepository) Delete(id uint) error {
	return r.db.Delete(&model.ResHotspot{}, id).Error
}

func (r *HotspotRepository) FindByID(id uint) (*model.ResHotspot, error) {
	var hotspot model.ResHotspot
	err := r.db.First(&hotspot, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("热点不存在")
		}
		return nil, err
	}
	return &hotspot, nil
}

func (r *HotspotRepository) FindByIDWithScene(id uint) (*model.ResHotspot, error) {
	var hotspot model.ResHotspot
	err := r.db.Preload("Scene").First(&hotspot, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("热点不存在")
		}
		return nil, err
	}
	return &hotspot, nil
}

func (r *HotspotRepository) FindByIDWithTarget(id uint) (*model.ResHotspot, error) {
	var hotspot model.ResHotspot
	err := r.db.Preload("Target").First(&hotspot, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("热点不存在")
		}
		return nil, err
	}
	return &hotspot, nil
}

type HotspotListQuery struct {
	Page     int
	PageSize int
	SceneID  uint
	Type     *uint8
	Status   *uint8
}

func (r *HotspotRepository) GetList(query *HotspotListQuery) ([]*model.ResHotspot, int64, error) {
	var hotspots []*model.ResHotspot
	var total int64

	db := r.db.Model(&model.ResHotspot{})

	if query.SceneID > 0 {
		db = db.Where("scene_id = ?", query.SceneID)
	}

	if query.Type != nil {
		db = db.Where("type = ?", *query.Type)
	}

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	err := db.Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&hotspots).Error

	if err != nil {
		return nil, 0, err
	}

	return hotspots, total, nil
}

func (r *HotspotRepository) GetListBySceneID(sceneID uint) ([]*model.ResHotspot, error) {
	var hotspots []*model.ResHotspot
	err := r.db.Where("scene_id = ?", sceneID).
		Order("sort_order ASC, created_at DESC").
		Find(&hotspots).Error
	if err != nil {
		return nil, err
	}
	return hotspots, nil
}

func (r *HotspotRepository) DeleteBySceneID(sceneID uint) error {
	return r.db.Where("scene_id = ?", sceneID).Delete(&model.ResHotspot{}).Error
}

func (r *HotspotRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *HotspotRepository) CountBySceneID(sceneID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.ResHotspot{}).Where("scene_id = ?", sceneID).Count(&count).Error
	return count, err
}

func (r *HotspotRepository) CountByType(sceneID uint, hotspotType uint8) (int64, error) {
	var count int64
	err := r.db.Model(&model.ResHotspot{}).
		Where("scene_id = ? AND type = ?", sceneID, hotspotType).
		Count(&count).Error
	return count, err
}
