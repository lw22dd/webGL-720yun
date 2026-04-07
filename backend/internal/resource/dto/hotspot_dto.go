package dto

import (
	"time"
	"webGL-720yun/internal/model"
	"webGL-720yun/pkg/utils"
)

type HotspotListRequest struct {
	Page     int    `form:"page" binding:"min=1" json:"page"`
	PageSize int    `form:"page_size" binding:"min=1,max=100" json:"page_size"`
	SceneID  uint   `form:"scene_id" json:"scene_id"`
	Type     *uint8 `form:"type" json:"type"`
	Status   *uint8 `form:"status" json:"status"`
}

type CreateHotspotRequest struct {
	SceneID          uint    `json:"scene_id" binding:"required"`
	Type             uint8   `json:"type" binding:"required,oneof=1 2 3"`
	TargetSceneID    *uint   `json:"target_scene_id"`
	Pitch            float64 `json:"pitch" binding:"required,min=-90,max=90"`
	Yaw              float64 `json:"yaw" binding:"required,min=-180,max=180"`
	Title            string  `json:"title" binding:"required,min=1,max=100"`
	IconURL          string  `json:"icon_url" binding:"omitempty,url"`
	Style            string  `json:"style" binding:"omitempty,max=50"`
	Content          string  `json:"content"`
	MediaType        string  `json:"media_type" binding:"omitempty,oneof=image video text"`
	MediaURL         string  `json:"media_url" binding:"omitempty,url"`
	Question         string  `json:"question"`
	Options          string  `json:"options"`
	Answer           string  `json:"answer"`
	Score            int     `json:"score" binding:"omitempty,min=0"`
	TransitionEffect string  `json:"transition_effect" binding:"omitempty,oneof=fade zoom none"`
	SortOrder        int     `json:"sort_order"`
}

type UpdateHotspotRequest struct {
	TargetSceneID    *uint   `json:"target_scene_id"`
	Pitch            float64 `json:"pitch" binding:"omitempty,min=-90,max=90"`
	Yaw              float64 `json:"yaw" binding:"omitempty,min=-180,max=180"`
	Title            string  `json:"title" binding:"omitempty,min=1,max=100"`
	IconURL          string  `json:"icon_url" binding:"omitempty,url"`
	Style            string  `json:"style" binding:"omitempty,max=50"`
	Content          string  `json:"content"`
	MediaType        string  `json:"media_type" binding:"omitempty,oneof=image video text"`
	MediaURL         string  `json:"media_url" binding:"omitempty,url"`
	Question         string  `json:"question"`
	Options          string  `json:"options"`
	Answer           string  `json:"answer"`
	Score            int     `json:"score" binding:"omitempty,min=0"`
	TransitionEffect string  `json:"transition_effect" binding:"omitempty,oneof=fade zoom none"`
	SortOrder        int     `json:"sort_order"`
	Status           *uint8  `json:"status" binding:"omitempty"`
}

type HotspotListResponse struct {
	utils.PageInfo `json:"page_info"`
	Hotspots       []*HotspotListItem `json:"hotspots"`
}

type HotspotListItem struct {
	ID               uint      `json:"id"`
	SceneID          uint      `json:"scene_id"`
	Type             uint8     `json:"type"`
	TargetSceneID    *uint     `json:"target_scene_id,omitempty"`
	Pitch            float64   `json:"pitch"`
	Yaw              float64   `json:"yaw"`
	Title            string    `json:"title"`
	IconURL          string    `json:"icon_url"`
	Style            string    `json:"style"`
	TransitionEffect string    `json:"transition_effect"`
	SortOrder        int       `json:"sort_order"`
	Status           uint8     `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type HotspotDetailResponse struct {
	ID               uint      `json:"id"`
	SceneID          uint      `json:"scene_id"`
	Type             uint8     `json:"type"`
	TargetSceneID    *uint     `json:"target_scene_id,omitempty"`
	Pitch            float64   `json:"pitch"`
	Yaw              float64   `json:"yaw"`
	Title            string    `json:"title"`
	IconURL          string    `json:"icon_url"`
	Style            string    `json:"style"`
	Content          string    `json:"content,omitempty"`
	MediaType        string    `json:"media_type,omitempty"`
	MediaURL         string    `json:"media_url,omitempty"`
	Question         string    `json:"question,omitempty"`
	Options          string    `json:"options,omitempty"`
	Answer           string    `json:"answer,omitempty"`
	Score            int       `json:"score"`
	TransitionEffect string    `json:"transition_effect"`
	SortOrder        int       `json:"sort_order"`
	Status           uint8     `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

func ToHotspotListItem(hotspot *model.ResHotspot) *HotspotListItem {
	return &HotspotListItem{
		ID:               hotspot.ID,
		SceneID:          hotspot.SceneID,
		Type:             hotspot.Type,
		TargetSceneID:    hotspot.TargetSceneID,
		Pitch:            hotspot.Pitch,
		Yaw:              hotspot.Yaw,
		Title:            hotspot.Title,
		IconURL:          hotspot.IconURL,
		Style:            hotspot.Style,
		TransitionEffect: hotspot.TransitionEffect,
		SortOrder:        hotspot.SortOrder,
		Status:           hotspot.Status,
		CreatedAt:        hotspot.CreatedAt,
	}
}

func ToHotspotDetailResponse(hotspot *model.ResHotspot) *HotspotDetailResponse {
	return &HotspotDetailResponse{
		ID:               hotspot.ID,
		SceneID:          hotspot.SceneID,
		Type:             hotspot.Type,
		TargetSceneID:    hotspot.TargetSceneID,
		Pitch:            hotspot.Pitch,
		Yaw:              hotspot.Yaw,
		Title:            hotspot.Title,
		IconURL:          hotspot.IconURL,
		Style:            hotspot.Style,
		Content:          hotspot.Content,
		MediaType:        hotspot.MediaType,
		MediaURL:         hotspot.MediaURL,
		Question:         hotspot.Question,
		Options:          hotspot.Options,
		Answer:           hotspot.Answer,
		Score:            hotspot.Score,
		TransitionEffect: hotspot.TransitionEffect,
		SortOrder:        hotspot.SortOrder,
		Status:           hotspot.Status,
		CreatedAt:        hotspot.CreatedAt,
	}
}
