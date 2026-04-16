package setup

import (
	"path/filepath"
	"webGL-720yun/pkg/utils"
)

type SceneSeed struct {
	Title     string
	SceneCode string
	FileName  string
}

type ScenicSpotSeed struct {
	Name        string
	Slug        string
	LocalPath   string
	CoverFile   string
	Description string
	Province    string
	City        string
	Longitude   float64
	Latitude    float64
	Scenes      []SceneSeed
}

// generateSceneCode 根据标题生成拼音语义化的 SceneCode
func generateSceneCode(title string) string {
	return utils.GenerateSceneCodeStable(title)
}

func GetScenicSpotSeeds() []ScenicSpotSeed {
	return []ScenicSpotSeed{
		{
			Name:        "都江堰景区",
			Slug:        "dujiangyan",
			LocalPath:   filepath.Join("..", "参考", "都江堰"),
			CoverFile:   "cover.jpg",
			Description: "世界文化遗产，古代水利工程杰作，始建于公元前256年，是全世界迄今为止年代最久、唯一留存、以无坝引水为特征的宏大水利工程。",
			Province:    "四川省",
			City:        "成都市",
			Longitude:   103.611379,
			Latitude:    31.001752,
			Scenes: []SceneSeed{
				{Title: "都江堰景区大门", SceneCode: generateSceneCode("都江堰景区大门"), FileName: "都江堰景区大门_3576761_sphere.jpg"},
				{Title: "伏龙观", SceneCode: generateSceneCode("伏龙观"), FileName: "伏龙观_3576621_sphere.jpg"},
				{Title: "鱼嘴分水堤1", SceneCode: generateSceneCode("鱼嘴分水堤1"), FileName: "鱼嘴分水堤1_3576420_sphere.jpg"},
				{Title: "安澜索桥", SceneCode: generateSceneCode("安澜索桥"), FileName: "安澜索桥_3576562_sphere.jpg"},
				{Title: "二王庙", SceneCode: generateSceneCode("二王庙"), FileName: "二王庙_3576708_sphere.jpg"},
				{Title: "禹王宫", SceneCode: generateSceneCode("禹王宫"), FileName: "禹王宫_3576391_sphere.jpg"},
				{Title: "秦堰楼", SceneCode: generateSceneCode("秦堰楼"), FileName: "秦堰楼_3576546_sphere.jpg"},
				{Title: "玉垒关观赏平台", SceneCode: generateSceneCode("玉垒关观赏平台"), FileName: "玉垒关观赏平台_3576333_sphere.jpg"},
				{Title: "都江堰后门-白景", SceneCode: generateSceneCode("都江堰后门-白景"), FileName: "都江堰后门-白景_3576741_sphere.jpg"},
				{Title: "飞沙堰", SceneCode: generateSceneCode("飞沙堰"), FileName: "飞沙堰_3576638_sphere.jpg"},
			},
		},
	}
}

