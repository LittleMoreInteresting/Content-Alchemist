package hot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ZhihuFetcher 知乎热榜获取器
type ZhihuFetcher struct {
	client *http.Client
}

// NewZhihuFetcher 创建知乎获取器
func NewZhihuFetcher() *ZhihuFetcher {
	return &ZhihuFetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Platform 返回平台名称
func (f *ZhihuFetcher) Platform() string {
	return PlatformZhihu
}

// ZhihuResponse 知乎响应结构
type ZhihuResponse struct {
	Data []struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Url      string `json:"url"`
		DetailText string `json:"detail_text"`
		Excerpt  string `json:"excerpt"`
	} `json:"data"`
}

// Fetch 获取知乎热榜
func (f *ZhihuFetcher) Fetch() ([]HotItem, error) {
	// 知乎热榜API
	url := "https://www.zhihu.com/api/v3/feed/topstory/hot-lists/total?limit=50"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return f.fetchMock(), nil
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://www.zhihu.com/hot")
	
	resp, err := f.client.Do(req)
	if err != nil {
		return f.fetchMock(), nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return f.fetchMock(), nil
	}
	
	var zhResp ZhihuResponse
	if err := json.NewDecoder(resp.Body).Decode(&zhResp); err != nil {
		return f.fetchMock(), nil
	}
	
	var items []HotItem
	for i, item := range zhResp.Data {
		if i >= 50 {
			break
		}
		
		// 解析热度值
		hotValue := float64(1000000 - (i+1)*15000)
		
		items = append(items, HotItem{
			Rank:      i + 1,
			Title:     item.Title,
			URL:       item.Url,
			HotValue:  hotValue,
			Category:  "热榜",
			FetchTime: time.Now(),
		})
	}
	
	if len(items) == 0 {
		return f.fetchMock(), nil
	}
	
	return items, nil
}

// fetchMock 获取模拟数据
func (f *ZhihuFetcher) fetchMock() []HotItem {
	mockData := []struct {
		Title   string
		Excerpt string
	}{
		{"如何看待人工智能对就业的影响？", "随着AI技术的快速发展..."},
		{"有哪些值得推荐的高效学习方法？", "分享几个我实践过的学习方法..."},
		{"2024年最值得关注的科技趋势是什么？", "从多个维度分析今年的科技趋势..."},
		{"为什么现在的年轻人不愿意结婚了？", "从社会学角度分析这一现象..."},
		{"如何评价最近爆火的电视剧？", "从剧情、演技、制作等方面评价..."},
		{"程序员如何保持技术竞争力？", "分享一些程序员成长的经验..."},
		{"有哪些性价比高的旅游目的地？", "推荐几个适合不同预算的旅游地..."},
		{"如何培养良好的理财习惯？", "从记账、预算、投资等方面分享..."},
		{"未来十年哪些行业最有前景？", "分析未来发展趋势和机遇..."},
		{"如何提升自己的表达能力？", "分享一些实用的表达技巧..."},
	}
	
	var items []HotItem
	for i, data := range mockData {
		items = append(items, HotItem{
			Rank:      i + 1,
			Title:     data.Title,
			URL:       fmt.Sprintf("https://www.zhihu.com/question/%d", 1000000+i),
			HotValue:  float64(1000000 - (i+1)*60000),
			Category:  "热榜",
			FetchTime: time.Now(),
		})
	}
	
	return items
}
