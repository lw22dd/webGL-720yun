package model

import (
	"time"

	"gorm.io/gorm"
)

// ResSpace 空间/景区表
type ResSpace struct {
	ID          uint           `gorm:"primaryKey;column:id" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;uniqueIndex;column:name;comment:空间名称" json:"name"`
	Slug        string         `gorm:"type:varchar(100);not null;index;column:slug;comment:空间标识(英文/拼音)" json:"slug"`
	CoverURL    string         `gorm:"type:varchar(500);column:cover_url;comment:封面URL" json:"cover_url"`
	Description string         `gorm:"type:text;column:description;comment:空间描述" json:"description"`
	Province    string         `gorm:"type:varchar(50);column:province;comment:省份" json:"province"`
	City        string         `gorm:"type:varchar(50);column:city;comment:城市" json:"city"`
	Longitude   float64        `gorm:"type:decimal(10,7);column:longitude;comment:经度" json:"longitude"`
	Latitude    float64        `gorm:"type:decimal(10,7);column:latitude;comment:纬度" json:"latitude"`
	ZoomLevel   int            `gorm:"default:12;column:zoom_level;comment:缩放级别" json:"zoom_level"`
	SortOrder   int            `gorm:"default:0;column:sort_order;comment:排序" json:"sort_order"`
	Status      uint8          `gorm:"default:1;column:status;comment:状态(1启用0禁用)" json:"status"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`
	CreatedBy   uint           `gorm:"index;column:created_by;comment:创建者ID" json:"created_by"`

	Scenes []ResScene `gorm:"foreignKey:SpaceID" json:"scenes,omitempty"`
}

func (ResSpace) TableName() string {
	return "res_space"
}

