package panorama

// FaceNameString 立方体面的标准字符串命名
const (
	FaceNamePX = "px"
	FaceNameNX = "nx"
	FaceNamePY = "py"
	FaceNameNY = "ny"
	FaceNamePZ = "pz"
	FaceNameNZ = "nz"
)

// FaceNameStrings 返回所有面的标准字符串名称列表
func FaceNameStrings() []string {
	return []string{FaceNamePX, FaceNameNX, FaceNamePY, FaceNameNY, FaceNamePZ, FaceNameNZ}
}

// TileXYToLonLat 将瓦片坐标转换为经纬度（用于前端定位）
func TileXYToLonLat(tileX, tileY, level int, tileSize int) (lon, lat float64) {
	tileCount := GetTileCount(level, tileSize)

	normX := float64(tileX) / float64(tileCount)
	normY := float64(tileY) / float64(tileCount)

	lon = (normX - 0.5) * 360
	lat = (0.5 - normY) * 180

	return lon, lat
}

// LonLatToTileXY 将经纬度转换为瓦片坐标
func LonLatToTileXY(lon, lat float64, level int, tileSize int) (tileX, tileY int) {
	tileCount := GetTileCount(level, tileSize)

	normX := (lon / 360) + 0.5
	normY := 0.5 - (lat / 180)

	tileX = int(normX * float64(tileCount))
	tileY = int(normY * float64(tileCount))

	tileX = tileX % tileCount
	if tileX < 0 {
		tileX += tileCount
	}

	tileY = tileY % tileCount
	if tileY < 0 {
		tileY += tileCount
	}

	return tileX, tileY
}

// GetLevelDimension 获取指定层级的像素尺寸
func GetLevelDimension(level int) int {
	return 512 << level
}

// GetTileCount 获取指定层级的瓦片数量
func GetTileCount(level int, tileSize int) int {
	dimension := GetLevelDimension(level)
	return (dimension + tileSize - 1) / tileSize
}
