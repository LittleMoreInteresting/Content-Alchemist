package hot

import (
	"fmt"
	"net/http"
	"time"
)

// ToutiaoFetcher 头条热榜获取器
type ToutiaoFetcher struct {
	client *http.Client
}

// NewToutiaoFetcher 创建头条获取器
func NewToutiaoFetcher() *ToutiaoFetcher {
	return &ToutiaoFetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Platform 返回平台名称
func (f *ToutiaoFetcher) Platform() string {
	return PlatformToutiao
}

// Fetch 获取头条热榜（由于头条API限制，返回模拟数据）
func (f *ToutiaoFetcher) Fetch() ([]HotItem, error) {
	// 头条热榜需要特殊处理，这里返回模拟数据
	return f.fetchMock(), nil
}

// fetchMock 获取模拟数据
func (f *ToutiaoFetcher) fetchMock() []HotItem {
	mockData := []struct {
		Title    string
		Category string
	}{
		{"今日头条独家报道", "新闻"},
		{"本地生活新鲜事", "本地"},
		{"美食探店推荐", "美食"},
		{"搞笑视频精选", "搞笑"},
		{"汽车测评报告", "汽车"},
		{"时尚穿搭指南", "时尚"},
		{"育儿经验分享", "育儿"},
		{"宠物可爱瞬间", "宠物"},
		{"职场技能提升", "职场"},
		{"历史故事揭秘", "历史"},
	}
	
	var items []HotItem
	for i, data := range mockData {
		items = append(items, HotItem{
			Rank:      i + 1,
			Title:     data.Title,
			URL:       fmt.Sprintf("https://www.toutiao.com/search/?keyword=%s", data.Title),
			HotValue:  float64(1000000 - (i+1)*52000),
			Category:  data.Category,
			FetchTime: time.Now(),
		})
	}
	
	return items
}
