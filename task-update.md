# 🚀 Content-Alchemist Pro 改版计划
## AI 运营工作流引擎重构方案

---

## 📋 改版概述

### 核心目标
将现有「单篇创作工具」升级为「AI 运营工作流引擎」，实现：
- **自动选题** → **自动创作** → **自动排版** → **发布草稿** 的全链路自动化
- 保留并强化「预览」和「公众号自动排版」核心能力
- 支持批量运营、定时发布、多账号管理

### 新旧架构对比

```
【旧架构】单篇创作模式
┌─────────────────────────────────────────────────────┐
│  用户输入标题 → AI生成大纲 → 人工编辑 → 复制发布    │
│                  ↓                                   │
│            单次交互，人工驱动                        │
└─────────────────────────────────────────────────────┘

【新架构】工作流引擎模式
┌──────────────────────────────────────────────────────────────────────┐
│  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐          │
│  │ 选题引擎 │ → │ 创作引擎 │ → │ 排版引擎 │ → │ 发布引擎 │          │
│  │          │   │          │   │          │   │          │          │
│  │•热点追踪 │   │•大纲生成 │   │•主题渲染 │   │•微信API  │          │
│  │•选题评分 │   │•分段创作 │   │•自动配图 │   │•定时发布 │          │
│  │•竞品分析 │   │•质量评估 │   │•预览生成 │   │•数据回传 │          │
│  └──────────┘   └──────────┘   └──────────┘   └──────────┘          │
│       ↑                                              ↓               │
│       └────────────── 人工审核/干预点 ────────────────┘               │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 🏗️ 一、架构重构设计

### 1.1 新架构分层

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Presentation Layer (UI)                            │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐       │
│  │ 工作流控制台  │ │ 选题中心     │ │ 创作编辑器   │ │ 发布管理     │       │
│  │              │ │              │ │ (保留现有)   │ │              │       │
│  │•工作流列表   │ │•热点看板     │ │              │ │•发布队列     │       │
│  │•运行监控     │ │•选题推荐     │ │              │ │•发布历史     │       │
│  │•调度配置     │ │•选题库       │ │              │ │•账号管理     │       │
│  └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘       │
├─────────────────────────────────────────────────────────────────────────────┤
│                           Workflow Engine Layer                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      Workflow Orchestrator                           │    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │    │
│  │  │ 状态机   │ │ 调度器   │ │ 重试器   │ │ 并行器   │ │ 钩子系统 │   │    │
│  │  │          │ │          │ │          │ │          │ │          │   │    │
│  │  │•草稿    │ │•定时    │ │•指数退避│ │•MapReduce│ │•前置钩子│   │    │
│  │  │•审核中  │ │•优先级  │ ││•死信队列│ │•Pipeline │ │•后置钩子│   │    │
│  │  │•发布中  │ │•并发控  │ ││•告警    │ │•Fork/Join│ │•人工审核│   │    │
│  │  │•已发布  │ │         │ │          │ │          │ │          │   │    │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────────────────────┤
│                            Agent Layer (Eino)                                │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐       │
│  │ TopicAgent   │ │ OutlineAgent │ │ WritingAgent │ │ ReviewAgent  │       │
│  │              │ │              │ │              │ │              │       │
│  │•热点分析     │ │•大纲生成     │ │•分段创作     │ │•质量评估     │       │
│  │•选题评分     │ │•结构调整     │ │•润色优化     │ │•合规检查     │       │
│  │•竞品分析     │ │•标题生成     │ │•扩写/精简   │ │•A/B测试标题 │       │
│  └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘       │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐                        │
│  │LayoutAgent   │ │ ImageAgent   │ │ PublishAgent │                        │
│  │              │ │              │ │              │                        │
│  │•公众号排版   │ │•配图生成     │ │•微信API     │                        │
│  │•主题适配     │ │•封面图       │ │•定时发布     │                        │
│  │•代码高亮     │ │•插图搜索     │ │•草稿管理     │                        │
│  └──────────────┘ └──────────────┘ └──────────────┘                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                           Infrastructure Layer                               │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐          │
│  │ SQLite   │ │ LocalFS  │ │ WeChatAPI│ │ HotAPI   │ │ LLM API  │          │
│  │          │ │          │ │          │ │          │ │          │          │
│  │•工作流状态│ │•素材存储 │ │•公众号  │ │•热榜接口 │ │•DeepSeek│          │
│  │•文章数据 │ │•图片缓存 │ │•草稿箱  │ │•舆情监控 │ │•OpenAI  │          │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘          │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 核心模块职责

| 模块 | 职责 | 关键技术 |
|------|------|----------|
| Workflow Engine | 工作流编排、状态管理、调度执行 | 状态机、Cron、事件驱动 |
| TopicAgent | 热点追踪、选题生成、竞品分析 | RSS/API采集、NLP分析 |
| ReviewAgent | 内容审核、质量评分、合规检查 | 规则引擎、AI评估 |
| LayoutAgent | 公众号排版、主题渲染、样式优化 | Markdown→HTML、CSS变量 |
| ImageAgent | 封面图生成、插图搜索、图片优化 | AI绘图、图库API |
| PublishAgent | 微信API对接、定时发布、发布队列 | 微信MP API、任务调度 |

---

## 📊 二、数据模型扩展

### 2.1 新增核心模型

```go
// ==================== 工作流相关 ====================

// Workflow 工作流定义
type Workflow struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    // 工作流配置
    Trigger     WorkflowTrigger   `json:"trigger"`     // 触发器配置
    Steps       []WorkflowStep    `json:"steps"`       // 执行步骤
    // 自动化设置
    AutoPublish bool              `json:"autoPublish"` // 是否自动发布
    NeedReview  bool              `json:"needReview"`  // 是否需要人工审核
    // 状态
    Status      string            `json:"status"`      // active, paused, archived
    CreatedAt   time.Time         `json:"createdAt"`
    UpdatedAt   time.Time         `json:"updatedAt"`
}

// WorkflowTrigger 工作流触发器
type WorkflowTrigger struct {
    Type      string            `json:"type"`      // manual, schedule, webhook, rss
    CronExpr  string            `json:"cronExpr"`  // 定时表达式
    RSSURL    string            `json:"rssUrl"`    // RSS源地址
    WebhookURL string           `json:"webhookUrl"` // Webhook地址
    Config    map[string]any    `json:"config"`    // 额外配置
}

// WorkflowStep 工作流步骤
type WorkflowStep struct {
    ID       string   `json:"id"`
    Name     string   `json:"name"`
    Type     string   `json:"type"`     // topic, outline, write, review, layout, image, publish
    Config   map[string]any `json:"config"` // 步骤配置
    NextStep string   `json:"nextStep"` // 下一步ID
    OnError  string   `json:"onError"`  // 错误处理: retry, skip, abort, manual
}

// WorkflowRun 工作流运行实例
type WorkflowRun struct {
    ID         string            `json:"id"`
    WorkflowID string            `json:"workflowId"`
    Status     string            `json:"status"`     // pending, running, paused, completed, failed
    CurrentStep string           `json:"currentStep"`
    Input      map[string]any    `json:"input"`      // 输入参数
    Output     map[string]any    `json:"output"`     // 输出结果
    Steps      []WorkflowRunStep `json:"steps"`      // 各步骤执行状态
    StartedAt  time.Time         `json:"startedAt"`
    CompletedAt *time.Time       `json:"completedAt"`
    Error      string            `json:"error"`
}

