package utils

import (
	"strings"
	"testing"
)

func TestToPinyin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"伏龙观", "fulongguan"},
		{"安澜索桥", "anlansuoqiao"},
		{"二王庙", "erwangmiao"},
		{"鱼嘴分水堤", "yuzuishuifenti"},
		{"飞沙堰", "feishayan"},
		{"宝瓶口", "baopingkou"},
		{"秦堰楼", "qinyanlou"},
		{"玉垒关", "yuleiguan"},
		{"离堆公园", "liduigongyuan"},
		{"都江堰景区大门", "dujiangyanjingqudamen"},
		{"", ""},
		{"ABC", ""}, // 非中文字符返回空
	}

	for _, test := range tests {
		result := ToPinyin(test.input)
		if result != test.expected {
			t.Errorf("ToPinyin(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestToPinyinWithSep(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"伏龙观", "fulongguan"},
		{"都江堰 景区", "dujiangyan-jingqu"},
		{"都江堰-景区", "dujiangyan-jingqu"},
		{"都江堰_景区", "dujiangyan-jingqu"},
		{"二王庙 入口", "erwangmiao-rukou"},
		{"", ""},
	}

	for _, test := range tests {
		result := ToPinyinWithSep(test.input)
		if result != test.expected {
			t.Errorf("ToPinyinWithSep(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestSanitizeSceneCode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Hello_World", "hello-world"},
		{"Hello-World", "hello-world"},
		{"Hello--World", "hello-world"},
		{"-Hello-World-", "hello-world"},
		{"Hello@#$World", "helloworld"},
		{"Hello123", "hello123"},
		{"", ""},
		{"ABC", "abc"},
		{"abc-123-xyz", "abc-123-xyz"},
	}

	for _, test := range tests {
		result := SanitizeSceneCode(test.input)
		if result != test.expected {
			t.Errorf("SanitizeSceneCode(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestGenerateSceneCode(t *testing.T) {
	tests := []struct {
		input         string
		expectedPrefix string
	}{
		{"伏龙观", "fulongguan-"},
		{"安澜索桥", "anlansuoqiao-"},
		{"二王庙", "erwangmiao-"},
		{"都江堰景区大门", "dujiangyanjingqudamen-"},
		{"", "scene-"}, // 空标题使用 scene 前缀
	}

	for _, test := range tests {
		result := GenerateSceneCode(test.input)
		
		// 验证前缀正确
		if !strings.HasPrefix(result, test.expectedPrefix) {
			t.Errorf("GenerateSceneCode(%q) = %q, expected prefix %q", 
				test.input, result, test.expectedPrefix)
		}
		
		// 验证格式：拼音-UUID（8位）
		parts := strings.Split(result, "-")
		if len(parts) < 2 {
			t.Errorf("GenerateSceneCode(%q) = %q, invalid format", test.input, result)
			continue
		}
		
		// 验证 UUID 部分长度为 8
		uuidPart := parts[len(parts)-1]
		if len(uuidPart) != 8 {
			t.Errorf("GenerateSceneCode(%q) = %q, UUID part length %d != 8", 
				test.input, result, len(uuidPart))
		}
		
		// 验证整体格式合法
		if !IsValidSceneCode(result) {
			t.Errorf("GenerateSceneCode(%q) = %q, is not valid scene code", test.input, result)
		}
	}
}

func TestGenerateSceneCodeWithUUID(t *testing.T) {
	title := "伏龙观"
	uuid := "544bfab9-bd61-4c03-bae3-9059d904cd43"
	
	result := GenerateSceneCodeWithUUID(title, uuid)
	expected := "fulongguan-544bfab9"
	
	if result != expected {
		t.Errorf("GenerateSceneCodeWithUUID(%q, %q) = %q, expected %q", 
			title, uuid, result, expected)
	}
}

func TestIsValidSceneCode(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"fulongguan-a1b2c3d4", true},
		{"dujiangyan-jingqu-rukou-255566cb", true},
		{"scene-12345678", true},
		{"abc-123", true},
		{"", false},                    // 空字符串
		{"a", false},                   // 太短
		{"-abc-123", false},            // 以连字符开头
		{"abc-123-", false},            // 以连字符结尾
		{"abc_123", false},             // 包含下划线
		{"ABC-123", false},             // 包含大写字母
		{"abc 123", false},             // 包含空格
		{"abc@123", false},             // 包含特殊字符
		{strings.Repeat("a", 101), false}, // 太长（超过100字符）
	}

	for _, test := range tests {
		result := IsValidSceneCode(test.input)
		if result != test.expected {
			t.Errorf("IsValidSceneCode(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestAddPolyphoneMapping(t *testing.T) {
	// 添加自定义多音字映射
	AddPolyphoneMapping("测试多音字", "ceshiduoyinzi")
	
	result := ToPinyin("测试多音字")
	expected := "ceshiduoyinzi"
	
	if result != expected {
		t.Errorf("After AddPolyphoneMapping, ToPinyin(%q) = %q, expected %q", 
			"测试多音字", result, expected)
	}
}

func TestPolyphoneMap(t *testing.T) {
	// 测试内置多音字映射
	tests := []struct {
		input    string
		expected string
	}{
		{"重庆", "chongqing"},
		{"长安", "changan"},
		{"西藏", "xizang"},
		{"银行", "yinhang"},
	}

	for _, test := range tests {
		result := ToPinyin(test.input)
		if result != test.expected {
			t.Errorf("ToPinyin(%q) = %q, expected %q (polyphone)", test.input, result, test.expected)
		}
	}
}

func TestLongTitle(t *testing.T) {
	// 测试超长标题截断
	longTitle := strings.Repeat("都江堰", 20) // 60 个字符
	result := GenerateSceneCode(longTitle)
	
	// 验证总长度不超过 60（50 + 1 + 8 + 一些连字符）
	if len(result) > 70 {
		t.Errorf("GenerateSceneCode with long title length %d > 70", len(result))
	}
	
	// 验证格式仍然正确
	if !IsValidSceneCode(result) {
		t.Errorf("GenerateSceneCode with long title = %q, is not valid", result)
	}
}
