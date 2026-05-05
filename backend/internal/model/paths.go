package model

import "fmt"

const (
	SpacePrefix = "spaces/%s"

	// 空间/景区封面图
	SpaceCoverPathFormat = SpacePrefix + "/covers/cover.jpg"

	// 原始全景图
	SourcePathFormat = SpacePrefix + "/sources/%s/source.jpg"

	// 预览图
	PreviewPathFormat = SpacePrefix + "/previews/%s/preview.jpg"

	// 瓦片图：spaces/{space_name}/tiles/{scene_code}/cubemap/{face}/level_{level}/tile_{x}_{y}.jpg
	TilePathFormat = SpacePrefix + "/tiles/%s/cubemap/%s/level_%d/tile_%d_%d.jpg"

	// 瓦片根目录（用于数据库记录）
	TileBaseFormat = SpacePrefix + "/tiles/%s/"
)

// GetSpaceCoverPath 获取空间封面图路径
func GetSpaceCoverPath(spaceSlug string) string {
	return fmt.Sprintf(SpaceCoverPathFormat, spaceSlug)
}

// GetSceneSourcePath 获取场景源文件路径
func GetSceneSourcePath(spaceSlug, fileID string) string {
	return fmt.Sprintf(SourcePathFormat, spaceSlug, fileID)
}

// GetScenePreviewPath 获取场景预览图路径
func GetScenePreviewPath(spaceSlug, sceneCode string) string {
	return fmt.Sprintf(PreviewPathFormat, spaceSlug, sceneCode)
}

// GetSceneTilePath 获取场景特定瓦片路径
// 注意：x, y 顺序必须与文件名保持一致
func GetSceneTilePath(spaceSlug, sceneCode, face string, level, x, y int) string {
	return fmt.Sprintf(TilePathFormat, spaceSlug, sceneCode, face, level, x, y)
}

// GetSceneTileBasePath 获取场景瓦片根目录
func GetSceneTileBasePath(spaceSlug, sceneCode string) string {
	return fmt.Sprintf(TileBaseFormat, spaceSlug, sceneCode)
}
