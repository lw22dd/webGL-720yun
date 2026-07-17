package setup

import (
	"webGL-720yun/internal/model"
	"webGL-720yun/pkg/logger"

	"gorm.io/gorm"
)

// HotspotSeed 热点种子数据
// 用 SourceTitle/TargetTitle 而非 SceneCode 关联场景：
// 因为 SceneCode 含 MD5 后缀，跨文件硬编码易出错；
// 这里通过 generateSceneCode() 动态生成，与 scenic_spots.go 保持一致。
type HotspotSeed struct {
	SourceTitle      string  // 源场景中文标题（热点所属场景）
	Title            string  // 热点标题
	Type             uint8   // 1=场景切换 2=教学 3=答题
	TargetTitle      string  // 目标场景中文标题（Type=1 时必填）
	Pitch            float64 // 俯仰角（弧度，PSV 内部直接使用）
	Yaw              float64 // 偏航角（弧度，PSV 内部直接使用）
	IconURL          string  // 自定义图标 URL，留空用 style 默认
	Style            string  // 预设图标 key（与前端 HOTSPOT_ICON_PRESETS 对齐）
	TransitionEffect string  // 转场效果：fade / zoom / none
}

// GetHotspotSeeds 默认热点种子列表
// 用于开发/演示环境初始化边数据，让小地图有向图可见
func GetHotspotSeeds() []HotspotSeed {
	return []HotspotSeed{
		{
			SourceTitle:      "安澜索桥",
			Title:            "前往都江堰后门",
			Type:             model.HotspotTypeSwitch,
			TargetTitle:      "都江堰后门-白景",
			Pitch:            -0.26,
			Yaw:              2.55,
			Style:            "arrow-blue",
			TransitionEffect: "fade",
		},
	}
}

// SeedHotspotsIfNeeded 为已存在的场景创建默认热点
// 幂等：使用 (scene_id, title) 联合唯一查重
// 调用前置：SeedResourcesIfNeeded 已完成
func SeedHotspotsIfNeeded(db *gorm.DB) error {
	logger.Info("🎯 开始初始化默认热点...")

	seeds := GetHotspotSeeds()

	// 1. 预加载所有场景的 Title -> ID 映射
	// 注意：用 Title 而非 SceneCode 关联，避免 hash 后缀导致匹配失败
	type sceneRow struct {
		ID    uint
		Title string
	}
	var sceneRows []sceneRow
	if err := db.Model(&model.ResScene{}).Select("id, title").Find(&sceneRows).Error; err != nil {
		return err
	}
	titleToID := make(map[string]uint, len(sceneRows))
	for _, s := range sceneRows {
		titleToID[s.Title] = s.ID
	}

	created := 0
	skipped := 0
	for _, seed := range seeds {
		sourceID, ok := titleToID[seed.SourceTitle]
		if !ok {
			logger.Warnf("⚠️  源场景 %s 不存在，跳过热点 %q", seed.SourceTitle, seed.Title)
			skipped++
			continue
		}

		// 幂等检查：同场景下同标题不重复
		var existing model.ResHotspot
		err := db.Where("scene_id = ? AND title = ?", sourceID, seed.Title).First(&existing).Error
		if err == nil {
			skipped++
			continue
		}
		if err != gorm.ErrRecordNotFound {
			logger.Errorf("查询热点失败: %v", err)
			continue
		}

		hotspot := model.ResHotspot{
			SceneID: sourceID,
			Type:    seed.Type,
			Title:   seed.Title,
			Pitch:   seed.Pitch,
			Yaw:     seed.Yaw,
			IconURL: seed.IconURL,
			Style:   seed.Style,
			Status:  1,
		}

		if seed.Type == model.HotspotTypeSwitch && seed.TargetTitle != "" {
			targetID, ok := titleToID[seed.TargetTitle]
			if !ok {
				logger.Warnf("⚠️  目标场景 %s 不存在，跳过热点 %q", seed.TargetTitle, seed.Title)
				skipped++
				continue
			}
			hotspot.TargetSceneID = &targetID
			hotspot.TransitionEffect = seed.TransitionEffect
		}

		if err := db.Create(&hotspot).Error; err != nil {
			logger.Errorf("创建热点 %q 失败: %v", seed.Title, err)
			continue
		}
		created++
		logger.Infof("✅ 创建热点: %s (源场景=%s, ID=%d)", seed.Title, seed.SourceTitle, sourceID)
	}

	logger.Infof("🎯 热点初始化完成：新建=%d, 跳过=%d", created, skipped)
	return nil
}
