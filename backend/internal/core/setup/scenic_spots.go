package setup

<<<<<<< HEAD
import "path/filepath"
=======
import (
	"path/filepath"
	"webGL-720yun/pkg/utils"
)
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503

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

<<<<<<< HEAD
=======
// generateSceneCode 根据标题生成拼音语义化的 SceneCode
func generateSceneCode(title string) string {
	return utils.GenerateSceneCodeStable(title)
}

>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
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
<<<<<<< HEAD
				{Title: "玉垒关观赏平台", SceneCode: "3576333", FileName: "玉垒关观赏平台_3576333_sphere.jpg"},
				{Title: "玉垒关", SceneCode: "3576350", FileName: "玉垒关_3576350_sphere.jpg"},
				{Title: "玉垒殿", SceneCode: "3576369", FileName: "玉垒殿_3576369_sphere.jpg"},
				{Title: "禹王宫", SceneCode: "3576391", FileName: "禹王宫_3576391_sphere.jpg"},
				{Title: "鱼嘴分水堤1", SceneCode: "3576420", FileName: "鱼嘴分水堤1_3576420_sphere.jpg"},
				{Title: "鱼嘴分水堤2", SceneCode: "3576400", FileName: "鱼嘴分水堤2_3576400_sphere.jpg"},
				{Title: "宣威门", SceneCode: "3576433", FileName: "宣威门_3576433_sphere.jpg"},
				{Title: "天下爱情第一桥", SceneCode: "3576465", FileName: "天下爱情第一桥_3576465_sphere.jpg"},
				{Title: "西关", SceneCode: "3576455", FileName: "西关_3576455_sphere.jpg"},
				{Title: "太极殿", SceneCode: "3576475", FileName: "太极殿_3576475_sphere.jpg"},
				{Title: "石林", SceneCode: "3576487", FileName: "石林_3576487_sphere.jpg"},
				{Title: "善水阁", SceneCode: "3576515", FileName: "善水阁_3576515_sphere.jpg"},
				{Title: "商铺街", SceneCode: "3576503", FileName: "商铺街_3576503_sphere.jpg"},
				{Title: "神木艺术馆", SceneCode: "3576494", FileName: "神木艺术馆_3576494_sphere.jpg"},
				{Title: "安澜索桥", SceneCode: "3576562", FileName: "安澜索桥_3576562_sphere.jpg"},
				{Title: "门犀亭", SceneCode: "3576567", FileName: "门犀亭_3576567_sphere.jpg"},
				{Title: "三官殿", SceneCode: "3576534", FileName: "三官殿_3576534_sphere.jpg"},
				{Title: "秦堰楼", SceneCode: "3576546", FileName: "秦堰楼_3576546_sphere.jpg"},
				{Title: "魁星点斗", SceneCode: "3576583", FileName: "魁星点斗_3576583_sphere.jpg"},
				{Title: "观景台", SceneCode: "3576591", FileName: "观景台_3576591_sphere.jpg"},
				{Title: "伏龙观内景", SceneCode: "3576597", FileName: "伏龙观内景_3576597_sphere.jpg"},
				{Title: "伏龙观观景台", SceneCode: "3576605", FileName: "伏龙观观景台_3576605_sphere.jpg"},
				{Title: "伏龙观", SceneCode: "3576621", FileName: "伏龙观_3576621_sphere.jpg"},
				{Title: "飞沙堰", SceneCode: "3576638", FileName: "飞沙堰_3576638_sphere.jpg"},
				{Title: "二王庙入口", SceneCode: "3576650", FileName: "二王庙入口_3576650_sphere.jpg"},
				{Title: "二王庙内景", SceneCode: "3576658", FileName: "二王庙内景_3576658_sphere.jpg"},
				{Title: "二王庙道观", SceneCode: "3576681", FileName: "二王庙道观_3576681_sphere.jpg"},
				{Title: "二王庙", SceneCode: "3576708", FileName: "二王庙_3576708_sphere.jpg"},
				{Title: "东苑", SceneCode: "3576750", FileName: "东苑_3576750_sphere.jpg"},
				{Title: "马王殿", SceneCode: "3576575", FileName: "马王殿_3576575_sphere.jpg"},
				{Title: "斗犀亭", SceneCode: "3576730", FileName: "斗犀亭_3576730_sphere.jpg"},
				{Title: "财神殿", SceneCode: "3576767", FileName: "财神殿_3576767_sphere.jpg"},
				{Title: "安澜桥", SceneCode: "3576775", FileName: "安澜桥_3576775_sphere.jpg"},
				{Title: "南桥-夜景", SceneCode: "3576780", FileName: "南桥-夜景_3576780_sphere.jpg"},
				{Title: "都江堰景区大门", SceneCode: "3576761", FileName: "都江堰景区大门_3576761_sphere.jpg"},
				{Title: "都江堰后门-白景", SceneCode: "3576741", FileName: "都江堰后门-白景_3576741_sphere.jpg"},
				{Title: "都江堰后门-夜景", SceneCode: "3576798", FileName: "都江堰后门-夜景_3576798_sphere.jpg"},
=======
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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
			},
		},
	}
}

