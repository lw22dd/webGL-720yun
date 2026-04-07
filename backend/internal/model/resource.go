package model

import (
	"time"

	"gorm.io/gorm"
)

// ResSpace 空间/景区表
type ResSpace struct {
	ID          uint           `gorm:"primaryKey;column:id" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;uniqueIndex;column:name" json:"name"`
	Slug        string         `gorm:"type:varchar(100);not null;uniqueIndex;column:slug" json:"slug"`
	CoverURL    string         `gorm:"type:varchar(500);column:cover_url" json:"cover_url"`
	Description string         `gorm:"type:text;column:description" json:"description"`
	Province    string         `gorm:"type:varchar(50);column:province" json:"province"`
	City        string         `gorm:"type:varchar(50);column:city" json:"city"`
	Longitude   float64        `gorm:"type:decimal(10,7);column:longitude" json:"longitude"`
	Latitude    float64        `gorm:"type:decimal(10,7);column:latitude" json:"latitude"`
	ZoomLevel   int            `gorm:"default:12;column:zoom_level" json:"zoom_level"`
	SortOrder   int            `gorm:"default:0;column:sort_order" json:"sort_order"`
	Status      uint8          `gorm:"default:1;column:status" json:"status"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`
	CreatedBy   uint           `gorm:"index;column:created_by" json:"created_by"`

	Scenes []ResScene `gorm:"foreignKey:SpaceID" json:"scenes,omitempty"`
}

func (ResSpace) TableName() string {
	return "res_space"
}

// ResScene 场景/景点表
type ResScene struct {
	ID           uint   `gorm:"primaryKey;column:id" json:"id"`
	SpaceID      uint   `gorm:"not null;index;column:space_id" json:"space_id"`
	Title        string `gorm:"type:varchar(200);not null;column:title" json:"title"`
	SceneCode    string `gorm:"type:varchar(50);not null;uniqueIndex;column:scene_code" json:"scene_code"`
	PanoramaType string `gorm:"type:varchar(20);default:equirectangular;column:panorama_type" json:"panorama_type"`

	SourceURL      string `gorm:"type:varchar(500);column:source_url" json:"source_url"`
	SourceWidth    int    `gorm:"column:source_width" json:"source_width"`
	SourceHeight   int    `gorm:"column:source_height" json:"source_height"`
	SourceFileSize int64  `gorm:"column:source_file_size" json:"source_file_size"`
	SourceFileMD5  string `gorm:"type:varchar(32);index;column:source_file_md5" json:"source_file_md5"`

	CubemapURL  string `gorm:"type:varchar(500);column:cubemap_url" json:"cubemap_url"`
	IsConverted bool   `gorm:"default:false;column:is_converted" json:"is_converted"`

	ThumbnailURL  string `gorm:"type:varchar(500);column:thumbnail_url" json:"thumbnail_url"`
	CoverImageURL string `gorm:"type:varchar(500);column:cover_image_url" json:"cover_image_url"`

	InitialFOV   float64 `gorm:"type:decimal(5,2);default:100;column:initial_fov" json:"initial_fov"`
	InitialPitch float64 `gorm:"type:decimal(5,2);default:0;column:initial_pitch" json:"initial_pitch"`
	InitialYaw   float64 `gorm:"type:decimal(5,2);default:0;column:initial_yaw" json:"initial_yaw"`
	NorthOffset  float64 `gorm:"type:decimal(5,2);default:0;column:north_offset" json:"north_offset"`

	Longitude float64 `gorm:"type:decimal(10,7);column:longitude" json:"longitude"`
	Latitude  float64 `gorm:"type:decimal(10,7);column:latitude" json:"latitude"`

	ViewCount int64          `gorm:"default:0;column:view_count" json:"view_count"`
	SortOrder int            `gorm:"default:0;column:sort_order" json:"sort_order"`
	Status    uint8          `gorm:"default:1;column:status" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	Space    ResSpace     `gorm:"foreignKey:SpaceID" json:"space,omitempty"`
	Hotspots []ResHotspot `gorm:"foreignKey:SceneID" json:"hotspots,omitempty"`
}

func (ResScene) TableName() string {
	return "res_scene"
}

// ResHotspot 热点表
type ResHotspot struct {
	ID            uint  `gorm:"primaryKey;column:id" json:"id"`
	SceneID       uint  `gorm:"not null;index;column:scene_id" json:"scene_id"`
	Type          uint8 `gorm:"not null;column:type" json:"type"`
	TargetSceneID *uint `gorm:"index;column:target_scene_id" json:"target_scene_id,omitempty"`

	Pitch float64 `gorm:"type:decimal(5,2);column:pitch" json:"pitch"`
	Yaw   float64 `gorm:"type:decimal(5,2);column:yaw" json:"yaw"`

	Title   string `gorm:"type:varchar(100);column:title" json:"title"`
	IconURL string `gorm:"type:varchar(500);column:icon_url" json:"icon_url"`
	Style   string `gorm:"type:varchar(50);default:default;column:style" json:"style"`

	Content   string `gorm:"type:text;column:content" json:"content,omitempty"`
	MediaType string `gorm:"type:varchar(20);column:media_type" json:"media_type,omitempty"`
	MediaURL  string `gorm:"type:varchar(500);column:media_url" json:"media_url,omitempty"`

	Question string `gorm:"type:text;column:question" json:"question,omitempty"`
	Options  string `gorm:"type:text;column:options" json:"options,omitempty"`
	Answer   string `gorm:"type:varchar(200);column:answer" json:"answer,omitempty"`
	Score    int    `gorm:"default:10;column:score" json:"score"`

	TransitionEffect string `gorm:"type:varchar(20);default:fade;column:transition_effect" json:"transition_effect"`

	SortOrder int       `gorm:"default:0;column:sort_order" json:"sort_order"`
	Status    uint8     `gorm:"default:1;column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`

	Scene  ResScene  `gorm:"foreignKey:SceneID" json:"scene,omitempty"`
	Target *ResScene `gorm:"foreignKey:TargetSceneID" json:"target,omitempty"`
}

func (ResHotspot) TableName() string {
	return "res_hotspot"
}

const (
	HotspotTypeSwitch uint8 = 1 // 场景切换热点
	HotspotTypeTeach  uint8 = 2 // 教学热点
	HotspotTypeQuiz   uint8 = 3 // 答题热点
)