// SceneCoordinates 景点坐标映射表（WGS84坐标系）
// 使用标题作为 key，因为 SceneCode 现在是动态生成的
var SceneCoordinates = map[string]struct {
	Longitude float64
	Latitude  float64
}{
	// 离堆/大门区域
	"都江堰景区大门": {Longitude: 103.615379, Latitude: 30.995792},
	"景区大门":    {Longitude: 103.615379, Latitude: 30.995792},
	"卧铁":      {Longitude: 103.614500, Latitude: 30.996500},
	"清溪园":     {Longitude: 103.613800, Latitude: 30.997200},
	"天府源茶馆":   {Longitude: 103.613200, Latitude: 30.997800},
	"堰功道":     {Longitude: 103.612800, Latitude: 30.998200},
	"张松银杏":    {Longitude: 103.612500, Latitude: 30.998600},
	"伏龙观":     {Longitude: 103.612200, Latitude: 30.999000},
	"伏龙观内景":   {Longitude: 0, Latitude: 0},
	"伏龙观观景台":  {Longitude: 0, Latitude: 0},
	"离堆":      {Longitude: 103.611800, Latitude: 30.999500},

	// 鱼嘴/索桥区域
	"飞沙堰":   {Longitude: 103.608000, Latitude: 31.001000},
	"金刚堤":   {Longitude: 103.607000, Latitude: 31.003000}, 
	"鱼嘴分水堤2": {Longitude: 103.607000, Latitude: 31.003000}, // 映射金刚堤
	"安澜索桥": {Longitude: 103.609000, Latitude: 31.006000},
	"安澜桥":   {Longitude: 0, Latitude: 0},
	"天下爱情第一桥": {Longitude: 0, Latitude: 0},
	"鱼嘴":     {Longitude: 103.605000, Latitude: 31.009000},
	"鱼嘴分水堤1": {Longitude: 103.605000, Latitude: 31.009000}, // 映射鱼嘴

	// 山上/庙宇区域
	"秦堰楼": {Longitude: 103.607000, Latitude: 31.010000},
	"观景台": {Longitude: 0, Latitude: 0},
	"二王庙": {Longitude: 103.610500, Latitude: 31.006000},
	"二王庙入口": {Longitude: 0, Latitude: 0},
	"二王庙内景": {Longitude: 0, Latitude: 0},
	"二王庙道观": {Longitude: 0, Latitude: 0},
	"禹王宫":  {Longitude: 103.611862, Latitude: 31.003596},

	// 古道/广场区域
	"灵动森林":  {Longitude: 103.611500, Latitude: 31.004500},
	"善水阁":   {Longitude: 103.611500, Latitude: 31.004500}, // 映射灵动森林
	"敬修之牌坊": {Longitude: 103.612500, Latitude: 31.003000},
	"商铺街":   {Longitude: 103.612500, Latitude: 31.003000}, // 映射敬修之牌坊
	"松茂古道":  {Longitude: 103.613500, Latitude: 31.001500},
	"神木艺术馆": {Longitude: 103.613500, Latitude: 31.001500}, // 映射松茂古道
	"玉垒阁":   {Longitude: 103.614500, Latitude: 31.000500},
	"玉垒关观赏平台": {Longitude: 103.614500, Latitude: 31.000500}, // 映射玉垒阁
	"玉垒关":   {Longitude: 0, Latitude: 0},
	"玉垒殿":   {Longitude: 0, Latitude: 0},
	"宣威门":   {Longitude: 0, Latitude: 0},
	"城隍庙":   {Longitude: 103.615000, Latitude: 31.000000},
	"十殿":     {Longitude: 103.615500, Latitude: 30.999500},
	"玉垒山广场": {Longitude: 103.616000, Latitude: 30.999000},
	"南桥-夜景": {Longitude: 0, Latitude: 0},
	"西关":    {Longitude: 103.612581, Latitude: 30.999904},

	// 其他保持0,0的条目
	"太极殿":   {Longitude: 0, Latitude: 0},
	"石林":    {Longitude: 0, Latitude: 0},
	"门犀亭":   {Longitude: 0, Latitude: 0},
	"三官殿":   {Longitude: 0, Latitude: 0},
	"魁星点斗": {Longitude: 0, Latitude: 0},
	"东苑":    {Longitude: 0, Latitude: 0},
	"马王殿":   {Longitude: 0, Latitude: 0},
	"斗犀亭":   {Longitude: 0, Latitude: 0},
	"财神殿":   {Longitude: 0, Latitude: 0},
	"都江堰后门-白景": {Longitude: 103.607628, Latitude: 31.008030},
	"都江堰后门-夜景": {Longitude: 0, Latitude: 0},
}

// GetSceneCoordinate 获取指定场景标题的坐标
func GetSceneCoordinate(title string) (longitude, latitude float64, found bool) {
	if coord, ok := SceneCoordinates[title]; ok {
		return coord.Longitude, coord.Latitude, true
	}
	return 0, 0, false
}
