package utils

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/mozillazg/go-pinyin"
)

var (
	// 多音字映射表，覆盖常见多音字
	polyphoneMap = map[string]string{
		"重庆": "chongqing",
		"长安": "changan",
		"西藏": "xizang",
		"藏族": "zangzu",
		"宝藏": "baozang",
		"行当": "hangdang",
		"银行": "yinhang",
		"行走": "xingzou",
	}

	// 非字母数字字符的正则表达式
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9-]`)
	// 多个连字符的正则表达式
	multipleHyphenRegex = regexp.MustCompile(`-+`)
)

// PinyinConfig 拼音转换配置
var PinyinConfig = pinyin.NewArgs()

func init() {
	// 配置拼音转换选项
	PinyinConfig.Style = pinyin.Normal      // 普通风格，不带声调
	PinyinConfig.Heteronym = false          // 不启用多音字（使用默认读音）
	PinyinConfig.Separator = ""             // 不使用分隔符
}

// ToPinyin 将中文转换为拼音
// 示例："伏龙观" -> "fulongguan"
func ToPinyin(chinese string) string {
	if chinese == "" {
		return ""
	}

	// 先检查多音字映射表
	if pinyin, ok := polyphoneMap[chinese]; ok {
		return pinyin
	}

	// 使用 go-pinyin 转换
	result := pinyin.Pinyin(chinese, PinyinConfig)
	
	var builder strings.Builder
	for _, py := range result {
		if len(py) > 0 {
			builder.WriteString(py[0])
		}
	}
	
	return strings.ToLower(builder.String())
}

// ToPinyinWithSep 将中文转换为拼音，保留分隔符
// 示例："都江堰 景区" -> "dujiangyan-jingqu"
func ToPinyinWithSep(chinese string) string {
	if chinese == "" {
		return ""
	}

	// 替换常见分隔符为空格
	chinese = strings.ReplaceAll(chinese, "_", " ")
	chinese = strings.ReplaceAll(chinese, "-", " ")
	chinese = strings.ReplaceAll(chinese, "·", " ")

	// 分段转换
	parts := strings.Fields(chinese)
	var pinyinParts []string
	
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		// 检查多音字映射
		if pinyin, ok := polyphoneMap[part]; ok {
			pinyinParts = append(pinyinParts, pinyin)
			continue
		}
		
		// 转换这一段
		result := pinyin.Pinyin(part, PinyinConfig)
		var builder strings.Builder
		for _, py := range result {
			if len(py) > 0 {
				builder.WriteString(py[0])
			}
		}
		if pinyinStr := strings.ToLower(builder.String()); pinyinStr != "" {
			pinyinParts = append(pinyinParts, pinyinStr)
		}
	}
	
	return strings.Join(pinyinParts, "-")
}

// SanitizeSceneCode 清理场景编码，确保 URL 安全
// 1. 转换为小写
// 2. 移除非字母数字字符（保留连字符）
// 3. 合并多个连字符
// 4. 移除首尾连字符
func SanitizeSceneCode(code string) string {
	if code == "" {
		return ""
	}
	
	// 转换为小写
	code = strings.ToLower(code)
	
	// 替换下划线和空格为连字符
	code = strings.ReplaceAll(code, "_", "-")
	code = strings.ReplaceAll(code, " ", "-")
	
	// 移除非字母数字字符（保留连字符）
	code = nonAlphanumericRegex.ReplaceAllString(code, "")
	
	// 合并多个连字符
	code = multipleHyphenRegex.ReplaceAllString(code, "-")
	
	// 移除首尾连字符
	code = strings.Trim(code, "-")
	
	return code
}

// GenerateSceneCode 生成场景编码：拼音-UUID前缀
// 示例："伏龙观" -> "fulongguan-a1b2c3d4"
func GenerateSceneCode(title string) string {
	if title == "" {
		return generateShortUUID()
	}
	pinyinPart := ToPinyinWithSep(title)
	pinyinPart = SanitizeSceneCode(pinyinPart)
	if pinyinPart == "" {
		pinyinPart = "scene"
	}
	if len(pinyinPart) > 50 {
		pinyinPart = pinyinPart[:50]
		pinyinPart = strings.TrimSuffix(pinyinPart, "-")
	}
	uuidPart := generateShortUUID()
	return pinyinPart + "-" + uuidPart
}

// GenerateSceneCodeStable 生成稳定计算的场景编码（用于种子数据，防止重启后编码变化导致重复导入）
func GenerateSceneCodeStable(title string) string {
	if title == "" {
		return "scene-default"
	}
	pinyinPart := ToPinyinWithSep(title)
	pinyinPart = SanitizeSceneCode(pinyinPart)
	if pinyinPart == "" {
		pinyinPart = "scene"
	}
	
	// 使用标题的 MD5 前 8 位作为后缀，确保稳定性
	hash := md5.Sum([]byte(title))
	stablePart := hex.EncodeToString(hash[:4])
	
	return pinyinPart + "-" + stablePart
}

// GenerateSceneCodeWithUUID 使用指定的 UUID 生成场景编码
// 用于确保唯一性时重复使用同一个 UUID
func GenerateSceneCodeWithUUID(title string, uuidStr string) string {
	if title == "" {
		return sanitizeUUID(uuidStr)
	}
	
	pinyinPart := ToPinyinWithSep(title)
	pinyinPart = SanitizeSceneCode(pinyinPart)
	
	if pinyinPart == "" {
		pinyinPart = "scene"
	}
	
	if len(pinyinPart) > 50 {
		pinyinPart = pinyinPart[:50]
		pinyinPart = strings.TrimSuffix(pinyinPart, "-")
	}
	
	uuidPart := sanitizeUUID(uuidStr)
	
	return pinyinPart + "-" + uuidPart
}

// generateShortUUID 生成短 UUID（8 位）
func generateShortUUID() string {
	u := uuid.New()
	return strings.ToLower(u.String()[:8])
}

// sanitizeUUID 清理 UUID，只保留前 8 位
func sanitizeUUID(uuidStr string) string {
	// 移除所有连字符
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	// 转小写
	uuidStr = strings.ToLower(uuidStr)
	// 取前 8 位
	if len(uuidStr) > 8 {
		return uuidStr[:8]
	}
	return uuidStr
}

// IsValidSceneCode 验证场景编码是否合法
// 规则：小写字母、数字、连字符，长度 2-100
func IsValidSceneCode(code string) bool {
	if code == "" {
		return false
	}
	
	if len(code) < 2 || len(code) > 100 {
		return false
	}
	
	// 只能包含小写字母、数字和连字符
	for _, r := range code {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	
	// 不能以连字符开头或结尾
	if code[0] == '-' || code[len(code)-1] == '-' {
		return false
	}
	
	return true
}

// AddPolyphoneMapping 添加多音字映射
// 用于扩展多音字映射表
func AddPolyphoneMapping(chinese, pinyin string) {
	polyphoneMap[chinese] = pinyin
}
