package hot

import (
	"time"
)

// Fetcher 热点获取器接口
type Fetcher interface {
	Platform() string
	Fetch() ([]HotItem, error)
}

// HotItem 热点条目
type HotItem struct {
	Rank      int       `json:"rank"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	HotValue  float64   `json:"hotValue"`
	Category  string    `json:"category"`
	FetchTime time.Time `json:"fetchTime"`
}

// Result 获取结果
type Result struct {
	Platform  string    `json:"platform"`
	Items     []HotItem `json:"items"`
	FetchTime time.Time `json:"fetchTime"`
	Error     string    `json:"error,omitempty"`
}

// Platform 平台常量
const (
	PlatformWeibo   = "weibo"
	PlatformZhihu   = "zhihu"
	PlatformBaidu   = "baidu"
	PlatformToutiao = "toutiao"
	PlatformDouyin  = "douyin"
)

// PlatformNames 平台名称映射
var PlatformNames = map[string]string{
	PlatformWeibo:   "微博",
	PlatformZhihu:   "知乎",
	PlatformBaidu:   "百度",
	PlatformToutiao: "今日头条",
	PlatformDouyin:  "抖音",
}

// Manager 热点管理器
type Manager struct {
	fetchers map[string]Fetcher
}

// NewManager 创建热点管理器
func NewManager() *Manager {
	return &Manager{
		fetchers: make(map[string]Fetcher),
	}
}

// Register 注册获取器
func (m *Manager) Register(fetcher Fetcher) {
	m.fetchers[fetcher.Platform()] = fetcher
}

// GetFetcher 获取获取器
func (m *Manager) GetFetcher(platform string) (Fetcher, bool) {
	f, ok := m.fetchers[platform]
	return f, ok
}

// GetPlatforms 获取所有平台
func (m *Manager) GetPlatforms() []string {
	platforms := make([]string, 0, len(m.fetchers))
	for p := range m.fetchers {
		platforms = append(platforms, p)
	}
	return platforms
}

// FetchAll 获取所有平台热点
func (m *Manager) FetchAll() []Result {
	var results []Result
	for _, fetcher := range m.fetchers {
		items, err := fetcher.Fetch()
		result := Result{
			Platform:  fetcher.Platform(),
			Items:     items,
			FetchTime: time.Now(),
		}
		if err != nil {
			result.Error = err.Error()
		}
		results = append(results, result)
	}
	return results
}

// FetchPlatforms 获取指定平台热点
func (m *Manager) FetchPlatforms(platforms []string) []Result {
	var results []Result
	for _, platform := range platforms {
		if fetcher, ok := m.fetchers[platform]; ok {
			items, err := fetcher.Fetch()
			result := Result{
				Platform:  platform,
				Items:     items,
				FetchTime: time.Now(),
			}
			if err != nil {
				result.Error = err.Error()
			}
			results = append(results, result)
		}
	}
	return results
}
