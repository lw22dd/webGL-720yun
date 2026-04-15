package dto

import (
	"webGL-720yun/internal/model"
)

type GraphDataResponse struct {
	SpaceInfo *SpaceInfoForGraph `json:"space_info"`
	Nodes     []*SceneNodeData   `json:"nodes"`
	Edges     []*EdgeData        `json:"edges"`
}

type SpaceInfoForGraph struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	ZoomLevel   int     `json:"zoom_level"`
}

type SceneNodeData struct {
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	SceneCode    string  `json:"scene_code"`
	ThumbnailURL string  `json:"thumbnail_url"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	HasPosition  bool    `json:"has_position"`
	ViewCount    int64   `json:"view_count"`
	Status       string  `json:"status"`
}

type EdgeData struct {
	ID           uint   `json:"id"`
	SourceID     uint   `json:"source_id"`
	TargetID     uint   `json:"target_id"`
	HotspotID    uint   `json:"hotspot_id"`
	HotspotTitle string `json:"hotspot_title"`
	Type         string `json:"type"`
}

type UpdatePositionRequest struct {
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
}

type BatchUpdatePositionRequest struct {
	Positions []ScenePosition `json:"positions" binding:"required,min=1"`
}

type ScenePosition struct {
	SceneID   uint    `json:"scene_id" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
}

func ToSpaceInfoForGraph(space *model.ResSpace) *SpaceInfoForGraph {
	return &SpaceInfoForGraph{
		ID:        space.ID,
		Name:      space.Name,
		Longitude: space.Longitude,
		Latitude:  space.Latitude,
		ZoomLevel: space.ZoomLevel,
	}
}

func ToSceneNodeData(scene *model.ResScene) *SceneNodeData {
	hasPosition := scene.Longitude != 0 || scene.Latitude != 0
	status := "pending"
	if hasPosition {
		status = "placed"
	}
	return &SceneNodeData{
		ID:           scene.ID,
		Title:        scene.Title,
		SceneCode:    scene.SceneCode,
		ThumbnailURL: scene.ThumbnailURL,
		Longitude:    scene.Longitude,
		Latitude:     scene.Latitude,
		HasPosition:  hasPosition,
		ViewCount:    scene.ViewCount,
		Status:       status,
	}
}

func ToEdgeData(hotspot *model.ResHotspot) *EdgeData {
	targetID := uint(0)
	if hotspot.TargetSceneID != nil {
		targetID = *hotspot.TargetSceneID
	}
	edgeType := "walk"
	return &EdgeData{
		ID:           hotspot.ID,
		SourceID:     hotspot.SceneID,
		TargetID:     targetID,
		HotspotID:    hotspot.ID,
		HotspotTitle: hotspot.Title,
		Type:         edgeType,
	}
}