// ResScene 场景/景点表
type ResScene struct {
	ID           uint   `gorm:"primaryKey;column:id" json:"id"`
	SpaceID      uint   `gorm:"not null;index;column:space_id;comment:空间ID" json:"space_id"`
	Title        string `gorm:"type:varchar(200);not null;column:title;comment:场景标题" json:"title"`
	SceneCode    string `gorm:"type:varchar(50);not null;uniqueIndex;column:scene_code;comment:场景编码，记录经纬度" json:"scene_code"`
	PanoramaType string `gorm:"type:varchar(20);default:equirectangular;column:panorama_type;comment:全景类型" json:"panorama_type"`

	FileID      string `gorm:"type:varchar(50);index;column:file_id;comment:关联的文件ID" json:"file_id"`
	SourceURL      string `gorm:"type:varchar(500);column:source_url;comment:源文件URL" json:"source_url"`
	SourceWidth    int    `gorm:"column:source_width;comment:源文件宽度" json:"source_width"`
	SourceHeight   int    `gorm:"column:source_height;comment:源文件高度" json:"source_height"`
	SourceFileSize int64  `gorm:"column:source_file_size;comment:源文件大小" json:"source_file_size"`
	SourceFileMD5  string `gorm:"type:varchar(32);index;column:source_file_md5;comment:源文件MD5" json:"source_file_md5"`

	TileURL      string `gorm:"type:varchar(500);column:tile_url;comment:瓦片目录URL" json:"tile_url"`
	PreviewURL   string `gorm:"type:varchar(500);column:preview_url;comment:预览图URL" json:"preview_url"`
	CubemapURL   string `gorm:"type:varchar(500);column:cubemap_url;comment:立方体纹理URL" json:"cubemap_url"`
	IsConverted  bool   `gorm:"default:false;column:is_converted;comment:是否已转换" json:"is_converted"`
	SliceStatus  string `gorm:"type:varchar(20);default:pending;column:slice_status;comment:切片状态" json:"slice_status"`
	TaskID       string `gorm:"type:varchar(50);column:task_id;comment:切片任务ID" json:"task_id"`

	ThumbnailURL  string `gorm:"type:varchar(500);column:thumbnail_url;comment:缩略图URL" json:"thumbnail_url"`
	CoverImageURL string `gorm:"type:varchar(500);column:cover_image_url;comment:封面图URL" json:"cover_image_url"`

	InitialFOV   float64 `gorm:"type:decimal(5,2);default:100;column:initial_fov;comment:初始视野角度" json:"initial_fov"`
	InitialPitch float64 `gorm:"type:decimal(5,2);default:0;column:initial_pitch;comment:初始俯仰角" json:"initial_pitch"`
	InitialYaw   float64 `gorm:"type:decimal(5,2);default:0;column:initial_yaw;comment:初始偏航角" json:"initial_yaw"`
	NorthOffset  float64 `gorm:"type:decimal(5,2);default:0;column:north_offset;comment:北方偏移角" json:"north_offset"`

	Longitude float64 `gorm:"type:decimal(10,7);column:longitude;comment:经度" json:"longitude"`
	Latitude  float64 `gorm:"type:decimal(10,7);column:latitude;comment:纬度" json:"latitude"`

	ViewCount int64          `gorm:"default:0;column:view_count;comment:浏览次数" json:"view_count"`
	SortOrder int            `gorm:"default:0;column:sort_order;comment:排序" json:"sort_order"`
	Status    uint8          `gorm:"default:1;column:status;comment:状态(1启用0禁用)" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	Space    ResSpace     `gorm:"foreignKey:SpaceID" json:"space,omitempty"`
	Hotspots []ResHotspot `gorm:"foreignKey:SceneID" json:"hotspots,omitempty"`
}

const (
	SliceStatusPending   = "pending"
	SliceStatusSlicing   = "slicing"
	SliceStatusReady     = "ready"
	SliceStatusFailed    = "failed"
)

func (ResScene) TableName() string {
	return "res_scene"
}

// ResHotspot 热点表
type ResHotspot struct {
	ID            uint  `gorm:"primaryKey;column:id" json:"id"`
	SceneID       uint  `gorm:"not null;index;column:scene_id;comment:场景ID" json:"scene_id"`
	Type          uint8 `gorm:"not null;column:type;comment:热点类型(1场景切换2教学3答题)" json:"type"`
	TargetSceneID *uint `gorm:"index;column:target_scene_id;comment:目标场景ID" json:"target_scene_id,omitempty"`

	Pitch float64 `gorm:"type:decimal(5,2);column:pitch;comment:俯仰角" json:"pitch"`
	Yaw   float64 `gorm:"type:decimal(5,2);column:yaw;comment:偏航角" json:"yaw"`

	Title   string `gorm:"type:varchar(100);column:title;comment:热点标题" json:"title"`
	IconURL string `gorm:"type:varchar(500);column:icon_url;comment:图标URL" json:"icon_url"`
	Style   string `gorm:"type:varchar(50);default:default;column:style;comment:样式" json:"style"`

	Content   string `gorm:"type:text;column:content;comment:内容" json:"content,omitempty"`
	MediaType string `gorm:"type:varchar(20);column:media_type;comment:媒体类型" json:"media_type,omitempty"`
	MediaURL  string `gorm:"type:varchar(500);column:media_url;comment:媒体URL" json:"media_url,omitempty"`

	Question string `gorm:"type:text;column:question;comment:题目" json:"question,omitempty"`
	Options  string `gorm:"type:text;column:options;comment:选项(JSON数组)" json:"options,omitempty"`
	Answer   string `gorm:"type:varchar(200);column:answer;comment:答案" json:"answer,omitempty"`
	Score    int    `gorm:"default:10;column:score;comment:分值" json:"score"`

	TransitionEffect string `gorm:"type:varchar(20);default:fade;column:transition_effect;comment:转场效果" json:"transition_effect"`

	SortOrder int       `gorm:"default:0;column:sort_order;comment:排序" json:"sort_order"`
	Status    uint8     `gorm:"default:1;column:status;comment:状态(1启用0禁用)" json:"status"`
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