// SceneCoordinates 景点坐标映射表（WGS84坐标系）
<<<<<<< HEAD
=======
// 使用标题作为 key，因为 SceneCode 现在是动态生成的
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
var SceneCoordinates = map[string]struct {
	Longitude float64
	Latitude  float64
}{
<<<<<<< HEAD
	// 1. 玉垒关观赏平台 | 103.611379,31.001752
	"3576333": {Longitude: 103.611379, Latitude: 31.001752},
	// 2. 玉垒关 | 103.611379,31.001752
	"3576350": {Longitude: 103.611379, Latitude: 31.001752},
	// 3. 玉垒殿 | 0,0 未找到
	"3576369": {Longitude: 0, Latitude: 0},
	// 4. 禹王宫 | 103.611862,31.003596
	"3576391": {Longitude: 103.611862, Latitude: 31.003596},
	// 5. 鱼嘴分水堤1 | 103.604640,31.008830
	"3576420": {Longitude: 103.604640, Latitude: 31.008830},
	// 6. 鱼嘴分水堤2 | 103.608435,31.001146
	"3576400": {Longitude: 103.608435, Latitude: 31.001146},
	// 7. 宣威门 | 0,0 未找到
	"3576433": {Longitude: 0, Latitude: 0},
	// 8. 天下爱情第一桥 | 103.608945,31.006044
	"3576465": {Longitude: 103.608945, Latitude: 31.006044},
	// 9. 西关 | 103.612581,30.999904
	"3576455": {Longitude: 103.612581, Latitude: 30.999904},
	// 10. 太极殿 | 0,0 未找到
	"3576475": {Longitude: 0, Latitude: 0},
	// 11. 石林 | 0,0 未找到
	"3576487": {Longitude: 0, Latitude: 0},
	// 12. 善水阁 | 0,0 未找到
	"3576515": {Longitude: 0, Latitude: 0},
	// 13. 商铺街 | 0,0 未找到
	"3576503": {Longitude: 0, Latitude: 0},
	// 14. 神木艺术馆 | 0,0 未找到
	"3576494": {Longitude: 0, Latitude: 0},
	// 15. 安澜索桥 | 103.609389,31.005997
	"3576562": {Longitude: 103.609389, Latitude: 31.005997},
	// 16. 门犀亭 | 0,0 未找到
	"3576567": {Longitude: 0, Latitude: 0},
	// 17. 三官殿 | 0,0 未找到
	"3576534": {Longitude: 0, Latitude: 0},
	// 18. 秦堰楼 | 103.606649,31.010009
	"3576546": {Longitude: 103.606649, Latitude: 31.010009},
	// 19. 魁星点斗 | 0,0 未找到
	"3576583": {Longitude: 0, Latitude: 0},
	// 20. 观景台 | 103.606649,31.010009
	"3576591": {Longitude: 103.606649, Latitude: 31.010009},
	// 21. 伏龙观内景 | 103.612983,30.997706
	"3576597": {Longitude: 103.612983, Latitude: 30.997706},
	// 22. 伏龙观观景台 | 103.612983,30.997706
	"3576605": {Longitude: 103.612983, Latitude: 30.997706},
	// 23. 伏龙观 | 103.612983,30.997706
	"3576621": {Longitude: 103.612983, Latitude: 30.997706},
	// 24. 飞沙堰 | 103.607951,31.001338
	"3576638": {Longitude: 103.607951, Latitude: 31.001338},
	// 25. 二王庙入口 | 103.610210,31.005773
	"3576650": {Longitude: 103.610210, Latitude: 31.005773},
	// 26. 二王庙内景 | 103.610210,31.005773
	"3576658": {Longitude: 103.610210, Latitude: 31.005773},
	// 27. 二王庙道观 | 103.610210,31.005773
	"3576681": {Longitude: 103.610210, Latitude: 31.005773},
	// 28. 二王庙 | 103.610210,31.005773
	"3576708": {Longitude: 103.610210, Latitude: 31.005773},
	// 29. 东苑 | 0,0 未找到
	"3576750": {Longitude: 0, Latitude: 0},
	// 30. 马王殿 | 0,0 未找到
	"3576575": {Longitude: 0, Latitude: 0},
	// 31. 斗犀亭 | 0,0 未找到
	"3576730": {Longitude: 0, Latitude: 0},
	// 32. 财神殿 | 0,0 未找到
	"3576767": {Longitude: 0, Latitude: 0},
	// 33. 安澜桥 | 103.609389,31.005997
	"3576775": {Longitude: 103.609389, Latitude: 31.005997},
	// 34. 南桥-夜景 | 103.613595,30.998393
	"3576780": {Longitude: 103.613595, Latitude: 30.998393},
	// 35. 都江堰景区大门 | 103.615379,30.995792
	"3576761": {Longitude: 103.615379, Latitude: 30.995792},
	// 36. 都江堰后门-白景 | 103.607628,31.008030
	"3576741": {Longitude: 103.607628, Latitude: 31.008030},
	// 37. 都江堰后门-夜景 | 103.607628,31.008030
	"3576798": {Longitude: 103.607628, Latitude: 31.008030},
}

// GetSceneCoordinate 获取指定场景码的坐标
func GetSceneCoordinate(sceneCode string) (longitude, latitude float64, found bool) {
	if coord, ok := SceneCoordinates[sceneCode]; ok {
=======
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
>>>>>>> 8146554307dc850256079e5aa35fb05bb5b6a503
		return coord.Longitude, coord.Latitude, true
	}
	return 0, 0, false
}
