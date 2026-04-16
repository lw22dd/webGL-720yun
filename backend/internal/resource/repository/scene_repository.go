package repository

import (
	"errors"

	"webGL-720yun/internal/model"

	"gorm.io/gorm"
)

type SceneRepository struct {
	db *gorm.DB
}

func NewSceneRepository(db *gorm.DB) *SceneRepository {
	return &SceneRepository{db: db}
}

func (r *SceneRepository) Create(scene *model.ResScene) error {
	return r.db.Create(scene).Error
}

func (r *SceneRepository) Update(scene *model.ResScene) error {
	return r.db.Save(scene).Error
}

func (r *SceneRepository) Delete(id uint) error {
	return r.db.Delete(&model.ResScene{}, id).Error
}

func (r *SceneRepository) FindByID(id uint) (*model.ResScene, error) {
	var scene model.ResScene
	err := r.db.First(&scene, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("场景不存在")
		}
		return nil, err
	}
	return &scene, nil
}

func (r *SceneRepository) FindByIDWithHotspots(id uint) (*model.ResScene, error) {
	var scene model.ResScene
	err := r.db.Preload("Hotspots").First(&scene, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("场景不存在")
		}
		return nil, err
	}
	return &scene, nil
}

func (r *SceneRepository) FindByIDWithSpace(id uint) (*model.ResScene, error) {
	var scene model.ResScene
	err := r.db.Preload("Space").First(&scene, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("场景不存在")
		}
		return nil, err
	}
	return &scene, nil
}

func (r *SceneRepository) FindBySceneCode(code string) (*model.ResScene, error) {
	var scene model.ResScene
	err := r.db.Where("scene_code = ?", code).First(&scene).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("场景不存在")
		}
		return nil, err
	}
	return &scene, nil
}

<<<<<<< HEAD
=======
// FindBySceneCodeWithSpace 通过 sceneCode 查找场景并预加载 Space 关联
// 用于瓦片读取时通过 sceneCode 解析出 space.slug（spaceName）
func (r *SceneRepository) FindBySceneCodeWithSpace(code string) (*model.ResScene, error) {
	var scene model.ResScene
	err := r.db.Where("scene_code = ?", code).Preload("Space").First(&scene).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("场景不存在")
		}
		return nil, err
	}
	return &scene, nil
}

>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
func (r *SceneRepository) FindByMD5(md5 string) (*model.ResScene, error) {
	var scene model.ResScene
	err := r.db.Where("source_file_md5 = ?", md5).First(&scene).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &scene, nil
}

type SceneListQuery struct {
	Page     int
	PageSize int
	SpaceID  uint
	Title    string
	Status   *uint8
	Keyword  string
}

func (r *SceneRepository) GetList(query *SceneListQuery) ([]*model.ResScene, int64, error) {
	var scenes []*model.ResScene
	var total int64

	db := r.db.Model(&model.ResScene{})

	if query.SpaceID > 0 {
		db = db.Where("space_id = ?", query.SpaceID)
	}

	if query.Title != "" {
		db = db.Where("title LIKE ?", "%"+query.Title+"%")
	}

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("title LIKE ? OR scene_code LIKE ?", keyword, keyword)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	err := db.Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&scenes).Error

	if err != nil {
		return nil, 0, err
	}

	return scenes, total, nil
}

func (r *SceneRepository) GetListBySpaceID(spaceID uint) ([]*model.ResScene, error) {
	var scenes []*model.ResScene
	err := r.db.Where("space_id = ?", spaceID).
		Order("sort_order ASC, created_at DESC").
		Find(&scenes).Error
	if err != nil {
		return nil, err
	}
	return scenes, nil
}

func (r *SceneRepository) BatchCreate(scenes []*model.ResScene) error {
	if len(scenes) == 0 {
		return nil
	}
	return r.db.Create(&scenes).Error
}

func (r *SceneRepository) CheckSceneCodeExists(code string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.ResScene{}).Where("scene_code = ?", code)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SceneRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *SceneRepository) DeleteHotspotsBySceneID(sceneID uint) error {
	return r.db.Where("scene_id = ?", sceneID).Delete(&model.ResHotspot{}).Error
}

func (r *SceneRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&model.ResScene{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *SceneRepository) FindBySpaceIDWithHotspots(spaceID uint) ([]*model.ResScene, error) {
	var scenes []*model.ResScene
	err := r.db.Where("space_id = ?", spaceID).
		Preload("Hotspots").
		Order("sort_order ASC, created_at DESC").
		Find(&scenes).Error
	if err != nil {
		return nil, err
	}
	return scenes, nil
}

func (r *SceneRepository) UpdatePosition(id uint, longitude, latitude float64) error {
	return r.db.Model(&model.ResScene{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"longitude": longitude,
			"latitude":  latitude,
		}).Error
}

type ScenePositionUpdate struct {
	ID        uint
	Longitude float64
	Latitude  float64
}

func (r *SceneRepository) BatchUpdatePosition(positions []ScenePositionUpdate) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, pos := range positions {
			if err := tx.Model(&model.ResScene{}).Where("id = ?", pos.ID).
				Updates(map[string]interface{}{
					"longitude": pos.Longitude,
					"latitude":  pos.Latitude,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
