package dto

import (
	"time"
	"webGL-720yun/internal/model"
	"webGL-720yun/pkg/utils"
)

type SpaceListRequest struct {
	Page     int    `form:"page" binding:"min=1" json:"page"`
	PageSize int    `form:"page_size" binding:"min=1,max=100" json:"page_size"`
	Name     string `form:"name" json:"name"`
	Status   *uint8 `form:"status" json:"status"`
	Keyword  string `form:"keyword" json:"keyword"`
}

type CreateSpaceRequest struct {
	Name        string  `form:"name" json:"name" binding:"required,min=2,max=100"`
	Slug        string  `form:"slug" json:"slug" binding:"omitempty,min=2,max=100"`
	Description string  `form:"description" json:"description" binding:"max=1000"`
	Province    string  `form:"province" json:"province" binding:"max=50"`
	City        string  `form:"city" json:"city" binding:"max=50"`
	Longitude   float64 `form:"longitude" json:"longitude" binding:"omitempty,min=-180,max=180"`
	Latitude    float64 `form:"latitude" json:"latitude" binding:"omitempty,min=-90,max=90"`
	ZoomLevel   int     `form:"zoom_level" json:"zoom_level" binding:"omitempty,min=0,max=20"`
	SortOrder   int     `form:"sort_order" json:"sort_order"`
}

type UpdateSpaceRequest struct {
	Name        string  `form:"name" json:"name" binding:"omitempty,min=2,max=100"`
	Slug        string  `form:"slug" json:"slug" binding:"omitempty,min=2,max=100"`
	Description string  `form:"description" json:"description" binding:"omitempty,max=1000"`
	Province    string  `form:"province" json:"province" binding:"omitempty,max=50"`
	City        string  `form:"city" json:"city" binding:"omitempty,max=50"`
	Longitude   float64 `form:"longitude" json:"longitude" binding:"omitempty,min=-180,max=180"`
	Latitude    float64 `form:"latitude" json:"latitude" binding:"omitempty,min=-90,max=90"`
	ZoomLevel   int     `form:"zoom_level" json:"zoom_level" binding:"omitempty,min=0,max=20"`
	SortOrder   int     `form:"sort_order" json:"sort_order"`
	Status      *uint8  `form:"status" json:"status" binding:"omitempty"`
}

type SpaceListResponse struct {
	utils.PageInfo `json:"page_info"`
	Spaces         []*SpaceListItem `json:"spaces"`
}

type SpaceListItem struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	CoverURL    string    `json:"cover_url"`
	Description string    `json:"description"`
	Province    string    `json:"province"`
	City        string    `json:"city"`
	Longitude   float64   `json:"longitude"`
	Latitude    float64   `json:"latitude"`
	ZoomLevel   int       `json:"zoom_level"`
	SortOrder   int       `json:"sort_order"`
	Status      uint8     `json:"status"`
	SceneCount  int       `json:"scene_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SpaceDetailResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	CoverURL    string         `json:"cover_url"`
	Description string         `json:"description"`
	Province    string         `json:"province"`
	City        string         `json:"city"`
	Longitude   float64        `json:"longitude"`
	Latitude    float64        `json:"latitude"`
	ZoomLevel   int            `json:"zoom_level"`
	SortOrder   int            `json:"sort_order"`
	Status      uint8          `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedBy   uint           `json:"created_by"`
	Scenes      []*SceneSimple `json:"scenes,omitempty"`
}

type SceneSimple struct {
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	SceneCode    string  `json:"scene_code"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	ThumbnailURL string  `json:"thumbnail_url"`
	ViewCount    int64   `json:"view_count"`
	SortOrder    int     `json:"sort_order"`
}

func ToSpaceListItem(space *model.ResSpace) *SpaceListItem {
	return &SpaceListItem{
		ID:          space.ID,
		Name:        space.Name,
		Slug:        space.Slug,
		CoverURL:    space.CoverURL,
		Description: space.Description,
		Province:    space.Province,
		City:        space.City,
		Longitude:   space.Longitude,
		Latitude:    space.Latitude,
		ZoomLevel:   space.ZoomLevel,
		SortOrder:   space.SortOrder,
		Status:      space.Status,
		SceneCount:  len(space.Scenes),
		CreatedAt:   space.CreatedAt,
		UpdatedAt:   space.UpdatedAt,
	}
}

func ToSpaceDetailResponse(space *model.ResSpace) *SpaceDetailResponse {
	response := &SpaceDetailResponse{
		ID:          space.ID,
		Name:        space.Name,
		Slug:        space.Slug,
		CoverURL:    space.CoverURL,
		Description: space.Description,
		Province:    space.Province,
		City:        space.City,
		Longitude:   space.Longitude,
		Latitude:    space.Latitude,
		ZoomLevel:   space.ZoomLevel,
		SortOrder:   space.SortOrder,
		Status:      space.Status,
		CreatedAt:   space.CreatedAt,
		UpdatedAt:   space.UpdatedAt,
		CreatedBy:   space.CreatedBy,
	}

	if len(space.Scenes) > 0 {
		response.Scenes = make([]*SceneSimple, len(space.Scenes))
		for i, scene := range space.Scenes {
			response.Scenes[i] = &SceneSimple{
				ID:           scene.ID,
				Title:        scene.Title,
				SceneCode:    scene.SceneCode,
				Longitude:    scene.Longitude,
				Latitude:     scene.Latitude,
				ThumbnailURL: scene.ThumbnailURL,
				ViewCount:    scene.ViewCount,
				SortOrder:    scene.SortOrder,
			}
		}
	}

	return response
}
