package repository

import (
	"errors"

	"webGL-720yun/internal/model"

	"gorm.io/gorm"
)

type SpaceRepository struct {
	db *gorm.DB
}

func NewSpaceRepository(db *gorm.DB) *SpaceRepository {
	return &SpaceRepository{db: db}
}

func (r *SpaceRepository) Create(space *model.ResSpace) error {
	return r.db.Create(space).Error
}

func (r *SpaceRepository) Update(space *model.ResSpace) error {
	return r.db.Save(space).Error
}

func (r *SpaceRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&model.ResSpace{}, id).Error
}

func (r *SpaceRepository) FindByID(id uint) (*model.ResSpace, error) {
	var space model.ResSpace
	err := r.db.First(&space, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("空间不存在")
		}
		return nil, err
	}
	return &space, nil
}

func (r *SpaceRepository) FindByIDWithScenes(id uint) (*model.ResSpace, error) {
	var space model.ResSpace
	err := r.db.Preload("Scenes").First(&space, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("空间不存在")
		}
		return nil, err
	}
	return &space, nil
}

type SpaceListQuery struct {
	Page     int
	PageSize int
	Name     string
	Status   *uint8
	Keyword  string
	CreatedBy uint
}

func (r *SpaceRepository) GetList(query *SpaceListQuery) ([]*model.ResSpace, int64, error) {
	var spaces []*model.ResSpace
	var total int64

	db := r.db.Model(&model.ResSpace{})

	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("name LIKE ? OR description LIKE ? OR province LIKE ? OR city LIKE ?",
			keyword, keyword, keyword, keyword)
	}

	if query.CreatedBy > 0 {
		db = db.Where("created_by = ?", query.CreatedBy)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	err := db.Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&spaces).Error

	if err != nil {
		return nil, 0, err
	}

	return spaces, total, nil
}

func (r *SpaceRepository) CheckNameExists(name string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.ResSpace{}).Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SpaceRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *SpaceRepository) DeleteScenesBySpaceID(spaceID uint) error {
	return r.db.Unscoped().Where("space_id = ?", spaceID).Delete(&model.ResScene{}).Error
}

func (r *SpaceRepository) DeleteHotspotsBySpaceID(spaceID uint) error {
	return r.db.Unscoped().Where("scene_id IN (SELECT id FROM res_scene WHERE space_id = ?)", spaceID).
		Delete(&model.ResHotspot{}).Error
}
