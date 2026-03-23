package model

import (
	"time"
)

// ==================== 选题相关模型 ====================

// Topic 选题
type Topic struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`    // 选题标题
	Category string   `json:"category"` // 分类
	Source   string   `json:"source"`   // 来源: ai, rss, manual, hot
	SourceURL string  `json:"sourceUrl"` // 来源链接

	// AI评估
	Score     float64 `json:"score"`     // 综合评分 0-100
	HotScore  float64 `json:"hotScore"`  // 热度评分
	CompScore float64 `json:"compScore"` // 竞争度评分
	FitScore  float64 `json:"fitScore"`  // 匹配度评分

	// 分析数据
	Keywords   []string `json:"keywords"`   // 关键词
	Summary    string   `json:"summary"`    // 内容摘要
	References []string `json:"references"` // 参考文章
	Angles     []string `json:"angles"`     // 切入角度

	// 状态
	Status        string `json:"status"`        // pending, approved, rejected, processing, published, archived
	WorkflowRunID string `json:"workflowRunId"` // 关联的工作流运行ID

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TopicSource 选题来源类型
const (
	TopicSourceAI     = "ai"
	TopicSourceRSS    = "rss"
	TopicSourceManual = "manual"
	TopicSourceHot    = "hot"
)

// TopicStatus 选题状态
const (
	TopicStatusPending   = "pending"
	TopicStatusApproved  = "approved"
	TopicStatusRejected  = "rejected"
	TopicStatusProcessing = "processing"
	TopicStatusPublished = "published"
	TopicStatusArchived  = "archived"
)

// HotTrend 热点趋势
type HotTrend struct {
	ID       string    `json:"id"`
	Platform string    `json:"platform"` // 平台: weibo, zhihu, baidu, toutiao, douyin
	Title    string    `json:"title"`
	URL      string    `json:"url"`
	HotRank  int       `json:"hotRank"`  // 热度排名
	HotValue float64   `json:"hotValue"` // 热度值
	Category string    `json:"category"`
	CreatedAt time.Time `json:"createdAt"`
}

// HotPlatform 热点平台常量
const (
	HotPlatformWeibo  = "weibo"
	HotPlatformZhihu  = "zhihu"
	HotPlatformBaidu  = "baidu"
	HotPlatformToutiao = "toutiao"
	HotPlatformDouyin = "douyin"
)

// TopicRequest AI生成选题请求
type TopicRequest struct {
	Keywords      []string `json:"keywords"`      // 关注领域
	ExcludeTopics []string `json:"excludeTopics"` // 排除的选题
	Limit         int      `json:"limit"`         // 生成数量
}

// TopicResult AI生成的选题结果
type TopicResult struct {
	Title    string   `json:"title"`
	Score    float64  `json:"score"`
	Reason   string   `json:"reason"`   // 推荐理由
	Angles   []string `json:"angles"`   // 切入角度建议
	Keywords []string `json:"keywords"`
}

// TopicAnalysis 选题深度分析
type TopicAnalysis struct {
	TopicID    string            `json:"topicId"`
	Title      string            `json:"title"`
	HotScore   float64           `json:"hotScore"`   // 热度评分
	CompScore  float64           `json:"compScore"`  // 竞争度评分
	FitScore   float64           `json:"fitScore"`   // 匹配度评分
	Potential  float64           `json:"potential"`  // 传播潜力
	Timeline   string            `json:"timeline"`   // 最佳发布时间

	// 竞品分析
	Competitors []CompetitorInfo `json:"competitors"`

	// 关键词分析
	Keywords []KeywordAnalysis `json:"keywords"`
}

// CompetitorInfo 竞品信息
type CompetitorInfo struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Platform    string `json:"platform"`
	ReadCount   int    `json:"readCount"`
	LikeCount   int    `json:"likeCount"`
	PublishTime string `json:"publishTime"`
}

// KeywordAnalysis 关键词分析
type KeywordAnalysis struct {
	Word       string  `json:"word"`
	SearchVol  int     `json:"searchVol"`  // 搜索量
	Competition float64 `json:"competition"` // 竞争度
}
