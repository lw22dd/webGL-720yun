package model

import "fmt"

const (
	SpacePrefix = "spaces/%s"

	SourcePathFormat  = SpacePrefix + "/sources/%s/source.jpg"
	PreviewPathFormat = SpacePrefix + "/previews/%s/preview.jpg"
	TilePathFormat    = SpacePrefix + "/tiles/%s/cubemap/%s/level_%d/tile_%d_%d.jpg"
	TileBaseFormat    = SpacePrefix + "/tiles/%s/"
)

func GetSceneSourcePath(spaceSlug, fileID string) string {
	return fmt.Sprintf(SourcePathFormat, spaceSlug, fileID)
}

func GetScenePreviewPath(spaceSlug, sceneCode string) string {
	return fmt.Sprintf(PreviewPathFormat, spaceSlug, sceneCode)
}

func GetSceneTilePath(spaceSlug, sceneCode, face string, level, x, y int) string {
	return fmt.Sprintf(TilePathFormat, spaceSlug, sceneCode, face, level, x, y)
}

func GetSceneTileBasePath(spaceSlug, sceneCode string) string {
	return fmt.Sprintf(TileBaseFormat, spaceSlug, sceneCode)
}

func ParseTilePath(path string) (spaceSlug, sceneCode, face string, level, row, col int, err error) {
	_, err = fmt.Sscanf(path, TilePathFormat,
		&spaceSlug, &sceneCode, &face, &level, &row, &col)
	return
}
