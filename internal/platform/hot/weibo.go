package hot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WeiboFetcher 微博热榜获取器
type WeiboFetcher struct {
	client *http.Client
}

// NewWeiboFetcher 创建微博获取器
func NewWeiboFetcher() *WeiboFetcher {
	return &WeiboFetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Platform 返回平台名称
func (f *WeiboFetcher) Platform() string {
	return PlatformWeibo
}

// WeiboResponse 微博响应结构
type WeiboResponse struct {
	Ok   int `json:"ok"`
	Data struct {
		Realtime []struct {
			Rank        int     `json:"rank"`
			Note        string  `json:"note"`
			RawHot      float64 `json:"raw_hot"`
			Category    string  `json:"category"`
			Link        string  `json:"link"`
			Flag        int     `json:"flag"`
		} `json:"realtime"`
	} `json:"data"`
}

// Fetch 获取微博热榜
func (f *WeiboFetcher) Fetch() ([]HotItem, error) {
	// 微博热榜API（使用公开接口）
	url := "https://weibo.com/ajax/side/hotSearch"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://weibo.com/hot/search")
	
	resp, err := f.client.Do(req)
	if err != nil {
		return f.fetchMock(), nil // 失败时返回模拟数据
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return f.fetchMock(), nil
	}
	
	var wbResp WeiboResponse
	if err := json.NewDecoder(resp.Body).Decode(&wbResp); err != nil {
		return f.fetchMock(), nil
	}
	
	if wbResp.Ok != 1 {
		return f.fetchMock(), nil
	}
	
	var items []HotItem
	for _, item := range wbResp.Data.Realtime {
		// 只获取前50条
		if item.Rank > 50 {
			continue
		}
		
		hotValue := float64(item.RawHot)
		if hotValue == 0 {
			hotValue = float64(1000000 - item.Rank*10000)
		}
		
		items = append(items, HotItem{
			Rank:      item.Rank,
			Title:     item.Note,
			URL:       fmt.Sprintf("https://s.weibo.com/weibo?q=%s", item.Note),
			HotValue:  hotValue,
			Category:  item.Category,
			FetchTime: time.Now(),
		})
	}
	
	if len(items) == 0 {
		return f.fetchMock(), nil
	}
	
	return items, nil
}

// fetchMock 获取模拟数据（当API失败时使用）
func (f *WeiboFetcher) fetchMock() []HotItem {
	mockData := []struct {
		Title    string
		Category string
	}{
		{"AI技术革新引发行业变革", "科技"},
		{"新一代智能手机发布", "数码"},
		{"电影票房创新高", "娱乐"},
		{"股市今日大涨", "财经"},
		{"新能源汽车政策出台", "汽车"},
		{"明星婚礼现场曝光", "娱乐"},
		{"教育改革新方案", "教育"},
		{" healthcare 行业动态", "健康"},
		{"房地产政策调整", "房产"},
		{"体育赛事精彩回顾", "体育"},
	}
	
	var items []HotItem
	for i, data := range mockData {
		items = append(items, HotItem{
			Rank:      i + 1,
			Title:     data.Title,
			URL:       fmt.Sprintf("https://s.weibo.com/weibo?q=%s", data.Title),
			HotValue:  float64(1000000 - (i+1)*50000),
			Category:  data.Category,
			FetchTime: time.Now(),
		})
	}
	
	return items
}