// WorkflowRunStep 工作流步骤执行状态
type WorkflowRunStep struct {
    StepID    string     `json:"stepId"`
    Status    string     `json:"status"`    // pending, running, completed, failed, skipped
    Input     any        `json:"input"`
    Output    any        `json:"output"`
    StartedAt time.Time  `json:"startedAt"`
    CompletedAt *time.Time `json:"completedAt"`
    Error     string     `json:"error"`
    RetryCount int       `json:"retryCount"`
}

// ==================== 选题相关 ====================

// Topic 选题
type Topic struct {
    ID          string   `json:"id"`
    Title       string   `json:"title"`       // 选题标题
    Category    string   `json:"category"`    // 分类
    Source      string   `json:"source"`      // 来源: ai, rss, manual, hot
    SourceURL   string   `json:"sourceUrl"`   // 来源链接
    
    // AI评估
    Score       float64  `json:"score"`       // 综合评分 0-100
    HotScore    float64  `json:"hotScore"`    // 热度评分
    CompScore   float64  `json:"compScore"`   // 竞争度评分
    FitScore    float64  `json:"fitScore"`    // 匹配度评分
    
    // 分析数据
    Keywords    []string `json:"keywords"`    // 关键词
    Summary     string   `json:"summary"`     // 内容摘要
    References  []string `json:"references"`  // 参考文章
    
    // 状态
    Status      string   `json:"status"`      // pending, approved, rejected, processing, published
    WorkflowRunID string `json:"workflowRunId"` // 关联的工作流运行ID
    
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

// HotTrend 热点趋势
type HotTrend struct {
    ID          string    `json:"id"`
    Platform    string    `json:"platform"`    // 平台: weibo, zhihu, baidu, toutiao
    Title       string    `json:"title"`
    URL         string    `json:"url"`
    HotRank     int       `json:"hotRank"`     // 热度排名
    HotValue    float64   `json:"hotValue"`    // 热度值
    Category    string    `json:"category"`
    CreatedAt   time.Time `json:"createdAt"`
}

// ==================== 发布相关 ====================

// PublishAccount 发布账号
type PublishAccount struct {
    ID           string    `json:"id"`
    Platform     string    `json:"platform"`     // wechat, zhihu, juejin等
    Name         string    `json:"name"`         // 账号名称
    AccountID    string    `json:"accountId"`    // 平台账号ID
    
    // 微信特有
    AppID        string    `json:"appId"`
    AppSecret    string    `json:"appSecret"`
    AccessToken  string    `json:"accessToken"`
    TokenExpiry  time.Time `json:"tokenExpiry"`
    
    // 状态
    Status       string    `json:"status"`       // active, expired, disabled
    LastUsedAt   time.Time `json:"lastUsedAt"`
    CreatedAt    time.Time `json:"createdAt"`
}

// PublishTask 发布任务
type PublishTask struct {
    ID           string    `json:"id"`
    ArticleID    string    `json:"articleId"`
    AccountID    string    `json:"accountId"`
    
    // 发布设置
    Title        string    `json:"title"`
    Content      string    `json:"content"`      // 发布内容(HTML)
    CoverImage   string    `json:"coverImage"`   // 封面图
    Summary      string    `json:"summary"`      // 摘要
    
    // 定时发布
    ScheduleType string    `json:"scheduleType"` // immediate, scheduled
    ScheduleAt   *time.Time `json:"scheduleAt"`
    
    // 发布状态
    Status       string    `json:"status"`       // pending, scheduled, publishing, published, failed
    PlatformID   string    `json:"platformId"`   // 平台返回的ID
    PlatformURL  string    `json:"platformUrl"`  // 发布后链接
    
    // 结果
    Error        string    `json:"error"`
    PublishedAt  *time.Time `json:"publishedAt"`
    CreatedAt    time.Time `json:"createdAt"`
}

// ==================== 文章扩展 ====================

// Article 扩展现有模型
type Article struct {
    ID        string        `json:"id"`
    Title     string        `json:"title"`
    Content   string        `json:"content"`
    Outline   []OutlineNode `json:"outline"`
    CreatedAt time.Time     `json:"createdAt"`
    UpdatedAt time.Time     `json:"updatedAt"`
    Status    string        `json:"status"` // draft, reviewing, ready, published, archived
    
    // 新增字段
    SourceType    string    `json:"sourceType"`    // manual, workflow
    WorkflowRunID string    `json:"workflowRunId"` // 来源工作流
    TopicID       string    `json:"topicId"`       // 关联选题
    
    // 质量评估
    QualityScore  float64   `json:"qualityScore"`
    ReadTime      int       `json:"readTime"`      // 阅读时长(分钟)
    WordCount     int       `json:"wordCount"`
    
    // 发布信息
    PublishTaskID string    `json:"publishTaskId"`
    PublishedAt   *time.Time `json:"publishedAt"`
}
```

### 2.2 数据库存储设计

```sql
-- 工作流表
CREATE TABLE workflows (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    trigger_config TEXT, -- JSON
    steps TEXT, -- JSON
    auto_publish BOOLEAN DEFAULT FALSE,
    need_review BOOLEAN DEFAULT TRUE,
    status TEXT DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 工作流运行表
CREATE TABLE workflow_runs (
    id TEXT PRIMARY KEY,
    workflow_id TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    current_step TEXT,
    input TEXT, -- JSON
    output TEXT, -- JSON
    steps TEXT, -- JSON
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error TEXT,
    FOREIGN KEY (workflow_id) REFERENCES workflows(id)
);

-- 选题表
CREATE TABLE topics (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    category TEXT,
    source TEXT,
    source_url TEXT,
    score REAL,
    hot_score REAL,
    comp_score REAL,
    fit_score REAL,
    keywords TEXT, -- JSON array
    summary TEXT,
    references TEXT, -- JSON array
    status TEXT DEFAULT 'pending',
    workflow_run_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 热点趋势表
CREATE TABLE hot_trends (
    id TEXT PRIMARY KEY,
    platform TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT,
    hot_rank INTEGER,
    hot_value REAL,
    category TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 发布账号表
CREATE TABLE publish_accounts (
    id TEXT PRIMARY KEY,
    platform TEXT NOT NULL,
    name TEXT NOT NULL,
    account_id TEXT,
    app_id TEXT,
    app_secret TEXT, -- 加密存储
    access_token TEXT, -- 加密存储
    token_expiry TIMESTAMP,
    status TEXT DEFAULT 'active',
    last_used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 发布任务表
CREATE TABLE publish_tasks (
    id TEXT PRIMARY KEY,
    article_id TEXT NOT NULL,
    account_id TEXT NOT NULL,
    title TEXT,
    content TEXT,
    cover_image TEXT,
    summary TEXT,
    schedule_type TEXT DEFAULT 'immediate',
    schedule_at TIMESTAMP,
    status TEXT DEFAULT 'pending',
    platform_id TEXT,
    platform_url TEXT,
    error TEXT,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES publish_accounts(id)
);
```

---

## 🤖 三、AI Agent 扩展

### 3.1 TopicAgent - 选题引擎

```go
// internal/agent/topic_agent.go
package agent

// TopicAgent 选题Agent
type TopicAgent struct {
    llm einoModel.ToolCallingChatModel
    hotAPI HotTrendAPI // 热榜API客户端
}

// TopicRequest 选题请求
type TopicRequest struct {
    Keywords    []string `json:"keywords"`    // 关注领域
    ExcludeTopics []string `json:"excludeTopics"` // 排除的选题
    Limit       int      `json:"limit"`       // 生成数量
}

// TopicResult 选题结果
type TopicResult struct {
    Title    string   `json:"title"`
    Score    float64  `json:"score"`
    Reason   string   `json:"reason"`    // 推荐理由
    Angles   []string `json:"angles"`    // 切入角度建议
    Keywords []string `json:"keywords"`
}

// GenerateTopics 基于热点生成选题
func (a *TopicAgent) GenerateTopics(ctx context.Context, req TopicRequest) ([]TopicResult, error) {
    // 1. 获取热点数据
    hotTrends := a.hotAPI.FetchTrends(ctx, req.Keywords)
    
    // 2. 构建Prompt进行AI分析
    prompt := fmt.Sprintf(`基于以下热点数据，生成%d个适合公众号的选题：

关注领域：%v
热点数据：%v

要求：
1. 每个选题需包含：标题、评分(0-100)、推荐理由
2. 提供3个不同的切入角度
3. 避免与已发布内容重复
4. 考虑时效性和传播性

输出JSON格式：
{
  "topics": [
    {
      "title": "选题标题",
      "score": 85,
      "reason": "推荐理由",
      "angles": ["角度1", "角度2", "角度3"],
      "keywords": ["关键词1", "关键词2"]
    }
  ]
}`, req.Limit, req.Keywords, hotTrends)
    
    // 3. 调用LLM生成
    messages := []*schema.Message{{Role: schema.User, Content: prompt}}
    resp, err := a.llm.Generate(ctx, messages, einoModel.WithTemperature(0.8))
    if err != nil {
        return nil, err
    }
    
    // 4. 解析结果
    return a.parseTopicResults(resp.Content)
}

// AnalyzeTopic 深度分析选题
func (a *TopicAgent) AnalyzeTopic(ctx context.Context, title string) (*TopicAnalysis, error) {
    // 分析选题的竞争度、热度、匹配度
    // 返回详细的分析报告
}
```

### 3.2 ReviewAgent - 质量审核

```go
// internal/agent/review_agent.go
package agent

// ReviewAgent 审核Agent
type ReviewAgent struct {
    llm einoModel.ToolCallingChatModel
}

// ReviewRequest 审核请求
type ReviewRequest struct {
    Title   string `json:"title"`
    Content string `json:"content"`
    Outline []model.OutlineNode `json:"outline"`
}

// ReviewResult 审核结果
type ReviewResult struct {
    Passed      bool             `json:"passed"`
    Score       float64          `json:"score"`
    Issues      []ReviewIssue    `json:"issues"`
    Suggestions []string         `json:"suggestions"`
}

// ReviewIssue 审核问题
type ReviewIssue struct {
    Type     string `json:"type"`     // error, warning, info
    Category string `json:"category"` // grammar, logic, compliance, style
    Message  string `json:"message"`
    Position string `json:"position"` // 问题位置
}

// Review 审核文章
func (a *ReviewAgent) Review(ctx context.Context, req ReviewRequest) (*ReviewResult, error) {
    prompt := fmt.Sprintf(`请对以下公众号文章进行全面审核：

标题：%s
大纲：%v
内容：%s

审核维度：
1. 内容质量：逻辑性、深度、原创性
2. 语言表达：流畅度、专业度、可读性
3. 合规检查：敏感词、违规内容
4. 排版检查：格式规范、图片说明
5. 传播潜力：标题吸引力、互动引导

输出JSON格式：
{
  "passed": true/false,
  "score": 85,
  "issues": [
    {"type": "warning", "category": "style", "message": "第三段过长", "position": "第3段"}
  ],
  "suggestions": ["建议1", "建议2"]
}`, req.Title, req.Outline, req.Content)
    
    messages := []*schema.Message{{Role: schema.User, Content: prompt}}
    resp, err := a.llm.Generate(ctx, messages)
    if err != nil {
        return nil, err
    }
    
    return a.parseReviewResult(resp.Content)
}
```

### 3.3 ImageAgent - 配图生成

```go
// internal/agent/image_agent.go
package agent

// ImageAgent 图片Agent
type ImageAgent struct {
    aiDrawAPI   AIDrawAPI   // AI绘图API
    imageSearch ImageSearchAPI // 图库搜索API
}

// GenerateCoverRequest 生成封面请求
type GenerateCoverRequest struct {
    Title    string   `json:"title"`
    Keywords []string `json:"keywords"`
    Style    string   `json:"style"`    // 风格：简约、商务、科技等
}

// GenerateCover 生成封面图
func (a *ImageAgent) GenerateCover(ctx context.Context, req GenerateCoverRequest) (string, error) {
    // 1. 生成提示词
    prompt := a.buildCoverPrompt(req)
    
    // 2. 调用AI绘图
    imageURL, err := a.aiDrawAPI.Generate(ctx, prompt)
    if err != nil {
        // 降级：使用图库搜索
        return a.imageSearch.Search(ctx, req.Keywords)
    }
    
    return imageURL, nil
}

// SearchIllustration 搜索插图
func (a *ImageAgent) SearchIllustration(ctx context.Context, keyword string) ([]string, error) {
    // 搜索文章相关插图
}
```

### 3.4 PublishAgent - 发布代理

```go
// internal/agent/publish_agent.go
package agent

// PublishAgent 发布Agent
type PublishAgent struct {
    wechatAPI *WechatAPI
}

// WechatAPI 微信公众号API封装
type WechatAPI struct {
    appID       string
    appSecret   string
    accessToken string
    tokenExpiry time.Time
}

// PublishDraftRequest 发布草稿请求
type PublishDraftRequest struct {
    Title      string   `json:"title"`
    Content    string   `json:"content"`    // HTML内容
    CoverImage string   `json:"coverImage"`
    Author     string   `json:"author"`
    Digest     string   `json:"digest"`     // 摘要
}

// PublishDraft 发布到微信草稿箱
func (a *PublishAgent) PublishDraft(ctx context.Context, req PublishDraftRequest) (string, error) {
    // 1. 确保access_token有效
    if err := a.ensureToken(); err != nil {
        return "", err
    }
    
    // 2. 上传封面图片到微信
    mediaID, err := a.uploadImage(ctx, req.CoverImage)
    if err != nil {
        return "", fmt.Errorf("upload cover failed: %w", err)
    }
    
    // 3. 上传文章中的图片
    content, err := a.processContentImages(ctx, req.Content)
    if err != nil {
        return "", fmt.Errorf("process content images failed: %w", err)
    }
    
    // 4. 创建草稿
    draftReq := WechatDraftRequest{
        Articles: []WechatArticle{
            {
                Title:              req.Title,
                Content:            content,
                ThumbMediaID:       mediaID,
                Author:             req.Author,
                Digest:             req.Digest,
                ShowCoverPic:       1,
                ContentSourceURL:   "",
                NeedOpenComment:    1,
                OnlyFansCanComment: 0,
            },
        },
    }
    
    resp, err := a.createDraft(ctx, draftReq)
    if err != nil {
        return "", fmt.Errorf("create draft failed: %w", err)
    }
    
    return resp.MediaID, nil
}

// SchedulePublish 定时发布
func (a *PublishAgent) SchedulePublish(ctx context.Context, mediaID string, publishAt time.Time) error {
    // 调用微信发布接口
}
```

---

## 🔄 四、工作流引擎设计

### 4.1 工作流定义DSL

```yaml
# 工作流配置示例
# workflows/content_production.yaml

name: "全自动内容生产"
description: "从选题到发布的全自动化工作流"

trigger:
  type: schedule
  cron: "0 9 * * *"  # 每天上午9点执行

steps:
  - id: fetch_trends
    name: "获取热点"
    type: hot_fetch
    config:
      platforms: [weibo, zhihu, toutiao]
      limit: 20
    next: generate_topics
    
  - id: generate_topics
    name: "生成选题"
    type: topic_generate
    config:
      limit: 5
      min_score: 70
    next: select_topic
    
  - id: select_topic
    name: "选择最佳选题"
    type: topic_select
    config:
      strategy: highest_score
    next: generate_outline
    
  - id: generate_outline
    name: "生成大纲"
    type: outline_generate
    config:
      style: "干货专业"
      sections: 5
    next: write_content
    
  - id: write_content
    name: "创作正文"
    type: content_write
    config:
      parallel: true
      word_count: 2000
    next: review
    
  - id: review
    name: "质量审核"
    type: content_review
    config:
      min_score: 80
    next: check_review
    
  - id: check_review
    name: "检查审核结果"
    type: condition
    config:
      condition: "review.score >= 80"
      true_next: generate_cover
      false_next: rewrite
    
  - id: rewrite
    name: "重写优化"
    type: content_rewrite
    config:
      focus: review.issues
    next: review
    
  - id: generate_cover
    name: "生成封面"
    type: image_generate
    config:
      image_type: cover
    next: layout
    
  - id: layout
    name: "公众号排版"
    type: layout_apply
    config:
      theme: "科技蓝"
    next: publish
    
  - id: publish
    name: "发布草稿"
    type: publish_draft
    config:
      account: "default"
      schedule: "next_day_8am"

# 错误处理
error_handling:
  default: retry
  max_retries: 3
  retry_delay: 60  # 秒
```

### 4.2 工作流引擎核心代码

```go
// internal/workflow/engine.go
package workflow

// Engine 工作流引擎
type Engine struct {
    store      WorkflowStore
    scheduler  Scheduler
    dispatcher Dispatcher
    hooks      Hooks
    
    // Agent集合
    agents     map[string]Agent
    
    // 执行状态
    runs       map[string]*RunContext
    mu         sync.RWMutex
}

// RunContext 运行上下文
type RunContext struct {
    Run      *model.WorkflowRun
    Workflow *model.Workflow
    Cancel   context.CancelFunc
}

// Start 启动工作流
func (e *Engine) Start(ctx context.Context, workflowID string, input map[string]any) (*model.WorkflowRun, error) {
    // 1. 加载工作流定义
    workflow, err := e.store.GetWorkflow(workflowID)
    if err != nil {
        return nil, err
    }
    
    // 2. 创建运行实例
    run := &model.WorkflowRun{
        ID:         generateID(),
        WorkflowID: workflowID,
        Status:     "running",
        Input:      input,
        Steps:      make([]model.WorkflowRunStep, len(workflow.Steps)),
        StartedAt:  time.Now(),
    }
    
    // 3. 初始化步骤状态
    for i, step := range workflow.Steps {
        run.Steps[i] = model.WorkflowRunStep{
            StepID: step.ID,
            Status: "pending",
        }
    }
    
    // 4. 保存运行状态
    if err := e.store.SaveRun(run); err != nil {
        return nil, err
    }
    
    // 5. 启动执行协程
    runCtx, cancel := context.WithCancel(ctx)
    e.mu.Lock()
    e.runs[run.ID] = &RunContext{
        Run:      run,
        Workflow: workflow,
        Cancel:   cancel,
    }
    e.mu.Unlock()
    
    go e.execute(runCtx, run, workflow)
    
    return run, nil
}

// execute 执行工作流
func (e *Engine) execute(ctx context.Context, run *model.WorkflowRun, workflow *model.Workflow) {
    defer func() {
        e.mu.Lock()
        delete(e.runs, run.ID)
        e.mu.Unlock()
    }()
    
    // 构建步骤映射
    stepMap := make(map[string]model.WorkflowStep)
    for _, step := range workflow.Steps {
        stepMap[step.ID] = step
    }
    
    // 找到开始步骤
    currentStepID := workflow.Steps[0].ID
    
    for currentStepID != "" {
        step := stepMap[currentStepID]
        
        // 执行步骤
        result, err := e.executeStep(ctx, run, step)
        
        if err != nil {
            // 错误处理
            handleErr := e.handleStepError(run, step, err)
            if handleErr == ErrorActionAbort {
                run.Status = "failed"
                run.Error = err.Error()
                e.store.SaveRun(run)
                return
            } else if handleErr == ErrorActionRetry {
                continue
            }
        }
        
        // 确定下一步
        currentStepID = e.getNextStep(step, result)
        
        // 触发钩子
        e.hooks.Trigger(AfterStep, run, step, result)
    }
    
    // 完成
    run.Status = "completed"
    now := time.Now()
    run.CompletedAt = &now
    e.store.SaveRun(run)
    
    // 触发完成钩子
    e.hooks.Trigger(AfterWorkflow, run, nil, nil)
}

// executeStep 执行单个步骤
func (e *Engine) executeStep(ctx context.Context, run *model.WorkflowRun, step model.WorkflowStep) (any, error) {
    // 更新步骤状态
    e.updateStepStatus(run, step.ID, "running", nil, nil)
    
    // 获取Agent
    agent, ok := e.agents[step.Type]
    if !ok {
        return nil, fmt.Errorf("unknown agent type: %s", step.Type)
    }
    
    // 构建输入
    input := e.buildStepInput(run, step)
    
    // 执行
    output, err := agent.Execute(ctx, input)
    
    if err != nil {
        e.updateStepStatus(run, step.ID, "failed", input, err)
        return nil, err
    }
    
    // 更新成功状态
    e.updateStepStatus(run, step.ID, "completed", input, output)
    
    // 更新运行状态
    run.CurrentStep = step.ID
    run.Output = mergeOutput(run.Output, output)
    e.store.SaveRun(run)
    
    return output, nil
}
```

### 4.3 内置Agent类型

```go
// internal/workflow/agents/builtin.go

// HotFetchAgent 热点获取Agent
type HotFetchAgent struct{}
func (a *HotFetchAgent) Execute(ctx context.Context, input any) (any, error) {
    config := input.(HotFetchConfig)
    // 调用各平台API获取热点
    trends := fetchFromPlatforms(config.Platforms)
    return map[string]any{"trends": trends}, nil
}

// TopicGenerateAgent 选题生成Agent
type TopicGenerateAgent struct {
    topicAgent *agent.TopicAgent
}
func (a *TopicGenerateAgent) Execute(ctx context.Context, input any) (any, error) {
    data := input.(map[string]any)
    trends := data["trends"].([]model.HotTrend)
    
    topics, err := a.topicAgent.GenerateTopics(ctx, agent.TopicRequest{
        Limit: data["limit"].(int),
    })
    if err != nil {
        return nil, err
    }
    
    return map[string]any{"topics": topics}, nil
}

// TopicSelectAgent 选题选择Agent
type TopicSelectAgent struct{}
func (a *TopicSelectAgent) Execute(ctx context.Context, input any) (any, error) {
    data := input.(map[string]any)
    topics := data["topics"].([]agent.TopicResult)
    
    // 选择评分最高的选题
    var best *agent.TopicResult
    for _, t := range topics {
        if best == nil || t.Score > best.Score {
            best = &t
        }
    }
    
    return map[string]any{"selected_topic": best}, nil
}

// OutlineGenerateAgent 大纲生成Agent
type OutlineGenerateAgent struct {
    outlineAgent *agent.OutlineAgent
}
func (a *OutlineGenerateAgent) Execute(ctx context.Context, input any) (any, error) {
    data := input.(map[string]any)
    topic := data["selected_topic"].(*agent.TopicResult)
    
    outline, err := a.outlineAgent.GenerateOutline(ctx, agent.OutlineRequest{
        Title: topic.Title,
    })
    if err != nil {
        return nil, err
    }
    
    return map[string]any{"outline": outline}, nil
}

// ContentWriteAgent 内容创作Agent
type ContentWriteAgent struct {
    writingAgent *agent.WritingAgent
}
func (a *ContentWriteAgent) Execute(ctx context.Context, input any) (any, error) {
    data := input.(map[string]any)
    outline := data["outline"].([]model.OutlineNode)
    
    // 并行生成各章节
    var wg sync.WaitGroup
    contents := make([]string, len(outline))
    
    for i, node := range outline {
        wg.Add(1)
        go func(idx int, n model.OutlineNode) {
            defer wg.Done()
            content, _ := a.writingAgent.GenerateSection(ctx, n)
            contents[idx] = content
        }(i, node)
    }
    
    wg.Wait()
    
    // 组装文章
    article := assembleArticle(outline, contents)
    
    return map[string]any{
        "article": article,
        "outline": outline,
    }, nil
}

// ContentReviewAgent 内容审核Agent
type ContentReviewAgent struct {
    reviewAgent *agent.ReviewAgent
}
func (a *ContentReviewAgent) Execute(ctx context.Context, input any) (any, error) {
    data := input.(map[string]any)
    article := data["article"].(*model.Article)
    
    review, err := a.reviewAgent.Review(ctx, agent.ReviewRequest{
        Title:   article.Title,
        Content: article.Content,
    })
    if err != nil {
        return nil, err
    }
    
    return map[string]any{"review": review}, nil
}

// LayoutApplyAgent 排版Agent
type LayoutApplyAgent struct {
    layoutAgent *agent.LayoutAgent
}
func (a *LayoutApplyAgent) Execute(ctx context.Context, input any) (any, error) {
    data := input.(map[string]any)
    article := data["article"].(*model.Article)
    theme := data["theme"].(string)
    
    html, err := a.layoutAgent.Render(ctx, article, theme)
    if err != nil {
        return nil, err
    }
    
    return map[string]any{
        "html": html,
        "article": article,
    }, nil
}

// PublishDraftAgent 发布Agent
type PublishDraftAgent struct {
    publishAgent *agent.PublishAgent
}
func (a *PublishDraftAgent) Execute(ctx context.Context, input any) (any, error) {
    data := input.(map[string]any)
    article := data["article"].(*model.Article)
    html := data["html"].(string)
    
    mediaID, err := a.publishAgent.PublishDraft(ctx, agent.PublishDraftRequest{
        Title:   article.Title,
        Content: html,
    })
    if err != nil {
        return nil, err
    }
    
    return map[string]any{
        "media_id": mediaID,
        "published": true,
    }, nil
}
```

---

## 📱 五、前端界面设计

### 5.1 新增页面规划

```
frontend/src/views/
├── WelcomeView.vue          # 欢迎引导（保留）
├── EditorView.vue           # 编辑器（保留，增强）
├── SettingsView.vue         # 设置（保留）
├── PublishView.vue          # 发布（保留，增强）
│
├── workflow/                # 【新增】工作流
│   ├── WorkflowList.vue     # 工作流列表
│   ├── WorkflowEditor.vue   # 工作流编辑器
│   ├── WorkflowMonitor.vue  # 运行监控
│   └── WorkflowLogs.vue     # 运行日志
│
├── topic/                   # 【新增】选题中心
│   ├── TopicCenter.vue      # 选题中心首页
│   ├── HotTrends.vue        # 热点趋势
│   ├── TopicLibrary.vue     # 选题库
│   └── TopicDetail.vue      # 选题详情
│
├── publish/                 # 【新增】发布管理
│   ├── PublishQueue.vue     # 发布队列
│   ├── PublishHistory.vue   # 发布历史
│   ├── AccountManage.vue    # 账号管理
│   └── ScheduleCalendar.vue # 发布日历
│
└── dashboard/               # 【新增】数据看板
    ├── Dashboard.vue        # 总览
    ├── ArticleStats.vue     # 文章统计
    └── WorkflowStats.vue    # 工作流统计
```

### 5.2 工作流编辑器设计

```vue
<!-- WorkflowEditor.vue -->
<template>
  <div class="workflow-editor">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <el-input v-model="workflow.name" placeholder="工作流名称" />
      <el-button @click="save">保存</el-button>
      <el-button type="primary" @click="testRun">测试运行</el-button>
    </div>
    
    <!-- 左侧：步骤组件库 -->
    <div class="component-panel">
      <h4>步骤组件</h4>
      <draggable v-model="components" :group="{ name: 'steps', pull: 'clone', put: false }">
        <div v-for="c in components" :key="c.type" class="component-item">
          <el-icon><component :is="c.icon" /></el-icon>
          <span>{{ c.name }}</span>
        </div>
      </draggable>
    </div>
    
    <!-- 中间：流程画布 -->
    <div class="flow-canvas" ref="canvas">
      <vue-flow
        v-model="elements"
        :default-zoom="1"
        :min-zoom="0.2"
        :max-zoom="4"
        @connect="onConnect"
        @node-click="onNodeClick"
      >
        <!-- 自定义节点 -->
        <template #node-step="{ data }">
          <div class="step-node" :class="data.type">
            <el-icon><component :is="data.icon" /></el-icon>
            <span>{{ data.name }}</span>
            <div class="step-config" v-if="data.config">
              <el-tag size="small">{{ formatConfig(data.config) }}</el-tag>
            </div>
          </div>
        </template>
        
        <Background pattern-color="#aaa" :gap="16" />
        <MiniMap />
        <Controls />
      </vue-flow>
    </div>
    
    <!-- 右侧：属性面板 -->
    <div class="property-panel" v-if="selectedNode">
      <h4>步骤配置</h4>
      <component 
        :is="getConfigComponent(selectedNode.type)" 
        v-model="selectedNode.config"
      />
    </div>
    
    <!-- 底部：运行日志 -->
    <div class="log-panel" v-if="isRunning">
      <el-timeline>
        <el-timeline-item 
          v-for="log in runLogs" 
          :key="log.id"
          :type="log.status"
          :timestamp="log.time"
        >
          {{ log.message }}
        </el-timeline-item>
      </el-timeline>
    </div>
  </div>
</template>

<script setup>
import { VueFlow, useVueFlow } from '@vue-flow/core'
import '@vue-flow/core/dist/style.css'

// 步骤组件库
const components = [
  { type: 'hot_fetch', name: '获取热点', icon: 'TrendCharts' },
  { type: 'topic_generate', name: '生成选题', icon: 'Lightbulb' },
  { type: 'outline_generate', name: '生成大纲', icon: 'List' },
  { type: 'content_write', name: '创作内容', icon: 'Edit' },
  { type: 'content_review', name: '质量审核', icon: 'View' },
  { type: 'image_generate', name: '生成配图', icon: 'Picture' },
  { type: 'layout_apply', name: '公众号排版', icon: 'Document' },
  { type: 'publish_draft', name: '发布草稿', icon: 'Upload' },
  { type: 'condition', name: '条件判断', icon: 'Share' },
  { type: 'delay', name: '延迟等待', icon: 'Timer' },
  { type: 'manual_review', name: '人工审核', icon: 'User' },
]

// 流程元素
const elements = ref([])
const selectedNode = ref(null)
const isRunning = ref(false)
const runLogs = ref([])

// 添加节点
function onDrop(event) {
  const type = event.dataTransfer.getData('componentType')
  const position = getEventPosition(event)
  
  elements.value.push({
    id: generateId(),
    type: 'step',
    position,
    data: {
      type,
      name: getComponentName(type),
      icon: getComponentIcon(type),
      config: getDefaultConfig(type),
    },
  })
}

// 连接节点
function onConnect(connection) {
  elements.value.push({
    id: `e-${connection.source}-${connection.target}`,
    source: connection.source,
    target: connection.target,
    type: 'smoothstep',
  })
}

// 测试运行
async function testRun() {
  isRunning.value = true
  runLogs.value = []
  
  const workflow = buildWorkflowFromElements(elements.value)
  
  // 调用后端API启动工作流
  const run = await StartWorkflow(workflow)
  
  // 监听运行状态
  watchWorkflowRun(run.id, (update) => {
    runLogs.value.push({
      id: Date.now(),
      status: update.status,
      message: update.message,
      time: new Date().toLocaleTimeString(),
    })
  })
}
</script>
```

### 5.3 选题中心设计

```vue
<!-- TopicCenter.vue -->
<template>
  <div class="topic-center">
    <!-- 顶部筛选 -->
    <div class="filter-bar">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="AI推荐" name="ai">
          <AIRecommend />
        </el-tab-pane>
        <el-tab-pane label="热点趋势" name="hot">
          <HotTrends />
        </el-tab-pane>
        <el-tab-pane label="我的选题库" name="library">
          <TopicLibrary />
        </el-tab-pane>
      </el-tabs>
      
      <el-button type="primary" @click="refreshTopics">
        <el-icon><Refresh /></el-icon> 刷新选题
      </el-button>
    </div>
    
    <!-- AI推荐内容 -->
    <div v-if="activeTab === 'ai'" class="ai-recommend">
      <div class="recommend-header">
        <h3>为您推荐 {{ topics.length }} 个选题</h3>
        <p>基于热点趋势和您的公众号定位智能生成</p>
      </div>
      
      <el-row :gutter="20">
        <el-col :span="8" v-for="topic in topics" :key="topic.id">
          <el-card class="topic-card" :class="{ 'high-score': topic.score >= 80 }">
            <div class="topic-header">
              <el-tag v-if="topic.score >= 80" type="danger">强烈推荐</el-tag>
              <el-tag v-else-if="topic.score >= 70" type="warning">推荐</el-tag>
              <span class="score">{{ topic.score }}分</span>
            </div>
            
            <h4 class="topic-title">{{ topic.title }}</h4>
            
            <p class="topic-reason">{{ topic.reason }}</p>
            
            <div class="topic-angles">
              <h5>切入角度：</h5>
              <el-tag v-for="angle in topic.angles" :key="angle" size="small">
                {{ angle }}
              </el-tag>
            </div>
            
            <div class="topic-keywords">
              <el-tag v-for="kw in topic.keywords" :key="kw" type="info" size="small">
                {{ kw }}
              </el-tag>
            </div>
            
            <div class="topic-actions">
              <el-button type="primary" @click="createArticle(topic)">
                立即创作
              </el-button>
              <el-button @click="saveToLibrary(topic)">
                加入选题库
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>
```

---

## 🔌 六、微信公众号API集成

### 6.1 微信API封装

```go
// internal/platform/wechat/api.go
package wechat

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

const (
    BaseURL = "https://api.weixin.qq.com/cgi-bin"
)

// Client 微信API客户端
type Client struct {
    appID       string
    appSecret   string
    accessToken string
    expiry      time.Time
    httpClient  *http.Client
}

// NewClient 创建微信客户端
func NewClient(appID, appSecret string) *Client {
    return &Client{
        appID:      appID,
        appSecret:  appSecret,
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }
}

// ensureAccessToken 确保access_token有效
func (c *Client) ensureAccessToken(ctx context.Context) error {
    if c.accessToken != "" && time.Now().Before(c.expiry) {
        return nil
    }
    
    url := fmt.Sprintf("%s/token?grant_type=client_credential&appid=%s&secret=%s",
        BaseURL, c.appID, c.appSecret)
    
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    var result struct {
        AccessToken string `json:"access_token"`
        ExpiresIn   int    `json:"expires_in"`
        ErrCode     int    `json:"errcode"`
        ErrMsg      string `json:"errmsg"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return err
    }
    
    if result.ErrCode != 0 {
        return fmt.Errorf("wechat error: %s", result.ErrMsg)
    }
    
    c.accessToken = result.AccessToken
    c.expiry = time.Now().Add(time.Duration(result.ExpiresIn-300) * time.Second)
    
    return nil
}

// UploadImage 上传图片素材
func (c *Client) UploadImage(ctx context.Context, imageURL string) (string, error) {
    if err := c.ensureAccessToken(ctx); err != nil {
        return "", err
    }
    
    // 下载图片
    resp, err := http.Get(imageURL)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    // 构建multipart请求
    var buf bytes.Buffer
    writer := multipart.NewWriter(&buf)
    part, err := writer.CreateFormFile("media", "image.jpg")
    if err != nil {
        return "", err
    }
    
    _, err = io.Copy(part, resp.Body)
    if err != nil {
        return "", err
    }
    writer.Close()
    
    // 上传
    url := fmt.Sprintf("%s/media/uploadimg?access_token=%s", BaseURL, c.accessToken)
    uploadResp, err := c.httpClient.Post(url, writer.FormDataContentType(), &buf)
    if err != nil {
        return "", err
    }
    defer uploadResp.Body.Close()
    
    var result struct {
        URL     string `json:"url"`
        ErrCode int    `json:"errcode"`
        ErrMsg  string `json:"errmsg"`
    }
    
    if err := json.NewDecoder(uploadResp.Body).Decode(&result); err != nil {
        return "", err
    }
    
    if result.ErrCode != 0 {
        return "", fmt.Errorf("upload failed: %s", result.ErrMsg)
    }
    
    return result.URL, nil
}

// DraftArticle 草稿文章
type DraftArticle struct {
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

// AddDraft 添加草稿
func (c *Client) AddDraft(ctx context.Context, articles []DraftArticle) (string, error) {
    if err := c.ensureAccessToken(ctx); err != nil {
        return "", err
    }
    
    payload := map[string]any{
        "articles": articles,
    }
    
    body, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }
    
    url := fmt.Sprintf("%s/draft/add?access_token=%s", BaseURL, c.accessToken)
    resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    var result struct {
        MediaID string `json:"media_id"`
        ErrCode int    `json:"errcode"`
        ErrMsg  string `json:"errmsg"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }
    
    if result.ErrCode != 0 {
        return "", fmt.Errorf("add draft failed: %s", result.ErrMsg)
    }
    
    return result.MediaID, nil
}

// Publish 发布文章
func (c *Client) Publish(ctx context.Context, mediaID string) error {
    if err := c.ensureAccessToken(ctx); err != nil {
        return err
    }
    
    payload := map[string]any{
        "media_id": mediaID,
    }
    
    body, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    
    url := fmt.Sprintf("%s/freepublish/submit?access_token=%s", BaseURL, c.accessToken)
    resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    var result struct {
        PublishID int64 `json:"publish_id"`
        ErrCode   int   `json:"errcode"`
        ErrMsg    string `json:"errmsg"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return err
    }
    
    if result.ErrCode != 0 {
        return fmt.Errorf("publish failed: %s", result.ErrMsg)
    }
    
    return nil
}
```

### 6.2 发布服务

```go
// internal/service/publish_service.go
package service

// PublishService 发布服务
type PublishService struct {
    db          *repository.Repository
    wechatAPI   *wechat.Client
    taskQueue   chan *model.PublishTask
}

// NewPublishService 创建发布服务
func NewPublishService(db *repository.Repository) *PublishService {
    svc := &PublishService{
        db:        db,
        taskQueue: make(chan *model.PublishTask, 100),
    }
    
    // 启动发布协程
    go svc.processQueue()
    
    return svc
}

// SchedulePublish 调度发布任务
func (s *PublishService) SchedulePublish(ctx context.Context, task *model.PublishTask) error {
    // 保存任务
    if err := s.db.CreatePublishTask(task); err != nil {
        return err
    }
    
    // 根据调度类型处理
    switch task.ScheduleType {
    case "immediate":
        // 立即发布
        s.taskQueue <- task
        
    case "scheduled":
        // 定时发布 - 使用调度器
        s.scheduleAt(task)
    }
    
    return nil
}

// processQueue 处理发布队列
func (s *PublishService) processQueue() {
    for task := range s.taskQueue {
        ctx := context.Background()
        
        // 更新状态
        task.Status = "publishing"
        s.db.UpdatePublishTask(task)
        
        // 获取账号
        account, err := s.db.GetPublishAccount(task.AccountID)
        if err != nil {
            s.markFailed(task, fmt.Sprintf("get account failed: %v", err))
            continue
        }
        
        // 初始化微信客户端
        wechatClient := wechat.NewClient(account.AppID, account.AppSecret)
        
        // 上传封面图
        coverMediaID, err := s.uploadCover(ctx, wechatClient, task.CoverImage)
        if err != nil {
            s.markFailed(task, fmt.Sprintf("upload cover failed: %v", err))
            continue
        }
        
        // 处理内容中的图片
        content, err := s.processContentImages(ctx, wechatClient, task.Content)
        if err != nil {
            s.markFailed(task, fmt.Sprintf("process images failed: %v", err))
            continue
        }
        
        // 创建草稿
        mediaID, err := wechatClient.AddDraft(ctx, []wechat.DraftArticle{
            {
                Title:              task.Title,
                Content:            content,
                ThumbMediaID:       coverMediaID,
                Author:             "", // 可从配置读取
                Digest:             task.Summary,
                ShowCoverPic:       1,
                NeedOpenComment:    1,
                OnlyFansCanComment: 0,
            },
        })
        
        if err != nil {
            s.markFailed(task, fmt.Sprintf("create draft failed: %v", err))
            continue
        }
        
        task.PlatformID = mediaID
        task.Status = "published"
        now := time.Now()
        task.PublishedAt = &now
        s.db.UpdatePublishTask(task)
    }
}

// markFailed 标记失败
func (s *PublishService) markFailed(task *model.PublishTask, err string) {
    task.Status = "failed"
    task.Error = err
    s.db.UpdatePublishTask(task)
}
```

---

## 📅 七、实施路线图

### 7.1 阶段划分

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           实施路线图                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Phase 1: 基础架构 (2周)                                                      │
│ ├── 数据库迁移：新增工作流、选题、发布相关表                                    │
│ ├── 核心模型定义：Workflow, Topic, PublishTask                               │
│ ├── 工作流引擎基础：状态机、调度器、事件系统                                    │
│ └── Agent框架扩展：支持新Agent类型注册                                        │
│                              ↓                                              │
│ Phase 2: 选题系统 (1周)                                                      │
│ ├── TopicAgent实现：热点获取、选题生成、选题评分                               │
│ ├── 热点API接入：微博、知乎、百度热榜                                        │
│ ├── 选题中心前端：热点看板、AI推荐、选题库                                     │
│ └── 选题→工作流联动：选中选题触发创作流程                                      │
│                              ↓                                              │
│ Phase 3: 工作流引擎 (2周)                                                    │
│ ├── 工作流CRUD API：增删改查、YAML解析                                       │
│ ├── 工作流可视化编辑器：Vue Flow集成                                          │
│ ├── 内置Agent实现：全部11个步骤组件                                          │
│ ├── 工作流运行监控：实时日志、状态追踪                                        │
│ └── 定时调度：Cron表达式解析、任务调度                                        │
│                              ↓                                              │
│ Phase 4: 微信发布 (1周)                                                      │
│ ├── 微信API封装：access_token管理、素材上传、草稿发布                          │
│ ├── PublishAgent实现：草稿发布、定时发布                                      │
│ ├── 账号管理：多账号支持、Token自动刷新                                       │
│ ├── 发布队列：发布任务调度、失败重试                                          │
│ └── 发布管理前端：发布队列、发布历史、发布日历                                  │
│                              ↓                                              │
│ Phase 5: 编辑器增强 (1周)                                                    │
│ ├── 集成工作流：编辑器可启动/暂停工作流                                        │
│ ├── 人工审核点：工作流暂停等待人工确认                                        │
│ ├── 预览增强：支持查看工作流生成的文章                                        │
│ └── 公众号排版保留：现有功能迁移                                               │
│                              ↓                                              │
│ Phase 6: 测试优化 (1周)                                                      │
│ ├── 端到端测试：完整工作流测试                                                │
│ ├── 性能优化：并发控制、缓存优化                                              │
│ ├── 错误处理：重试机制、降级策略、告警通知                                     │
│ └── 文档完善：使用文档、API文档                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 详细任务分解

#### Phase 1: 基础架构 (2周)

| 任务 | 工时 | 负责人 | 产出 |
|------|------|--------|------|
| 数据库迁移脚本 | 1d | - | migration.sql |
| 数据模型定义 | 2d | - | model/workflow.go, model/topic.go, model/publish.go |
| 工作流引擎核心 | 3d | - | workflow/engine.go, workflow/store.go |
| 状态机实现 | 2d | - | workflow/state_machine.go |
| 事件系统 | 2d | - | workflow/events.go |
| Agent注册框架 | 2d | - | workflow/agent_registry.go |
| 单元测试 | 2d | - | *_test.go |

#### Phase 2: 选题系统 (1周)

| 任务 | 工时 | 产出 |
|------|------|------|
| 热点API接口定义 | 1d | platform/hot/ |
| 微博热榜采集 | 1d | platform/hot/weibo.go |
| 知乎热榜采集 | 1d | platform/hot/zhihu.go |
| TopicAgent实现 | 2d | agent/topic_agent.go |
| 选题中心API | 1d | api/topic_api.go |
| 选题中心前端 | 2d | views/topic/TopicCenter.vue |

#### Phase 3: 工作流引擎 (2周)

| 任务 | 工时 | 产出 |
|------|------|------|
| 工作流CRUD API | 2d | api/workflow_api.go |
| YAML解析器 | 1d | workflow/yaml_parser.go |
| Vue Flow集成 | 2d | views/workflow/WorkflowEditor.vue |
| 11个内置Agent | 5d | workflow/agents/*.go |
| 运行监控系统 | 2d | workflow/monitor.go, WorkflowMonitor.vue |
| Cron调度器 | 2d | workflow/scheduler.go |

#### Phase 4: 微信发布 (1周)

| 任务 | 工时 | 产出 |
|------|------|------|
| 微信API封装 | 2d | platform/wechat/api.go |
| Token管理 | 1d | platform/wechat/token.go |
| PublishAgent | 1d | agent/publish_agent.go |
| 发布服务 | 2d | service/publish_service.go |
| 发布管理前端 | 2d | views/publish/*.vue |

#### Phase 5: 编辑器增强 (1周)

| 任务 | 工时 | 产出 |
|------|------|------|
| 编辑器与工作流集成 | 2d | EditorView.vue 增强 |
| 人工审核点实现 | 2d | workflow/hooks/manual_review.go |
| 公众号排版迁移 | 2d | utils/wechatRenderer.ts 保留 |
| 预览功能保留 | 1d | PreviewPanel.vue 保留 |

#### Phase 6: 测试优化 (1周)

| 任务 | 工时 | 产出 |
|------|------|------|
| 端到端测试 | 2d | e2e/workflow_test.go |
| 性能测试 | 1d | benchmark/ |
| 错误处理完善 | 2d | error_handler.go |
| 文档编写 | 2d | docs/ |

---

## 📁 八、文件变更清单

### 8.1 新增文件

```
internal/
├── workflow/
│   ├── engine.go              # 工作流引擎核心
│   ├── store.go               # 工作流存储接口
│   ├── state_machine.go       # 状态机
│   ├── scheduler.go           # 调度器
│   ├── events.go              # 事件系统
│   ├── hooks.go               # 钩子系统
│   ├── yaml_parser.go         # YAML解析
│   ├── agent_registry.go      # Agent注册中心
│   └── agents/                # 内置Agent
│       ├── hot_fetch.go
│       ├── topic_generate.go
│       ├── topic_select.go
│       ├── outline_generate.go
│       ├── content_write.go
│       ├── content_review.go
│       ├── content_rewrite.go
│       ├── image_generate.go
│       ├── layout_apply.go
│       └── publish_draft.go
│
├── agent/
│   ├── topic_agent.go         # 选题Agent
│   ├── review_agent.go        # 审核Agent
│   ├── image_agent.go         # 图片Agent
│   ├── publish_agent.go       # 发布Agent
│   └── layout_agent.go        # 排版Agent
│
├── model/
│   ├── workflow.go            # 工作流模型
│   ├── topic.go               # 选题模型
│   └── publish.go             # 发布模型
│
├── platform/
│   ├── wechat/
│   │   ├── api.go             # 微信API
│   │   ├── token.go           # Token管理
│   │   └── types.go           # 类型定义
│   └── hot/
│       ├── types.go           # 热点类型
│       ├── weibo.go           # 微博热榜
│       ├── zhihu.go           # 知乎热榜
│       ├── baidu.go           # 百度热榜
│       └── toutiao.go         # 头条热榜
│
└── service/
    ├── workflow_service.go    # 工作流服务
    ├── topic_service.go       # 选题服务
    └── publish_service.go     # 发布服务

frontend/src/
├── views/
│   ├── workflow/
│   │   ├── WorkflowList.vue
│   │   ├── WorkflowEditor.vue
│   │   ├── WorkflowMonitor.vue
│   │   └── WorkflowLogs.vue
│   ├── topic/
│   │   ├── TopicCenter.vue
│   │   ├── HotTrends.vue
│   │   ├── TopicLibrary.vue
│   │   └── TopicDetail.vue
│   └── publish/
│       ├── PublishQueue.vue
│       ├── PublishHistory.vue
│       ├── AccountManage.vue
│       └── ScheduleCalendar.vue
│
├── stores/
│   ├── workflow.ts
│   ├── topic.ts
│   └── publish.ts
│
├── api/
│   ├── workflow.ts
│   ├── topic.ts
│   └── publish.ts
│
└── components/
    └── workflow/
        ├── StepNode.vue
        ├── FlowCanvas.vue
        └── PropertyPanel.vue
```

### 8.2 修改文件

```
internal/
├── model/
│   └── article.go             # 扩展Article模型
│
├── repository/
│   └── sqlite.go              # 新增表初始化
│
└── service/
    └── article_service.go     # 集成工作流

main.go                        # 注册新服务

frontend/src/
├── views/
│   └── EditorView.vue         # 集成工作流入口
│
├── router/
│   └── index.ts               # 新增路由
│
└── types/
    └── index.ts               # 新增类型定义
```

---

## ⚠️ 九、风险与应对

| 风险 | 影响 | 应对策略 |
|------|------|----------|
| 微信API审核不通过 | 高 | 提前申请测试账号，准备备用方案（生成复制内容） |
| 工作流引擎复杂度 | 中 | 采用迭代开发，先实现串行执行，再支持并行 |
| 热点API限制 | 中 | 接入多个数据源，实现降级机制 |
| AI生成质量不稳定 | 中 | 增加人工审核点，支持A/B测试和重写 |
| 现有功能兼容性 | 低 | 保持现有API不变，新增功能并行开发 |

---

## 📝 十、验收标准

### 10.1 功能验收

- [ ] 工作流可以正常创建、编辑、保存
- [ ] 工作流可以从YAML配置加载
- [ ] 工作流可以手动触发运行
- [ ] 工作流支持定时触发
- [ ] 选题中心可以获取并展示热点
- [ ] AI可以基于热点生成选题
- [ ] 工作流可以自动完成选题→创作→排版→发布
- [ ] 发布任务可以成功发布到微信公众号草稿箱
- [ ] 定时发布功能正常工作
- [ ] 运行日志可以实时查看

### 10.2 性能验收

- [ ] 工作流启动响应时间 < 1s
- [ ] 选题生成时间 < 10s
- [ ] 文章创作时间 < 60s（2000字）
- [ ] 并发工作流执行稳定（5个同时运行）
- [ ] 发布队列处理速度 > 1篇/分钟

### 10.3 兼容性验收

- [ ] 现有编辑器功能不受影响
- [ ] 现有文章数据可以正常访问
- [ ] 预览功能正常工作
- [ ] 公众号排版效果保持一致

---

## 🎯 总结

本次改版将 Content-Alchemist 从「单篇创作工具」升级为「AI 运营工作流引擎」，核心能力包括：

1. **自动选题**：基于热点趋势AI生成优质选题
2. **自动创作**：工作流驱动的大纲生成、内容创作、质量审核
3. **自动排版**：保留并优化公众号排版能力
4. **自动发布**：对接微信公众号API，支持定时发布

通过可视化的工作流编辑器，用户可以灵活配置自动化流程，在关键节点设置人工审核，实现"无人值守"的内容运营。
