package dto

import (
	"time"
	"webGL-720yun/internal/model"
	"webGL-720yun/pkg/utils"
)

type SceneListRequest struct {
	Page     int    `form:"page" binding:"min=1" json:"page"`
	PageSize int    `form:"page_size" binding:"min=1,max=100" json:"page_size"`
	SpaceID  uint   `form:"space_id" json:"space_id"`
	Title    string `form:"title" json:"title"`
	Status   *uint8 `form:"status" json:"status"`
	Keyword  string `form:"keyword" json:"keyword"`
}

type CreateSceneRequest struct {
	SpaceID      uint    `json:"space_id" binding:"required"`
	Title        string  `json:"title" binding:"required,min=2,max=200"`
	SceneCode    string  `json:"scene_code" binding:"omitempty,min=2,max=100"`
	FileID       string  `json:"file_id" binding:"omitempty"`
	PanoramaType string  `json:"panorama_type" binding:"omitempty,oneof=equirectangular cubemap"`
	InitialFOV   float64 `json:"initial_fov" binding:"omitempty,min=30,max=150"`
	InitialPitch float64 `json:"initial_pitch" binding:"omitempty,min=-90,max=90"`
	InitialYaw   float64 `json:"initial_yaw" binding:"omitempty,min=-180,max=180"`
	NorthOffset  float64 `json:"north_offset" binding:"omitempty,min=-180,max=180"`
	Longitude    float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	Latitude     float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	SortOrder    int     `json:"sort_order"`
}

type CreateSceneResponse struct {
	SceneID     uint   `json:"scene_id"`
	Title       string `json:"title"`
	SceneCode   string `json:"scene_code"`
	SliceStatus string `json:"slice_status"`
	TaskID      string `json:"task_id"`
}

type UpdateSceneRequest struct {
	Title        string  `json:"title" binding:"omitempty,min=2,max=200"`
	PanoramaType string  `json:"panorama_type" binding:"omitempty,oneof=equirectangular cubemap"`
	FileID       string  `json:"file_id" binding:"omitempty"`
	InitialFOV   float64 `json:"initial_fov" binding:"omitempty,min=30,max=150"`
	InitialPitch float64 `json:"initial_pitch" binding:"omitempty,min=-90,max=90"`
	InitialYaw   float64 `json:"initial_yaw" binding:"omitempty,min=-180,max=180"`
	NorthOffset  float64 `json:"north_offset" binding:"omitempty,min=-180,max=180"`
	Longitude    float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	Latitude     float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	SortOrder    int     `json:"sort_order"`
	Status       *uint8  `json:"status" binding:"omitempty"`
}

type SceneListResponse struct {
	utils.PageInfo `json:"page_info"`
	Scenes         []*SceneListItem `json:"scenes"`
}

type SceneListItem struct {
	ID           uint      `json:"id"`
	SpaceID      uint      `json:"space_id"`
	Title        string    `json:"title"`
	SceneCode    string    `json:"scene_code"`
	PanoramaType string    `json:"panorama_type"`
	ThumbnailURL string    `json:"thumbnail_url"`
	SourceWidth  int       `json:"source_width"`
	SourceHeight int       `json:"source_height"`
	ViewCount    int64     `json:"view_count"`
	SliceStatus  string    `json:"slice_status"`
	SortOrder    int       `json:"sort_order"`
	Status       uint8     `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SceneDetailResponse struct {
	ID             uint             `json:"id"`
	SpaceID        uint             `json:"space_id"`
	Title          string           `json:"title"`
	SceneCode      string           `json:"scene_code"`
	PanoramaType   string           `json:"panorama_type"`
	FileID         string           `json:"file_id"`
	SourceURL      string           `json:"source_url"`
	SourceWidth    int              `json:"source_width"`
	SourceHeight   int              `json:"source_height"`
	SourceFileSize int64            `json:"source_file_size"`
	SourceFileSHA256  string           `json:"source_file_sha256"`
	TileURL        string           `json:"tile_url"`
	PreviewURL     string           `json:"preview_url"`
	CubemapURL     string           `json:"cubemap_url"`
	IsConverted    bool             `json:"is_converted"`
	SliceStatus    string           `json:"slice_status"`
	TaskID         string           `json:"task_id"`
	ThumbnailURL   string           `json:"thumbnail_url"`
	CoverImageURL  string           `json:"cover_image_url"`
	InitialFOV     float64          `json:"initial_fov"`
	InitialPitch   float64          `json:"initial_pitch"`
	InitialYaw     float64          `json:"initial_yaw"`
	NorthOffset    float64          `json:"north_offset"`
	Longitude      float64          `json:"longitude"`
	Latitude       float64          `json:"latitude"`
	ViewCount      int64            `json:"view_count"`
	SortOrder      int              `json:"sort_order"`
	Status         uint8            `json:"status"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	Hotspots       []*HotspotSimple `json:"hotspots,omitempty"`
}

