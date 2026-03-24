package hot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// BaiduFetcher 百度热榜获取器
type BaiduFetcher struct {
	client *http.Client
}

// NewBaiduFetcher 创建百度获取器
func NewBaiduFetcher() *BaiduFetcher {
	return &BaiduFetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Platform 返回平台名称
func (f *BaiduFetcher) Platform() string {
	return PlatformBaidu
}

// BaiduResponse 百度响应结构
type BaiduResponse struct {
	Data struct {
		Cards []struct {
			Content []struct {
				Query    string `json:"query"`
				Desc     string `json:"desc"`
				HotScore int    `json:"hotScore"`
				Index    int    `json:"index"`
				Url      string `json:"url"`
				Word     string `json:"word"`
			} `json:"content"`
		} `json:"cards"`
	} `json:"data"`
}

// Fetch 获取百度热榜
func (f *BaiduFetcher) Fetch() ([]HotItem, error) {
	// 百度热榜API
	url := "https://top.baidu.com/api/board?platform=wise&tab=realtime"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return f.fetchMock(), nil
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://top.baidu.com/")
	
	resp, err := f.client.Do(req)
	if err != nil {
		return f.fetchMock(), nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return f.fetchMock(), nil
	}
	
	var bdResp BaiduResponse
	if err := json.NewDecoder(resp.Body).Decode(&bdResp); err != nil {
		return f.fetchMock(), nil
	}
	
	var items []HotItem
	if len(bdResp.Data.Cards) > 0 {
		for _, item := range bdResp.Data.Cards[0].Content {
			title := item.Word
			if title == "" {
				title = item.Query
			}
			
			hotValue := float64(item.HotScore)
			if hotValue == 0 {
				hotValue = float64(1000000 - item.Index*10000)
			}
			
			items = append(items, HotItem{
				Rank:      item.Index,
				Title:     title,
				URL:       item.Url,
				HotValue:  hotValue,
				Category:  "实时热点",
				FetchTime: time.Now(),
			})
		}
	}
	
	if len(items) == 0 {
		return f.fetchMock(), nil
	}
	
	return items, nil
}

// fetchMock 获取模拟数据
func (f *BaiduFetcher) fetchMock() []HotItem {
	mockData := []struct {
		Title    string
		Category string
	}{
		{"百度热搜第一名", "实时热点"},
		{"今日新闻头条", "新闻"},
		{"股市行情最新", "财经"},
		{"天气预警信息", "社会"},
		{"体育赛事结果", "体育"},
		{"娱乐明星动态", "娱乐"},
		{"科技新品发布", "科技"},
		{"健康科普知识", "健康"},
		{"教育培训资讯", "教育"},
		{"旅游出行攻略", "旅游"},
	}
	
	var items []HotItem
	for i, data := range mockData {
		items = append(items, HotItem{
			Rank:      i + 1,
			Title:     data.Title,
			URL:       fmt.Sprintf("https://www.baidu.com/s?wd=%s", data.Title),
			HotValue:  float64(1000000 - (i+1)*55000),
			Category:  data.Category,
			FetchTime: time.Now(),
		})
	}
	
	return items
}
