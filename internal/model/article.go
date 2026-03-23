package model

import "time"

// Article 文章模型
type Article struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	Outline   []OutlineNode `json:"outline"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Status    string        `json:"status"` // draft, reviewing, ready, published, archived

	// 新增字段 - 工作流支持
	SourceType    string     `json:"sourceType"`    // manual, workflow
	WorkflowRunID string     `json:"workflowRunId"` // 来源工作流
	TopicID       string     `json:"topicId"`       // 关联选题

	// 质量评估
	QualityScore float64 `json:"qualityScore"`
	ReadTime     int     `json:"readTime"`  // 阅读时长(分钟)
	WordCount    int     `json:"wordCount"`

	// 发布信息
	PublishTaskID string     `json:"publishTaskId"`
	PublishedAt   *time.Time `json:"publishedAt"`
}

// OutlineNode 大纲节点
type OutlineNode struct {
	ID       string `json:"id"`
	Level    int    `json:"level"`    // 1,2,3
	Title    string `json:"title"`
	Content  string `json:"content"`  // 该节点对应正文
	ParentID string `json:"parentId"`
	Status   string `json:"status"`   // empty, draft, done
	WordCount int   `json:"wordCount"`
	TargetWords int `json:"targetWords"`
}

// Config 应用配置
type Config struct {
	ID          string   `json:"id"`
	APIBaseURL  string   `json:"apiBaseUrl"`
	APIKey      string   `json:"apiKey"`
	Model       string   `json:"model"`
	Temperature float64  `json:"temperature"`
	StyleTags   []string `json:"styleTags"`
	Audience    string   `json:"audience"`
	Persona     string   `json:"persona"` // 我/我们/小编/笔者
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Material 素材
type Material struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"` // snippet, data, quote, history
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Tags       []string  `json:"tags"`
	Source     string    `json:"source"`
	CreatedAt  time.Time `json:"createdAt"`
	UsageCount int       `json:"usageCount"`
}

// ArticleStatus 文章状态常量
const (
	ArticleStatusDraft     = "draft"
	ArticleStatusReviewing = "reviewing"
	ArticleStatusReady     = "ready"
	ArticleStatusPublished = "published"
	ArticleStatusArchived  = "archived"
)

// ArticleSource 文章来源类型
const (
	ArticleSourceManual   = "manual"
	ArticleSourceWorkflow = "workflow"
)

// Version 文章版本
type Version struct {
	ID        string    `json:"id"`
	ArticleID string    `json:"articleId"`
	Content   string    `json:"content"`
	Snapshot  string    `json:"snapshot"` // 前100字摘要
	CreatedAt time.Time `json:"createdAt"`
	Type      string    `json:"type"` // auto, manual
}
