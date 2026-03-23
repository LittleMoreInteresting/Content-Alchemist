package model

import (
	"time"
)

// ==================== 发布相关模型 ====================

// PublishAccount 发布账号
type PublishAccount struct {
	ID          string    `json:"id"`
	Platform    string    `json:"platform"`    // wechat, zhihu, juejin等
	Name        string    `json:"name"`        // 账号名称
	AccountID   string    `json:"accountId"`   // 平台账号ID
	AccountType string    `json:"accountType"` // 账号类型：个人/企业

	// 微信特有配置
	AppID       string    `json:"appId"`
	AppSecret   string    `json:"appSecret"`   // 加密存储
	AccessToken string    `json:"accessToken"` // 加密存储
	TokenExpiry time.Time `json:"tokenExpiry"`
	RefreshToken string   `json:"refreshToken"` // 刷新令牌

	// 其他平台配置
	Config map[string]string `json:"config"` // 其他平台配置

	// 状态
	Status     string    `json:"status"`     // active, expired, disabled
	LastUsedAt time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time `json:"createdAt"`
}

// PublishTask 发布任务
type PublishTask struct {
	ID        string `json:"id"`
	ArticleID string `json:"articleId"`
	AccountID string `json:"accountId"`

	// 发布设置
	Title      string `json:"title"`
	Content    string `json:"content"`    // 发布内容(HTML)
	CoverImage string `json:"coverImage"` // 封面图URL或MediaID
	Summary    string `json:"summary"`    // 摘要
	Author     string `json:"author"`     // 作者

	// 其他发布选项
	Tags        []string `json:"tags"`        // 标签
	Category    string   `json:"category"`    // 分类
	OriginalURL string   `json:"originalUrl"` // 原文链接

	// 定时发布
	ScheduleType string     `json:"scheduleType"` // immediate, scheduled
	ScheduleAt   *time.Time `json:"scheduleAt"`

	// 发布状态
	Status      string `json:"status"`      // pending, scheduled, publishing, published, failed, cancelled
	PlatformID  string `json:"platformId"`  // 平台返回的ID
	PlatformURL string `json:"platformUrl"` // 发布后链接

	// 结果
	Error       string     `json:"error"`
	PublishedAt *time.Time `json:"publishedAt"`
	RetryCount  int        `json:"retryCount"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// PublishPlatform 发布平台常量
const (
	PlatformWechat = "wechat"
	PlatformZhihu  = "zhihu"
	PlatformJuejin = "juejin"
	PlatformCSDN   = "csdn"
	PlatformJianshu = "jianshu"
)

// PublishStatus 发布状态常量
const (
	PublishStatusPending    = "pending"
	PublishStatusScheduled  = "scheduled"
	PublishStatusPublishing = "publishing"
	PublishStatusPublished  = "published"
	PublishStatusFailed     = "failed"
	PublishStatusCancelled  = "cancelled"
)

// ScheduleType 调度类型
const (
	ScheduleTypeImmediate  = "immediate"
	ScheduleTypeScheduled  = "scheduled"
)

// AccountStatus 账号状态
const (
	AccountStatusActive   = "active"
	AccountStatusExpired  = "expired"
	AccountStatusDisabled = "disabled"
)

// WechatDraftRequest 微信草稿请求
type WechatDraftRequest struct {
	Articles []WechatArticle `json:"articles"`
}

// WechatArticle 微信文章
type WechatArticle struct {
	Title              string `json:"title"`
	Author             string `json:"author"`
	Digest             string `json:"digest"`
	Content            string `json:"content"`
	ContentSourceURL   string `json:"content_source_url"`
	ThumbMediaID       string `json:"thumb_media_id"`
	ShowCoverPic       int    `json:"show_cover_pic"`
	NeedOpenComment    int    `json:"need_open_comment"`
	OnlyFansCanComment int    `json:"only_fans_can_comment"`
}

// WechatPublishResponse 微信发布响应
type WechatPublishResponse struct {
	PublishID int64  `json:"publish_id"`
	MsgID     int64  `json:"msg_id"`
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

// PublishStats 发布统计
type PublishStats struct {
	TotalCount    int64 `json:"totalCount"`
	SuccessCount  int64 `json:"successCount"`
	FailedCount   int64 `json:"failedCount"`
	PendingCount  int64 `json:"pendingCount"`
	ScheduleCount int64 `json:"scheduleCount"`
}

// PublishQueueItem 发布队列项
type PublishQueueItem struct {
	TaskID       string    `json:"taskId"`
	Priority     int       `json:"priority"`     // 优先级
	ScheduledAt  time.Time `json:"scheduledAt"`
	RetryCount   int       `json:"retryCount"`
	MaxRetries   int       `json:"maxRetries"`
}
