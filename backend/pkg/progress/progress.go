package progress

import (
	"fmt"
	"os"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// Manager 管理多个进度条
type Manager struct {
	bars map[string]*progressbar.ProgressBar
	mu   sync.RWMutex
}

// NewManager 创建新的进度条管理器
func NewManager() *Manager {
	return &Manager{
		bars: make(map[string]*progressbar.ProgressBar),
	}
}

// CreateBar 创建通用进度条
// max: 最大值（字节数或100等）
// description: 描述文字
// showBytes: 是否显示字节格式（上传/下载用true，处理任务用false）
func (m *Manager) CreateBar(id string, max int64, description string, showBytes bool) *progressbar.ProgressBar {
	opts := []progressbar.Option{
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionClearOnFinish(),
	}

	if showBytes {
		opts = append(opts,
			progressbar.OptionShowBytes(true),
			progressbar.OptionShowCount(),
		)
	}

	bar := progressbar.NewOptions64(max, opts...)

	m.mu.Lock()
	m.bars[id] = bar
	m.mu.Unlock()

	return bar
}

// GetBar 获取已存在的进度条
func (m *Manager) GetBar(id string) *progressbar.ProgressBar {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.bars[id]
}

// RemoveBar 移除进度条
func (m *Manager) RemoveBar(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if bar, ok := m.bars[id]; ok {
		bar.Close()
		delete(m.bars, id)
	}
}

// PrintStart 打印任务开始信息
func PrintStart(taskType string, name string, details string) {
	fmt.Printf("\n[%s] 开始%s: %s %s\n", taskType, getActionName(taskType), name, details)
}

// PrintComplete 打印任务完成信息
func PrintComplete(taskType string, name string) {
	fmt.Printf("[%s] %s完成: %s\n", taskType, getActionName(taskType), name)
}

// PrintInfo 打印任务信息
func PrintInfo(taskType string, message string) {
	fmt.Printf("[%s] %s\n", taskType, message)
}

func getActionName(taskType string) string {
	switch taskType {
	case "上传":
		return "上传"
	case "切片":
		return "处理"
	case "下载":
		return "下载"
	default:
		return "处理"
	}
}

// GlobalManager 全局进度条管理器实例
var GlobalManager = NewManager()