type HotspotSimple struct {
	ID      uint    `json:"id"`
	Type    uint8   `json:"type"`
	Title   string  `json:"title"`
	Pitch   float64 `json:"pitch"`
	Yaw     float64 `json:"yaw"`
	IconURL string  `json:"icon_url"`
}

type BatchImportRequest struct {
	SpaceID uint `json:"space_id" binding:"required"`
}

type BatchImportItem struct {
	Title     string  `json:"title" binding:"required"`
	SceneCode string  `json:"scene_code" binding:"omitempty"`
	FileName  string  `json:"file_name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type BatchImportResponse struct {
	SuccessCount int                 `json:"success_count"`
	FailedCount  int                 `json:"failed_count"`
	Results      []BatchImportResult `json:"results"`
	Errors       []BatchImportError  `json:"errors,omitempty"`
}

type BatchImportResult struct {
	Index     int    `json:"index"`
	SceneCode string `json:"scene_code"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type BatchImportError struct {
	Index     int    `json:"index"`
	SceneCode string `json:"scene_code"`
	Error     string `json:"error"`
}

func ToSceneListItem(scene *model.ResScene) *SceneListItem {
	return &SceneListItem{
		ID:           scene.ID,
		SpaceID:      scene.SpaceID,
		Title:        scene.Title,
		SceneCode:    scene.SceneCode,
		PanoramaType: scene.PanoramaType,
		ThumbnailURL: scene.ThumbnailURL,
		SourceWidth:  scene.SourceWidth,
		SourceHeight: scene.SourceHeight,
		ViewCount:    scene.ViewCount,
		SliceStatus:  scene.SliceStatus,
		SortOrder:    scene.SortOrder,
		Status:       scene.Status,
		CreatedAt:    scene.CreatedAt,
		UpdatedAt:    scene.UpdatedAt,
	}
}

func ToSceneDetailResponse(scene *model.ResScene) *SceneDetailResponse {
	response := &SceneDetailResponse{
		ID:             scene.ID,
		SpaceID:        scene.SpaceID,
		Title:          scene.Title,
		SceneCode:      scene.SceneCode,
		PanoramaType:   scene.PanoramaType,
		FileID:         scene.FileID,
		SourceURL:      scene.SourceURL,
		SourceWidth:    scene.SourceWidth,
		SourceHeight:   scene.SourceHeight,
		SourceFileSize: scene.SourceFileSize,
		SourceFileSHA256:  scene.SourceFileSHA256,
		TileURL:        scene.TileURL,
		PreviewURL:     scene.PreviewURL,
		CubemapURL:     scene.CubemapURL,
		IsConverted:    scene.IsConverted,
		SliceStatus:    scene.SliceStatus,
		TaskID:         scene.TaskID,
		ThumbnailURL:   scene.ThumbnailURL,
		CoverImageURL:  scene.CoverImageURL,
		InitialFOV:     scene.InitialFOV,
		InitialPitch:   scene.InitialPitch,
		InitialYaw:     scene.InitialYaw,
		NorthOffset:    scene.NorthOffset,
		Longitude:      scene.Longitude,
		Latitude:       scene.Latitude,
		ViewCount:      scene.ViewCount,
		SortOrder:      scene.SortOrder,
		Status:         scene.Status,
		CreatedAt:      scene.CreatedAt,
		UpdatedAt:      scene.UpdatedAt,
	}

	if len(scene.Hotspots) > 0 {
		response.Hotspots = make([]*HotspotSimple, len(scene.Hotspots))
		for i, hotspot := range scene.Hotspots {
			response.Hotspots[i] = &HotspotSimple{
				ID:      hotspot.ID,
				Type:    hotspot.Type,
				Title:   hotspot.Title,
				Pitch:   hotspot.Pitch,
				Yaw:     hotspot.Yaw,
				IconURL: hotspot.IconURL,
			}
		}
	}

	return response
}
