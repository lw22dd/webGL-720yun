package setup

import "path/filepath"

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
			Longitude:   103.5956,
			Latitude:    30.9388,
			Scenes: []SceneSeed{
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
			},
		},
	}
}
